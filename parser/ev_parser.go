package parser

//--------解析器实现---------

import (
	uuid "github.com/satori/go.uuid"
	"regexp"
)

var transactionCacheMap = make(map[string]interface{})

// StartTransaction 为了不频繁传输大文本，采取缓存
func StartTransaction(inputString string) string {
	var tId = uuid.NewV4().String()
	transactionCacheMap[tId] = inputString
	return tId
}

// InjectArg 当前会话注入参数
func InjectArg(tId, aKey, aValue string) {
	//todo 基本上是给js应用的，暂时不考虑
}

func CacheTransactionData(tId string, data interface{}) {
	transactionCacheMap[tId] = data
}

// ParseRuleRaw 解析返回新的会话id?,完全是为了批量解析
func ParseRuleRaw(tId string, rule string) string {
	var newId = uuid.NewV4().String()
	factoryTid(tId, rule).parse(rule, false, newId)
	return newId
}

// ParseRuleStr 解析返回字符串
func ParseRuleStr(tId string, rule string) []string {
	return factoryTid(tId, rule).parse(rule, true, "")
}

// ParseRuleStrForParent 解析返回字符串
func ParseRuleStrForParent(tId string, rule string, index int) []string {
	return factoryForParent(tId, rule, index).parse(rule, true, "")
}

func QueryBatchResultSize(tId string) int {
	return len(transactionCacheMap[tId].(BatchResult).cacheData)
}

// EndTransaction 结束会话，清空缓存,在批量操作的时候，还要注意清除批量的会话
func EndTransaction(tId string) {
	delete(transactionCacheMap, tId)
}

func factory(input interface{}, rule string, parentQueryType string) ActionParser {
	var action Action
	switch judgeQueryType(rule, parentQueryType) {
	case ACTION_TYPE_JSOUP:
		action = JsoupAction{}
	case ACTION_TYPE_CSS:
		action = CSSAction{}
	default:
		action = JsoupAction{}
	}
	return ActionParser{
		action:    action,
		inputData: input,
	}

}

func factoryTid(tId string, rule string) ActionParser {
	return factory(transactionCacheMap[tId], rule, "")
}

func factoryForParent(tId string, rule string, index int) ActionParser {
	var batchResult = transactionCacheMap[tId].(BatchResult)
	return factory(batchResult.cacheData[index], rule, batchResult.queryType)
}

func judgeQueryType(rule string, parentQueryType string) string {
	//todo 根据类型指派解析器,先判断自己属于什么类型,然后结合父类型
	var matchCss, _ = regexp.MatchString(PARSER_TYPE_CSS, rule)
	if matchCss {
		return ACTION_TYPE_CSS
	}
	return ACTION_TYPE_JSOUP
}
