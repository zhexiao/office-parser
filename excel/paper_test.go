package excel

import (
	"io/ioutil"
	"testing"
)

var paperTestFile1 = "../_testdata/paper_test.xlsx"

func readPaperFile(filepath string, t *testing.T) *CT_Excel {
	ct := NewCT_Excel()

	fileBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		t.Error(err)
	}

	err = ct.read(fileBytes)
	if err != nil {
		t.Error(err)
	}

	excelData, err := ct.GetExcelData()
	if err != nil {
		t.Error(err)
	}

	return excelData
}
func TestParsePaper(t *testing.T) {
	e := readPaperFile(paperTestFile1, t)
	paper, err := ParsePaper(e)
	if err != nil {
		t.Error(err)
	}

	if len(paper.Questions) != 2 {
		t.Error("解析出错，大题数对不上")
	}

	for i, sb := range paper.Questions {
		if i == 0 && len(sb.Question) != 5 {
			t.Error("第一个大题题目解析出错")

			for j, q := range sb.Question {
				if j == 0 && len(q.Child) != 4 {
					t.Error("第一大题的第一小题的子题数读取失败")
				}
			}
		} else if i == 1 && len(sb.Question) != 6 {
			t.Error("第二个大题题目解析出错")

			for j, q := range sb.Question {
				if j == 1 && len(q.Child) != 2 {
					t.Error("第二大题的第二小题的子题数读取失败")
				}
			}
		}
	}
}
