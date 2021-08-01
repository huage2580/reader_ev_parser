package test

import (
	"fmt"
	"reader_ev_parser/parser"
)

func Jsoup() {
	tId := parser.StartTransaction(DATA)
	result := parser.ParseRuleStr(tId, "tag.p.0@text")
	result2 := parser.ParseRuleStr(tId, `href@js:result+',{webView:\"true\"}'`)
	fmt.Println(result, result2)
	parser.EndTransaction(tId)
}
