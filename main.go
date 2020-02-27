package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/zhexiao/office-parser/bases"
	"github.com/zhexiao/office-parser/excel"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type CT_YamlSettings struct {
	Qiniu struct {
		AccessKey    string `yaml:"accessKey"`
		SecretKey    string `yaml:"secretKey"`
		Bucket       string `yaml:"bucket"`
		Zone         string `yaml:"zone"`
		Domain       string `yaml:"domain"`
		UploadPrefix string `yaml:"uploadPrefix"`
	}

	Wmf struct {
		Uri string `yaml:"uri"`
	}
}

func readSettings() *CT_YamlSettings {
	//读取配置文件
	yamlSettings := new(CT_YamlSettings)

	yamlFile, err := ioutil.ReadFile("settings.yaml")
	if err != nil {
		return nil
	}

	err = yaml.Unmarshal(yamlFile, yamlSettings)
	if err != nil {
		return nil
	}

	return yamlSettings
}

func init() {
	//读取配置文件
	yamlSettings := readSettings()

	//如果默认没有配置，则可以作为第三方库支持用户自定义配置
	if yamlSettings != nil {
		//初始化七牛的配置
		bases.OpQiniu = &bases.CT_Qiniu{
			//七牛key
			AccessKey: yamlSettings.Qiniu.AccessKey,
			//七牛secret
			SecretKey: yamlSettings.Qiniu.SecretKey,
			//七牛存储的bucket
			Bucket: yamlSettings.Qiniu.Bucket,
			//所属区域
			Zone: yamlSettings.Qiniu.Zone,
			//访问的域名地址
			Domain: yamlSettings.Qiniu.Domain,
			//路径
			UploadPrefix: yamlSettings.Qiniu.UploadPrefix,
		}

		//初始化WMF的配置
		bases.OpWmf = &bases.CT_WmfCfg{
			Uri: yamlSettings.Wmf.Uri,
		}
	}
}

func main() {
	//var (
	//	filepath string
	//	eduType  string
	//	data     interface{}
	//)
	//
	//app := cli.NewApp()
	//app.Name = "Office Parser"
	//app.Usage = "Convert Word、Excel to json data"
	//app.Version = "2.0"
	//app.EnableBashCompletion = true
	//
	//app.Flags = []cli.Flag{
	//	cli.StringFlag{
	//		Name:        "filepath, f",
	//		Usage:       "filepath",
	//		Destination: &filepath,
	//	},
	//	cli.StringFlag{
	//		Name:        "eduType, t",
	//		Usage:       "The type of the document belong to!",
	//		Destination: &eduType,
	//	},
	//}
	//
	//app.Action = func(c *cli.Context) error {
	//	if filepath == "" {
	//		log.Panic("缺少必要的文件路径")
	//	}
	//
	//	switch eduType {
	//	case "question":
	//		data = word.ConvertFromFile(filepath)
	//	case "word":
	//		data = word.ConvertPaperFromFile(filepath)
	//	case "paper":
	//		data = excel.ConvertFromFile(filepath, "paper")
	//	case "book":
	//		data = excel.ConvertFromFile(filepath, "book")
	//	case "outline":
	//		data = excel.ConvertFromFile(filepath, "outline")
	//	case "cognition_map":
	//		data = excel.ConvertFromFile(filepath, "cognition_map")
	//	case "cognition_sp":
	//		data = excel.ConvertFromFile(filepath, "cognition_sp")
	//	default:
	//		log.Panicf("不支持的解析类型：%s", eduType)
	//	}
	//
	//	jsonBytes, err := json.Marshal(data)
	//	if err != nil {
	//		log.Panicf("json转换失败: %s", err)
	//	}
	//
	//	//saveIntoFile(string(jsonBytes))
	//	fmt.Println("=================================================")
	//	fmt.Println("=================================================")
	//	fmt.Println(string(jsonBytes))
	//	fmt.Println("=================================================")
	//	fmt.Println("=================================================")
	//	return nil
	//}
	//
	////执行命令行
	//err := app.Run(os.Args)
	//if err != nil {
	//	log.Panic(err)
	//}

	//注释上面，运行测试
	test()
}

func test() {
	//data, err := word.ConvertFromFile("./test/question-fill-201903011.docx")
	//data, err := word.ConvertPaperFromFile("./test/ywgs.docx")
	data, err := excel.ConvertFromFile("./_testdata/cognition_map_test.xlsx", "cognition_map")
	if err != nil {
		log.Panicf("失败: %s", err)
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Panicf("json转换失败: %s", err)
	}
	fmt.Println(string(jsonBytes))

}

//把数据写入文件
func saveIntoFile(data string) {
	filename := fmt.Sprintf("%d%d%d%d%d%d.json", time.Now().Year(), int(time.Now().Month()), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())
	file, err := os.Create(filename)
	if err != nil {
		log.Panicf("打开文件失败err=%s", err)
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		log.Panicf("写入文件失败err=%s", err)
	}
}
