package excel

import (
	"fmt"
	"strings"
)

type CT_Book struct {
	BookId     string          `json:"book_id"`
	BookIndexs []*CT_BookIndex `json:"book_indexs"`
}

type CT_BookIndex struct {
	Level   string `json:"level"`
	Name    string `json:"name"`
	PaperId string `json:"paper_id"`
}

func NewCT_Book() *CT_Book {
	return &CT_Book{}
}

func NewCT_BookIndex() *CT_BookIndex {
	return &CT_BookIndex{}
}

func ParseBook(e *CT_Excel) (*CT_Book, error) {
	var (
		contentIdx    int
		levelStrArray []string
	)

	//记录每一个level已经有多少个数据了
	levelCount := make(map[int]int)

	//循环数据
	book := NewCT_Book()
	for idx, row := range e.RowsData {
		if idx == 0 {
			continue
		}

		//实例化
		bIdx := NewCT_BookIndex()

		for n, val := range row.Content {
			levelIdx := n + 1
			content := strings.Trim(val, " ")

			//作为目录内容
			if contentIdx == 0 && content != "" {
				bIdx.Name = content

				//每一个level数据加1
				levelCount[levelIdx] += 1

				for m := 1; m <= levelIdx; m++ {
					levelStrArray = append(levelStrArray, fmt.Sprintf("%03d", levelCount[m]))
				}
				bIdx.Level = strings.Join(levelStrArray, ".")

				//设置目录idx
				contentIdx = levelIdx
			} else if content != "" {
				//如果还循环到存在内容，并且不是第一次循环到内容，则视为paper_id
				bIdx.PaperId = content
			}
		}

		//数据重置
		contentIdx = 0
		levelStrArray = nil

		book.BookIndexs = append(book.BookIndexs, bIdx)
	}

	return book, nil
}
