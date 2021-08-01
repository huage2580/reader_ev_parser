package parser

import (
	"fmt"
	"regexp"
	"strings"
)

/**
这里实现通用的解析流程
*/

type Action interface {
	parseEach(rule string, needFilterString bool) []string
	formatRule(rule string) string
}

type ActionParser struct {
	action    Action
	inputData interface{}
}

func (parser ActionParser) parse(rule string, needFilterString bool) []string {
	var result = make([]string, 0, 99999)
	//处理倒叙
	var needReverse = strings.HasPrefix(rule, "-")
	// 过滤分类用的前缀
	var ruleWithoutPrefix = parser.action.formatRule(rule)
	//过滤js规则，在最后面
	var ruleWithoutJS, jsScript = formatJs(ruleWithoutPrefix)
	var s = fmt.Sprintf("ruleWithoutJS->%s | jsScript->%s", ruleWithoutJS, jsScript)
	fmt.Println(s)
	//过滤正则净化，标记后面处理
	var ruleWithoutRegexp, regexpRule = formatRegexp(ruleWithoutJS)
	fmt.Printf("ruleWithoutRegexp->%s | regexpRule->%s\n", ruleWithoutRegexp, regexpRule)
	//切割操作符 && || %% 组合条件
	var ruleList, opMode = splitOperator(ruleWithoutRegexp)
	//单独执行规则
	var resultList = make([][]string, 2)
	for index, ruleEach := range ruleList {
		var resultEach = parser.action.parseEach(ruleEach, needFilterString)
		resultList[index] = resultEach
		//或的操作，有数据不执行后面的
		if len(resultEach) > 0 && opMode == OPERATOR_OR {
			break
		}
	}
	//合并结果集
	var resultComb = CombineResultEach(resultList, opMode)
	//正则净化结果
	var resultAfterRegexp = regexpFilter(resultComb, regexpRule)
	//maybe 执行js？
	//反转列表
	if needReverse {
		result = reverseArray(resultAfterRegexp)
	}
	return result
}

// 按照组合操作符分割规则
func splitOperator(input string) ([]string, string) {
	indexAnd := strings.Index(input, OPERATOR_AND)
	indexOr := strings.Index(input, OPERATOR_OR)
	indexMerge := strings.Index(input, OPERATOR_MERGE)
	if indexAnd > 0 {
		return strings.Split(input, OPERATOR_AND), OPERATOR_AND
	}
	if indexOr > 0 {
		return strings.Split(input, OPERATOR_OR), OPERATOR_OR
	}
	if indexMerge > 0 {
		return strings.Split(input, OPERATOR_MERGE), OPERATOR_MERGE
	}
	return []string{input}, ""
}

// CombineResultEach 组合结果
func CombineResultEach(input [][]string, opMode string) []string {
	switch opMode {
	case OPERATOR_AND:
		return append(input[0], input[1]...)
	case OPERATOR_OR:
		return append(input[0], input[1]...)
	case OPERATOR_MERGE:
		var result []string
		var maxLength = 0
		var length1 = len(input[0])
		var length2 = len(input[1])
		if length1 > length2 {
			maxLength = length1
		} else {
			maxLength = length2
		}
		for i := 0; i < maxLength; i++ {
			if length1 > i {
				result = append(result, input[0][i])
			}
			if length2 > i {
				result = append(result, input[1][i])
			}
		}
		return result

	}
	return input[0]
}

// 正则净化结果
func regexpFilter(input []string, regexpRule string) []string {
	if regexpRule == "" {
		return input
	}
	//todo
	return nil
}

// 反转数组
func reverseArray(x []string) []string {
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}
	return x
}

//过滤js脚本
func formatJs(input string) (string, string) {
	re, _ := regexp.Compile(RE_JS_TOKEN)
	var output = re.ReplaceAllString(input, "")
	var jsScript = re.FindStringSubmatch(input)
	var js string
	if len(jsScript) > 1 {
		js = jsScript[1]
	}
	return output, js
}

//过滤##正则净化规则
func formatRegexp(input string) (string, string) {
	index := strings.Index(input, RE_REPLACE)
	if index > 0 {
		var output = input[0:index]
		var regexpRule = input[index+len(RE_REPLACE):]
		return output, regexpRule
	}
	return input, ""
}
