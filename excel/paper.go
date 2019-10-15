package excel

import (
	"fmt"
	"log"
	"office-parser/utils"
	"strconv"
	"strings"
)

type Paper struct {
	Paper     PaperAttr `json:"'paper'"`
	Questions []*SubtypeAttr
}

type PaperAttr struct {
	Labels      []string `json:"labels"`
	ExcelUri    string   `json:"excel_uri"`
	Name        string   `json:"name"`
	Area        string   `json:"area"`
	Grade       int      `json:"grade"`
	LabelString string   `json:"label_string"`
	Note        string   `json:"note"`
	Score       float64  `json:"score"`
	OutlineNums string   `json:"outline_nums"`
	Time        int      `json:"time"`
	ResUsage    string   `json:"res_usage"`
	Year        int      `json:"year"`
}

type SubtypeAttr struct {
	Name     string          `json:"name"`
	Note     string          `json:"note"`
	Question []*QuestionAttr `json:"question"`
}

type QuestionAttr struct {
	Qid   string          `json:"qid"`
	Score float64         `json:"score"`
	Child []*QuestionAttr `json:"child"`
}

func NewPaper() *Paper {
	return &Paper{}
}

func NewPaperAttr() *PaperAttr {
	return &PaperAttr{}
}

func NewQuestionAttr() *QuestionAttr {
	return &QuestionAttr{}
}

func NewSubtypeAttr() *SubtypeAttr {
	return &SubtypeAttr{}
}

func ParsePaper(e *Excel) {
	paper := NewPaper()

	//解析数据
	paper.parseRow(e)
}

func (paper Paper) parseRow(e *Excel) {
	pAttr := NewPaperAttr()

	//循环数据
	for idx, row := range e.RowsData {
		//根据当前行数读取数据
		switch idx {
		case 0, 2, 4:
			continue
		case 1:
			resUsageNum, err := strconv.Atoi(row.Content[0])
			if err != nil {
				log.Fatalf("应用类型 解析失败 %s", err)
			}
			resUsage := utils.ResUsage(resUsageNum).Val()

			labelStringNum, err := strconv.Atoi(row.Content[1])
			if err != nil {
				log.Fatalf("试卷描述类型 解析失败 %s", err)
			}
			labelString := utils.PaperLabelString(labelStringNum).Val()

			note := row.Content[2]
			name := row.Content[3]

			pAttr.ResUsage = resUsage
			pAttr.LabelString = labelString
			pAttr.Labels = []string{labelString}
			pAttr.Note = note
			pAttr.Name = name
		case 3:
			outlineNumsContent := row.Content[0]
			nums := strings.Join(utils.ReadNum(outlineNumsContent), ",")

			grade, err := strconv.Atoi(row.Content[1])
			if err != nil {
				log.Fatalf("年级 解析失败 %s", err)
			}

			time, err := strconv.Atoi(row.Content[2])
			if err != nil {
				log.Fatalf("时间 解析失败 %s", err)
			}

			totalScore, err := strconv.ParseFloat(row.Content[3], 2)
			if err != nil {
				log.Fatalf("总分 解析失败 %s", err)
			}

			year, err := strconv.Atoi(row.Content[4])
			if err != nil {
				log.Fatalf("年度 解析失败 %s", err)
			}

			area := row.Content[5]

			pAttr.OutlineNums = nums
			pAttr.Grade = grade
			pAttr.Time = time
			pAttr.Score = totalScore
			pAttr.Year = year
			pAttr.Area = area
		default:
			fmt.Println(row)
		}
	}

	fmt.Printf("%#v \n", pAttr)
}
