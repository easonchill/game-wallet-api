package env

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

//Config
var Wtoken, Port, Mode string

//Config
var MysqlConfig string
var MongodbConfig string
var TGurl string
var OutputLogToTG bool

func init() {

	//設定讀取路徑
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	//設定預設值
	viper.SetDefault("gin.port", 8080)
	viper.SetDefault("gin.wtoken", "1234")
	viper.SetDefault("gin.mode", gin.DebugMode)

	err := viper.ReadInConfig()

	if err != nil {
		panic("讀取設定檔出現錯誤，原因為：" + err.Error())
	}

	Wtoken = viper.GetString("gin.wtoken")
	Port = ":" + viper.GetString("gin.port")
	Mode = viper.GetString("gin.mode")

	//載入config的DB設定

	MysqlConfig = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"),
	)

	MongodbConfig = fmt.Sprintf("mongodb://%s:%s@%s:%d",
		viper.GetString("mongodb.user"),
		viper.GetString("mongodb.password"),
		viper.GetString("mongodb.host"),
		viper.GetInt("mongodb.port"),
	)

	OutputLogToTG = viper.GetBool("other.outputLogToTG")
	TGurl = viper.GetString("other.TGurl")
	fmt.Println("config load success!")
}
