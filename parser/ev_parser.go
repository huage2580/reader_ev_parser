package parser

//--------解析器实现---------

import uuid "github.com/satori/go.uuid"

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

// ParseRuleRaw 解析返回html文本
func ParseRuleRaw(tId string, rule string) []string {
	return factoryTid(tId).parse(rule, false)
}

// ParseRuleStr 解析返回字符串
func ParseRuleStr(tId string, rule string) []string {
	return factoryTid(tId).parse(rule, true)
}

// EndTransaction 结束会话，清空网页缓存
func EndTransaction(tId string) {
	delete(transactionCacheMap, tId)
}

func factory(input string) ActionParser {
	//todo 根据类型指派解析器
	return ActionParser{
		action:    JsoupAction{},
		inputData: input,
	}
}

func factoryTid(tId string) ActionParser {
	return factory(transactionCacheMap[tId].(string))
}
