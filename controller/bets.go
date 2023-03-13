package controller

import (
	"context"
	"fmt"
	"game-wallet-api/module"
	"game-wallet-api/structs"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Bets(c *gin.Context) {
	var betsTotal float64
	var betsEvent []structs.Event
	var transactionMgoLog structs.TransactionMgoLog

	data := structs.BetsReq{} //注意該結構接受的内容
	c.BindJSON(&data)

	crt := nowTime()

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
	opts := options.FindOne().SetSkip(0)

	for i := 0; i < len(data.Data); i++ {
		findMtcode := structs.FindByMtcode{Mtcode: data.Data[i].Mtcode}
		err = collection.FindOne(
			context.TODO(),
			findMtcode,
			opts,
		).Decode(&transactionMgoLog)

		if err != mongo.ErrNoDocuments {
			wrapResponse(c, 200, nil, "2009")
			return
		}

		betsTotal = betsTotal + data.Data[i].Amount

		//檢查本次陣列裡是否有重覆的mtcode，有的話回傳錯誤然後結束
		for _, v := range betsEvent {
			if data.Data[i].Mtcode == v.Mtcode {
				wrapResponse(c, 200, nil, "2009")
				return
			}
		}

		betsEvent = append(betsEvent, structs.Event{
			Mtcode:    data.Data[i].Mtcode,
			Amount:    data.Data[i].Amount,
			Even_time: data.Data[i].Eventime,
			Status:    "success",
		})
	}

	//透過GORM連接MySQL，先取得MySQL裡玩家目前餘額
	testdb, _ := module.GetMysql()
	defer testdb.Close()
	account := data.Account
	var GormUser = new(User)

	testdb.Where("account = ?", account).Find(&GormUser)

	//開啟GORM交易處理(事務)，將新的餘額寫回MySQL
	tx := testdb.Begin()
	beforeBalance := GormUser.Balance
	if GormUser.Balance < betsTotal {
		wrapResponse(c, 200, nil, "1005")
	} else {
		//扣掉本次下注，取到新的餘額
		NewBalance := beforeBalance - betsTotal

		if err := tx.Model(&GormUser).Where("account = ?", account).Update("balance", NewBalance).Error; err != nil {
			tx.Rollback()
			fmt.Println(err)
			wrapResponse(c, 200, err, "1100")
		}
		//交易提交
		tx.Commit()
		ent := nowTime()

		//準備寫入MongoDB的交易紀錄ㄖ
		transactionMgoLog = structs.TransactionMgoLog{
			Action: "bets",
			Target: structs.Target{
				Account: data.Account,
			},
			Status: structs.Status{
				Create_time: crt,
				End_time:    ent,
				Status:      "success",
				Msg:         "success",
			},
			Before:   beforeBalance,
			Balance:  NewBalance,
			Currency: GormUser.Currency,
			Event:    betsEvent,
		}

		//把交易紀錄塞進mongoDB
		_, err := collection.InsertOne(context.TODO(), transactionMgoLog)
		if err != nil {
			panic(err)
		}

		// id := iResult.InsertedID.(primitive.ObjectID)
		// fmt.Println("自動增加ID", id.Hex())
		wrapResponse(c, 200, structs.BetsResp{Balance: NewBalance, Currency: GormUser.Currency}, "0")
	}

}
