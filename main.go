package main

import (
	"office-parser/word"
)

func main() {
	w := word.Word{}
	w.Parser("./test/question-fill.docx")

	q := word.Question{}
	q.Parser(&w)

	//s1 := "abc 1   "
	//fmt.Println(strings.Trim(s1, " "))
}