package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-hbase/hbase"
	"os"
	"strconv"
	"time"
)

func MD5(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

type UserQuery struct {
	UID  *string  `json:"uid"`
	UIDS []string `json:"uids"`
}

func main() {
	domain := os.Getenv("domain")
	if domain == "" {
		domain = "127.0.0.1:9090"
	}
	maxLength := os.Getenv("maxLength")
	if maxLength == "" {
		maxLength = "100"
	}
	fmt.Println("hbase domain:", domain)
	protocolFactory := thrift.NewTBinaryProtocolFactoryConf(&thrift.TConfiguration{
		ConnectTimeout: 3 * time.Second,
		SocketTimeout:  3 * time.Second,
	})
	transport := thrift.NewTSocketConf(domain, &thrift.TConfiguration{})
	client := hbase.NewTHBaseServiceClientFactory(transport, protocolFactory)
	if err := transport.Open(); err != nil {
		panic(err)
	}
	userTable := []byte("USER_TABLE_V3")

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Token"},
		AllowCredentials: true,
	}))

	r.POST("/user", func(c *gin.Context) {
		ctx := c.Request.Context()
		var userQuery UserQuery
		err := c.ShouldBindJSON(&userQuery)
		if err != nil {
			c.JSON(200, gin.H{
				"msg": err.Error(),
			})
			return
		}
		if userQuery.UID != nil {
			result, err := client.Get(ctx, userTable, &hbase.TGet{Row: []byte(MD5(*userQuery.UID))})
			if err != nil {
				c.JSON(200, gin.H{
					"msg": err.Error(),
				})
				return
			}
			data := make(map[string]string)
			for _, column := range result.ColumnValues {
				data[string(column.Qualifier)] = string(column.Value)
			}
			c.JSON(200, gin.H{
				"code": 1,
				"data": data,
			})
		} else if userQuery.UIDS != nil {
			max, err := strconv.Atoi(maxLength)
			if err != nil {
				c.JSON(200, gin.H{
					"msg": "maxLength必须为数字",
				})
				return
			}
			if len(userQuery.UIDS) > max {
				c.JSON(200, gin.H{
					"msg": "uids最多" + maxLength + "个",
				})
				return
			}
			var gets []*hbase.TGet
			for _, uid := range userQuery.UIDS {
				gets = append(gets, &hbase.TGet{Row: []byte(MD5(uid))})
			}
			results, err := client.GetMultiple(ctx, userTable, gets)
			if err != nil {
				c.JSON(200, gin.H{
					"msg": err.Error(),
				})
				return
			}
			var datas []map[string]string
			for _, result := range results {
				maps := make(map[string]string)
				for _, column := range result.ColumnValues {
					maps[string(column.Qualifier)] = string(column.Value)
				}
				datas = append(datas, maps)
			}
			c.JSON(200, gin.H{
				"code": 1,
				"data": datas,
			})
		} else {
			c.JSON(200, gin.H{
				"msg": "uid或者uids必填其中一个参数",
			})
		}
	})

	r.Run(":8000")
}
