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

func Amends(c *gin.Context) {

	data := structs.AmendsReq{} //注意該結構接受的内容
	err := c.ShouldBindJSON(&data)

	if err != nil {
		wrapResponse(c, 200, err.Error(), "0")
		return
	}

	var success []structs.AmendsSuccessList
	var fail []structs.AmendsFailList
	var playerdata []structs.AmendsList
	var amendsSum float64
	var amendsEvent []structs.AmendEvent

OuterLoop:
	for i := 0; i < len(data.List); i++ {
		playerdata = append(playerdata, data.List[i])
		amendsSum = 0
		amendsTranRecord := structs.AmendRecord{}

		for k := 0; k < len(playerdata[i].Event); k++ {

			if playerdata[i].Event[k].Amount < 0 {
				fail = append(fail, structs.AmendsFailList{
					Account: data.List[i].Account,
					Code:    "1003",
					Message: "amount不得為負值",
					Ucode:   data.List[i].Ucode,
				})
				continue OuterLoop
			}

			amendsEvent = append(amendsEvent, structs.AmendEvent{
				Mtcode:    data.List[i].Event[k].Mtcode,
				Amount:    data.List[i].Event[k].Amount,
				Action:    data.List[i].Event[k].Action,
				Eventtime: data.List[i].Event[k].Eventtime,
			})
		}

		if data.List[i].Action == "credit" {
			amendsSum = data.List[i].Amount
		} else if data.List[i].Action == "debit" {
			amendsSum = -data.List[i].Amount
		} else {
			wrapResponse(c, 200, nil, "1003")
			return
		}

		testdb, _ := module.GetMysql()
		defer testdb.Close()
		var GormUser = new(User)

		result := testdb.Where("account = ?", data.List[i].Account).Find(&GormUser)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fail = append(fail, structs.AmendsFailList{
				Account: data.List[i].Account,
				Code:    "1006",
				Message: "該玩家不存在",
				Ucode:   data.List[i].Ucode,
			})
			continue
		}

		beforeBalance := GormUser.Balance
		newBalance := amendsSum + beforeBalance

		if newBalance < 0 {
			fail = append(fail, structs.AmendsFailList{
				Account: data.List[i].Account,
				Code:    "1003",
				Message: "餘額不足",
				Ucode:   data.List[i].Ucode,
			})
			continue
		}

		success = append(success, structs.AmendsSuccessList{Account: data.List[i].Account,
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

		amendsTranRecord = structs.AmendRecord{
			Action: "amends",
			Target: structs.AmendTarget{
				Account: data.List[i].Account,
			},
			Status: structs.AmendStatus{
				Createtime: nowTime(),
				Endtime:    data.List[i].Eventtime,
				Status:     "success",
				Message:    "success",
			},
			Before:   beforeBalance,
			Balance:  newBalance,
			Currency: GormUser.Currency,
			Event:    amendsEvent,
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

		_, err = collection.InsertOne(context.TODO(), amendsTranRecord)

		if err != nil {
			wrapResponse(c, 200, err, "1100")
			return
		}

	}

	AmendsRespOutput := structs.AmendsResp{
		Success: success,
		Failed:  fail,
	}
	wrapResponse(c, 200, AmendsRespOutput, "0")
}
