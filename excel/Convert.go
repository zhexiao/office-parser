package excel

import "io"

func ConvertFromFile(filepath string, _type string) interface{} {
	//得到excel数据
	e := Open(filepath)
	return convert(e, _type)

}

func ConvertFromData(r io.ReaderAt, size int64, _type string) interface{} {
	//得到excel数据
	e := Read(r, size)
	return convert(e, _type)
}

func convert(e *Excel, _type string) interface{} {
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
	default:
		return nil
	}
}
