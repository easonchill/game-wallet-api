package module

import (
	"context"
	"game-wallet-api/structs"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CheckMtcodeRecode(mtc string) (amount float64, currency string, CheckMtcodeRecodeErr error) {
	var tx structs.TransactionMgoLog
	mongoClient, err := GetMgoCli()
	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	if err != nil {
		return 0, "", err
	}

	collection := mongoClient.Database("transaction").Collection("transaction_log")
	opts := options.FindOne().SetSkip(0)

	findMtcode := structs.FindByMtcode{Mtcode: mtc}
	err = collection.FindOne(
		context.TODO(),
		findMtcode,
		opts,
	).Decode(&tx)

	//先查詢該mtcode是否不存在，是的話回餘額amout=0,err報錯，讓controller繼續動作
	if err == mongo.ErrNoDocuments {
		log.Println("amend：該mtcode是新的一筆")
		return 0, "", err

	}
	log.Println("amend：該mtcode已存在，讀取mongodb裡的紀錄")
	return tx.Balance, tx.Currency, nil

}
