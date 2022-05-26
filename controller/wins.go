package controller

import (
	"Basil/module"
	"Basil/structs"
	"context"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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
	var WinsEvent []event

OuterLoop:
	for i := 0; i < len(data.List); i++ {
		playerdata = append(playerdata, data.List[i])
		winsSum = 0
		t_log := transaction_log{}
		//playerevent := []structs.WinsEvent{}
		for k := 0; k < len(playerdata[i].Event); k++ {
			//playerevent = append(playerevent, playerdata[i].Event[k])

			if playerdata[i].Event[k].Amount < 0 {
				winFail = append(winFail, structs.WinsFailList{
					Account: data.List[i].Account,
					Code:    "1003",
					Message: "amount不得為負值",
					Ucode:   data.List[i].Ucode,
				})
				continue OuterLoop
			}

			winsSum = winsSum + playerdata[i].Event[k].Amount
			WinsEvent = append(WinsEvent, event{
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

		t_log = transaction_log{
			Action: "wins",
			Target: struct {
				Account string "json:\"account\" bson:\"account\""
			}{
				Account: data.List[i].Account,
			},

			Status: struct {
				Create_time string "json:\"createtime\" bson:\"createtime\""
				End_time    string "json:\"endtime\"  bson:\"endtime\""
				Status      string "json:\"status\" bson:\"status\""
				Msg         string "json:\"message\" bson:\"message\""
			}{
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
	wrapResponse(c, 200, o, "0")
}
