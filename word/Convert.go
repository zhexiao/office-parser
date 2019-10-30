package word

import "io"

//word上传试题
func ConvertFromFile(filepath string) *Question {
	//得到word数据
	doc := Open(filepath)
	w := QuestionWord(doc)

	//解析试题结构
	return ParseQuestion(w)
}

//word上传试题
func ConvertFromData(r io.ReaderAt, size int64) *Question {
	//得到excel数据
	doc := Read(r, size)
	w := QuestionWord(doc)

	//解析试题结构
	return ParseQuestion(w)
}

//word上传试卷
func ConvertPaperFromFile(filepath string) *CT_PureWord {
	//得到word数据
	doc := Open(filepath)
	w := PaperWord(doc)

	//解析试卷
	return ParsePaper(w)
}
