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
