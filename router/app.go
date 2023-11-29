package router

import (
	"github.com/gin-gonic/gin"
	"go-server/controllers"
	"go-server/services"
)

func Router() *gin.Engine {
	r := gin.Default()
	//apiGroup := r.Group("/api")
	apiGroup := r.Group("")
	{
		/* ----- 后台管理 ----- */
		apiGroup.POST("/login", services.AdminLogin)
		apiGroup.GET("/getTest", services.DownRolesHandler) //测试 Get 接口
		/* ----- 后台管理 ----- */

		/* ----- 菜单路由相关路由组 ----- */
		menuGroup := apiGroup.Group("/menu")
		{
			menuGroup.PUT("/update", controllers.UpdateMenu)
			menuGroup.POST("/create", controllers.CreateMenu)
			menuGroup.GET("/info", controllers.GetMenuInfo)
			menuGroup.GET("/list", controllers.GetMenuList)
			//menuGroup.POST("/login", services.Login)
		}
		/* ----- 菜单路由相关路由组 ----- */

		/* ----- 用户相关路由组 ----- */
		userGroup := apiGroup.Group("/user")
		{
			userGroup.GET("/export", services.ExportUserList) // 导出用户列表
			userGroup.PUT("/update", services.UpdateUserInfo) // 修改用户个人信息
			userGroup.POST("/create", services.CreateUser)    // 获取用户列表
			userGroup.GET("/info", services.GetUserInfo)      // 获取用户个人信息
			userGroup.GET("/list", services.GetUserList)      // 获取用户列表
			userGroup.POST("/login", services.Login)          // 用户存在登录 不存在注册并登录
		}
		/* ----- 用户相关路由组 ----- */

		/* ----- 上传文件 ----- */
		apiGroup.POST("/upload", services.UploadFile) // 上传文件
		/* ----- 上传文件 ----- */
	}

	return r
}
