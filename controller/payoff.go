package controller

import (
	"fmt"
	"game-wallet-api/module"
	"game-wallet-api/structs"

	"github.com/gin-gonic/gin"
)

func Payoff(c *gin.Context) {
	data := structs.PayoffReq{}

	c.ShouldBind(&data)

	nb, err := module.AddMoney(data.Account, data.Amount)

	fmt.Println(err)
	if err != nil {
		wrapResponse(c, 200, err, "1100")
		fmt.Println(err)
		return
	}
	wrapResponse(c, 200, nb, "0")
	return
}
