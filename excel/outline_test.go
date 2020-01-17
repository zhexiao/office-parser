package excel

import (
	"io/ioutil"
	"testing"
)

var outlineTestFile1 = "../_testdata/outline_test.xlsx"
var outlineTestFile2 = "../_testdata/outline_test_multiple_root.xlsx"

var outlineTestFile3 = "../_testdata/outline_test_missing_faculty.xlsx"
var outlineTestFile4 = "../_testdata/outline_test_missing_grade.xlsx"
var outlineTestFile5 = "../_testdata/outline_test_missing_subject.xlsx"
var outlineTestFile6 = "../_testdata/outline_test_missing_year.xlsx"

func TestNewCT_Outline(t *testing.T) {
	var b interface{}
	b1 := NewCT_Outline()

	b = b1

	_, ok := b.(*CT_Outline)
	if !ok {
		t.Error("实例化有误")
	}
}

func TestNewCT_OutlineBook(t *testing.T) {
	var b interface{}
	b1 := NewCT_OutlineBook()

	b = b1

	_, ok := b.(*CT_OutlineBook)
	if !ok {
		t.Error("实例化有误")
	}
}

func TestNewCT_OutlineAttr(t *testing.T) {
	var b interface{}
	b1 := NewCT_OutlineAttr()

	b = b1

	_, ok := b.(*CT_OutlineAttr)
	if !ok {
		t.Error("实例化有误")
	}
}

func TestParseOutline(t *testing.T) {
	e := NewCT_Excel()

	fileBytes, err := ioutil.ReadFile(outlineTestFile1)
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

	data, err := ParseOutline(excel)
	if err != nil {
		t.Error(err)
	}

	if len(data.Outline) != 8 {
		t.Error("数据解析有错误")
	}
}

func TestParseOutlineWithMultipleRoot(t *testing.T) {
	e := NewCT_Excel()

	fileBytes, err := ioutil.ReadFile(outlineTestFile2)
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

	_, err = ParseOutline(excel)
	if err == nil {
		t.Error("存在多个根目录，应该报错")
	}
}

func readData(filepath string, t *testing.T) (*CT_Outline, error) {
	e := NewCT_Excel()

	fileBytes, err := ioutil.ReadFile(filepath)
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

	data, err := ParseOutline(excel)
	return data, err
}

func TestParseOutlineWithMissing(t *testing.T) {
	_, err := readData(outlineTestFile3, t)
	if err == nil {
		t.Error("学段 不存在，应该报错")
	}

	_, err = readData(outlineTestFile4, t)
	if err == nil {
		t.Error("年级 不存在，应该报错")
	}

	_, err = readData(outlineTestFile5, t)
	if err == nil {
		t.Error("学科 不存在，应该报错")
	}

	_, err = readData(outlineTestFile6, t)
	if err == nil {
		t.Error("年份 不存在，应该报错")
	}
}
