package main

import (
	"Basil/controller"
	"Basil/routers"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	gin.SetMode(controller.Mode)
	//啟動gin

	//f := controller.LogTest{}

	//gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	r := gin.Default()
	routers.SetRouter(r)

	r.Run(controller.Port)

}
