package excel

import (
	"bytes"
	"fmt"
	"gitee.com/zhexiao/unioffice/spreadsheet"
	"github.com/zhexiao/office-parser/bases"
	"strconv"
)

type CT_RowData struct {
	Content []string
}

type CT_Excel struct {
	Uri      string
	excel    *spreadsheet.Workbook
	RowsData []*CT_RowData
}

func NewCT_Excel() *CT_Excel {
	return &CT_Excel{}
}

func NewCT_RowData() *CT_RowData {
	return &CT_RowData{}
}

//解析excel
func (e *CT_Excel) GetExcelData() (*CT_Excel, error) {
	//读取第一个sheet的数据
	if err := e.getSheetData(0); err != nil {
		return nil, err
	}

	return e, nil
}

//读取excel
func (e *CT_Excel) read(fileBytes []byte) error {
	workbook, err := spreadsheet.Read(bytes.NewReader(fileBytes), int64(len(fileBytes)))
	if err != nil {
		return bases.NewOpError(bases.NormalError, err.Error())
	}

	e.excel = workbook
	return nil
}

//读取sheet表的所有rows
func (e *CT_Excel) getSheetData(n int) error {
	//读取sheet
	sheetCount := e.excel.SheetCount()
	if n > sheetCount {
		return bases.NewOpError(bases.NormalError, "传入的数字大于excel最大的sheet")
	}
	sheets := e.excel.Sheets()
	sheet := sheets[n]

	//读取行列数
	colStart, rowStart, colEnd, rowEnd := sheet.ExtentsIndex()
	colStartInt := []byte(colStart)[0]
	colEndInt := []byte(colEnd)[0]

	//读取每一行的数据
	for i := rowStart; i <= rowEnd; i++ {
		rowData := NewCT_RowData()

		//读取每个单元的数据
		for n := colStartInt; n <= colEndInt; n++ {
			cellRef := fmt.Sprintf("%s%s", string([]byte{n}), strconv.Itoa(int(i)))
			val := sheet.Cell(cellRef).GetString()

			rowData.Content = append(rowData.Content, val)
		}

		e.RowsData = append(e.RowsData, rowData)
	}

	return nil
}
