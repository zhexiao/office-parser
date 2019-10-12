package word

func Convert(filepath string) *Question {
	//得到word数据
	w := Parser(filepath)

	//根据word数据解析试题数据
	q := ParseQuestion(w)

	return q
}
