package excel

import (
	"fmt"
	"gitee.com/zhexiao/unioffice/spreadsheet"
	"io"
	"log"
	"strconv"
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

//直接读文件内容
func Read(r io.ReaderAt, size int64) *Excel {
	workbook, err := spreadsheet.Read(r, size)
	if err != nil {
		log.Fatal(err)
	}

	return parser(workbook)
}

//打开文件内容
func Open(filepath string) *Excel {
	workbook, err := spreadsheet.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	return parser(workbook)
}

//解析excel
func parser(workbook *spreadsheet.Workbook) *Excel {
	//解析word
	e := NewExcel()
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

	rowStart, colStart, rowEnd, colEnd := sheet.ExtentsIndex()
	rowStartInt := []byte(rowStart)[0]
	rowEndInt := []byte(rowEnd)[0]

	//读取每一行的数据
	for i := colStart; i <= colEnd; i++ {
		rowData := NewRowData()

		//读取每个单元的数据
		for n := rowStartInt; n <= rowEndInt; n++ {
			cellRef := fmt.Sprintf("%s%s", string([]byte{n}), strconv.Itoa(int(i)))
			val := sheet.Cell(cellRef).GetString()

			rowData.Content = append(rowData.Content, val)
		}

		e.RowsData = append(e.RowsData, rowData)
	}
}
