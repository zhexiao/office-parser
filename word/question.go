package word

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Question struct {
	basicType       string
	resUsage        int
	year            int
	author          string
	labelString     int
	grade           int
	questionAppType int
	oftenTest       int
	autoGrade       int
	note            string
	estimated_time  int
	diff            float64
	identify        float64
	guess           float64
	model_type      string
	stem            string
}

//解析word数据到试题结构
func (q *Question) Parser(w *Word) {
	//读取基本类型
	firstRow := w.Content[0]
	q.basicType = strings.Trim(firstRow.Content[0], " ")

	q.parseTable(w)

	//打印数据
	fmt.Printf("%#v", q)
}

func (q *Question) parseTable(w *Word) {
	//读取试题数据
	for _, row := range w.Content {
		firstText := strings.Trim(row.Content[0], " ")

		switch {
		case strings.Contains(firstText, "应用类型"):
			resUsage, err := strconv.Atoi(row.Content[1])
			if err != nil {
				log.Fatalf("应用类型 解析失败 %s", err)
			}

			q.resUsage = resUsage
		case strings.Contains(firstText, "题库年度"):
			year, err := strconv.Atoi(row.Content[1])
			if err != nil {
				log.Fatalf("题库年度 解析失败 %s", err)
			}

			q.year = year
			q.author = row.Content[3]
		case strings.Contains(firstText, "试题描述类型"):
			labelString, err := strconv.Atoi(row.Content[1])
			if err != nil {
				log.Fatalf("试题描述类型 解析失败 %s", err)
			}

			q.labelString = labelString
		case strings.Contains(firstText, "试用年级"):
			grade, err := strconv.Atoi(row.Content[1])
			if err != nil {
				log.Fatalf("试用年级 解析失败 %s", err)
			}

			q.grade = grade
		case strings.Contains(firstText, "学科题型"):
			questionAppType, err := strconv.Atoi(row.Content[1])
			if err != nil {
				log.Fatalf("学科题型 解析失败 %s", err)
			}

			q.questionAppType = questionAppType
		case strings.Contains(firstText, "常考题"):
			oftenTest, err := strconv.Atoi(row.Content[1])
			if err != nil {
				log.Fatalf("常考题 解析失败 %s", err)
			}

			autoGrade, err := strconv.Atoi(row.Content[4])
			if err != nil {
				log.Fatalf("自动批改 解析失败 %s", err)
			}

			q.oftenTest = oftenTest
			q.autoGrade = autoGrade
		case strings.Contains(firstText, "试题备注"):
			q.note = row.Content[1]
		case strings.Contains(firstText, "解题时间"):
			estimated_time, err := strconv.Atoi(row.Content[1])
			if err != nil {
				log.Fatalf("解题时间 解析失败 %s", err)
			}

			diff, err := strconv.ParseFloat(row.Content[3], 2)
			if err != nil {
				log.Fatalf("困难度 解析失败 %s", err)
			}

			identify := row.Content[5]
			if identify == "" {
				q.identify = 0
			} else {
				identify, err := strconv.ParseFloat(identify, 2)
				if err != nil {
					log.Fatalf("鉴别度 解析失败 %s", err)
				}
				q.identify = identify
			}

			guess := row.Content[7]
			if guess == "" {
				q.guess = 0
			} else {
				guess, err := strconv.ParseFloat(guess, 2)
				if err != nil {
					log.Fatalf("猜度 解析失败 %s", err)
				}
				q.guess = guess
			}

			q.estimated_time = estimated_time
			q.diff = diff
		case strings.Contains(firstText, "版型"):
			q.model_type = row.Content[1]
		case strings.Contains(firstText, "题目文字"):
			q.stem = row.Content[1]
		}
	}
}
