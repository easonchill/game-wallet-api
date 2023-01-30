package module

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mysqlConfig string

type User struct {
	Account     string  `json:"account"`
	Password    string  `json:"password"`
	Balance     float64 `json:"balance"`
	Currency    string  `json:"currency"`
	Status      bool    `json:"status"`
	Last_time   string  `json:"last_time"`
	Create_time string  `json:"create_time"`
}

func init() {
	//載入config的DB設定

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()

	if err != nil {
		panic("讀取設定檔出現錯誤，原因為：" + err.Error())
	}

	mysqlConfig = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"),
	)
	fmt.Println("config load success!")
}

func GetMysql() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", mysqlConfig)
	//defer db.Close()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	return db, err
}

func GetMgoCli() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:55000"))
	return client, err
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
