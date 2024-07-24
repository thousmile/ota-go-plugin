# golang plugin

## 注意

1.不支持 Windows 系统
2.不支持 调试模式
3.不支持 go-delve

## linux
```shell
mkdir plugins
go build -buildmode=plugin -o v1.so ./extensions/v1
go build -buildmode=plugin -o v2.so ./extensions/v2
go build -o main ./primary

./main

2024/07/24 21:12:35 INFO launcher listen http server  addr=:18086
2024/07/24 21:12:35 INFO v1 listen http server  addr=:18081

```
