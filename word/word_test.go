package word

import (
	"io/ioutil"
	"testing"
)

var t1 = "../_testdata/t1.docx"
var t2 = "../_testdata/t2.docx"

func TestNewCT_Word(t *testing.T) {
	var t1 interface{}
	t2 := NewCT_Word()

	t1 = t2
	_, ok := t1.(*CT_Word)
	if !ok {
		t.Errorf("初始化方法有错误")
	}
}

func TestWordRead(t *testing.T) {
	bytes, err := ioutil.ReadFile(t1)
	if err != nil {
		t.Error(err)
	}

	ct := NewCT_Word()
	err = ct.read(bytes)
	if err != nil {
		t.Error(err)
	}

	if ct.doc == nil {
		t.Error("无法读取doc的数据")
	}
}

func TestParseOle(t *testing.T) {
	bytes, err := ioutil.ReadFile(t1)
	if err != nil {
		t.Error(err)
	}

	ct := NewCT_Word()
	err = ct.read(bytes)
	if err != nil {
		t.Error(err)
	}

	ct.parseOle()
	ct.oles.Range(func(key, value interface{}) bool {
		if key == nil || value == nil {
			t.Error("无法识别公式数据")
		}
		return true
	})
}

func TestParseImage(t *testing.T) {
	bytes, err := ioutil.ReadFile(t1)
	if err != nil {
		t.Error(err)
	}

	ct := NewCT_Word()
	err = ct.read(bytes)
	if err != nil {
		t.Error(err)
	}

	ct.parseImage()
	ct.images.Range(func(key, value interface{}) bool {
		if key == nil || value == nil {
			t.Error("无法识别图片数据")
		}
		return true
	})
}

func TestParseOrder(t *testing.T) {
	bytes, err := ioutil.ReadFile(t2)
	if err != nil {
		t.Error(err)
	}

	ct := NewCT_Word()
	err = ct.read(bytes)
	if err != nil {
		t.Error(err)
	}

	ct.parseOrder()
	if len(ct.numIdMapAbNumId) == 0 || len(ct.numData) ==0{
		t.Error("无法识别序号数据")
	}

	if len(ct.numIdMapAbNumId) != 1 {
		t.Error("序号个数不正确")
	}
}

func TestTableData(t *testing.T) {
	bytes, err := ioutil.ReadFile(t2)
	if err != nil {
		t.Error(err)
	}

	ct := NewCT_Word()
	err = ct.read(bytes)
	if err != nil {
		t.Error(err)
	}

	ct.getTableData()
	if len(ct.Tables) != 1{
		t.Error("读取表格数据有误")
	}

	for _,value := range ct.Tables{
		row := value.Rows
		if len(row) != 3 {
			t.Error("行数据读取失败")
		}
	}
}