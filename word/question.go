package word

import (
	"fmt"
	"log"
	"office-parser/utils"
	"strconv"
	"strings"
)

type QuestionCognitionMap struct {
	cognitionMapNum string
}

type QuestionOutline struct {
	outlineNum string
}

type QuestionCognitionSp struct {
	cognitionSpNum string
}

type QuestionResolve struct {
	resolve        string
	mimeType       int
	mimeUri        string
	defaultResolve int
}

type QuestionAnswer struct {
	answer           string
	autoCorrect      int
	cognitionMapNums string
	cognitionSpNums  string
	assessment       string
}

type QuestionHint struct {
	hint string
}

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
	estimatedTime   int
	diff            float64
	identify        float64
	guess           float64
	modelType       string
	stem            string
	image           string
	qCognitionMap   []QuestionCognitionMap
	qCognitionSp    []QuestionCognitionSp
	qOutline        []QuestionOutline
	qHint           []QuestionHint
	qResolve        []QuestionResolve
	qAnswer         []QuestionAnswer
}

//解析word数据到试题结构
func (q *Question) Parser(w *Word) {
	//读取基本类型
	firstRow := w.Content[0]
	q.basicType = strings.Trim(firstRow.Content[0], " ")

	q.parseTable(w)

	//打印数据
	fmt.Printf("%#v \n", q)
}

func (q *Question) parseTable(w *Word) {
	q.parseMeta(w)
	q.parseAddon(w)
}

//试题基础属性
func (q *Question) parseMeta(w *Word) {
	for _, row := range w.Content {
		title := strings.Trim(row.Content[0], " ")

		switch {
		case strings.Contains(title, "应用类型"):
			resUsage, err := strconv.Atoi(row.Content[1])
			if err != nil {
				log.Fatalf("应用类型 解析失败 %s", err)
			}

			q.resUsage = resUsage
		case strings.Contains(title, "题库年度"):
			year, err := strconv.Atoi(row.Content[1])
			if err != nil {
				log.Fatalf("题库年度 解析失败 %s", err)
			}

			q.year = year
			q.author = row.Content[3]
		case strings.Contains(title, "试题描述类型"):
			labelString, err := strconv.Atoi(row.Content[1])
			if err != nil {
				log.Fatalf("试题描述类型 解析失败 %s", err)
			}

			q.labelString = labelString
		case strings.Contains(title, "试用年级"):
			grade, err := strconv.Atoi(row.Content[1])
			if err != nil {
				log.Fatalf("试用年级 解析失败 %s", err)
			}

			q.grade = grade
		case strings.Contains(title, "学科题型"):
			questionAppType, err := strconv.Atoi(row.Content[1])
			if err != nil {
				log.Fatalf("学科题型 解析失败 %s", err)
			}

			q.questionAppType = questionAppType
		case strings.Contains(title, "常考题"):
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
		case strings.Contains(title, "试题备注"):
			q.note = row.Content[1]
		case strings.Contains(title, "解题时间"):
			estimatedTime, err := strconv.Atoi(row.Content[1])
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

			q.estimatedTime = estimatedTime
			q.diff = diff
		case strings.Contains(title, "版型"):
			q.modelType = row.Content[1]
		case strings.Contains(title, "题目文字"):
			q.stem = row.Content[1]
		case strings.Contains(title, "题目图片"):
			q.image = "需要通过上传到七牛保存图片地址"
		}
	}
}

//试题附加属性
func (q *Question) parseAddon(w *Word) {
	var answerTable bool
	var hintTable bool

	for _, row := range w.Content {
		title := strings.Trim(row.Content[0], " ")

		switch {
		case strings.Contains(title, "知识地图"):
			content := row.Content[1]
			q.parseCognitionMap(content)
		case strings.Contains(title, "教材目录"):
			content := row.Content[1]
			q.parseOutline(content)
		case strings.Contains(title, "特异性知识点"):
			content := row.Content[1]
			q.parseCognitionSp(content)
		case strings.Contains(title, "详解"):
			content := row.Content[1]
			q.parseResolve(content)
		case strings.Contains(title, "#ANSWER"):
			answerTable = true
			continue
		case strings.Contains(title, "#HINT"):
			answerTable = false
			hintTable = true
			continue
		}

		if answerTable {
			content := row.Content[0]
			cognitionMapContent := row.Content[1]
			cognitionSpContent := row.Content[2]

			if strings.Contains(content, "答案内容") {
				continue
			}

			q.parseAnswer(content, cognitionMapContent, cognitionSpContent)
		}

		if hintTable {
			content := row.Content[1]
			q.parseHint(content)
		}
	}
}

//试题普适性知识点
func (q *Question) parseCognitionMap(content string) {
	numList := utils.ReadNum(content)

	for _, num := range numList {
		numObj := QuestionCognitionMap{
			cognitionMapNum: num,
		}

		q.qCognitionMap = append(q.qCognitionMap, numObj)
	}
}

//试题目录
func (q *Question) parseOutline(content string) {
	numList := utils.ReadNum(content)

	for _, num := range numList {
		numObj := QuestionOutline{
			outlineNum: num,
		}

		q.qOutline = append(q.qOutline, numObj)
	}
}

//试题特异性知识点
func (q *Question) parseCognitionSp(content string) {
	numList := utils.ReadNum(content)

	for _, num := range numList {
		numObj := QuestionCognitionSp{
			cognitionSpNum: num,
		}

		q.qCognitionSp = append(q.qCognitionSp, numObj)
	}
}

//试题解析
func (q *Question) parseResolve(content string) {
	//todo 带图片的文本需要处理
	resolveObj := QuestionResolve{
		resolve:        content,
		mimeType:       0,
		mimeUri:        "",
		defaultResolve: 1,
	}

	q.qResolve = append(q.qResolve, resolveObj)
}

//试题答案
func (q *Question) parseAnswer(content string, maps string, sps string) {
	//数组转字符串
	mapNumList := utils.ReadNum(maps)
	spNumlist := utils.ReadNum(sps)

	mapNums := strings.Join(mapNumList, ",")
	spNums := strings.Join(spNumlist, ",")

	//todo 带图片的文本需要处理
	answerObj := QuestionAnswer{
		answer:           content,
		autoCorrect:      0,
		cognitionMapNums: mapNums,
		cognitionSpNums:  spNums,
		assessment:       "",
	}

	q.qAnswer = append(q.qAnswer, answerObj)
}

//试题提示
func (q *Question) parseHint(content string) {
	//todo 带图片的文本需要处理
	hintObj := QuestionHint{
		hint: content,
	}

	q.qHint = append(q.qHint, hintObj)
}
