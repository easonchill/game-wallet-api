package controller

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"Basil/module"
	"Basil/structs"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

//Config
var wtoken, Port, Mode string

func init() {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.SetDefault("gin.port", 8080)
	viper.SetDefault("gin.wtoken", "1234")
	viper.SetDefault("gin.mode", gin.DebugMode)

	err := viper.ReadInConfig()

	if err != nil {
		panic("讀取設定檔出現錯誤，原因為：" + err.Error())
	}

	wtoken = viper.GetString("gin.wtoken")
	Port = ":" + viper.GetString("gin.port")
	Mode = viper.GetString("gin.mode")
}

type LoginInfo struct {
	UserID     int64     `json:"userId"`
	ClientIP   string    `json:"clientIP"`
	LoginState string    `json:"loginState"`
	LoginTime  time.Time `json:"loginTime"`
}

type User struct {
	Account     string  `json:"account"`
	Password    string  `json:"password"`
	Balance     float64 `json:"balance"`
	Currency    string  `json:"currency"`
	Status      bool    `json:"status"`
	Last_time   string  `json:"last_time"`
	Create_time string  `json:"create_time"`
}

func NewLoginInfo(id int64, clientIP string, loginState string) *LoginInfo {
	return &LoginInfo{
		UserID:     id,
		ClientIP:   clientIP,
		LoginState: loginState,
		LoginTime:  time.Now(),
	}
}

func CheckPlayer(c *gin.Context) {

	account := c.Param("input")

	// if account == "20220523sporttest" {
	// 	time.Sleep(3 * time.Second)
	// }

	testdb, _ := module.GetMysql()

	var GormUser = new(User)

	err := testdb.Where("account = ?", account).Find(&GormUser).Error

	testdb.Close()

	accountData := false

	if err != nil {
		accountData = false
	} else {
		accountData = true
	}

	c.JSON(200, gin.H{
		"data": accountData,
		"status": gin.H{
			"code":     "0",
			"message":  "Success",
			"datetime": nowTime(),
		},
	})
}

func Balance(c *gin.Context) {

	account := c.Param("input")

	testdb, _ := module.GetMysql()
	var GormUser = new(User)

	testdb.Where("account = ?", account).Find(&GormUser)
	testdb.Close()

	// time.Sleep(10 * time.Second)

	wrapResponse(c, 200, gin.H{
		"balance":  GormUser.Balance,
		"currency": GormUser.Currency,
	}, "0")
}

func CreateUser(c *gin.Context) {
	var GormUser = new(User)

	data := structs.CreateUserReq{}

	err := c.ShouldBind(&data)

	if err != nil {
		wrapResponse(c, 200, err.Error(), "1100")
		return
	}

	testdb, _ := module.GetMysql()
	defer testdb.Close()

	//先檢查該帳號是否存在
	result := testdb.Where("account = ?", data.Account).First(&GormUser)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		wrapResponse(c, 200, "account already exists!", "1100")
		return
	}

	data.Status = true
	data.LastTime = nowTime()
	data.CreateTime = nowTime()
	data.Password = passwordHash(data.Account, data.Password)
	if err := testdb.Create(data).Error; err != nil {
		wrapResponse(c, 200, err, "1100")
		return
	} else {
		wrapResponse(c, 200, gin.H{"Account Create Success": data.Account}, "0")
		return
	}
}

func Version(c *gin.Context) {
	wrapResponse(c, 200, gin.H{"Basil Version": "2.0.0_test60658"}, "0")
}

// 中間件
func CheckTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//取得header的token值
		token := c.GetHeader("wtoken")
		if token == "" {
			//沒有token的話就不執行後面的handler
			c.Abort()
			//回傳錯誤值
			wrapResponse(c, 403, map[string]string{"msg": "未帶入token"}, "1003")
		} else {
			if token == wtoken {
				c.Next()
			} else {
				c.Abort()
				//回傳錯誤值
				wrapResponse(c, 403, map[string]string{"msg": "token 錯誤"}, "1003")
			}
		}

	}
}

//統一設置成json格式回覆
func SetContentTypeJson() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Content-Type", "application/json")
		c.Next()
	}
}

func nowTime() (nowTime string) {
	return time.Now().In(time.FixedZone("myTimeZone", -4*60*60)).Format(time.RFC3339Nano) //設定現在時間，並調整為UTC-4時區
}

//統一回應的格式
func wrapResponse(c *gin.Context, statusCode int, data interface{}, code string) {

	type status struct {
		Code     string `json:"code"`
		Message  string `json:"message"`
		Datetime string `json:"dateime"`
	}

	type ret struct {
		Data   interface{} `json:"data"`
		Status status      `json:"status"`
	}

	s := status{
		Code:     code,
		Message:  "Success",
		Datetime: nowTime(),
	}

	d := ret{
		Data:   data,
		Status: s,
	}

	switch code {
	case "1002":
		d.Status.Message = "執行的action錯誤"
	case "1003":
		d.Status.Message = "Bad parameter！"
	case "1004":
		d.Status.Message = "時間格式錯誤"
	case "1005":
		d.Status.Message = "餘額不足"
	case "1006":
		d.Status.Message = "查無該玩家"
	case "1014":
		d.Status.Message = "未查詢到紀錄"
	case "1015":
		d.Status.Message = "該mtcode已經被refund"
	case "1100":
		d.Status.Message = "伺服器錯誤"
	case "2009":
		d.Status.Message = "mtcode重覆"
	}

	c.JSON(statusCode, d)
}

func passwordHash(account, pwd string) string {
	mac := hmac.New(sha256.New, []byte(account))
	mac.Write([]byte(pwd))
	return hex.EncodeToString(mac.Sum(nil))
}

//同步Log寫到TG
type LogTest struct {
	s []byte
}

func (l LogTest) Write(p []byte) (n int, err error) {
	l.s = p
	d := string(l.s)
	d = strings.Replace(d, "\n", "", -1)
	url := "https://api.telegram.org/bot5397186861:AAF_6jiONQg4KYvICRy0OOU9MGRI1OnOTnA/sendMessage?&parse_mode=html&chat_id=527687006&text=" + d
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	return len(p), nil
}
