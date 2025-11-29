package db

import (
	"context"
	"github.com/StephenChristianW/go-movies-open/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"time"
)

const (
	dbLog = "[DB] "
)

var (
	clientInstance *mongo.Client
	clientOnce     sync.Once
)

// DefaultTimeout 默认超时时间
const DefaultTimeout = 10 * time.Second

// GetCtx 返回一个带超时的 context
func GetCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultTimeout)
}

// MongoClient 获取全局 MongoDB 客户端
func MongoClient() func() {
	clientOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(config.DBUrl)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatalf(dbLog+"MongoDB 连接失败: %v", err)
		}

		// 测试连通性
		if err := client.Ping(ctx, nil); err != nil {
			log.Fatalf(dbLog+"MongoDB Ping 失败: %v", err)
		}

		clientInstance = client
		log.Println(dbLog + "MongoDB 连接成功")
	})
	return func() {
		if clientInstance == nil {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := clientInstance.Disconnect(ctx); err != nil {
			log.Printf(dbLog+"关闭 MongoDB 失败: %v\n", err)
		} else {
			log.Println(dbLog + "MongoDB 已关闭")
		}
	}
}

// GetStackBuilderCollection 获取StackBuilder的集合
func GetStackBuilderCollection(collName string) *mongo.Collection {
	return clientInstance.Database(config.DBName).Collection(collName)
}

// CloseMongoClient 关闭全局 MongoDB 客户端
