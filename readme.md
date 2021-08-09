## reader_ev_parser

使用go写阅读3.0格式的书源解析  

### 为什么是go？  
为了完全的跨端，调研JS和GO的生态，最后为了多应用一门语言，决定采用go。  

为啥不是dart，dart的生态实在太差劲，很多东西我没办法一个人完成。  

### 目标 
1. 服务端的可用性校验  
2. 通过go编译出ios,android,flutter,server,windows通用  
3. 为了保证可用兼容性，放弃对JS的支持,做到剩下60%书源的兼容  

##编译
### Android arm64(aarch64)
```shell
$env:GOOS="android"
$env:GOARCH="arm64"
$env:CGO_ENABLED="1"
$env:CC="C:\Users\hua\AppData\Local\Android\Sdk\ndk\22.0.7026061\toolchains\llvm\prebuilt\windows-x86_64\bin\aarch64-linux-android21-clang.cmd"
$env:CXX="C:\Users\hua\AppData\Local\Android\Sdk\ndk\22.0.7026061\toolchains\llvm\prebuilt\windows-x86_64\bin\aarch64-linux-android21-clang++.cmd"
go build -buildmode=c-shared -o libevparser.so
# 查看编译结果
readelf -h evparser.so
```
### windows
```shell
$env:GOOS=""
$env:GOARCH=""
$env:CGO_ENABLED="1"
$env:CC=""
$env:CXX=""
go build -buildmode=c-shared -o evparser.dll
```

### IOS
```shell
export GOOS="ios"
export GOARCH="arm64"
export CGO_ENABLED="1"
go build -buildmode=c-archive -o libevparser.a
# 查看便衣结果
lipo -info libevparser.a

```