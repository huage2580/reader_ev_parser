package parser

import (
	"regexp"
	"strings"
)

/**
这里实现通用的解析流程
*/

// BatchResult 缓存的结果集,用来批量操作的，用于解析目录列表和搜索列表页
type BatchResult struct {
	queryType string
	cacheData []interface{}
}

type Action interface {
	parseEach(input interface{}, rule string, needFilterString bool) []interface{}
	formatRule(rule string) string
	getType() string
}

type ActionParser struct {
	action    Action
	inputData interface{}
}

//[cacheId] 表示把结果存储到这里
func (parser ActionParser) parse(rule string, needFilterString bool, cacheId string) []string {
	//处理倒叙
	var needReverse = strings.HasPrefix(rule, "-")
	// 过滤分类用的前缀
	var ruleWithoutPrefix = parser.action.formatRule(rule)
	//过滤js规则，在最后面
	var ruleWithoutJS, _ = formatJs(ruleWithoutPrefix)

	//过滤正则净化，标记后面处理
	var ruleWithoutRegexp, regexpRule = formatRegexp(ruleWithoutJS)
	//fmt.Printf("ruleWithoutRegexp->%s | regexpRule->%s\n", ruleWithoutRegexp, regexpRule)

	//切割操作符 && || %% 组合条件
	var ruleList, opMode = splitOperator(ruleWithoutRegexp)
	//单独执行规则
	var resultList = make([][]interface{}, len(ruleList))
	for index, ruleEach := range ruleList {
		var resultEach = parser.action.parseEach(parser.inputData, ruleEach, needFilterString)
		resultList[index] = resultEach
		//或的操作，有数据不执行后面的
		if len(resultEach) > 0 && opMode == OPERATOR_OR {
			break
		}
	}
	//合并结果集
	var resultComb = CombineResultEach(resultList, opMode)
	//缓存结果集,在作为批量操作的时候
	if cacheId != "" {
		CacheTransactionData(cacheId, BatchResult{queryType: parser.action.getType(), cacheData: resultComb})
	}
	//正则净化结果
	var resultAfterRegexp = make([]string, 0)
	if needFilterString {
		var resultStr = make([]string, len(resultComb))
		for i, item := range resultComb {
			resultStr[i] = item.(string)
		}
		resultAfterRegexp = RegexpFilter(resultStr, regexpRule)
	}
	//maybe 执行js？
	//反转列表
	if needReverse {
		resultAfterRegexp = reverseArray(resultAfterRegexp)
	}
	//只有要求输出文本的才有结果集,原生的缓存起来等待批量解析
	return resultAfterRegexp
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
func CombineResultEach(input [][]interface{}, opMode string) []interface{} {
	//
	var result = input[0]
	switch opMode {
	case OPERATOR_AND:
		for _, item := range input {
			result = append(result, item...)
		}
		return result
	case OPERATOR_OR:
		for _, item := range input {
			result = append(result, item...)
		}
		return result
	case OPERATOR_MERGE:
		var maxLength = 0
		for _, item := range input {
			len1 := len(item)
			if len1 > maxLength {
				maxLength = len1
			}
		}
		for i := 0; i < maxLength; i++ {
			for ind, i2 := range input {
				if len(i2) > i {
					result = append(result, input[ind][i])
				}
			}
		}
		return result

	}
	return result
}

// RegexpFilter 正则净化结果
func RegexpFilter(input []string, regexpRule string) []string {
	if regexpRule == "" {
		return input
	}
	reList := strings.Split(regexpRule, RE_REPLACE)
	length := len(reList)
	var rule, replace string
	var onlyOne = false
	rule = reList[1]
	if length == 2 { //
		replace = ""
	} else if length == 3 { //
		replace = reList[2]
	} else if length == 4 { //
		replace = reList[2]
		onlyOne = true
	}
	//把 $1 换成 ${1}的格式，go对前者兼容很差
	realReplace := regexp.MustCompile(`\$(\d*)`).ReplaceAllString(replace, `${$1}`)
	for i, str := range input {
		re := regexp.MustCompile(rule)
		onlyOneFlag := false
		output := re.ReplaceAllStringFunc(str, func(a string) string {
			if onlyOneFlag && onlyOne {
				return a
			}
			onlyOneFlag = true
			return re.ReplaceAllString(a, realReplace)
		})
		input[i] = output
	}
	return input
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
	if index >= 0 {
		var output = input[0:index]
		var regexpRule = input[index:]
		return output, regexpRule
	}
	return input, ""
}
