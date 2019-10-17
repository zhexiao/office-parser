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
		//根据 excel 数据解析试卷
		paper := ParsePaper(e)

		jsonBytes, err := json.Marshal(paper)
		if err != nil {
			log.Fatalf("json转换失败: %s", err)
		}

		fmt.Println(string(jsonBytes))
	}

}
