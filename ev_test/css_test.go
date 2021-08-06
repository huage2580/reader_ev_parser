package ev_test

import (
	"fmt"
	"reader_ev_parser/parser"
)

func CSS() {
	testCssList("@css:h1@text")
	cssBatch("@css:#list dd", "@css:a@href", "text")
}

func testCssList(rule string) {
	tId := parser.StartTransaction(DATA)
	result := parser.ParseRuleStr(tId, rule)
	fmt.Printf("css rule-> [%s] result-> %s\n", rule, result)
	parser.EndTransaction(tId)
}

func cssBatch(rule string, rule1 string, rule2 string) {
	tId := parser.StartTransaction(DATA)
	bId := parser.ParseRuleRaw(tId, rule)
	size := parser.QueryBatchResultSize(bId)
	fmt.Printf("css rule-> [%s] result-> %d\n", rule, size)
	for i := 0; i < size; i++ {
		r1 := parser.ParseRuleStrForParent(bId, rule1, i)
		fmt.Printf("css rule-> [%s] result-> %s\n", rule1, r1)
		r2 := parser.ParseRuleStrForParent(bId, rule2, i)
		fmt.Printf("css rule-> [%s] result-> %s\n", rule2, r2)
	}
	parser.EndTransaction(bId)
	parser.EndTransaction(tId)
}
