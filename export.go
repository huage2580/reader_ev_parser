package main

/*

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"log"
	"reader_ev_parser/parser"
	"reflect"
	"unsafe"
)

//export StartTransaction
func StartTransaction(inputString *C.char) *C.char {
	r := parser.StartTransaction(toGoString(inputString))
	return toCString(r)
}

//export ParseRuleRaw
func ParseRuleRaw(tId *C.char, rule *C.char) *C.char {
	var r = ""
	TryWithLog(func() {
		r = parser.ParseRuleRaw(toGoString(tId), toGoString(rule))
	})
	return toCString(r)
}

//export ParseRuleStr
func ParseRuleStr(tId *C.char, rule *C.char) **C.char {
	var r []string
	TryWithLog(func() {
		r = parser.ParseRuleStr(toGoString(tId), toGoString(rule))
	})
	return stringSliceToC(r)
}

//export ParseRuleStrForParent
func ParseRuleStrForParent(tId *C.char, rule *C.char, index int) **C.char {
	var r []string
	TryWithLog(func() {
		r = parser.ParseRuleStrForParent(toGoString(tId), toGoString(rule), index)
	})
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
func stringSliceToC(input []string) **C.char {
	arr := make([]*C.char, len(input))
	for i, s := range input {
		arr[i] = C.CString(s)
	}
	ptr := (*reflect.SliceHeader)(unsafe.Pointer(&arr))
	p1 := unsafe.Pointer(ptr.Data)
	//copy 一份给c用，slice的数据在大量操作的时候，会触发go的gc，导致数据被回收
	var sizeOfPoint = unsafe.Sizeof(p1)
	var size = (int(sizeOfPoint)) * (len(arr) + 1) //+1是为了后面放个空指针，数据完结
	var sizeLong = C.size_t(size)
	p2 := C.malloc(sizeLong)
	C.memcpy(p2, p1, sizeLong)
	//C.test_print((**C.char)(p2))
	return (**C.char)(p2)
}

//实现 try catch

func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

func TryWithLog(fun func()) {
	Try(fun, func(i interface{}) {
		log.Fatal(i)
	})
}
