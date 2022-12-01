package routers

import (
	"Basil/controller"
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

	c.HTML(http.StatusOK, "index.html", nil)

}
