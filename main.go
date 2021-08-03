package main

/**
* 用于导出函
 */

import (
	"reader_ev_parser/parser"
	"reader_ev_parser/test"
	"reflect"
	"unsafe"
)

import "C"

func main() {
	//test.Regexp()
	//test.CSS()
	test.Jsoup()
}

//export StartTransaction
func StartTransaction(inputString *C.char) *C.char {
	r := parser.StartTransaction(toGoString(inputString))
	return toCString(r)
}

//export ParseRuleRaw
func ParseRuleRaw(tId *C.char, rule *C.char) *C.char {
	r := parser.ParseRuleRaw(toGoString(tId), toGoString(rule))
	return toCString(r)
}

//export ParseRuleStr
func ParseRuleStr(tId *C.char, rule *C.char) uintptr {
	r := parser.ParseRuleStr(toGoString(tId), toGoString(rule))
	return stringSliceToC(r)
}

//export ParseRuleStrForParent
func ParseRuleStrForParent(tId *C.char, rule *C.char, index int) uintptr {
	r := parser.ParseRuleStrForParent(toGoString(tId), toGoString(rule), index)
	return stringSliceToC(r)
}

//export QueryBatchResultSize
func QueryBatchResultSize(tId *C.char) int {
	return parser.QueryBatchResultSize(toGoString(tId))
}

//export EndTransaction
func EndTransaction(tId *C.char) {
	parser.EndTransaction(toGoString(tId))
}

//------------------------------------------------------------------------
func toCString(input string) *C.char {
	return C.CString(input)
}

func toGoString(input *C.char) string {
	return C.GoString(input)
}

//先转换GoString为CString,再拼接数组,返回的是指针,对应dart Point<Point<Utf8>> **char
func stringSliceToC(input []string) uintptr {
	arr := make([]*C.char, len(input))
	for i, s := range input {
		arr[i] = C.CString(s)
	}
	ptr := (*reflect.SliceHeader)(unsafe.Pointer(&arr))
	return ptr.Data
}

//----------------------------------------------------------------------
