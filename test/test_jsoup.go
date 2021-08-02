package test

import (
	"fmt"
	"reader_ev_parser/parser"
)

func Jsoup() {
	jsoupListStr("tag.h1.0@ownText")
	jsoupListStr("text.甲方@html")
	jsoupListStr("id.list.0@tag.dd.0:1:2@text&&id.list.0@tag.dd.12:13@text&&id.list.0@tag.dd.12:13@html")
	//jsoupListStr("tag.h1.0@html")
	//jsoupListStr("id.list.0@tag.dd.0@all")
	//jsoupListStr("id.list.0@tag.dd.0@html")
	//jsoupListStr("id.list.0@tag.dd.0@tag.a.0@href")

	//jsoupBatch("id.list.0@tag.dd.0:1:2:3:4", "tag.a@href", "text")

}

func jsoupListStr(rule string) {
	tId := parser.StartTransaction(DATA)
	result := parser.ParseRuleStr(tId, rule)
	fmt.Printf("jsoup rule-> [%s] result-> %s\n", rule, result)
	parser.EndTransaction(tId)
}

func jsoupBatch(rule string, rule1 string, rule2 string) {
	tId := parser.StartTransaction(DATA)
	bId := parser.ParseRuleRaw(tId, rule)
	size := parser.QueryBatchResultSize(bId)
	fmt.Printf("jsoup rule-> [%s] result-> %d\n", rule, size)
	for i := 0; i < size; i++ {
		r1 := parser.ParseRuleStrForParent(bId, rule1, i)
		fmt.Printf("jsoup rule-> [%s] result-> %s\n", rule1, r1)
		r2 := parser.ParseRuleStrForParent(bId, rule2, i)
		fmt.Printf("jsoup rule-> [%s] result-> %s\n", rule2, r2)
	}
	parser.EndTransaction(bId)
	parser.EndTransaction(tId)
}

func Regexp() {
	var test = parser.RegexpFilter([]string{"测试内容,测试内容，测试"}, "##测试")
	fmt.Println(test)
	var test2 = parser.RegexpFilter([]string{"测试内容,测屁内容，测试"}, "##测(.)##体$1")
	fmt.Println(test2)
	var test3 = parser.RegexpFilter([]string{"测试内容,测试内容，测试"}, "##(.)试##替换$1的###")
	fmt.Println(test3)
}
