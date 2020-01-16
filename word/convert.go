package word

import (
	"io/ioutil"
)

//word上传试题
func ConvertFromFile(filepath string) (*Question, error) {
	fileBytes, _ := ioutil.ReadFile(filepath)

	//解析试题结构
	return ConvertFromData(fileBytes)
}

//word上传试题
func ConvertFromData(fileBytes []byte) (*Question, error) {
	q, err := ParseQuestion(fileBytes)
	if err != nil {
		return nil, err
	}
	return q, nil
}

//word上传试卷
func ConvertPaperFromFile(filepath string) (*CT_PureWord, error) {
	fileBytes, _ := ioutil.ReadFile(filepath)

	//解析试题结构
	return ConvertPaperFromData(fileBytes)
}

//word上传试卷
func ConvertPaperFromData(fileBytes []byte) (*CT_PureWord, error) {
	paper, err := ParsePaper(fileBytes)
	if err != nil {
		return nil, err
	}

	return paper, nil
}
