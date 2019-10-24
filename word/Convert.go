package word

import (
	"encoding/json"
	"log"
)

func Convert(filepath string) string {
	//得到word数据
	w := Parser(filepath)

	//根据word数据解析试题数据
	q := ParseQuestion(w)

	jsonBytes, err := json.Marshal(q)
	if err != nil {
		log.Fatalf("json转换失败: %s", err)
	}
	return string(jsonBytes)
}
