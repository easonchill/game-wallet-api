package module

import (
	"Basil/structs"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func CheckRefundsMtcodeRecode(mtc string) (amount float64, currency string, CheckMtcodeRecodStatus bool) {
	var tx structs.TransactionMgoLog
	var mtcodeRecord structs.MtcodeBalanceMgoLog

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

	collection = mongoClient.Database("transaction").Collection("mtcodeRecord")

	opts := options.FindOne().SetSkip(0)

	findMtcode := structs.FindByMtcodeRecord{Mtcode: mtc}

	err = collection.FindOne(
		context.TODO(),
		findMtcode,
		opts,
	).Decode(&mtcodeRecord)

	return mtcodeRecord.Balance, tx.Currency, true

}

func SaveMtcodeRecord(mtcode string, action string, balance float64) (e error) {

	t := structs.MtcodeBalanceMgoLog{
		Id:      primitive.NewObjectID(),
		Mtcode:  mtcode,
		Action:  action,
		Balance: balance,
	}

	mongoClient, err := GetMgoCli()

	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	if err != nil {
		panic(err)
	}

	collection := mongoClient.Database("transaction").Collection("mtcodeRecord")
	_, err = collection.InsertOne(context.TODO(), t)

	if err != nil {
		panic(err)
	}

	return err
}
