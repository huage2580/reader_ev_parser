package parser

import (
	"golang.org/x/net/html"
	"log"
	"strconv"
	"strings"
)
import "github.com/PuerkitoBio/goquery"

// JsoupAction Jsoup解析
type JsoupAction struct {
	Action
}

func (action JsoupAction) parseEach(input interface{}, rule string, needFilterString bool) []interface{} {
	switch input.(type) {
	case string:
		if rule == "" {
			return []interface{}{input}
		}
		return parseHtmlJsoup(input.(string), rule, needFilterString)
	case *goquery.Selection:
		var node = input.(*goquery.Selection)
		return parseDocumentJsoup(node, rule, needFilterString)
	case *html.Node: //预留给xpath的
		var nde = input.(*html.Node)
		return parseDocumentJsoup(goquery.NewDocumentFromNode(nde).Find("body").Children().First(), rule, needFilterString)
	default:
		//not support
	}
	return []interface{}{}
}

func (action JsoupAction) formatRule(rule string) string {
	if strings.HasPrefix(rule, "-") || strings.HasPrefix(rule, "+") {
		return rule[1:]
	}
	return rule
}

func (action JsoupAction) getType() string {
	return ACTION_TYPE_JSOUP
}

func parseHtmlJsoup(html string, rule string, needFilterString bool) []interface{} {
	document, e := goquery.NewDocumentFromReader(strings.NewReader(html))
	if e != nil {
		log.Fatal(e)
		return []interface{}{}
	}
	return parseDocumentJsoup(document.Selection, rule, needFilterString)
}

func parseDocumentJsoup(node *goquery.Selection, rule string, needFilterString bool) []interface{} {
	//分割文本类型,如果是获取文本，最后一个就是要求返回的文本结果，除了常用匹配，就是返回attr属性值
	var ruleList = strings.Split(rule, DELIMITER)
	var stringFilterValue string
	var execRuleList = ruleList
	if needFilterString { //需要返回文本的情况
		stringFilterValue = ruleList[len(ruleList)-1]
		execRuleList = ruleList[0 : len(ruleList)-1]
	}

	var result = make([]*goquery.Selection, 0)
	var currentNodes = []*goquery.Selection{node}
	//逐条执行命令
	for _, r := range execRuleList {
		cssQuery, includeIndex, excludeIndex := remapToCssQuery(r)
		var nodesInRound = make([]*goquery.Selection, 0)
		for _, currentNode := range currentNodes {
			nodesInRound = make([]*goquery.Selection, 0) //清空结果
			var find *goquery.Selection
			if cssQuery == JSOUP_SUPPORT_CHILD { //单独处理获取子元素
				find = currentNode.Children()
			} else {
				find = currentNode.Find(cssQuery)
			}
			var pResults = find.Nodes
			for i := range pResults {
				var s = find.Eq(i)
				nodesInRound = append(nodesInRound, s)
			}
		}
		//保留或过滤指定序号
		currentNodes = filterIndex(nodesInRound, includeIndex, excludeIndex)
	}
	result = currentNodes

	return selectionToInterface(result, needFilterString, stringFilterValue)
}

//转换结果为interface，或者string
func selectionToInterface(result []*goquery.Selection, needFilterString bool, stringFilterValue string) []interface{} {
	//结果集
	var output = make([]interface{}, len(result))

	//过滤需要的文本
	if needFilterString {
		temp := filterText(result, stringFilterValue)
		for i, item := range temp {
			output[i] = item
		}
	} else {
		//转换interface
		for i, r := range result {
			output[i] = r
		}
	}
	return output
}

func remapToCssQuery(r string) (string, []int, []int) {
	var queryWithIndex = strings.Split(r, JSOUP_EXCLUDE_CHAR)
	var querySplit = strings.Split(queryWithIndex[0], JSOUP_SPLIT)

	var excludeStr = ""
	if len(queryWithIndex) > 1 {
		excludeStr = queryWithIndex[1]
	}
	//---------------------------------
	var css = ""
	var aType = querySplit[0] //类型
	var aValue = ""           // 值
	var includeStr = ""       //筛选序号 1:2:3:4:5
	if len(querySplit) > 1 {
		aValue = querySplit[1]
	}
	if len(querySplit) > 2 {
		includeStr = querySplit[2]
	}
	//根据类型生成css
	switch aType {
	case JSOUP_SUPPORT_CHILD:
		css = JSOUP_SUPPORT_CHILD
	case JSOUP_SUPPORT_CLASS:
		css = "." + aValue
	case JSOUP_SUPPORT_TAG:
		css = aValue
	case JSOUP_SUPPORT_ID:
		css = "#" + aValue
	case JSOUP_SUPPORT_TEXT:
		css = ":containsOwn(" + aValue + ")"
	default:
		css = aType
	}
	return css, indexStringToArray(includeStr), indexStringToArray(excludeStr)
}

func filterText(nodes []*goquery.Selection, stringFilterValue string) []string {
	var result = make([]string, len(nodes))
	for i, node := range nodes {
		doc := mapText(node, stringFilterValue)
		result[i] = doc
	}
	return result
}

// 获取指定的文本结果
func mapText(node *goquery.Selection, clazz string) string {
	switch clazz {
	case FILTER_HTML:
		var out = ""
		node.Contents().Each(func(i int, s *goquery.Selection) {
			if goquery.NodeName(s) == "br" {
				out = out + "\n"
			} else {
				out = out + s.Text()
			}
		})
		return out
	case FILTER_TEXT:
		return strings.TrimSpace(node.Text())
	case FILTER_OWN_TEXT:
		var out = ""
		node.Contents().Each(func(i int, s *goquery.Selection) {
			if goquery.NodeName(s) == "#text" {
				out = out + s.Text()
			}
		})
		return out
	case FILTER_TEXT_NODE:
		var out = ""
		node.Contents().Each(func(i int, s *goquery.Selection) {
			out = out + strings.TrimSpace(s.Text()) + "\n"
		})
		return out
	case FILTER_ALL:
		out, _ := node.Html()
		return out
	case FILTER_TAG_NAME:
		return goquery.NodeName(node)
	default:
		return node.AttrOr(clazz, "null")
	}
}

func filterIndex(nodes []*goquery.Selection, include []int, exclude []int) []*goquery.Selection {
	var result = make([]*goquery.Selection, 0)
	if len(nodes) == 0 {
		return result
	}
	//copy
	if len(exclude) == 0 && len(include) == 0 {
		for _, node := range nodes {
			result = append(result, node)
		}
	}
	for i, node := range nodes {
		//--排除
		if len(exclude) > 0 && contains(exclude, negative(i, len(nodes))) {
			continue
		} else if len(exclude) > 0 {
			result = append(result, node)
		}

	}
	//指定
	for _, c := range include {
		result = append(result, nodes[negative(c, len(nodes))])
	}
	return result
}

func contains(s []int, str int) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

//转换负数,取正数
func negative(i int, len int) int {
	if i >= 0 {
		return i
	} else {
		return len + i
	}
}

// 1:2:3:4:5 转换数组
func indexStringToArray(input string) []int {
	if input == "" {
		return []int{}
	}
	var sp = strings.Split(input, JSOUP_EXCLUDE_INT)
	var result = make([]int, len(sp))
	for i, s := range sp {
		result[i], _ = strconv.Atoi(s)
	}
	return result
}
