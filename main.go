package main

// @title Go-server API
// @version 1.0.1
// @contact.name silenceLamb
// @contact.url http://www.swagger.io/support
// @contact.email ooooooooooos@163.com
// @host localhost:9090
// @BasePath /

import (
	"bytes"
	"fmt"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "go-server/docs"
	"go-server/router"
	"os"
	"os/exec"
)

func runCommand() {
	cmd := exec.Command("swag", "init")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	runCommand()
	r := router.Router()
	r.Static("/static", "./static")
	// 注册Swagger接口文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":9090")
}
