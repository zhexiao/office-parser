package main

import (
	"fmt"
	"github.com/unidoc/unioffice/document"
)

func main() {
	doc, err := document.Open("./aaa.docx")
	if err != nil {
		fmt.Println(err)
	}

	pars := doc.Paragraphs()

	for _, p := range pars {
		fmt.Println("======================")
		runs := p.Runs()
		styles := p.Style()
		prop := p.Properties()

		fmt.Println(styles)
		fmt.Println(prop.Style())

		for _, r := range runs {
			fmt.Println(r.Text())
			rs := r.Properties()

			fmt.Println(rs.X().Color, rs.X().U)
			fmt.Println(rs.IsBold())
			fmt.Println(rs.IsItalic())

			rcolor := rs.Color()
			fmt.Println(rcolor.X().ValAttr)
		}
	}

}
