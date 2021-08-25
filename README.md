# ipTools
纯真ip数据库go语言解析器

0.复制qqwry.dat到项目data目录下

1.帮助
```
输入例示:
  -i string
        input you ip
  -w string
        web port
```
2.演示:
```
运行：
go run main.go -i 119.147.146.189
输出：
{"ip":"119.147.146.189","country":"中国","province":"广东省","city":"广东市","county":"广东县","isp":"","area":"中国广东省东莞市省东莞市"}

```

3.函数测试：
```cassandraql
goos: windows
goarch: amd64
pkg: github.com/Bmixo/ipTools
BenchmarkAll-8              8594            130904 ns/op
PASS
ok      github.com/Bmixo/ipTools        1.190s

```