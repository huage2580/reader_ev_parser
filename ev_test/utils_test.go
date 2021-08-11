package ev

import (
	"fmt"
	"reader_ev_parser/parser"
)

func JsoupBatchInput(input string, rule string, rule1 []string) {
	tId := parser.StartTransaction(input)
	bId := parser.ParseRuleRaw(tId, rule)
	size := parser.QueryBatchResultSize(bId)
	fmt.Printf("jsoup rule-> [%s] result-> %d\n", rule, size)
	for i := 0; i < size; i++ {
		for _, r := range rule1 {
			r1 := parser.ParseRuleStrForParent(bId, r, i)
			fmt.Printf("%s", r1)
		}
		fmt.Println("")
	}
	parser.EndTransaction(bId)
	parser.EndTransaction(tId)
}

func JsoupStrInput(input string, rule string) {
	tId := parser.StartTransaction(input)
	result := parser.ParseRuleStr(tId, rule)
	fmt.Printf("jsoup rule-> [%s] result-> \n", rule)
	for i, s := range result {
		fmt.Printf("%d->%s\n", i, s)
	}
	parser.EndTransaction(tId)
}
