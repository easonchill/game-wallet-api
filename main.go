package main

import (
	"game-wallet-api/controller"
	"game-wallet-api/env"
	"game-wallet-api/routers"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	//設定為true，會將log同步輸出到設定好的tg頻道
	gin.SetMode(env.Mode)

	if env.OutputLogToTG {
		f := controller.LogTest{}
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	}

	//啟動gin
	r := gin.Default()

	routers.SetRouter(r)

	r.Run(env.Port)

}
