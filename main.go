package main

// @title 快递代取11 API
// @version 1.0.1
// @contact.name silenceLamb
// @contact.url http://www.swagger.io/support
// @contact.email ooooooooooos@163.com
// @host localhost:9090
// @BasePath /

import (
	"bytes"
	_ "express-service/docs"
	"express-service/router"
	"fmt"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
	"os/exec"
)

func runCommand() {
	cmd := exec.Command("swag", "init")
	fmt.Println("Cmd", cmd.Args)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("String", out.String())
}

func main() {
	runCommand()
	r := router.Router()
	r.Static("/static", "./static")
	// 注册Swagger接口文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":9090")
}
