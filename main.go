package main

import (
	"encoding/json"
	"fmt"
	"log"
	"office-parser/word"
)

func main() {
	//解析word
	w := word.Word{}
	w.Parser("./test/question-fill.docx")

	//解析试题数据
	q := word.Question{}
	q.Parser(&w)

	//转为json数据
	jsonBytes, err := json.Marshal(q)
	if err != nil {
		log.Fatalf("json转换失败: %s", err)
	}

	fmt.Println(string(jsonBytes))
}
