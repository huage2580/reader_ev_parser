dart FFI 参考
```
package main
import 'dart:ffi' as ffi;
import 'package:ffi/ffi.dart';
import 'dart:io';
import 'package:path/path.dart' as path;
import 'parser_bindings.dart' as binding;
void main(){
  print("hello world");
  var libraryPath = path.join(Directory.current.path,'share','godll.dll');
  final dylib = ffi.DynamicLibrary.open(libraryPath);
  var lib = binding.EVParser(dylib);
  lib.HelloWorld();
  var s =lib.Add(6, 3);
  var dartStringPointer = ffi.Pointer<Utf8>.fromAddress(s.address);
  print(dartStringPointer.toDartString());
  //not need free? calloc.free(s);
  var sIn = 'plain old C string'.toNativeUtf8().cast<ffi.Int8>();
  lib.PrintText(sIn);
  calloc.free(sIn);
  var pArray = lib.SliceTest();
  var p1 = pArray.elementAt(0).value.toDartString();
  var p2 = pArray.elementAt(1).value.toDartString();
  var pp = pArray.elementAt(2).value.address;//address == 0就是没数据了
  print("${p1},${p2},${pp}");
}
```

```go
package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

import "C"

func main() {

}
//export HelloWorld
func HelloWorld(){
	fmt.Println("hello world from golang")
}

//export PrintText
func PrintText(cstr *C.char) {
	var s = C.GoString(cstr)
	fmt.Println("golang receipt->" + s)
}

//export Add
func Add(a int, b int) *C.char {
	var r = fmt.Sprintf("%d+%d=%d",a,b,a+b)
	return C.CString(r)
}

//export SliceTest
func SliceTest() uintptr {
	return stringSliceToC([]string{"hello","golang"})
}

//先转换GoString为CString,再拼接数组
func stringSliceToC(input []string) uintptr{
	arr := make([]*C.char,len(input))
	for i, s := range input {
		arr[i] = C.CString(s)
	}
	ptr := (*reflect.SliceHeader)(unsafe.Pointer(&arr))
	//ptr := unsafe.Pointer(&arr[0])
	return ptr.Data
}
```

```
# dart run ffigen

ffigen:
  output: 'parser_bindings.dart'
  name: 'EVParser'
  headers:
    entry-points:
      - 'share/godll.h'
```

