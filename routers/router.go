package routers

import (
	"fmt"
	"game-wallet-api/controller"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetRouter(r *gin.Engine) {
	//r.Use(controller.CheckTokenMiddleware(), controller.SetContentTypeJson())
	r.LoadHTMLGlob("views/*")

	r.GET("/player/check/:input", controller.CheckPlayer)
	r.GET("/version", controller.Version)
	r.GET("/", index)
	t := r.Group("/transaction")
	{
		t.GET("/balance/:input", controller.Balance)
		t.GET("/record/:mtcode", controller.Record)

		t.POST("/game/bets", controller.Bets)
		t.POST("/game/wins", controller.Wins)
		t.POST("/game/refunds", controller.Refunds)
		t.POST("/game/cancel", controller.Cancel)
		t.POST("/game/amend", controller.Amend)
		t.POST("/game/amends", controller.Amends)
		t.POST("/user/payoff", controller.Payoff)
	}

	u := r.Group("/swclient/test")
	{
		//新增帳號
		u.POST("/setaccount", controller.CreateUser)
		//u.GET("/fetchUser", fetchUser)
	}

}

func index(c *gin.Context) {

	url := "http://10.30.5.79:88/version"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("wtoken", "9527")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{"title": string(body)})

}
