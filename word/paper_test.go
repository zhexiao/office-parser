package word

import (
	"io/ioutil"
	"testing"
)

var p1 = "../_testdata/t1.docx"

func TestNewCT_PureWord(t *testing.T) {
	var t1 interface{}
	t2 := NewCT_PureWord()

	t1 = t2
	_, ok := t1.(*CT_PureWord)
	if !ok {
		t.Errorf("初始化方法有错误")
	}
}

func TestParsePaper(t *testing.T) {
	fBytes, err := ioutil.ReadFile(p1)
	if err != nil {
		t.Error(err)
	}

	data, err := ParsePaper(fBytes)
	if err != nil {
		t.Error(err)
	}

	if data.WordText == "" {
		t.Error("试卷文本数据解析失败")
	}
}
