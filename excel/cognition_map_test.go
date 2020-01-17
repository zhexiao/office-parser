package excel

import (
	"io/ioutil"
	"testing"
)

var cognitionMapTestFile = "../_testdata/cognition_map_test.xlsx"
var cognitionMapTestFile1 = "../_testdata/cognition_map_test_multiple_root.xlsx"
var cognitionMapTestFile2 = "../_testdata/cognition_map_test_missing_faculty.xlsx"
var cognitionMapTestFile3 = "../_testdata/cognition_map_test_missing_subject.xlsx"

func TestNewCT_CognitionMap(t *testing.T) {
	var b interface{}
	b1 := NewCT_CognitionMap()

	b = b1

	_, ok := b.(*CT_CognitionMap)
	if !ok {
		t.Error("实例化有误")
	}
}

func TestParseCognitionMapWithMultipleRoot(t *testing.T) {
	e := NewCT_Excel()

	fileBytes, err := ioutil.ReadFile(cognitionMapTestFile1)
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

	_, err = ParseCognitionMap(excel)
	if err == nil {
		t.Error("存在多个根目录，应该报错")
	}
}

func TestParseCognitionMap(t *testing.T) {
	e := NewCT_Excel()

	fileBytes, err := ioutil.ReadFile(cognitionMapTestFile)
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

	data, err := ParseCognitionMap(excel)
	if err != nil {
		t.Error(err)
	}

	if len(data) != 8 {
		t.Error("数据解析有错误")
	}
}

func TestParseCognitionMapMissingSubject(t *testing.T) {
	e := NewCT_Excel()

	fileBytes, err := ioutil.ReadFile(cognitionMapTestFile3)
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

	_, err = ParseCognitionMap(excel)
	if err == nil {
		t.Error("缺少学科，应该报错")
	}
}

func TestParseCognitionMapMissingFaculty(t *testing.T) {
	e := NewCT_Excel()

	fileBytes, err := ioutil.ReadFile(cognitionMapTestFile2)
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

	_, err = ParseCognitionMap(excel)
	if err == nil {
		t.Error("缺少学段，应该报错")
	}
}
