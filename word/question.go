package word

import (
	"log"
	"office-parser/utils"
	"strconv"
	"strings"
)

type QuestionCognitionMap struct {
	CognitionMapNum string `json:"cognition_map_num"`
}

type QuestionOutline struct {
	OutlineNum string `json:"outline_num"`
}

type QuestionCognitionSp struct {
	CognitionSpNum string `json:"cognition_sp_num"`
}

type QuestionResolve struct {
	Resolve        string `json:"resolve"`
	MimeType       int    `json:"mime_type"`
	MimeUri        string `json:"mime_uri"`
	DefaultResolve int    `json:"default_resolve"`
}

type QuestionAnswer struct {
	Answer           string `json:"answer"`
	AutoCorrect      int    `json:"auto_correct"`
	CognitionMapNums string `json:"cognition_map_nums"`
	CognitionSpNums  string `json:"cognition_sp_nums"`
	Assessment       string `json:"assessment"`
}

type QuestionHint struct {
	Hint string `json:"hint"`
}

type Question struct {
	BasicType       string                 `json:"basic_type"`
	ResUsage        int                    `json:"res_usage"`
	Year            int                    `json:"year"`
	Author          string                 `json:"author"`
	LabelString     int                    `json:"label_string"`
	Grade           int                    `json:"grade"`
	QuestionAppType int                    `json:"question_app_type"`
	OftenTest       int                    `json:"often_test"`
	AutoGrade       int                    `json:"auto_grade"`
	Note            string                 `json:"note"`
	EstimatedTime   int                    `json:"estimated_time"`
	Diff            float64                `json:"diff"`
	Identify        float64                `json:"identify"`
	Guess           float64                `json:"guess"`
	ModelType       string                 `json:"model_type"`
	Stem            string                 `json:"stem"`
	Image           string                 `json:"image"`
	QCognitionMap   []QuestionCognitionMap `json:"q_cognition_map"`
	QCognitionSp    []QuestionCognitionSp  `json:"q_cognition_sp"`
	QOutline        []QuestionOutline      `json:"q_outline"`
	QHint           []QuestionHint         `json:"q_hint"`
	QResolve        []QuestionResolve      `json:"q_resolve"`
	QAnswer         []QuestionAnswer       `json:"q_answer"`
}

//解析word数据到试题结构
func (q *Question) Parser(w *Word) {
	//读取基本类型
	firstRow := w.Content[0]
	q.BasicType = strings.Trim(firstRow.Content[0], " ")

	q.parseTable(w)

	//打印结构体数据
	//fmt.Printf("%#v \n", q)
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

			q.ResUsage = resUsage
		case strings.Contains(title, "题库年度"):
			year, err := strconv.Atoi(row.Content[1])
			if err != nil {
				log.Fatalf("题库年度 解析失败 %s", err)
			}

			q.Year = year
			q.Author = row.Content[3]
		case strings.Contains(title, "试题描述类型"):
			labelString, err := strconv.Atoi(row.Content[1])
			if err != nil {
				log.Fatalf("试题描述类型 解析失败 %s", err)
			}

			q.LabelString = labelString
		case strings.Contains(title, "试用年级"):
			grade, err := strconv.Atoi(row.Content[1])
			if err != nil {
				log.Fatalf("试用年级 解析失败 %s", err)
			}

			q.Grade = grade
		case strings.Contains(title, "学科题型"):
			questionAppType, err := strconv.Atoi(row.Content[1])
			if err != nil {
				log.Fatalf("学科题型 解析失败 %s", err)
			}

			q.QuestionAppType = questionAppType
		case strings.Contains(title, "常考题"):
			oftenTest, err := strconv.Atoi(row.Content[1])
			if err != nil {
				log.Fatalf("常考题 解析失败 %s", err)
			}

			autoGrade, err := strconv.Atoi(row.Content[4])
			if err != nil {
				log.Fatalf("自动批改 解析失败 %s", err)
			}

			q.OftenTest = oftenTest
			q.AutoGrade = autoGrade
		case strings.Contains(title, "试题备注"):
			q.Note = row.Content[1]
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
				q.Identify = 0
			} else {
				identify, err := strconv.ParseFloat(identify, 2)
				if err != nil {
					log.Fatalf("鉴别度 解析失败 %s", err)
				}
				q.Identify = identify
			}

			guess := row.Content[7]
			if guess == "" {
				q.Guess = 0
			} else {
				guess, err := strconv.ParseFloat(guess, 2)
				if err != nil {
					log.Fatalf("猜度 解析失败 %s", err)
				}
				q.Guess = guess
			}

			q.EstimatedTime = estimatedTime
			q.Diff = diff
		case strings.Contains(title, "版型"):
			q.ModelType = row.Content[1]
		case strings.Contains(title, "题目文字"):
			q.Stem = row.Content[1]
		case strings.Contains(title, "题目图片"):
			q.Image = "需要通过上传到七牛保存图片地址"
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

		//读取答案数据
		if answerTable {
			content := row.Content[0]
			cognitionMapContent := row.Content[1]
			cognitionSpContent := row.Content[2]

			if strings.Contains(content, "答案内容") {
				continue
			}

			q.parseAnswer(content, cognitionMapContent, cognitionSpContent)
		}

		//读取提示数据
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
			CognitionMapNum: num,
		}

		q.QCognitionMap = append(q.QCognitionMap, numObj)
	}
}

//试题目录
func (q *Question) parseOutline(content string) {
	numList := utils.ReadNum(content)

	for _, num := range numList {
		numObj := QuestionOutline{
			OutlineNum: num,
		}

		q.QOutline = append(q.QOutline, numObj)
	}
}

//试题特异性知识点
func (q *Question) parseCognitionSp(content string) {
	numList := utils.ReadNum(content)

	for _, num := range numList {
		numObj := QuestionCognitionSp{
			CognitionSpNum: num,
		}

		q.QCognitionSp = append(q.QCognitionSp, numObj)
	}
}

//试题解析
func (q *Question) parseResolve(content string) {
	//todo 带图片的文本需要处理
	resolveObj := QuestionResolve{
		Resolve:        content,
		MimeType:       0,
		MimeUri:        "",
		DefaultResolve: 1,
	}

	q.QResolve = append(q.QResolve, resolveObj)
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
		Answer:           content,
		AutoCorrect:      0,
		CognitionMapNums: mapNums,
		CognitionSpNums:  spNums,
		Assessment:       "",
	}

	q.QAnswer = append(q.QAnswer, answerObj)
}

//试题提示
func (q *Question) parseHint(content string) {
	if content == "" {
		return
	}

	//todo 带图片的文本需要处理
	hintObj := QuestionHint{
		Hint: content,
	}

	q.QHint = append(q.QHint, hintObj)
}
