package word

import (
	"bytes"
	"fmt"
	"gitee.com/zhexiao/unioffice/common"
	"gitee.com/zhexiao/unioffice/document"
	"github.com/zhexiao/mtef-go/eqn"
	"github.com/zhexiao/office-parser/bases"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

type CT_AbNumILvl struct {
	ilvl   int64
	numFmt string
	text   string
	start  int64
}

type CT_AbNumData struct {
	numId int64
	iLvl  []*CT_AbNumILvl
}

type CT_RowData struct {
	Content     []string
	HtmlContent []string
}

type CT_TableData struct {
	Rows []*CT_RowData
}

type CT_Word struct {
	Uri    string
	Tables []*CT_TableData
	doc    *document.Document

	//公式对象 RID:LATEX 的对应关系
	//oles map[string]*string
	oles *sync.Map

	//无法转换的公式转换成图片
	//olesImages map[string]string
	olesImages *sync.Map

	//图片 RID:七牛地址 的对应关系
	//images map[string]string
	images *sync.Map

	//自动序号相关
	numIdMapAbNumId map[int64]int64
	numData         []*CT_AbNumData
}

func NewCT_Word() *CT_Word {
	return &CT_Word{}
}

func (w *CT_Word) read(fileByte []byte) error {
	doc, err := document.Read(bytes.NewReader(fileByte), int64(len(fileByte)))
	if err != nil {
		return bases.NewOpError(bases.LibError, err.Error())
	}

	w.doc = doc
	return nil
}

//解析整个word数据
func (w *CT_Word) getWordData() (string, error) {
	//解析公式和图片
	w.parseOle()
	w.parseImage()

	//解析自动序号
	w.parseOrder()

	//得到word的所有文本信息
	data, err := w.getPureText()
	if err != nil {
		return "", err
	}

	return data, nil
}

//解析word表格数据
func (w *CT_Word) getWordTableData() ([]*CT_TableData, error) {
	//解析公式和图片
	w.parseOle()
	w.parseImage()

	//读取table数据
	w.getTableData()

	return w.Tables, nil
}

//把ole对象文件转为latex字符串
func (w *CT_Word) parseOle() {
	w.oles = &sync.Map{}
	w.olesImages = &sync.Map{}

	//使用 WaitGroup 来跟踪 goroutine 的工作是否完成
	var wg sync.WaitGroup
	wg.Add(len(w.doc.OleObjectPaths))

	for _, ole := range w.doc.OleObjectPaths {
		//goroutine 运行
		go func(word *CT_Word, oleObjPath document.OleObjectPath) {
			defer wg.Done()

			//调用解析库解析公式
			latex := eqn.Convert(oleObjPath.Path())
			if latex == "" {
				//无法解析的公式，转图片
				wmfObj := w.doc.OleObjectWmfPath[0]
				imageName := fmt.Sprintf("%s_%s", strconv.Itoa(int(time.Now().UnixNano())), wmfObj.Rid())
				bases.WmfConvert(wmfObj.Path(), imageName)

				word.olesImages.Store(wmfObj.Rid(), fmt.Sprintf("%s/%s/%s.jpg", bases.OpQiniu.Domain, bases.OpQiniu.UploadPrefix, imageName))
			} else {
				//成功解析的公式，替换$$为[ 或 ]
				latex = strings.Replace(latex, "$$", "\\[", 1)
				latex = strings.Replace(latex, "$$", "\\]", 1)

				w.oles.Store(oleObjPath.Rid(), &latex)
			}
		}(w, ole)
	}

	wg.Wait()
}

//把图片上传到七牛
func (w *CT_Word) parseImage() {
	w.images = &sync.Map{}

	//使用 WaitGroup 来跟踪 goroutine 的工作是否完成
	var wg sync.WaitGroup
	wg.Add(len(w.doc.Images))

	for _, img := range w.doc.Images {
		//goroutine 运行
		go func(word *CT_Word, image common.ImageRef) {
			defer wg.Done()

			//调用图片上传
			localFile := image.Path()
			key := fmt.Sprintf("%s_%s.%s", strconv.Itoa(int(time.Now().UnixNano())), image.RelID(), image.Format())

			//上传到七牛
			uri, _ := bases.UploadFileToQiniu(key, localFile)

			word.images.Store(image.RelID(), uri)
		}(w, img)
	}

	wg.Wait()
}

//执行自动序号数据读取
func (w *CT_Word) parseOrder() {
	if w.doc.Numbering.X() != nil {
		//读取序号数据
		for _, df := range w.doc.Numbering.Definitions() {
			abData := &CT_AbNumData{}
			abData.numId = df.AbstractNumberID()
			for _, lv := range df.X().Lvl {
				abData.iLvl = append(abData.iLvl, &CT_AbNumILvl{
					ilvl:   lv.IlvlAttr,
					numFmt: lv.NumFmt.ValAttr.String(),
					text:   *lv.LvlText.ValAttr,
					start:  lv.Start.ValAttr,
				})
			}

			w.numData = append(w.numData, abData)
		}

		//numId与abstractNumId的映射关系
		numIdMapAbNumId := make(map[int64]int64)
		for _, nu := range w.doc.Numbering.X().Num {
			numIdMapAbNumId[nu.NumIdAttr] = nu.AbstractNumId.ValAttr
		}

		w.numIdMapAbNumId = numIdMapAbNumId
	}
}

//得到纯解析的word文本数据
func (w *CT_Word) getPureText() (string, error) {
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
				//fmt.Println(pString)
				//fmt.Printf("%#v \n", paragraph.X().PPr.Jc)
				//fmt.Println("====================================")
				paragraphStyle = fmt.Sprintf(" align='%s' ", paragraph.X().PPr.Jc.ValAttr.String())
			}

			//段落自动编号样式读取
			//参考文档：http://c-rex.net/projects/samples/ooxml/e1/Part4/OOXML_P4_DOCX_Numbering_topic_ID0EN6IU.html
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
				ivlData, err := w.readAbNumData(paragraphSortNumId, 0)
				if err != nil {
					return "", err
				}

				numFmt := ivlData.numFmt
				numText := ivlData.text

				var numVal string
				if numFmt == "decimal" {
					numVal = NUM_Decimal(paragraphSortNum).String()
				} else if numFmt == "decimalEnclosedCircle" {
					numVal = NUM_DecimalEnclosedCircle(paragraphSortNum).String()
				} else if numFmt == "japaneseCounting" || numFmt == "chineseCountingThousand" {
					numVal = NUM_Counting(paragraphSortNum).String()
				} else if numFmt == "upperLetter" {
					numVal = NUM_UpperLetter(paragraphSortNum).String()
				} else if numFmt == "upperRoman" {
					numVal = NUM_UpperRoman(paragraphSortNum).String()
				} else {
					numVal = string(paragraphSortNum)
					log.Printf("暂时不支持的自动序号,numFmt=%s,text=%s", numFmt, numText)
				}

				//替换数据
				paragraphSortNumText = strings.Replace(numText, "%1", numVal, -1)

				//写入自动编号
				pString = fmt.Sprintf("%s %s", paragraphSortNumText, pString)
			}

			//段落缩进
			//https://docs.microsoft.com/zh-cn/dotnet/api/documentformat.openxml.wordprocessing.indentation?view=openxml-2.8.1
			if paragraph.X().PPr.Ind != nil {
				if paragraph.X().PPr.Ind.FirstLineCharsAttr != nil {
					indentNum := int(math.Round(float64(*(paragraph.X().PPr.Ind.FirstLineCharsAttr)) / 50))
					indentNbsp := strings.Repeat("&nbsp;", indentNum)

					//fmt.Println(paragraph.X().PPr.Ind.FirstLineCharsAttr)
					pString = fmt.Sprintf("%s%s", indentNbsp, pString)
				}
			}
		}

		//写入段落样式
		pString = fmt.Sprintf("<p %s>%s</p>", paragraphStyle, pString)

		//保存内容
		res.WriteString(pString)
	}

	return res.String(), nil
}

//读取段落数据
func (w *CT_Word) getParagraphData(paragraph document.Paragraph) string {
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

					if run.X().RPr != nil {
						//加粗
						if run.X().RPr.B != nil {
							rubyText = fmt.Sprintf("<b>%s</b>", rubyText)
						}
					}

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
			//parser_underline_wavyDouble 双波浪线
			//parser_em_zhuozhong	着重符
			if run.X().RPr != nil {
				//删除线
				if run.X().RPr.Strike != nil {
					text = fmt.Sprintf("<del>%s</del>", text)
				}

				//背景色
				if run.X().RPr.Highlight != nil {
					text = fmt.Sprintf("<span style='display:inline-block;background-color:%s;'>%s</span>", run.X().RPr.Highlight.ValAttr.String(), text)
				}

				//加粗
				if run.X().RPr.B != nil {
					text = fmt.Sprintf("<b>%s</b>", text)
				}

				//下划线、波浪线
				if run.X().RPr.U != nil {
					uVal := run.X().RPr.U.ValAttr.String()

					switch uVal {
					case "single":
						//下划线
						text = fmt.Sprintf("<span style='border-bottom:1px solid black;'>%s</span>", text)
					case "double":
						//双下划线
						text = fmt.Sprintf("<span style='border-bottom:3px double black;'>%s</span>", text)
					default:
						text = fmt.Sprintf("<span class='parser_underline_%s'>%s</span>", run.X().RPr.U.ValAttr.String(), text)
					}
				}

				//斜体
				if run.X().RPr.I != nil {
					text = fmt.Sprintf("<span style='font-style:italic'>%s</span>", text)
				}

				//着重符号
				if run.X().RPr.Em != nil {
					text = fmt.Sprintf("<span class='parser_em_zhuozhong'>%s</span>", text)
				}

				//颜色
				if run.X().RPr.Color != nil {
					colorVal := run.X().RPr.Color.ValAttr.String()

					//todo 默认不处理黑色背景，word组卷不希望出现
					if colorVal != "000" && colorVal != "000000" {
						text = fmt.Sprintf("<span style='color:#%s'>%s</span>", colorVal, text)
					}
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

//读取图片数据
func (w *CT_Word) readImage(images []document.InlineDrawing) string {
	var imageUri string
	for _, di := range images {
		imf, _ := di.GetImage()
		uri, _ := w.images.Load(imf.RelID())

		imageUri = fmt.Sprintf("<img src='%s' style='width:%s;height:%s'/>", uri, di.X().Extent.Size().Width, di.X().Extent.Size().Height)
	}

	return imageUri
}

//读取公式数据
func (w *CT_Word) readOles(ole []document.OleObject) string {
	var latexStr string
	for _, ole := range ole {
		latexPtr, ok := w.oles.Load(ole.OleRid())
		if ok {
			latexStr = *latexPtr.(*string)
		} else {
			oleImg, ok := w.olesImages.Load(ole.ImagedataRid())
			if ok {
				latexStr = fmt.Sprintf("<img src='%s' style='%s' />", oleImg, *ole.Shape().StyleAttr)
			}
		}
	}

	return latexStr
}

//读取自动序号的数据
func (w *CT_Word) readAbNumData(numId int64, ilvl int64) (*CT_AbNumILvl, error) {
	abNumId, ok := w.numIdMapAbNumId[numId]
	if !ok {
		return nil, bases.NewOpError(bases.LibError, fmt.Sprintf("自动序号解析失败，找不到numId=%d", numId))
	}

	//读取 abstractNum 数据
	var tmpAbData *CT_AbNumData
	for _, abData := range w.numData {
		abDataNumId := abData.numId

		if abDataNumId == abNumId {
			tmpAbData = abData
		}
	}

	if tmpAbData == nil {
		return nil, bases.NewOpError(bases.LibError, fmt.Sprintf("找不到AbNum实例数据，abNumId=%d", abNumId))
	}

	//读取 lvl 数据
	var tmpAbLvl *CT_AbNumILvl
	for _, abLvl := range tmpAbData.iLvl {
		abLvlVal := abLvl.ilvl

		if abLvlVal == ilvl {
			tmpAbLvl = abLvl
		}
	}

	if tmpAbLvl == nil {
		return nil, bases.NewOpError(bases.LibError, fmt.Sprintf("找不到ilvl实例数据，ilvl=%d", ilvl))
	}

	return tmpAbLvl, nil
}

//读取表格数据
func (w *CT_Word) getTableData() {
	tables := w.doc.Tables()
	for _, table := range tables {
		//读取一个表单里面的所有行
		rows := table.Rows()

		//读取行里面的数据
		td := &CT_TableData{}
		for _, row := range rows {
			cells := row.Cells()
			rowData := &CT_RowData{}

			for _, cell := range cells {
				rawText, htmlText := w.getCellText(&cell)
				rowData.Content = append(rowData.Content, rawText)
				rowData.HtmlContent = append(rowData.HtmlContent, htmlText)
			}

			td.Rows = append(td.Rows, rowData)
		}

		w.Tables = append(w.Tables, td)
	}
}

//读取行里面每一个单元的数据
func (w *CT_Word) getCellText(cell *document.Cell) (string, string) {
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
