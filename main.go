package main

import (
	"encoding/json"
	"fmt"
	"log"
	"office-parser/word"
)

func main() {
	//解析数据
	//q := word.Convert("./test/question-tzt-201903011.docx")
	q := word.Convert("./test/question-fill.docx")

	//转为json数据
	jsonBytes, err := json.Marshal(q)
	if err != nil {
		log.Fatalf("json转换失败: %s", err)
	}

	fmt.Println(string(jsonBytes))
}
