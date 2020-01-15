package word

import (
	"io/ioutil"
	"testing"
)

var t1 = "../_testdata/t1.docx"

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
