package module

import (
	"Basil/structs"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func CheckRefundsMtcodeRecode(mtc string) (amount float64, currency string, CheckMtcodeRecodStatus bool) {
	var tx structs.Transaction_mgolog
	mongoClient, err := GetMgoCli()
	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	if err != nil {
		panic(err)
	}

	collection := mongoClient.Database("transaction").Collection("transaction_log")
	fmt.Println(mtc)
	opts1 := options.FindOne().SetProjection((bson.D{{"event.$", 1}, {"target.account", 1}, {"balance", 1}, {"currency", 1}}))

	err = collection.FindOne(
		context.TODO(),
		bson.D{{"event.mtcode", mtc}},
		opts1,
	).Decode(&tx)

	//先查詢該mtcode是否不存在，不存在的話回餘額amout=0，讓controller繼續動作
	if err == mongo.ErrNoDocuments {
		log.Println("Refunds：該mtcode是新的一筆")
		return 0, "", false
	}

	//先查詢該mtcode是否已被退款，不是的話回餘額amout=0，讓controller繼續動作
	if tx.Event[0].Status != "refund" {
		log.Println("Refunds：該mtcode是新的一筆")
		return 0, "", false

	}
	log.Println("Refunds：該mtcode已被refund，讀取mongodb裡的紀錄")
	log.Println(tx.Balance)
	return tx.Balance, tx.Currency, true

}
