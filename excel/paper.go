package excel

import (
	"fmt"
	"github.com/zhexiao/office-parser/utils"
	"log"
	"strconv"
	"strings"
)

type CT_Paper struct {
	Paper     *PaperAttr     `json:"paper"`
	Questions []*SubtypeAttr `json:"questions"`
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

func NewCT_Paper() *CT_Paper {
	return &CT_Paper{}
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

func ParsePaper(e *Excel) *CT_Paper {
	paper := NewCT_Paper()
	paperAttr := NewPaperAttr()
	var subtypeAttrs []*SubtypeAttr

	var (
		//大题索引
		currentSubtypeIdx = -1
		//小题索引
		currentQuestionIdx int
		//子题索引
		currentChildQIdx int
	)

	//循环数据
	for idx, row := range e.RowsData {
		//根据当前行数读取数据
		switch idx {
		case 0, 2, 4:
			continue
		case 1:
			//读取试卷本身属性
			resUsageNum, err := strconv.Atoi(row.Content[0])
			if err != nil {
				log.Panicf("应用类型 解析失败 %s", err)
			}
			resUsage := utils.ResUsage(resUsageNum).Val()

			labelStringNum, err := strconv.Atoi(row.Content[1])
			if err != nil {
				log.Panicf("试卷描述类型 解析失败 %s", err)
			}
			labelString := utils.PaperLabelString(labelStringNum).Val()

			note := row.Content[2]
			name := row.Content[3]

			paperAttr.ResUsage = resUsage
			paperAttr.LabelString = labelString
			paperAttr.Labels = []string{labelString}
			paperAttr.Note = note
			paperAttr.Name = name
		case 3:
			//读取试卷本身属性
			outlineNumsContent := row.Content[0]
			nums := strings.Join(utils.ReadNum(outlineNumsContent), ",")

			grade, err := strconv.Atoi(row.Content[1])
			if err != nil {
				log.Panicf("年级 解析失败 %s", err)
			}

			time, err := strconv.Atoi(row.Content[2])
			if err != nil {
				log.Panicf("时间 解析失败 %s", err)
			}

			totalScore, err := strconv.ParseFloat(row.Content[3], 2)
			if err != nil {
				log.Panicf("总分 解析失败 %s", err)
			}

			year, err := strconv.Atoi(row.Content[4])
			if err != nil {
				log.Panicf("年度 解析失败 %s", err)
			}

			area := row.Content[5]

			paperAttr.OutlineNums = nums
			paperAttr.Grade = grade
			paperAttr.Time = time
			paperAttr.Score = totalScore
			paperAttr.Year = year
			paperAttr.Area = area
		default:
			rowLen := len(row.Content)
			subtypeName := strings.Trim(row.Content[0], " ")

			//存在即认为是一个大题的开始
			if subtypeName != "" {
				currentSubtypeIdx += 1

				//大题开始，重置下面的子题索引
				currentQuestionIdx = -1

				if rowLen < 5 {
					log.Panic("表单结构有误，找不到大题描述")
				}

				subtypNote := strings.Trim(row.Content[5], " ")

				//大题属性
				subAttr := NewSubtypeAttr()
				subAttr.Name = subtypeName
				subAttr.Note = subtypNote

				//把大题插入到大题列表中
				subtypeAttrs = append(subtypeAttrs, subAttr)
			} else {
				//认为是小题
				qId := strings.Trim(row.Content[1], " ")
				if qId != "" {
					//重置子题索引
					currentChildQIdx = 1

					//当前试题索引
					currentQuestionIdx += 1

					//读取试题属性
					qAttr := NewQuestionAttr()
					qAttr.Qid = strings.ToLower(qId)
					if rowLen < 4 {
						log.Panic("表结构有误，找不到分数列")
					}

					qScore, err := strconv.ParseFloat(row.Content[3], 2)
					if err != nil {
						log.Panicf("分数转为浮点类型出错 %s %s", err, qId)
					}
					qAttr.Score = qScore

					//插入到对应的所属大题
					subtypeAttrs[currentSubtypeIdx].Question = append(
						subtypeAttrs[currentSubtypeIdx].Question, qAttr)

				} else {
					//	认为是子题
					childQId := strings.Trim(row.Content[2], " ")
					if childQId == "" {
						continue
					}

					//按规范重构子题ID为母题id_n
					parentQ := subtypeAttrs[currentSubtypeIdx].Question[currentQuestionIdx]
					newChildQid := fmt.Sprintf("%s_%s", parentQ.Qid, strconv.Itoa(currentChildQIdx))

					//读取试题属性
					childQAttr := NewQuestionAttr()
					childQAttr.Qid = strings.ToLower(newChildQid)
					if rowLen < 5 {
						log.Panic("表结构有误，找不到分数列")
					}

					childQScore, err := strconv.ParseFloat(row.Content[4], 2)
					if err != nil {
						log.Panicf("子题分数 转为浮点类型出错 %s", err)
					}
					childQAttr.Score = childQScore

					//插入到对应的所属母题
					parentQ.Child = append(parentQ.Child, childQAttr)

					//子题索引 + 1
					currentChildQIdx += 1
				}
			}
		}
	}

	paper.Paper = paperAttr
	paper.Questions = subtypeAttrs

	return paper
}
