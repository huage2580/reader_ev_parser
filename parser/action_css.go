package parser

import (
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"log"
	"regexp"
	"strings"
)

type CSSAction struct {
	Action
}

func (action CSSAction) parseEach(input interface{}, rule string, needFilterString bool) []interface{} {
	switch input.(type) {
	case string:
		if rule == "" {
			return []interface{}{input}
		}
		return parseHtmlCss(input.(string), rule, needFilterString)
	case *goquery.Selection:
		var node = input.(*goquery.Selection)
		return parseDocumentCss(node, rule, needFilterString)
	case *html.Node: //预留给xpath的
		var nde = input.(*html.Node)
		return parseDocumentCss(goquery.NewDocumentFromNode(nde).Find("body").Children().First(), rule, needFilterString)
	default:
		//not support
	}
	return []interface{}{}
}

func (action CSSAction) formatRule(rule string) string {
	re, _ := regexp.Compile(PARSER_TYPE_CSS)
	return re.ReplaceAllString(rule, "")
}

func (action CSSAction) getType() string {
	return ACTION_TYPE_CSS
}

func parseHtmlCss(html string, rule string, needFilterString bool) []interface{} {
	document, e := goquery.NewDocumentFromReader(strings.NewReader(html))
	if e != nil {
		log.Fatal(e)
		return []interface{}{}
	}
	return parseDocumentCss(document.Selection, rule, needFilterString)
}

func parseDocumentCss(node *goquery.Selection, rule string, needFilterString bool) []interface{} {
	//分割文本类型,如果是获取文本，最后一个就是要求返回的文本结果，除了常用匹配，就是返回attr属性值
	var ruleList = strings.Split(rule, DELIMITER)
	var stringFilterValue string
	var execRuleList = ruleList
	if needFilterString { //需要返回文本的情况
		stringFilterValue = ruleList[len(ruleList)-1]
		execRuleList = ruleList[0 : len(ruleList)-1]
	}

	var result = make([]*goquery.Selection, 0)
	//执行查询
	var find = node.Find(execRuleList[0])
	var pResults = find.Nodes
	for i := range pResults {
		var s = find.Eq(i)
		result = append(result, s)
	}

	return selectionToInterface(result, needFilterString, stringFilterValue)

}
