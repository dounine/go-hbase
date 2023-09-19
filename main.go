package main

import (
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"go-hbase/hbase"
	"time"
)

func main() {
	ctx := context.Background()
	protocolFactory := thrift.NewTBinaryProtocolFactoryConf(&thrift.TConfiguration{
		ConnectTimeout: 3 * time.Second,
		SocketTimeout:  3 * time.Second,
	})
	transport := thrift.NewTSocketConf("47.110.140.35:9090", &thrift.TConfiguration{})
	client := hbase.NewTHBaseServiceClientFactory(transport, protocolFactory)
	if err := transport.Open(); err != nil {
		panic(err)
	}
	userTable := []byte("USER_TABLE_V3")

	result, err := client.Get(ctx, userTable, &hbase.TGet{Row: []byte("85eebccd895c279a1d9f5853c20fe874")})
	if err != nil {
		panic(err)
	}
	for _, column := range result.ColumnValues {
		fmt.Printf("%s:%s\n", column.Qualifier, column.Value)
	}
}
