package excel

func Convert(filepath string, _type string) {
	//得到excel数据
	e := Parser(filepath)

	switch _type {
	case "paper":
		//根据 excel 数据解析试卷
		ParsePaper(e)
	}
}
