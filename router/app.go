package router

import (
	"express-service/service"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	//apiGroup := r.Group("/api")
	apiGroup := r.Group("")
	{
		/* ----- 后台管理 ----- */
		apiGroup.POST("/admin/login", service.AdminLogin)
		apiGroup.GET("/getTest", service.DownRolesHandler) //测试 Get 接口
		/* ----- 后台管理 ----- */

		/* ----- 菜单路由相关路由组 ----- */
		menuGroup := apiGroup.Group("/menu")
		{
			//menuGroup.PUT("/update", service.UpdateUserInfo)
			menuGroup.POST("/create", service.CreateMenu) //
			//menuGroup.GET("/info", service.GetUserInfo)
			menuGroup.GET("/list", service.GetMenuList)
			//menuGroup.POST("/login", service.Login)
		}
		/* ----- 菜单路由相关路由组 ----- */

		/* ----- 用户相关路由组 ----- */
		userGroup := apiGroup.Group("/user")
		{
			userGroup.GET("/export", service.ExportUserList) // 导出用户列表
			userGroup.PUT("/update", service.UpdateUserInfo) // 修改用户个人信息
			userGroup.POST("/create", service.CreateUser)    // 获取用户列表
			userGroup.GET("/info", service.GetUserInfo)      // 获取用户个人信息
			userGroup.GET("/list", service.GetUserList)      // 获取用户列表
			userGroup.POST("/login", service.Login)          // 用户存在登录 不存在注册并登录
		}
		/* ----- 用户相关路由组 ----- */

		/* ----- 上传文件 ----- */
		apiGroup.POST("/upload", service.UploadFile) // 上传文件
		/* ----- 上传文件 ----- */
	}

	return r
}
