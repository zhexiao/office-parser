package excel

func Convert(filepath string, _type string) interface{} {
	//得到excel数据
	e := Parser(filepath)

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
