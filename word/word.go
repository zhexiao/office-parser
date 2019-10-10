package word

import (
	"bytes"
	"fmt"
	"github.com/unidoc/unioffice/document"
	"log"
)

type RowData struct {
	Content []string
}

type Word struct {
	Uri     string
	Content []*RowData
	doc     *document.Document
}

//解析word
func (w *Word) Parser(filepath string) {
	doc, err := document.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	//得到doc指针数据
	w.doc = doc

	//todo 得到文档的所有图片一次性上传到七牛
	//fmt.Println(w.doc.Images)

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
					fmt.Println(imf.Path())
				}
			} else if r.OleObjects() != nil {
				for _, ole := range r.OleObjects() {
					fmt.Printf("%#v", *ole.OleObject().IdAttr)
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
