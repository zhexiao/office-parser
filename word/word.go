package word

import (
	"bytes"
	"github.com/unidoc/unioffice/document"
	"log"
)

type RowData struct {
	Content []string
}

type Word struct {
	Uri     string
	Content []*RowData
}

//解析word
func (w *Word) Parser(filepath string) {
	doc, err := document.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	w.getTableData(doc)
}

//读取表单数据
func (w *Word) getTableData(doc *document.Document) {
	tables := doc.Tables()
	for _, table := range tables {
		rows := table.Rows()

		for _, row := range rows {
			rowData := w.getRowText(row)
			w.Content = append(w.Content, &rowData)
		}
	}
}

//读取每一行的数据
func (w *Word) getRowText(row document.Row) RowData {
	cells := row.Cells()
	rowData := RowData{}

	for _, cell := range cells {
		text := w.getCellText(cell)
		rowData.Content = append(rowData.Content, text)
	}

	return rowData
}

//读取每一个单元的数据
func (w *Word) getCellText(cell document.Cell) string {
	paras := cell.Paragraphs()

	resText := bytes.Buffer{}
	for _, p := range paras {
		runs := p.Runs()

		for _, r := range runs {
			text := r.Text()
			resText.WriteString(text)
		}
	}

	return resText.String()
}
