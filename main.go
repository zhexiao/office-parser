package main

import "office-parser/excel"

func main() {
	//解析word数据
	//q := word.Convert("./test/question-fill-201903011.docx")
	//
	////转为json数据
	//jsonBytes, err := json.Marshal(q)
	//if err != nil {
	//	log.Fatalf("json转换失败: %s", err)
	//}
	//
	//fmt.Println(string(jsonBytes))

	//解析excel数据
	//excel.Convert("./test/paper.xlsx", "paper")
	excel.Convert("./test/book-index-201903011.xlsx", "book")
	//excel.Convert("./test/cognition-map-201903011.xlsx", "cognition_map")
	//excel.Convert("./test/cognition-sp-201903011.xlsx", "cognition_sp")
	//excel.Convert("./test/outline-201903011.xlsx", "outline")
}
