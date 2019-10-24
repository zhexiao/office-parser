package main

import (
	"encoding/json"
	"fmt"
	"log"
	"office-parser/excel"
)

func main() {
	dest := `C:\Users\zhexiao\Downloads\真实excel数据`

	//解析word数据
	//d1 := word.Convert("./test/question-fill-201903011.docx")

	//解析excel数据
	//d1 := excel.Convert(fmt.Sprintf("%s/%s", dest, "paper_20190702.xlsx"), "paper")
	//d1 := excel.Convert(fmt.Sprintf("%s/%s", dest, "book.xlsx"), "book")
	//d1 := excel.Convert(fmt.Sprintf("%s/%s", dest, "outline.xlsx"), "outline")
	//d1 := excel.Convert(fmt.Sprintf("%s/%s", dest, "cognition_map.xlsx"), "cognition_map")
	d1 := excel.Convert(fmt.Sprintf("%s/%s", dest, "cognition_sp.xlsx"), "cognition_sp")

	jsonBytes, err := json.Marshal(d1)
	if err != nil {
		log.Fatalf("json转换失败: %s", err)
	}
	fmt.Println(string(jsonBytes))
}
