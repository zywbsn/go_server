package service

import (
	"express-service/define"
	"express-service/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// @Tags 菜单
// @Summary 修改用户信息
// @Description 修改用户信息接口
// @Router /user/info [put]
// @Param id formData  int true "用户 id"
// @Param identity formData  string true "用户唯一标识"
// @Param nickname formData string true "昵称"
// @Param username formData string true "账号"
// @Param password formData string true "密码"
// @Param phone formData string true "手机号"
// @Param rule formData string true "权限"
// @Produce application/json
// @Success 200 {string} string
func UpdateMenu(c *gin.Context) {
	Id, _ := strconv.Atoi(c.PostForm("id"))
	Identity := c.PostForm("identity")
	NickName := c.PostForm("nickname")
	UserName := c.PostForm("username")
	Password := c.PostForm("password")
	Phone := c.PostForm("phone")
	Rule := c.PostForm("rule")

	currentTime := time.Now()

	_, err := models.GetUserInfo(Identity)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"error":   "User Info err:" + err.Error(),
			"message": "用户不存在",
		})
		return
	}

	info := &models.UserList{
		Id:         Id,
		Identity:   Identity,
		NickName:   NickName,
		UserName:   UserName,
		Password:   Password,
		Phone:      Phone,
		Rule:       Rule,
		UpdateTime: currentTime,
	}
	err = models.DB.Model(new(models.UserList)).Where("identity = ?", Identity).
		Updates(info).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"error":   "User Info err:" + err.Error(),
			"message": "请求失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"data":    info,
		"message": "请求成功",
	})
}

// @Tags 菜单
// @Summary 新增菜单
// @Description 这是一个新增菜单接口
// @Router /menu/create [post]
// @Param key formData string true "key"
// @Param component formData string true "component"
// @Param label formData string true "label"
// @Param parent_id formData string true "parent_id"
// @Param sort formData string true "sort"
// @Produce application/json
// @Success 200 {string} string
func CreateMenu(c *gin.Context) {
	Path := c.PostForm("key")
	Component := c.PostForm("component")
	Label := c.PostForm("label")
	ParentId, pIdErr := strconv.Atoi(c.PostForm("parent_id"))
	Sort, sortErr := strconv.Atoi(c.PostForm("sort"))
	//Identity := helper.GetUUID()
	currentTime := time.Now()
	if Path == "" || Component == "" || Label == "" || pIdErr != nil || sortErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"error":   "参数为空",
			"message": "新增失败",
		})
		return
	}
	data := &models.MenuList{
		Path:       Path,
		Component:  Component,
		Label:      Label,
		ParentId:   ParentId,
		Sort:       Sort,
		CreateTime: currentTime,
	}
	err := models.DB.Create(data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"error":   err.Error(),
			"message": "新增失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"data":    data,
		"message": "创建成功",
	})
}

// @Tags 菜单
// @Summary 菜单路由列表
// @Description 菜单路由列表接口
// @Router /menu/list [get]
// @Param page query string true "page"
// @Param size query string true "size"
// @Param label query string false "label"
// @Param key query string false "key"
// @Param component query string false "component"
// @Produce application/json
// @Success 200 {string} string
func GetMenuList(c *gin.Context) {
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, _ := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	label := c.DefaultQuery("label", "")
	key := c.DefaultQuery("key", "")
	component := c.DefaultQuery("component", "")

	flag := label == "" && key == "" && component == ""

	if page == -1 || flag {
		GetMenu(c, page, size)
		return
	}

	var count int64

	list := make([]models.MenuList, 0)
	tx := models.GetMenuList()

	err := tx.Debug().Where("label LIKE  ? AND path LIKE  ?", "%"+label+"%", "%"+key+"%").Omit("content").Offset((page - 1) * size).Limit(size).Find(&list).Count(&count).Error

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"error":   err,
			"message": "请求错误",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data": map[string]interface{}{
			"list": list,
		},
		"message": "请求成功",
	})
}

//children menu list
func GetMenu(c *gin.Context, page int, size int) {
	list := make([]models.MenuList, 0)
	tx := models.GetMenuList()

	err := tx.Omit("content").Find(&list).Error

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"error":   err,
			"message": "请求错误",
		})
		return
	}

	returnList := models.GetMenus(list, 0)

	if page == -1 {
		c.JSON(http.StatusOK, gin.H{
			"status": 200,
			"data": map[string]interface{}{
				"list": returnList,
			},
			"message": "请求成功",
		})
		return
	}

	start := (page - 1) * size
	stop := 0
	if page*size > len(returnList) {
		stop = len(returnList)
	} else {
		stop = page * size
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data": map[string]interface{}{
			"list": returnList[start:stop],
		},
		"message": "请求成功",
	})

}
