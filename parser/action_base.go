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

type IActionParser struct {
	action    Action
	inputData interface{}
}

func (parser IActionParser) parse(rule string, needFilterString bool) []string {
	var result = make([]string, 0, 99999)
	//处理倒叙
	var needReverse = strings.HasPrefix(rule, "-")
	//过滤js规则，在最后面
	var ruleWithoutJS, jsScript = formatJs(rule)
	var s = fmt.Sprintf("ruleWithoutJS->%s | jsScript->%s", ruleWithoutJS, jsScript)
	fmt.Println(s)
	//todo 过滤正则净化，标记后面处理
	//todo 切割操作符 && || %% 组合条件
	//todo 单独执行规则
	//todo 合并结果集
	//todo 正则净化结果
	//todo 执行js？
	//反转列表
	if needReverse {
		reverseArray(result)
	}
	return result
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
