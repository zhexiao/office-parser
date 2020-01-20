package word

import (
	"testing"
)

var questionFillTest1 = "../_testdata/question-fill-test.docx"

func TestConvertFromFile(t *testing.T) {
	_, err := ConvertFromFile("")
	if err == nil {
		t.Error("传入空文件，应该报错")
	}
}

func TestConvertFromData(t *testing.T) {
	_, err := ConvertFromData([]byte{})
	if err == nil {
		t.Error("传入空字节，应该报错")
	}
}

func TestConvertPaperFromFile(t *testing.T) {
	_, err := ConvertPaperFromFile("")
	if err == nil {
		t.Error("传入空文件，应该报错")
	}
}

func TestConvertPaperFromData(t *testing.T) {
	_, err := ConvertPaperFromData([]byte{})
	if err == nil {
		t.Error("传入空字节，应该报错")
	}
}

func TestConvertFile(t *testing.T) {
	data, err := ConvertFromFile(questionFillTest1)
	if err != nil {
		t.Error(err)
	}

	if len(data.QCognitionMap) <= 0 {
		t.Error("解析出来的数据有误，未匹配上 知识点")
	}

	if len(data.QOutline) <= 0 {
		t.Error("解析出来的数据有误，未匹配上 目录")
	}

	if len(data.QAnswer) <= 0 {
		t.Error("解析出来的数据有误，未匹配上 答案")
	}
}

func TestConvertPaper(t *testing.T) {
	data, err := ConvertPaperFromFile(t1)
	if err != nil {
		t.Error(err)
	}

	if len(data.WordText) <= 0 {
		t.Error("数据解析有误")
	}
}
