package controller

import (
	"Basil/module"
	"Basil/structs"
	"context"
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"go.mongodb.org/mongo-driver/mongo"
)

func Wins(c *gin.Context) {

	data := structs.WinsReq{} //注意該結構接受的内容
	err := c.ShouldBindJSON(&data)

	if err != nil {
		wrapResponse(c, 200, err.Error(), "0")
		return
	}
	var winSuccess []structs.WinsSuccessList
	var winFail []structs.WinsFailList
	var playerdata []structs.WinsReqList
	var winsSum float64
	var WinsEvent []structs.Event

OuterLoop:
	for i := 0; i < len(data.List); i++ {
		playerdata = append(playerdata, data.List[i])
		winsSum = 0
		WinsEvent = nil
		t_log := structs.TransactionMgoLog{}

		for k := 0; k < len(playerdata[i].Event); k++ {

			//檢查該玩家該筆資料是否為負，是的話加入錯誤列表
			if playerdata[i].Event[k].Amount < 0 {
				winFail = append(winFail, structs.WinsFailList{
					Account: data.List[i].Account,
					Code:    "1003",
					Message: "amount不得為負值",
					Ucode:   data.List[i].Ucode,
				})
				continue OuterLoop
			}

			recordBalance, recorCurrency, err := module.CheckWinsMtcodeRecode(playerdata[i].Event[k].Mtcode)

			if err != mongo.ErrNoDocuments {

				winSuccess = append(winSuccess, structs.WinsSuccessList{Account: data.List[i].Account,
					Balance:  recordBalance,
					Currency: recorCurrency,
					Ucode:    data.List[i].Ucode,
				})
				continue OuterLoop
			}

			//加入成功列表
			winsSum = winsSum + playerdata[i].Event[k].Amount
			WinsEvent = append(WinsEvent, structs.Event{
				Mtcode:    data.List[i].Event[k].Mtcode,
				Amount:    data.List[i].Event[k].Amount,
				Even_time: data.List[i].Event[k].Eventime,
				Status:    "success",
			})

		}

		testdb, _ := module.GetMysql()
		defer testdb.Close()
		var GormUser = new(User)

		result := testdb.Where("account = ?", data.List[i].Account).Find(&GormUser)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			winFail = append(winFail, structs.WinsFailList{
				Account: data.List[i].Account,
				Code:    "1006",
				Message: "該玩家不存在",
				Ucode:   data.List[i].Ucode,
			})
			continue
		}

		beforeBalance := GormUser.Balance
		newBalance := winsSum + beforeBalance

		winSuccess = append(winSuccess, structs.WinsSuccessList{Account: data.List[i].Account,
			Balance:  newBalance,
			Currency: GormUser.Currency,
			Ucode:    data.List[i].Ucode,
		})

		tx := testdb.Begin()
		if err := tx.Model(&GormUser).Where("account = ?", data.List[i].Account).Update("balance", newBalance).Error; err != nil {
			tx.Rollback()
			panic(err)
		}
		//交易提交
		tx.Commit()

		t_log = structs.TransactionMgoLog{
			Action: "wins",
			Target: structs.Target{
				Account: data.List[i].Account,
			},

			Status: structs.Status{
				Create_time: nowTime(),
				End_time:    data.List[i].Eventtime,
				Status:      "success",
				Msg:         "success",
			},
			Before:   beforeBalance,
			Balance:  newBalance,
			Currency: GormUser.Currency,
			Event:    WinsEvent,
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
		//opts := options.FindOne().SetSkip(0)

		_, err = collection.InsertOne(context.TODO(), t_log)

		if err != nil {
			wrapResponse(c, 200, err, "1100")
			return
		}

	}

	o := structs.WinsResp{
		Success: winSuccess,
		Failed:  winFail,
	}

	time.Sleep(60 * time.Second)
	wrapResponse(c, 200, o, "0")
}
