package module

import (
	"errors"
	"game-wallet-api/env"

	"github.com/jinzhu/gorm"
)

type User struct {
	Account     string  `json:"account"`
	Password    string  `json:"password"`
	Balance     float64 `json:"balance"`
	Currency    string  `json:"currency"`
	Status      bool    `json:"status"`
	Last_time   string  `json:"last_time"`
	Create_time string  `json:"create_time"`
}

func GetMysql() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", env.MysqlConfig)
	//defer db.Close()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	return db, err
}

func AddMoney(account string, money float64) (newBalance float64, err error) {
	dbs, _ := GetMysql()
	defer dbs.Close()

	var GormUser = new(User)

	result := dbs.Where("account = ?", account).Find(&GormUser)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return GormUser.Balance, gorm.ErrRecordNotFound
	}
	//開啟交易
	tx := dbs.Begin()
	beforeBalance := GormUser.Balance
	newBalance = money + beforeBalance

	if err := tx.Model(&GormUser).Where("account = ?", account).Update("balance", newBalance).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	//交易提交
	tx.Commit()

	return GormUser.Balance, nil
}
