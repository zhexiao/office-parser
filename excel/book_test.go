package excel

import (
	"io/ioutil"
	"testing"
)

var bookTestFile = "../_testdata/book_test.xlsx"

func TestNewCT_Book(t *testing.T) {
	var b interface{}
	b1 := NewCT_Book()

	b = b1

	_, ok := b.(*CT_Book)
	if !ok {
		t.Error("实例化有误")
	}
}

func TestNewCT_BookIndex(t *testing.T) {
	var b interface{}
	b1 := NewCT_BookIndex()

	b = b1

	_, ok := b.(*CT_BookIndex)
	if !ok {
		t.Error("实例化有误")
	}
}

func TestParseBook(t *testing.T) {
	e := NewCT_Excel()

	fileBytes, err := ioutil.ReadFile(bookTestFile)
	if err != nil {
		t.Error(err)
	}

	err = e.read(fileBytes)
	if err != nil {
		t.Error(err)
	}

	excel, err := e.GetExcelData()
	if err != nil {
		t.Error(err)
	}

	data, err := ParseBook(excel)
	if err != nil {
		t.Error(err)
	}

	if len(data.BookIndexs) != 18 {
		t.Error("解析出来的数据有误，未匹配上数量")
	}
}
