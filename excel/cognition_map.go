package excel

import (
	"fmt"
	"github.com/zhexiao/office-parser/bases"
	"strconv"
	"strings"
)

type CT_CognitionMap struct {
	Num       string `json:"num"`
	ParentNum string `json:"parent_num"`
	Level     int    `json:"level"`
	Sort      int    `json:"sort"`
	Name      string `json:"name"`
	PreNum    string `json:"pre_num"`
	ExtendNum string `json:"extend_num"`
	Weight    string `json:"weight"`
	Faculty   int    `json:"faculty"`
	Subject   int    `json:"subject"`
}

func NewCT_CognitionMap() *CT_CognitionMap {
	return &CT_CognitionMap{}
}

func ParseCognitionMap(e *CT_Excel) ([]*CT_CognitionMap, error) {
	var (
		//获得节点结束的列，从0开始
		nodeEndCol int
		faculty    int
		subject    int

		//检查是否出现过根节点，一个知识点只允许出现一个根节点
		rootNode bool
	)

	//记录每一个level已经有多少个数据了
	levelCount := make(map[int]int)

	//记录每一级的最后一个num
	levelNum := make(map[int]string)

	var cogs []*CT_CognitionMap
	for idx, row := range e.RowsData {
		if idx == 0 {
			continue
		} else if idx == 1 {
			//得到学科学段
			facultyTmp, err := strconv.Atoi(row.Content[0])
			if err != nil {
				return nil, bases.NewOpError(bases.NormalError, fmt.Sprintf("解析学段失败 %s", err))
			}

			subjectTmp, err := strconv.Atoi(row.Content[1])
			if err != nil {
				return nil, bases.NewOpError(bases.NormalError, fmt.Sprintf("解析学科失败 %s", err))
			}

			subject = subjectTmp
			faculty = facultyTmp
		} else if idx == 2 {
			//得到当前#节点的结束位置
			for n, v := range row.Content {
				if strings.Contains(v, "#节点") {
					nodeEndCol = n
				}
			}
		} else {

			//	下面的数据为节点数据
			preNumStr := strings.Trim(row.Content[nodeEndCol+1], " ")
			extendNumStr := strings.Trim(row.Content[nodeEndCol+2], " ")

			//实例化
			cognitionMap := NewCT_CognitionMap()
			cognitionMap.Faculty = faculty
			cognitionMap.Subject = subject
			cognitionMap.PreNum = strings.Join(bases.ReadNum(preNumStr), ",")
			cognitionMap.ExtendNum = strings.Join(bases.ReadNum(extendNumStr), ",")

			for m, v := range row.Content {
				if m <= nodeEndCol {

					num := strings.Join(bases.ReadNum(v), ",")
					name := bases.ReadText(v)

					if num != "" && name != "" {
						if rootNode == false {
							//设置根节点已找到
							if m == 1 {
								rootNode = true
							}
						} else {
							if m == 1 {
								return nil, bases.NewOpError(bases.NormalError, "存在多个根目录")
							}
						}

						//num转大写
						num = strings.ToUpper(num)

						//记录每一级的数量
						levelCount[m] += 1

						//记录每一级的最后一个num
						levelNum[m] = num
						if m > 0 {
							cognitionMap.ParentNum = levelNum[m-1]
						}

						cognitionMap.Name = name
						cognitionMap.Num = num
						cognitionMap.Level = m
						cognitionMap.Sort = levelCount[m]

						//插入列表
						cogs = append(cogs, cognitionMap)
					}
				}
			}
		}
	}

	return cogs, nil
}
