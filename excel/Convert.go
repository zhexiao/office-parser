package excel

import (
	"encoding/json"
	"fmt"
	"log"
)

func Convert(filepath string, _type string) {
	//得到excel数据
	e := Parser(filepath)

	switch _type {
	case "paper":
		paper := ParsePaper(e)
		jsonBytes, err := json.Marshal(paper)
		if err != nil {
			log.Fatalf("json转换失败: %s", err)
		}

		fmt.Println(string(jsonBytes))
	case "book":
		book := ParseBook(e)
		jsonBytes, err := json.Marshal(book)
		if err != nil {
			log.Fatalf("json转换失败: %s", err)
		}

		fmt.Println(string(jsonBytes))
	case "cognition_map":
		cognitionMap := ParseCognitionMap(e)
		jsonBytes, err := json.Marshal(cognitionMap)
		if err != nil {
			log.Fatalf("json转换失败: %s", err)
		}

		fmt.Println(string(jsonBytes))
	case "cognition_sp":
		cognitionSp := ParseCognitionSp(e)
		jsonBytes, err := json.Marshal(cognitionSp)
		if err != nil {
			log.Fatalf("json转换失败: %s", err)
		}

		fmt.Println(string(jsonBytes))
	case "outline":
		outline := ParseOutline(e)
		jsonBytes, err := json.Marshal(outline)
		if err != nil {
			log.Fatalf("json转换失败: %s", err)
		}

		fmt.Println(string(jsonBytes))
	}

}
