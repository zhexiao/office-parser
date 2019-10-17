package utils

import (
	"regexp"
)

//返回需要读取的字符串列表
func ReadNum(text string) []string {
	//text := "abc[KMPS5-1-2-1]平的速度（多个以逗号隔开）[KMPS5-1]平抛运动"

	//这里必须使用非贪婪模式
	reg := regexp.MustCompile(`\[([\w\W\d-]+?)\]`)
	data := reg.FindAllStringSubmatch(text, -1)

	var res []string
	for _, dt := range data {
		res = append(res, dt[1])
	}

	return res
}

//返回需要读取的文本文字
func ReadText(text string) string {
	//text := "[KP1]生命科学"

	//这里必须使用非贪婪模式
	reg := regexp.MustCompile(`\[([\w\W\d-]+?)\](.*)`)
	data := reg.FindAllStringSubmatch(text, -1)

	var res string
	for _, dt := range data {
		res = dt[2]
	}

	return res
}
