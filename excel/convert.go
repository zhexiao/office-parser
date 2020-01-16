package excel

import (
	"github.com/zhexiao/office-parser/bases"
	"io/ioutil"
)

func ConvertFromFile(filepath string, _type string) (interface{}, error) {
	fileBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, bases.NewOpError(bases.NormalError, err.Error())
	}

	return ConvertFromData(fileBytes, _type)
}

func ConvertFromData(fileBytes []byte, _type string) (interface{}, error) {
	ct := NewCT_Excel()
	err := ct.read(fileBytes)
	if err != nil {
		return nil, bases.NewOpError(bases.NormalError, err.Error())
	}

	excelData, err := ct.GetExcelData()
	if err != nil {
		return nil, bases.NewOpError(bases.NormalError, err.Error())
	}

	data, err := convert(excelData, _type)
	if err != nil {
		return nil, bases.NewOpError(bases.NormalError, err.Error())
	}

	return data, nil
}

func convert(e *CT_Excel, _type string) (interface{}, error) {
	switch _type {
	case "paper":
		return ParsePaper(e)
	case "book":
		return ParseBook(e)
	case "cognition_map":
		return ParseCognitionMap(e)
	case "cognition_sp":
		return ParseCognitionSp(e)
	case "outline":
		return ParseOutline(e)
	}

	return nil, bases.NewOpError(bases.NormalError, "excel支持的类型不存在")
}
