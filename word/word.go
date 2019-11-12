package word

import (
	"bytes"
	"fmt"
	"gitee.com/zhexiao/unioffice/common"
	"gitee.com/zhexiao/unioffice/document"
	"github.com/zhexiao/mtef-go/eqn"
	"github.com/zhexiao/office-parser/utils"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type RowData struct {
	Content     []string
	HtmlContent []string
}

type TableData struct {
	Rows []*RowData
}

type Word struct {
	Uri    string
	Tables []*TableData
	doc    *document.Document

	//公式对象 RID:LATEX 的对应关系
	oles map[string]*string
	//无法转换的公式转换成图片
	olesImages map[string]string

	//图片 RID:七牛地址 的对应关系
	images map[string]string
}

func NewWord() *Word {
	return &Word{}
}

//直接读文件内容
func Read(r io.ReaderAt, size int64) *document.Document {
	doc, err := document.Read(r, size)
	if err != nil {
		log.Panic(err)
	}

	return doc
}

//打开文件内容
func Open(filepath string) *document.Document {
	doc, err := document.Open(filepath)
	if err != nil {
		log.Panic(err)
	}

	return doc
}

//解析word试题
func QuestionWord(doc *document.Document) *Word {
	//得到doc指针数据
	w := NewWord()
	w.doc = doc
	//解析公式和图片
	w.parseOle(w.doc.OleObjectPaths)
	w.parseImage(w.doc.Images)

	//todo 得到文档的所有公式WMF图片一次性解析
	//fmt.Println(w.doc.OleObjectWmfPath)

	//读取table数据
	w.getTableData()

	return w
}

//解析word试卷
func PaperWord(doc *document.Document) *Word {
	//得到doc指针数据
	w := NewWord()
	w.doc = doc
	//解析公式和图片
	w.parseOle(w.doc.OleObjectPaths)
	w.parseImage(w.doc.Images)

	return w
}

//得到纯解析的word文本数据
func (w *Word) getPureText() string {
	res := bytes.Buffer{}

	//p数据，段落自动编号当前值
	var (
		paragraphSortNum   int8
		paragraphSortNumId int64
	)
	for _, paragraph := range w.doc.Paragraphs() {
		var (
			//段落样式
			paragraphStyle string
			//段落自动编号应该呈现的值
			paragraphSortNumText string
		)

		//读取段落数据
		pString := w.getParagraphData(paragraph)

		//读取段落样式
		if paragraph.X().PPr != nil {
			//段落居中、居右
			if paragraph.X().PPr.Jc != nil {
				//if paragraph.X().PPr.Jc.ValAttr.String() == "center" {
				//	pString = fmt.Sprintf("<center>%s</center>", pString)
				//} else if paragraph.X().PPr.Jc.ValAttr.String() == "right" {
				//	pString = fmt.Sprintf("<span style='float:right;'>%s</span>", pString)
				//} else if paragraph.X().PPr.Jc.ValAttr.String() == "left" {
				//	pString = fmt.Sprintf("<span style='float:left;'>%s</span>", pString)
				//}

				paragraphStyle = fmt.Sprintf(" align='%s' ", paragraph.X().PPr.Jc.ValAttr.String())
			}

			//段落自动编号样式读取
			if paragraph.X().PPr.NumPr != nil {
				//初始化没有编号ID
				if paragraph.X().PPr.NumPr.NumId.ValAttr != paragraphSortNumId {
					//设置编号ID
					paragraphSortNumId = paragraph.X().PPr.NumPr.NumId.ValAttr

					//设置当前起始值为1
					paragraphSortNum = 1
				} else {
					//存在当前编号，当前值+1
					paragraphSortNum += 1
				}
			} else {
				//重置整个排序编号值
				paragraphSortNum = 0
			}

			if paragraphSortNum != 0 {
				switch paragraphSortNumId {
				//圆圈编号
				case 1000:
					switch paragraphSortNum {
					case 1:
						paragraphSortNumText = "①"
					case 2:
						paragraphSortNumText = "②"
					case 3:
						paragraphSortNumText = "③"
					case 4:
						paragraphSortNumText = "④"
					case 5:
						paragraphSortNumText = "⑤"
					case 6:
						paragraphSortNumText = "⑥"
					case 7:
						paragraphSortNumText = "⑦"
					case 8:
						paragraphSortNumText = "⑧"
					case 9:
						paragraphSortNumText = "⑨"
					case 10:
						paragraphSortNumText = "⑩"
					}
				//1. 2. 编号
				case 1100:
					paragraphSortNumText = fmt.Sprintf("%d.", paragraphSortNum)
				//(1) (2) 编号
				case 1200:
					paragraphSortNumText = fmt.Sprintf("(%d)", paragraphSortNum)
				default:
					paragraphSortNumText = fmt.Sprintf("(%d)", paragraphSortNum)
				}

				//写入自动编号
				pString = fmt.Sprintf("%s %s", paragraphSortNumText, pString)
			}
		}

		//写入段落样式
		pString = fmt.Sprintf("<p %s>%s</p>", paragraphStyle, pString)
		//pString = fmt.Sprintf("<p>%s</p>", pString)

		//保存内容
		res.WriteString(pString)
	}

	return res.String()
}

//读取段落数据
func (w *Word) getParagraphData(paragraph document.Paragraph) string {
	//存储run数据
	paragraphBuffer := bytes.Buffer{}

	//段落下面的每个单元文本数据
	for _, run := range paragraph.Runs() {
		//段落下面的每个单元文本数据
		var text string

		if run.DrawingInline() != nil {
			//图片数据
			text = w.readImage(run.DrawingInline())
		} else if run.OleObjects() != nil {
			//公式数据
			text = w.readOles(run.OleObjects())
		} else if len(run.Ruby().Rt) > 0 && len(run.Ruby().RubyBase) > 0 {
			//拼音数据
			if len(run.Ruby().Rt) != len(run.Ruby().RubyBase) {
				log.Println("拼音文本数据长度对不上")
			} else {
				for idx, rt := range run.Ruby().Rt {
					rubyText := run.Ruby().RubyBase[idx]
					text = fmt.Sprintf("<ruby>%s<rt>%s</rt></ruby>", rubyText, rt)
				}
			}
		} else {
			//	文本数据
			text = run.Text()

			//把空格替换成&nbsp;
			if strings.Contains(text, " ") {
				text = strings.Replace(text, " ", "&nbsp;", -1)
			}

			//检查文本样式
			//parser_underline_wave 波浪线
			//parser_em_zhuozhong	着重符
			if run.X().RPr != nil {
				//加粗
				if run.X().RPr.B != nil {
					text = fmt.Sprintf("<b>%s</b>", text)
				}

				//下划线、波浪线
				if run.X().RPr.U != nil {
					if run.X().RPr.U.ValAttr.String() == "single" {
						//下划线
						text = fmt.Sprintf("<u>%s</u>", text)
					} else {
						//波浪线
						text = fmt.Sprintf("<span class='parser_underline_%s'>%s</span>", run.X().RPr.U.ValAttr.String(), text)
					}
				}

				//斜体
				if run.X().RPr.I != nil {
					text = fmt.Sprintf("<span style='font-style:italic'>%s</span>", text)
				}

				//着重符号
				if run.X().RPr.Em != nil {
					text = fmt.Sprintf("<span class='em_zhuozhong'>%s</span>", text)
				}

				//颜色
				if run.X().RPr.Color != nil {
					text = fmt.Sprintf("<span style='color:#%s'>%s</span>", run.X().RPr.Color.ValAttr.String(), text)
				}

				//上标，下标
				if run.X().RPr.VertAlign != nil {
					switch run.X().RPr.VertAlign.ValAttr.String() {
					case "superscript":
						text = fmt.Sprintf("<sup>%s</sup>", text)
					case "subscript":
						text = fmt.Sprintf("<sub>%s</sub>", text)
					}
				}
			}
		}

		paragraphBuffer.WriteString(text)
	}

	return paragraphBuffer.String()
}

//读取表格数据
func (w *Word) getTableData() {
	tables := w.doc.Tables()
	for _, table := range tables {
		//读取一个表单里面的所有行
		rows := table.Rows()

		//读取行里面的数据
		tableData := w.getRowsData(&rows)
		w.Tables = append(w.Tables, &tableData)
	}
}

//读取所有行的数据
func (w *Word) getRowsData(rows *[]document.Row) TableData {
	var td TableData
	for _, row := range *rows {
		rowData := w.getRowText(&row)
		td.Rows = append(td.Rows, &rowData)
	}

	return td
}

//读取每一行的数据
func (w *Word) getRowText(row *document.Row) RowData {
	cells := row.Cells()
	rowData := RowData{}

	for _, cell := range cells {
		rawText, htmlText := w.getCellText(&cell)
		rowData.Content = append(rowData.Content, rawText)
		rowData.HtmlContent = append(rowData.HtmlContent, htmlText)
	}

	return rowData
}

//读取行里面每一个单元的数据
func (w *Word) getCellText(cell *document.Cell) (string, string) {
	paragraphs := cell.Paragraphs()

	resText := bytes.Buffer{}
	htmlResText := bytes.Buffer{}

	//循环每一个P标签数据
	for paragIdx, paragraph := range paragraphs {
		runs := paragraph.Runs()

		for _, r := range runs {
			var text string

			//图片数据
			if r.DrawingInline() != nil {
				text = w.readImage(r.DrawingInline())
				//	公式数据
			} else if r.OleObjects() != nil {
				text = w.readOles(r.OleObjects())
				//	文本数据
			} else {
				text = r.Text()
			}

			resText.WriteString(text)
			htmlResText.WriteString(text)
		}

		//新的段落换行
		if paragIdx < len(paragraphs)-1 {
			htmlResText.WriteString("<br/>")
		}
	}

	return resText.String(), htmlResText.String()
}

//读取图片
func (w *Word) readImage(images []document.InlineDrawing) string {
	var imageUri string
	for _, di := range images {
		imf, _ := di.GetImage()
		uri := w.images[imf.RelID()]

		imageUri = fmt.Sprintf("<img src='%s' style='width:%s;height:%s'/>", uri, di.X().Extent.Size().Width, di.X().Extent.Size().Height)
	}

	return imageUri
}

//读取公式
func (w *Word) readOles(ole []document.OleObject) string {
	var latexStr string
	for _, ole := range ole {
		latexPtr, ok := w.oles[ole.OleRid()]
		if ok {
			latexStr = *latexPtr
		} else {
			oleImg, ok := w.olesImages[ole.ImagedataRid()]
			if ok {
				latexStr = fmt.Sprintf("<img src='%s' style='%s' />", oleImg, *ole.Shape().StyleAttr)
			}
		}
	}

	return latexStr
}

//把ole对象文件转为latex字符串
func (w *Word) parseOle(olePaths []document.OleObjectPath) {
	w.oles = make(map[string]*string)
	w.olesImages = make(map[string]string)

	//使用 WaitGroup 来跟踪 goroutine 的工作是否完成
	var wg sync.WaitGroup
	wg.Add(len(olePaths))

	//循环数据
	var mutex sync.Mutex
	for _, ole := range olePaths {
		//goroutine 运行
		go func(word *Word, oleObjPath document.OleObjectPath) {
			// 在函数退出时调用 Done
			defer wg.Done()

			//调用解析库解析公式
			latex := eqn.Convert(oleObjPath.Path())
			if latex == "" {
				//无法解析的公式，转图片
				wmfObj := w.doc.OleObjectWmfPath[0]
				imageName := fmt.Sprintf("%s_%s", strconv.Itoa(int(time.Now().UnixNano())), wmfObj.Rid())
				utils.WmfConvert(wmfObj.Path(), imageName)

				//map并发问题
				mutex.Lock()
				word.olesImages[wmfObj.Rid()] = fmt.Sprintf("%s/%s/%s.jpg", utils.OfficeParserQiniuCfg.Domain, utils.OfficeParserQiniuCfg.UploadPrefix, imageName)
				mutex.Unlock()
			} else {
				//成功解析的公式，替换$$为[ 或 ]
				latex = strings.Replace(latex, "$$", "[", 1)
				latex = strings.Replace(latex, "$$", "]", 1)

				//map并发问题
				mutex.Lock()
				word.oles[oleObjPath.Rid()] = &latex
				mutex.Unlock()
			}
		}(w, ole)
	}

	wg.Wait()
}

//把图片上传到七牛
func (w *Word) parseImage(images []common.ImageRef) {
	w.images = make(map[string]string)

	//使用 WaitGroup 来跟踪 goroutine 的工作是否完成
	var wg sync.WaitGroup
	wg.Add(len(images))

	var mutex sync.Mutex
	for _, img := range images {
		//goroutine 运行
		go func(word *Word, image common.ImageRef) {
			// 在函数退出时调用 Done
			defer wg.Done()

			//调用图片上传
			localFile := image.Path()
			key := fmt.Sprintf("%s_%s.%s", strconv.Itoa(int(time.Now().UnixNano())), image.RelID(), image.Format())

			//上传到七牛
			uri := utils.UploadFileToQiniu(key, localFile)

			//map并发问题
			mutex.Lock()
			word.images[image.RelID()] = uri
			mutex.Unlock()
		}(w, img)
	}

	wg.Wait()
}
