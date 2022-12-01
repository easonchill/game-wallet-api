package controller

import (
	"Basil/module"
	"Basil/structs"
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Amend(c *gin.Context) {
	var amendTotal float64 = 0
	var amendEvent []structs.AmendEvent
	var amendResp structs.AmendResp
	var tx structs.TransactionMgoLog

	CheckMtcode := make(map[string]int)
	// startTime := time.Now()

	data := structs.AmendReq{} //注意該結構接受的内容

	if err := c.BindJSON(&data); err != nil {

		wrapResponse(c, c.Writer.Status(), err.Error(), "1003")

		return
	}

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
		//先檢查該mtcode，存在的話就拿原本餘額回傳後結束
		amendResp.Balance, amendResp.Currency, err = module.CheckMtcodeRecode(data.Data[i].Mtcode)

		if err != mongo.ErrNoDocuments {
			wrapResponse(c, 200, amendResp, "0")
			return
		}

		findMtcode := structs.FindByMtcode{Mtcode: data.Data[i].Mtcode}

		err = collection.FindOne(
			context.TODO(),
			findMtcode,
			opts,
		).Decode(&tx)

		//先查詢該mtcode是否存在，是的話就報錯結束
		if err != mongo.ErrNoDocuments {
			wrapResponse(c, 200, nil, "2009")
			return
		}

		if data.Data[i].Action == "credit" {
			amendTotal += data.Data[i].Amount
		} else if data.Data[i].Action == "debit" {
			amendTotal += -data.Data[i].Amount
		} else {
			wrapResponse(c, 200, nil, "1003")
			return
		}

		_, ok := CheckMtcode[data.Data[i].Mtcode]

		if ok {
			wrapResponse(c, 200, nil, "2009")
			return
		}

		CheckMtcode[data.Data[i].Mtcode] = 1

		// //檢查本次陣列裡是否有重覆的mtcode，有的話回傳錯誤然後結束
		// for _, v := range amendEvent {
		// 	if data.Data[i].Mtcode == v.Mtcode {
		// 		wrapResponse(c, 200, nil, "2009")
		// 		return
		// 	}
		// }

		amendEvent = append(amendEvent, structs.AmendEvent{
			Mtcode:    data.Data[i].Mtcode,
			Amount:    data.Data[i].Amount,
			Action:    data.Data[i].Action,
			Eventtime: data.Data[i].Eventtime,
		})
	}

	//透過GORM連接MySQL，先取得MySQL裡玩家目前餘額
	testdb, _ := module.GetMysql()
	account := data.Account
	var GormUser = new(User)

	testdb.Where("account = ?", account).Find(&GormUser)
	defer testdb.Close()

	beforeBalance := GormUser.Balance
	if GormUser.Balance < -(amendTotal) {
		wrapResponse(c, 200, nil, "1005")
	} else {
		//加上總共的amend金額，取到新的餘額
		NewBalance := beforeBalance + amendTotal
		//開啟GORM交易處理(事務)，將新的餘額寫回MySQL
		tx := testdb.Begin()
		if err := tx.Model(&GormUser).Where("account = ?", account).Update("balance", NewBalance).Error; err != nil {
			tx.Rollback()
			fmt.Println(err)
			wrapResponse(c, 200, err, "1100")
		}
		//交易提交
		tx.Commit()
		ent := nowTime()

		//準備寫入MongoDB的交易紀錄ㄖ
		amendTranRecord := structs.AmendRecord{
			Action: "amend",
			Target: structs.AmendTarget{
				Account: data.Account,
			},

			Status: structs.AmendStatus{
				Createtime: crt,
				Endtime:    ent,
				Status:     "success",
				Message:    "success",
			},
			Before:   beforeBalance,
			Balance:  NewBalance,
			Currency: GormUser.Currency,
			Event:    amendEvent,
		}

		//把交易紀錄塞進mongoDB
		_, err := collection.InsertOne(context.TODO(), amendTranRecord)
		if err != nil {
			panic(err)
		}

		amendResp.Balance = NewBalance
		amendResp.Currency = GormUser.Currency
		// id := iResult.InsertedID.(primitive.ObjectID)
		// fmt.Println("自動增加ID", id.Hex())

		// endTime := time.Now()
		// duration := endTime.Sub(startTime)
		// fmt.Println(duration)

		wrapResponse(c, 200, amendResp, "0")
	}

}
