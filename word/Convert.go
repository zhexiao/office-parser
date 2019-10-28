package word

import "io"

func ConvertFromFile(filepath string) *Question {
	//得到word数据
	w := Open(filepath)
	return convert(w)

}

func ConvertFromData(r io.ReaderAt, size int64) *Question {
	//得到excel数据
	w := Read(r, size)
	return convert(w)
}

func convert(w *Word) *Question {
	//根据word数据解析试题数据
	return ParseQuestion(w)
}
