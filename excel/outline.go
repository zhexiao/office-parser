package excel

import (
	"github.com/zhexiao/office-parser/utils"
	"log"
	"strconv"
	"strings"
)

type CT_OutlineBook struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Term        int    `json:"term"`
	SourceUri   string `json:"source_uri"`
	Isbn        string `json:"isbn"`
}

type CT_OutlineAttr struct {
	Num              string   `json:"num"`
	ParentNum        string   `json:"parent_num"`
	Name             string   `json:"name"`
	Level            int      `json:"level"`
	CognitionMapNums []string `json:"cognition_map_nums"`
	Sort             int      `json:"sort"`
	Description      string   `json:"description"`
}

type CT_Outline struct {
	Grade         int               `json:"grade"`
	PublishYear   int               `json:"publish_year"`
	Faculty       int               `json:"faculty"`
	Subject       int               `json:"subject"`
	PublisherName string            `json:"publisher_name"`
	Outline       []*CT_OutlineAttr `json:"outline"`
	OutlineBook   *CT_OutlineBook   `json:"outline_book"`
}

func NewCT_OutlineBook() *CT_OutlineBook {
	return &CT_OutlineBook{}
}

func NewCT_Outline() *CT_Outline {
	return &CT_Outline{}
}

func NewCT_OutlineAttr() *CT_OutlineAttr {
	return &CT_OutlineAttr{}
}

func ParseOutline(e *Excel) *CT_Outline {
	var (
		//所属目录的结束节点位置
		nodeEndCol int

		outAttrs []*CT_OutlineAttr
	)

	//记录每一个level已经有多少个数据了
	levelCount := make(map[int]int)

	//记录每一级的最后一个num
	levelNum := make(map[int]string)

	outlineBook := NewCT_OutlineBook()
	outline := NewCT_Outline()
	for idx, row := range e.RowsData {
		if idx == 0 {
			continue
		} else if idx == 1 {
			outlineName := strings.Trim(row.Content[0], " ")

			//得到学科学段
			faculty, err := strconv.Atoi(row.Content[1])
			if err != nil {
				log.Panicf("解析学段失败 %s", err)
			}

			subject, err := strconv.Atoi(row.Content[2])
			if err != nil {
				log.Panicf("解析学科失败 %s", err)
			}

			publisherName := strings.Trim(row.Content[3], " ")

			year, err := strconv.Atoi(row.Content[4])
			if err != nil {
				log.Panicf("解析审核年份失败 %s", err)
			}

			grade, err := strconv.Atoi(row.Content[5])
			if err != nil {
				log.Panicf("解析适用年级失败 %s", err)
			}

			termName := strings.Trim(row.Content[6], " ")
			var term int
			switch termName {
			case "上":
				term = 1
			case "下":
				term = 2
			default:
				term = 0
			}

			isbn := strings.Trim(row.Content[7], " ")

			outlineBook.Name = outlineName
			outlineBook.Term = term
			outlineBook.Isbn = isbn

			outline.Grade = grade
			outline.PublishYear = year
			outline.Faculty = faculty
			outline.Subject = subject
			outline.PublisherName = publisherName
		} else if idx == 2 {
			//得到当前 #节点 的结束位置
			for n, v := range row.Content {
				if strings.Contains(v, "#目录节点") {
					nodeEndCol = n
				}
			}
		} else {
			//	下面的数据为节点数据
			mapsStr := strings.Trim(row.Content[nodeEndCol+1], " ")

			//实例化
			outAttr := NewCT_OutlineAttr()
			outAttr.CognitionMapNums = utils.ReadNum(mapsStr)

			for m, v := range row.Content {
				if m <= nodeEndCol {
					numArr := utils.ReadNum(v)
					name := utils.ReadText(v)

					if len(numArr) > 0 && name != "" {
						num := numArr[0]

						//num转大写
						num = strings.ToUpper(num)

						//记录每一级的数量
						levelCount[m] += 1

						//记录每一级的最后一个num
						levelNum[m] = num
						if m > 0 {
							outAttr.ParentNum = levelNum[m-1]
						}

						outAttr.Name = name
						outAttr.Num = num
						outAttr.Level = m
						outAttr.Sort = levelCount[m]

						//插入列表
						outAttrs = append(outAttrs, outAttr)
					}
				}
			}
		}
	}

	outline.OutlineBook = outlineBook
	outline.Outline = outAttrs
	return outline
}
