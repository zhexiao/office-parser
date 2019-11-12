package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"github.com/zhexiao/office-parser/excel"
	"github.com/zhexiao/office-parser/word"
	"log"
	"os"
)

func init() {

}

func main() {
	var (
		filepath string
		eduType  string
		data     interface{}
	)

	app := cli.NewApp()
	app.Name = "Office Parser"
	app.Usage = "Convert Word、Excel to json data"
	app.Version = "2.0"
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "filepath, f",
			Usage:       "filepath",
			Destination: &filepath,
		},
		cli.StringFlag{
			Name:        "eduType, t",
			Usage:       "The type of the document belong to!",
			Destination: &eduType,
		},
	}

	app.Action = func(c *cli.Context) error {
		if filepath == "" {
			log.Panic("缺少必要的文件路径")
		}

		switch eduType {
		case "question":
			data = word.ConvertFromFile(filepath)
		case "word":
			data = word.ConvertPaperFromFile(filepath)
		case "paper":
			data = excel.ConvertFromFile(filepath, "paper")
		case "book":
			data = excel.ConvertFromFile(filepath, "book")
		case "outline":
			data = excel.ConvertFromFile(filepath, "outline")
		case "cognition_map":
			data = excel.ConvertFromFile(filepath, "cognition_map")
		case "cognition_sp":
			data = excel.ConvertFromFile(filepath, "cognition_sp")
		default:
			log.Panicf("不支持的解析类型：%s", eduType)
		}

		jsonBytes, err := json.Marshal(data)
		if err != nil {
			log.Panicf("json转换失败: %s", err)
		}
		fmt.Println(string(jsonBytes))

		return nil
	}

	//执行命令行
	err := app.Run(os.Args)
	if err != nil {
		log.Panic(err)
	}

	//注释上面，运行测试
	//test()
}

func test() {
	//data := word.ConvertFromFile("./test/question-fill-201903011.docx")
	data := word.ConvertPaperFromFile("./test/text1.docx")

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Panicf("json转换失败: %s", err)
	}
	fmt.Println(string(jsonBytes))
}
