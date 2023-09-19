# go-hbase
go 读取hbase数据
# hbase thrift2运行
```shell
bin/hbase-daemon.sh start thrift2
```
# 打包
## linux
```shell
CGO_ENABLED=1 go build -ldflags '-linkmode "external" -extldflags "-static"' main.go
```
## docker
打包
```shell
docker build . -t dounine/go-hbase
```
docker运行
```shell
docker run -e domain=server:9090 -e maxLength=100 -p 8000:8000 dounine/go-hbase
```
# 访问
```shell
curl --location --request GET 'localhost:8000/user' \
--header 'Content-Type: application/json' \
--data-raw '{
    "uid":"",
    "uids":["1","2"]
}'
```