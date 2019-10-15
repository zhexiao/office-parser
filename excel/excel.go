package excel

import (
	"github.com/unidoc/unioffice/spreadsheet"
	"log"
)

type RowData struct {
	Content []string
}

type Excel struct {
	Uri      string
	excel    *spreadsheet.Workbook
	RowsData []*RowData
}

func NewExcel() *Excel {
	return &Excel{}
}

func NewRowData() *RowData {
	return &RowData{}
}

//解析excel
func Parser(filepath string) *Excel {
	workbook, err := spreadsheet.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	//解析word
	e := NewExcel()
	e.Uri = filepath
	e.excel = workbook

	//读取第一个sheet的数据
	e.getSheetData(0)

	return e
}

//读取第n个sheet表
func (e *Excel) getSheet(n int) *spreadsheet.Sheet {
	sheetCount := e.excel.SheetCount()
	if n > sheetCount {
		log.Fatal("传入的数字大于excel最大的sheet")
	}

	sheets := e.excel.Sheets()
	return &sheets[n]
}

//读取sheet表的所有rows
func (e *Excel) getSheetData(n int) {
	sheet := e.getSheet(n)

	e.getRowsData(sheet.Rows())
}

//读取多行里面的数据
func (e *Excel) getRowsData(rows []spreadsheet.Row) {
	for _, row := range rows {
		rowData := e.getRowData(row.Cells())
		e.RowsData = append(e.RowsData, rowData)
	}
}

//读取每一行的数据
func (e *Excel) getRowData(cells []spreadsheet.Cell) *RowData {
	rowData := NewRowData()
	for _, cell := range cells {
		rowData.Content = append(rowData.Content, e.getCellData(cell))
	}

	return rowData
}

//读取单元格里面的数据值
func (e *Excel) getCellData(cell spreadsheet.Cell) string {
	return cell.GetString()
}
