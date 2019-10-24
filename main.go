package main

import (
	"fmt"
	"office-parser/excel"
	"office-parser/word"
)

func main() {
	//解析word数据
	d1 := word.Convert("./test/question-fill-201903011.docx")
	fmt.Println(d1)

	//解析excel数据
	//excel.Convert("./test/paper.xlsx", "paper")
	//excel.Convert("./test/book-index-201903011.xlsx", "book")
	//excel.Convert("./test/cognition-map-201903011.xlsx", "cognition_map")
	excel.Convert("./test/cognition-sp-201903011.xlsx", "cognition_sp")
	//excel.Convert("./test/outline-201903011.xlsx", "outline")
}
