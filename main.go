package main

import (
	"encoding/json"
	"fmt"
	"log"
	"office-parser/word"
)

func main() {
	w := word.Word{}
	w.Parser("./test/question-fill.docx")

	q := word.Question{}
	q.Parser(&w)

	jsonBytes, err := json.Marshal(q)
	if err != nil {
		log.Fatalf("json转换失败: %s", err)
	}

	fmt.Println(string(jsonBytes))
}
