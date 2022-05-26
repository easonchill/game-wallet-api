package routers

import (
	"Basil/controller"

	"github.com/gin-gonic/gin"
)

func SetRouter(r *gin.Engine) {
	r.Use(controller.CheckTokenMiddleware(), controller.SetContentTypeJson())

	r.GET("/player/check/:input", controller.CheckPlayer)

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

	// u := r.Group("/user")
	// {
	// 	u.POST("/createUser", controller.CreateUser)
	// 	//u.GET("/fetchUser", fetchUser)
	// }

	u := r.Group("/swclient/test")
	{
		//新增帳號
		u.POST("/setaccount", controller.CreateUser)
		//u.GET("/fetchUser", fetchUser)
	}

	// }

	// // func hello(c *gin.Context) {
	// // 	c.JSON(200, gin.H{
	// // 		"message": "Hello world!!",
	// // 	})
	// // }

	// func fetchUser(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "Hello world!!",
	// 	})
}
