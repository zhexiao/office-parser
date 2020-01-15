package excel

import (
	"github.com/zhexiao/office-parser/bases"
	"log"
	"strconv"
	"strings"
)

type CT_CognitionSp struct {
	Num              string `json:"num"`
	ParentNum        string `json:"parent_num"`
	Level            int    `json:"level"`
	Sort             int    `json:"sort"`
	Name             string `json:"name"`
	PreNum           string `json:"pre_num"`
	ExtendNum        string `json:"extend_num"`
	Weight           string `json:"weight"`
	Faculty          int    `json:"faculty"`
	Subject          int    `json:"subject"`
	SpType           int    `json:"sp_type"`
	CognitionMapNums string `json:"cognition_map_nums"`
}

func NewCT_CognitionSp() *CT_CognitionSp {
	return &CT_CognitionSp{}
}

func ParseCognitionSp(e *Excel) []*CT_CognitionSp {
	var (
		//获得节点结束的列，从0开始
		nodeEndCol int
		faculty    int
		subject    int
		spType     int
	)

	//记录每一个level已经有多少个数据了
	levelCount := make(map[int]int)

	//记录每一级的最后一个num
	levelNum := make(map[int]string)

	var cogs []*CT_CognitionSp
	for idx, row := range e.RowsData {
		if idx == 0 {
			continue
		} else if idx == 1 {
			//得到学科学段
			facultyTmp, err := strconv.Atoi(row.Content[0])
			if err != nil {
				log.Panicf("解析学段失败 %s", err)
			}

			subjectTmp, err := strconv.Atoi(row.Content[1])
			if err != nil {
				log.Panicf("解析学科失败 %s", err)
			}

			spTypeTmp, err := strconv.Atoi(row.Content[3])
			if err != nil {
				log.Panicf("解析认知点类型失败 %s", err)
			}

			subject = subjectTmp
			faculty = facultyTmp
			spType = spTypeTmp
		} else if idx == 2 {
			//得到当前 #节点 的结束位置
			for n, v := range row.Content {
				if strings.Contains(v, "#节点") {
					nodeEndCol = n
				}
			}
		} else {
			//	下面的数据为节点数据
			preNumStr := strings.Trim(row.Content[nodeEndCol+1], " ")
			extendNumStr := strings.Trim(row.Content[nodeEndCol+2], " ")
			mapsStr := strings.Trim(row.Content[nodeEndCol+4], " ")

			//实例化
			cog := NewCT_CognitionSp()
			cog.Faculty = faculty
			cog.Subject = subject
			cog.SpType = spType
			cog.PreNum = strings.Join(bases.ReadNum(preNumStr), ",")
			cog.ExtendNum = strings.Join(bases.ReadNum(extendNumStr), ",")
			cog.CognitionMapNums = strings.Join(bases.ReadNum(mapsStr), ",")

			for m, v := range row.Content {
				if m <= nodeEndCol {
					num := strings.Join(bases.ReadNum(v), ",")
					name := bases.ReadText(v)

					if num != "" && name != "" {
						//num转大写
						num = strings.ToUpper(num)

						//记录每一级的数量
						levelCount[m] += 1

						//记录每一级的最后一个num
						levelNum[m] = num
						if m > 0 {
							cog.ParentNum = levelNum[m-1]
						}

						cog.Name = name
						cog.Num = num
						cog.Level = m
						cog.Sort = levelCount[m]

						//插入列表
						cogs = append(cogs, cog)
					}
				}
			}
		}
	}

	return cogs
}
