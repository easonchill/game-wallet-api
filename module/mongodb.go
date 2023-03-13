package module

import (
	"context"
	"fmt"
	"game-wallet-api/env"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	client, _ := GetMgoCli()
	fmt.Println("進行MongoDB連線檢查")
	if err := client.Ping(context.TODO(), nil); err != nil {
		panic(err.Error())
	}
	fmt.Println("MongoDB連線檢查完成")
}

func GetMgoCli() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(env.MongodbConfig))

	return client, err
}
