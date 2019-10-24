package excel

import (
	"encoding/json"
	"log"
)

func Convert(filepath string, _type string) string {
	//得到excel数据
	e := Parser(filepath)

	var data interface{}
	switch _type {
	case "paper":
		data = ParsePaper(e)
	case "book":
		data = ParseBook(e)
	case "cognition_map":
		data = ParseCognitionMap(e)
	case "cognition_sp":
		data = ParseCognitionSp(e)
	case "outline":
		data = ParseOutline(e)
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("json转换失败: %s", err)
	}

	return string(jsonBytes)
}
