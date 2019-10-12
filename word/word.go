package word

import (
	"bytes"
	"fmt"
	"github.com/unidoc/unioffice/common"
	"github.com/unidoc/unioffice/document"
	"github.com/zhexiao/mtef-go/eqn"
	"log"
	"office-parser/utils"
	"strconv"
	"time"
)

type RowData struct {
	Content []string
}

type Word struct {
	Uri     string
	Content []*RowData
	doc     *document.Document

	//公式对象 RID:LATEX 的对应关系
	oles map[string]*string
	//图片 RID:七牛地址 的对应关系
	images map[string]string
}

//解析word
func (w *Word) Parser(filepath string) {
	doc, err := document.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	//得到doc指针数据
	w.doc = doc
	w.parseOle(w.doc.OleObjectPaths)
	w.parseImage(w.doc.Images)

	//todo 得到文档的所有公式一次性解析
	//fmt.Println(w.doc.OleObjectWmfPath)

	//读取table数据
	w.getTableData()
}

//读取表单数据
func (w *Word) getTableData() {
	tables := w.doc.Tables()
	for _, table := range tables {
		rows := table.Rows()
		w.getRowsData(&rows)
	}
}

//读取所有行的数据
func (w *Word) getRowsData(rows *[]document.Row) {
	for _, row := range *rows {
		rowData := w.getRowText(&row)
		w.Content = append(w.Content, &rowData)
	}
}

//读取每一行的数据
func (w *Word) getRowText(row *document.Row) RowData {
	cells := row.Cells()
	rowData := RowData{}

	for _, cell := range cells {
		cellText := w.getCellText(&cell)
		rowData.Content = append(rowData.Content, cellText)
	}

	return rowData
}

//读取行里面每一个单元的数据
func (w *Word) getCellText(cell *document.Cell) string {
	paras := cell.Paragraphs()

	resText := bytes.Buffer{}
	for _, p := range paras {
		runs := p.Runs()

		for _, r := range runs {
			var text string

			//图片数据
			if r.DrawingInline() != nil {
				for _, di := range r.DrawingInline() {
					imf, _ := di.GetImage()
					uri := w.images[imf.RelID()]
					fmt.Println(uri)
				}
			} else if r.OleObjects() != nil {
				for _, ole := range r.OleObjects() {
					latex := w.oles[ole.OleRid()]
					fmt.Println(*latex)
				}
			} else {
				//文本数据
				text = r.Text()
			}

			resText.WriteString(text)
		}
	}

	return resText.String()
}

//把ole对象文件转为latex字符串
func (w *Word) parseOle(olePaths []document.OleObjectPath) {
	w.oles = make(map[string]*string)

	for _, ole := range olePaths {
		//调用解析库解析公式
		latex := eqn.Convert(ole.Path())
		w.oles[ole.Rid()] = &latex
	}
}

//把图片上传到七牛
func (w *Word) parseImage(images []common.ImageRef) {
	w.images = make(map[string]string)

	for _, img := range images {
		localFile := img.Path()
		key := fmt.Sprintf("%s.%s", strconv.Itoa(int(time.Now().UnixNano())), img.Format())

		//上传到七牛
		uri := utils.UploadFileToQiniu(key, localFile)
		w.images[img.RelID()] = uri
	}
}
