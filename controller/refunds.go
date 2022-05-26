package controller

import (
	"Basil/module"
	"Basil/structs"
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Refunds(c *gin.Context) {
	var NewBalance float64 = 0
	data := structs.RefundsReq{}
	tx := transaction_mgolog{}
	refundResp := structs.RefundsResp{}

	if err := c.BindJSON(&data); err != nil {
		wrapResponse(c, 200, err.Error(), "1003")
		return
	}

	//初始化mongoDB
	mongoClient, err := module.GetMgoCli()
	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		wrapResponse(c, 200, "MongoDB connect error", "1100")
		return
	}

	collection := mongoClient.Database("transaction").Collection("transaction_log")

	if err != nil {
		wrapResponse(c, 200, "MongoDB connect error", "1100")
		return
	}

	//設定mongo搜尋結果，欄位只要帳號跟符合mtcode的event
	opts1 := options.FindOne().SetProjection((bson.D{{"event.$", 1}, {"target.account", 1}}))
	//先取出到該mtcode的帳號
	for _, v := range data.Mtcode {

		err = collection.FindOne(
			context.TODO(),
			bson.D{{"event.mtcode", v}},
			opts1,
		).Decode(&tx)

		//找不到該mtcode就報錯結束這回合
		if err != nil {
			if err == mongo.ErrNoDocuments {
				wrapResponse(c, 200, nil, "1014")
				return
			}
		}
		//發現該mtcode已經refund就報錯結束
		if tx.Event[0].Status == "refund" {
			wrapResponse(c, 200, nil, "1015")
			return
		}
		//加上這張單金額
		NewBalance += tx.Event[0].Amount
	}

	//開啟Mysql
	testdb, _ := module.GetMysql()
	defer testdb.Close()
	account := tx.Target.Account

	var GormUser = new(User)

	testdb.Where("account = ?", account).Find(&GormUser)

	beforeBalance := GormUser.Balance

	NewBalance += beforeBalance

	MysqlTx := testdb.Begin()
	if err := MysqlTx.Model(&GormUser).Where("account = ?", account).Update("balance", NewBalance).Error; err != nil {
		MysqlTx.Rollback()
		fmt.Println(err)
		wrapResponse(c, 200, err, "1100")
	}
	//交易提交
	MysqlTx.Commit()

	opts2 := options.FindOneAndUpdate().SetUpsert(false)

	for _, v := range data.Mtcode {

		findMtcode := structs.UpdateByMtcode{Mtcode: v}
		filter := findMtcode

		update := bson.D{{"$set", bson.D{{"event.$.status", "refund"}}}}
		err := collection.FindOneAndUpdate(
			context.TODO(),
			filter,
			update,
			opts2,
		).Decode(&tx)

		if err != nil {
			wrapResponse(c, 200, nil, "1014")
			return
		}

	}

	refundResp.Balance = NewBalance
	refundResp.Currency = GormUser.Currency
	wrapResponse(c, 200, refundResp, "0")

}
