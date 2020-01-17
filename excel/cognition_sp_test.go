package excel

import (
	"io/ioutil"
	"testing"
)

var cognitionSpTestFile = "../_testdata/cognition_sp_test.xlsx"
var cognitionSpTestFile1 = "../_testdata/cognition_sp_test_multiple_root.xlsx"
var cognitionSpTestFile2 = "../_testdata/cognition_sp_test_missing_faculty.xlsx"
var cognitionSpTestFile3 = "../_testdata/cognition_sp_test_missing_subject.xlsx"
var cognitionSpTestFile4 = "../_testdata/cognition_sp_test_missing_type.xlsx"

func TestNewCT_CognitionSp(t *testing.T) {
	var b interface{}
	b1 := NewCT_CognitionSp()

	b = b1

	_, ok := b.(*CT_CognitionSp)
	if !ok {
		t.Error("实例化有误")
	}
}

func TestParseCognitionSp(t *testing.T) {
	e := NewCT_Excel()

	fileBytes, err := ioutil.ReadFile(cognitionSpTestFile)
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

	data, err := ParseCognitionSp(excel)
	if err != nil {
		t.Error(err)
	}

	if len(data) != 8 {
		t.Error("数据解析有错误")
	}
}

func TestParseCognitionSpWithMultipleRoot(t *testing.T) {
	e := NewCT_Excel()

	fileBytes, err := ioutil.ReadFile(cognitionSpTestFile1)
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

	_, err = ParseCognitionSp(excel)
	if err == nil {
		t.Error("存在多个根目录，应该报错")
	}
}

func TestParseCognitionSpMissingSubject(t *testing.T) {
	e := NewCT_Excel()

	fileBytes, err := ioutil.ReadFile(cognitionSpTestFile3)
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

	_, err = ParseCognitionSp(excel)
	if err == nil {
		t.Error("缺少学科，应该报错")
	}
}

func TestParseCognitionSpMissingFaculty(t *testing.T) {
	e := NewCT_Excel()

	fileBytes, err := ioutil.ReadFile(cognitionSpTestFile2)
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

	_, err = ParseCognitionSp(excel)
	if err == nil {
		t.Error("缺少学段，应该报错")
	}
}

func TestParseCognitionSpMissingType(t *testing.T) {
	e := NewCT_Excel()

	fileBytes, err := ioutil.ReadFile(cognitionSpTestFile4)
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

	_, err = ParseCognitionSp(excel)
	if err == nil {
		t.Error("缺少类型，应该报错")
	}
}
