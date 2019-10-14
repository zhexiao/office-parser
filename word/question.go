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

type QuestionChoice struct {
	Content          string `json:"content"`
	Correct          bool   `json:"correct"`
	CognitionMapNums string `json:"cognition_map_nums"`
	CognitionSpNums  string `json:"cognition_sp_nums"`
	Assessment       string `json:"assessment"`
}

type QuestionHint struct {
	Hint string `json:"hint"`
}

type Question struct {
	BasicType       string                  `json:"basic_type"`
	ResUsage        string                  `json:"res_usage"`
	Year            int                     `json:"year"`
	Author          string                  `json:"author"`
	LabelString     string                  `json:"label_string"`
	Grade           int                     `json:"grade"`
	QuestionAppType int                     `json:"question_app_type"`
	OftenTest       int                     `json:"often_test"`
	Note            string                  `json:"note"`
	EstimatedTime   int                     `json:"estimated_time"`
	Diff            int                     `json:"diff"`
	DiffDisplay     float64                 `json:"diff_display"`
	IdentifyDisplay float64                 `json:"identify_display"`
	Identify        int                     `json:"identify"`
	GuessDisplay    float64                 `json:"guess_display"`
	Guess           int                     `json:"guess"`
	ModelType       string                  `json:"model_type"`
	Stem            string                  `json:"stem"`
	Image           string                  `json:"image"`
	HasHint         int                     `json:"has_hint"`
	StructureString string                  `json:"structure_string"`
	Subject         int                     `json:"subject"`
	Uploader        string                  `json:"uploader"`
	SourceType      int                     `json:"source_type"`
	SourceUri       string                  `json:"source_uri"`
	QBasicType      []map[string]string     `json:"q_basic_type"`
	QLabelString    []map[string]string     `json:"q_label_string"`
	QCognitionMap   []*QuestionCognitionMap `json:"q_cognition_map"`
	QCognitionSp    []*QuestionCognitionSp  `json:"q_cognition_sp"`
	QOutline        []*QuestionOutline      `json:"q_outline"`
	QHint           []*QuestionHint         `json:"q_hint"`
	QResolve        []*QuestionResolve      `json:"q_resolve"`
	QAnswer         []*QuestionAnswer       `json:"q_answer"`
	QChoice         []*QuestionChoice       `json:"q_choice"`
	QRelation       []*Question             `json:"q_relation"`
}

func NewQuestion() *Question {
	return &Question{}
}

//把word里面的试题数据解析出来
func ParseQuestion(w *Word) *Question {
	q := NewQuestion()
	for idx, table := range w.Tables {
		//读取基本类型
		firstRow := table.Rows[0]
		basicType := utils.BasicType(
			strings.Trim(firstRow.Content[0], " ")).Val()

		//基本类型，如果是选择题，则区分单选和多选
		if basicType == "XZT" {
			xztType, err := strconv.Atoi(firstRow.Content[2])
			if err != nil {
				log.Fatal("选择题类型转换失败")
			}
			switch xztType {
			case 1:
				basicType = utils.BasicType("单选题").Val()
			case 2:
				basicType = utils.BasicType("多选题").Val()
			default:
				log.Fatalf("选择题 类型数据有错误")
			}
		}

		//structuring string读取
		if basicType == "TZT" {
			tztType, err := strconv.Atoi(firstRow.Content[2])
			if err != nil {
				log.Fatal("题组题 类型转换失败")
			}

			q.StructureString = utils.StructuringString(tztType).Val()
		} else {
			q.StructureString = utils.StructuringString(1).Val()
		}

		//第一个表格解析
		if idx == 0 {
			q.BasicType = basicType
			//解析表单
			q.parseTable(table)
			//	后面的子题
		} else {
			if q.BasicType == "TZT" {
				childQuestion := NewQuestion()
				childQuestion.BasicType = basicType

				//解析表单
				childQuestion.parseTable(table)

				//一些基本属性继承于母题
				childQuestion.ResUsage = q.ResUsage
				childQuestion.Year = q.Year
				childQuestion.Author = q.Author
				childQuestion.LabelString = q.LabelString
				childQuestion.Grade = q.Grade
				childQuestion.OftenTest = q.OftenTest

				//子题插入到母题
				q.QRelation = append(q.QRelation, childQuestion)
			}
		}
	}

	return q
}

func (q *Question) parseTable(t *TableData) {
	//基础题型数据也插入数组中
	q.QBasicType = append(q.QBasicType, map[string]string{
		"type": q.BasicType,
	})

	//固定数据
	q.SourceType = 1

	//解析基础数据
	q.parseMeta(t)

	//解析附加属性
	q.parseAddon(t)
}

//试题基础属性
func (q *Question) parseMeta(t *TableData) {
	for _, row := range t.Rows {
		title := strings.Trim(row.Content[0], " ")

		switch {
		case strings.Contains(title, "应用类型"):
			resUsage, err := strconv.Atoi(row.Content[1])
			if err != nil {
				log.Fatalf("应用类型 解析失败 %s", err)
			}

			q.ResUsage = utils.ResUsage(resUsage).Val()
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

			q.LabelString = utils.QuestionLabelString(labelString).Val()
			q.QLabelString = append(q.QLabelString, map[string]string{
				"label": q.LabelString,
			})
		case strings.Contains(title, "适用年级"):
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

			////选择题没有自动批改，填空等题型才有，废弃
			//var autoGrade int
			//if len(row.Content) >= 4 {
			//	autoGrade, err = strconv.Atoi(row.Content[4])
			//	if err != nil {
			//		log.Fatalf("自动批改 解析失败 %s", err)
			//	}
			//} else {
			//	autoGrade = 0
			//}

			q.OftenTest = oftenTest
		case strings.Contains(title, "试题备注"):
			q.Note = row.Content[1]
		case strings.Contains(title, "解题时间"):
			estimatedTime, err := strconv.Atoi(row.Content[1])
			if err != nil {
				log.Fatalf("解题时间 解析失败 %s", err)
			}

			diffDisplay, err := strconv.ParseFloat(row.Content[3], 2)
			if err != nil {
				log.Fatalf("困难度 解析失败 %s", err)
			}

			identify := row.Content[5]
			if identify == "" {
				q.Identify = 0
			} else {
				identifyDisplay, err := strconv.ParseFloat(identify, 2)
				if err != nil {
					log.Fatalf("鉴别度 解析失败 %s", err)
				}
				q.IdentifyDisplay = identifyDisplay
				q.Identify = int(identifyDisplay * 100)
			}

			guess := row.Content[7]
			if guess == "" {
				q.Guess = 0
			} else {
				guessDisplay, err := strconv.ParseFloat(guess, 2)
				if err != nil {
					log.Fatalf("猜度 解析失败 %s", err)
				}
				q.GuessDisplay = guessDisplay
				q.Guess = int(guessDisplay * 100)
			}

			q.EstimatedTime = estimatedTime
			q.DiffDisplay = diffDisplay
			q.Diff = int(diffDisplay * 100)
		case strings.Contains(title, "版型"):
			q.ModelType = row.Content[1]
		case strings.Contains(title, "题目文字"):
			q.Stem = row.Content[1]
		case strings.Contains(title, "题目图片"):
			q.Image = row.Content[1]
		}
	}
}

//试题附加属性
func (q *Question) parseAddon(t *TableData) {
	var (
		answerTable bool
		hintTable   bool
	)

	for _, row := range t.Rows {
		title := strings.Trim(row.Content[0], " ")

		switch {
		case strings.Contains(title, "知识地图"):
			q.parseCognitionMap(row)
		case strings.Contains(title, "教材目录"):
			q.parseOutline(row)
		case strings.Contains(title, "特异性知识点"):
			q.parseCognitionSp(row)
		case strings.Contains(title, "详解"):
			q.parseResolve(row)
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

			//非选择题的标题
			if strings.Contains(content, "答案内容") {
				continue
			}

			q.parseAnswer(row)
		}

		//读取提示数据
		if hintTable {
			q.parseHint(row)
		}
	}
}

//试题普适性知识点
func (q *Question) parseCognitionMap(row *RowData) {
	content := row.Content[1]
	numList := utils.ReadNum(content)

	for _, num := range numList {
		numObj := QuestionCognitionMap{
			CognitionMapNum: num,
		}

		q.QCognitionMap = append(q.QCognitionMap, &numObj)
	}
}

//试题目录
func (q *Question) parseOutline(row *RowData) {
	content := row.Content[1]
	numList := utils.ReadNum(content)

	for _, num := range numList {
		numObj := QuestionOutline{
			OutlineNum: num,
		}

		q.QOutline = append(q.QOutline, &numObj)
	}
}

//试题特异性知识点
func (q *Question) parseCognitionSp(row *RowData) {
	content := row.Content[1]
	numList := utils.ReadNum(content)

	for _, num := range numList {
		numObj := QuestionCognitionSp{
			CognitionSpNum: num,
		}

		q.QCognitionSp = append(q.QCognitionSp, &numObj)
	}
}

//试题解析
func (q *Question) parseResolve(row *RowData) {
	content := row.Content[1]
	if content == "" {
		return
	}

	resolveObj := QuestionResolve{
		Resolve:        content,
		MimeType:       0,
		MimeUri:        "",
		DefaultResolve: 1,
	}

	q.QResolve = append(q.QResolve, &resolveObj)
}

//试题答案(答案的数据读取需要区分不同的题型)
func (q *Question) parseAnswer(row *RowData) {
	//选择题的属性
	var (
		isChoice = false
		correct  bool
	)

	var (
		content string
		maps    string
		sps     string
	)

	switch q.BasicType {
	case "TKT":
		content = row.Content[0]
		maps = row.Content[1]
		sps = row.Content[2]
	case "JDT":
		content = row.Content[1]
	case "PDT":
		content = row.Content[1]
	case "DANXT", "DUOXT":
		isChoice = true
		correctText := row.Content[0]

		//选择题的标题
		if strings.Contains(correctText, "是否正确") {
			return
		}

		if strings.EqualFold(correctText, "v") {
			correct = true
		} else {
			correct = false
		}

		content = row.Content[1]
		maps = row.Content[2]
		sps = row.Content[3]
	}

	var (
		mapNums = strings.Join(utils.ReadNum(maps), ",")
		spNums  = strings.Join(utils.ReadNum(sps), ",")
	)

	if isChoice {
		choiceObj := QuestionChoice{
			Content:          content,
			Correct:          correct,
			CognitionMapNums: mapNums,
			CognitionSpNums:  spNums,
			Assessment:       "",
		}

		q.QChoice = append(q.QChoice, &choiceObj)
	} else {
		answerObj := QuestionAnswer{
			Answer:           content,
			AutoCorrect:      0,
			CognitionMapNums: mapNums,
			CognitionSpNums:  spNums,
			Assessment:       "",
		}

		q.QAnswer = append(q.QAnswer, &answerObj)
	}

}

//试题提示
func (q *Question) parseHint(row *RowData) {
	content := row.Content[1]
	if content == "" {
		return
	}

	hintObj := QuestionHint{
		Hint: content,
	}

	q.HasHint = 1
	q.QHint = append(q.QHint, &hintObj)
}
