package bases

import (
	"testing"
)

func TestError(t *testing.T) {
	e := NewOpError(NormalError, "error")

	if e.Error() != "error" {
		t.Errorf("错误信息不一致")
	}
}

func TestErrorType(t *testing.T) {
	e1 := &OpError{}
	var e2 interface{} = e1

	_, ok := e2.(error)
	if !ok {
		t.Errorf("没有实现error的方法")
	}
}
