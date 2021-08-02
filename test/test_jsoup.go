package test

import (
	"fmt"
	"reader_ev_parser/parser"
)

func Jsoup() {
	//tId := parser.StartTransaction(DATA)
	//result := parser.ParseRuleStr(tId, "tag.p.0@text##hello")
	//result2 := parser.ParseRuleStr(tId, `href@js:result+',{webView:\"true\"}'`)
	//fmt.Println(result, result2)
	//parser.EndTransaction(tId)
	var test = parser.CombineResultEach([][]string{
		{"1", "2", "3"},
		{"4", "5", "6", "7", "8"},
	}, parser.OPERATOR_MERGE)
	fmt.Println(test)

}

func Regexp() {
	var test = parser.RegexpFilter([]string{"测试内容,测试内容，测试"}, "##测试")
	fmt.Println(test)
	var test2 = parser.RegexpFilter([]string{"测试内容,测屁内容，测试"}, "##测(.)##体$1")
	fmt.Println(test2)
	var test3 = parser.RegexpFilter([]string{"测试内容,测试内容，测试"}, "##(.)试##替换$1的###")
	fmt.Println(test3)
}
