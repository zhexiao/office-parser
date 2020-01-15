package bases

import (
	"testing"
)

//返回需要读取的字符串列表
func TestReadNum(t *testing.T) {
	text := "abc[KMPS5-1-2-1]平的速度（多个以逗号隔开）[KMPS5-1]平抛运动"
	data := ReadNum(text)
	if len(data) != 2 {
		t.Error("应该读取到2个数据，读取失败")
	}
}

//返回需要读取的字符串列表
func TestReadText(t *testing.T) {
	text := "[KP1]生命科学"
	data := ReadText(text)
	if data != "生命科学" {
		t.Error("文本读取失败")
	}
}
