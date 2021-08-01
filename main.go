package main

/**
* 用于导出函数
todo 入参、出参、是不是要转换成C的
*/

//import "C"
import (
	"reader_ev_parser/parser"
	"reader_ev_parser/test"
)

func main() {
	test.Jsoup()
}

func StartTransaction(inputString string) string {
	return parser.StartTransaction(inputString)
}

func InjectArg(tId, aKey, aValue string) {
	parser.InjectArg(tId, aKey, aValue)
}

func ParseRuleRaw(tId string, rule string) []string {
	return parser.ParseRuleRaw(tId, rule)
}

func ParseRuleStr(tId string, rule string) []string {
	return parser.ParseRuleStr(tId, rule)
}

func EndTransaction(tId string) {
	parser.EndTransaction(tId)
}
