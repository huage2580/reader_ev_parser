package parser

import "strings"

/**
Jsoup解析
*/

type JsoupAction struct {
	Action
}

func (action JsoupAction) parseEach(rule string, needFilterString bool) []string {
	return []string{"hello", "world"}
}

func (action JsoupAction) formatRule(rule string) string {
	if strings.HasPrefix(rule, "-") || strings.HasPrefix(rule, "+") {
		return rule[1:]
	}
	return rule
}
