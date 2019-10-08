package main

import (
	"office-parser/word"
)

func main() {
	w := word.Word{}
	w.Parser("./test/question-fill.docx")

	q := word.Question{}
	q.Parser(&w)
}
