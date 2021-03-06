package excel

import (
	"testing"
)

func TestConvertFromFile(t *testing.T) {
	_, err := ConvertFromFile("", "book")
	if err == nil {
		t.Error("传入空文件，应该报错")
	}
}

func TestConvertFromData(t *testing.T) {
	_, err := ConvertFromData([]byte{}, "book")
	if err == nil {
		t.Error("传入空字节，应该报错")
	}
}

func TestConvertBook(t *testing.T) {
	data, err := ConvertFromFile(bookTestFile, "book")
	if err != nil {
		t.Error(err)
	}

	data1 := data.(*CT_Book)
	if len(data1.BookIndexs) != 18 {
		t.Error("解析出来的数据有误，未匹配上数量")
	}
}

func TestConvertErrorType(t *testing.T) {
	_, err := ConvertFromFile(cognitionMapTestFile, "cognition")
	if err == nil {
		t.Error("错误类型，应该报错")
	}
}

func TestConvertCognitionMap(t *testing.T) {
	data, err := ConvertFromFile(cognitionMapTestFile, "cognition_map")
	if err != nil {
		t.Error(err)
	}

	data1 := data.([]*CT_CognitionMap)
	if len(data1) != 8 {
		t.Error("数据解析有错误")
	}
}

func TestConvertCognitionSp(t *testing.T) {
	data, err := ConvertFromFile(cognitionSpTestFile, "cognition_sp")
	if err != nil {
		t.Error(err)
	}

	data1 := data.([]*CT_CognitionSp)
	if len(data1) != 8 {
		t.Error("数据解析有错误")
	}
}

func TestConvertOutline(t *testing.T) {
	data, err := ConvertFromFile(outlineTestFile1, "outline")
	if err != nil {
		t.Error(err)
	}

	data1 := data.(*CT_Outline)
	if len(data1.Outline) != 8 {
		t.Error("数据解析有错误")
	}
}
