# go-hbase
go 读取hbase数据
# hbase thrift2运行
```shell
bin/hbase-daemon.sh start thrift2
```
# 打包
```shell
CGO_ENABLED=1 go build -ldflags '-linkmode "external" -extldflags "-static"' main.go
```