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
	"sync"
	"unsafe"
)

var holder = sync.Map{} // 并发问题

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
func ParseRuleStr(tId *C.char, rule *C.char) uintptr {
	var r []string
	TryWithLog(func() {
		r = parser.ParseRuleStr(toGoString(tId), toGoString(rule))
	})
	return stringSliceToC(r)
}

//export ParseRuleStrForParent
func ParseRuleStrForParent(tId *C.char, rule *C.char, index int) uintptr {
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

//export FreeStr
func FreeStr(address *C.char) {
	C.free(unsafe.Pointer(address))
}

//export FreeSlice
func FreeSlice(address uintptr) {
	//fmt.Println(address)
	var v, ok = holder.Load(address)
	if ok {
		var a = v.([]*C.char)
		for _, char := range a {
			C.free(unsafe.Pointer(char))
		}
	}
	holder.Delete(address)
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
	//fmt.Println("go->",input)
	arr := make([]*C.char, len(input)+1)
	for i, s := range input {
		arr[i] = C.CString(s)
	}
	ptr := (*reflect.SliceHeader)(unsafe.Pointer(&arr))
	//做持有操作，避免被回收，后续调用free
	holder.Store(ptr.Data, arr)
	return ptr.Data
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
