package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-server/define"
	"go-server/helper"
	"go-server/models"
	"net/http"
	"strconv"
	"time"
)

//UpdateMenu 更新菜单
func UpdateMenu(c *gin.Context, query models.MenuList) {
	err := models.DB.
		Model(new(models.MenuList)).
		Where("id = ?", query.Id).
		Updates(query).Error
	if err != nil {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  -1,
				"error":   "User Info err:" + err.Error(),
				"message": "操作失败",
			},
		)
		return
	}
	c.JSON(
		http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"data":    query,
			"message": "操作成功",
		},
	)
}

//GetMenuInfo 菜单详情
func GetMenuInfo(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))

	if err != nil {
		ctx.JSON(
			http.StatusOK, gin.H{
				"status":  -1,
				"error":   err.Error(),
				"message": "请求失败",
			},
		)
		return
	}

	info := models.GetMenuInfo(ctx, id)

	ctx.JSON(
		http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"info":    info,
			"message": "请求成功",
		},
	)
}

//CreateMenu 新增菜单
func CreateMenu(ctx *gin.Context) {
	data := models.FormatQuery(helper.PostJson(ctx))
	currentTime := time.Now()
	if data.Path == "" || data.Component == "" || data.Label == "" {
		ctx.JSON(
			http.StatusOK, gin.H{
				"status":  -1,
				"message": "参数为空",
			},
		)
		return
	}

	info := new(models.MenuList)

	err := models.DB.Where("path = ?", data.Path).First(&info).Error

	if err == nil {
		fmt.Println("bug3")
		ctx.JSON(
			http.StatusOK, gin.H{
				"status":  -1,
				"message": "菜单路径已存在",
			},
		)
		return
	}

	data.CreateTime = currentTime
	Data := &models.MenuList{
		Path:       data.Path,
		Component:  data.Component,
		Label:      data.Label,
		ParentId:   data.ParentId,
		Sort:       data.Sort,
		CreateTime: data.CreateTime,
	}

	err = models.DB.Debug().Create(Data).Error

	if err != nil {
		ctx.JSON(
			http.StatusOK, gin.H{
				"status":  -1,
				"error":   err.Error(),
				"message": "操作失败",
			},
		)
		return
	}
	ctx.JSON(
		http.StatusOK, gin.H{
			"status":  200,
			"data":    Data,
			"message": "操作成功",
		},
	)
}

//GetMenus 不分组菜单
func GetMenus(ctx *gin.Context) {
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", define.DefaultSize))
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", define.DefaultPage))
	label := ctx.DefaultQuery("label", "")
	key := ctx.DefaultQuery("key", "")
	component := ctx.DefaultQuery("component", "")

	var count int64

	list := make([]models.MenuList, 0)
	tx := models.GetMenuList()

	err := tx.Where(
		"label LIKE  ? AND path LIKE ? AND component LIKE  ?",
		"%"+label+"%",
		"%"+key+"%",
		"%"+component+"%",
	).
		Omit("content").
		Count(&count).
		Offset((page - 1) * size).
		Limit(size).
		Find(&list).
		Error

	if err != nil {
		ctx.JSON(
			http.StatusOK, gin.H{
				"status":  -1,
				"error":   err,
				"message": "请求失败",
			},
		)
		return
	}
	ctx.JSON(
		http.StatusOK, gin.H{
			"status": 200,
			"data": map[string]interface{}{
				"page":  page,
				"size":  size,
				"list":  list,
				"total": count,
			},
			"message": "请求成功",
		},
	)
}

//GetMenu 分组菜单
func GetMenu(ctx *gin.Context) {
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", define.DefaultSize))
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", define.DefaultPage))

	list := make([]models.MenuList, 0)
	tx := models.GetMenuList()

	err := tx.Omit("content").Find(&list).Error

	if err != nil {
		ctx.JSON(
			http.StatusOK, gin.H{
				"status":  -1,
				"error":   err,
				"message": "请求失败",
			},
		)
		return
	}

	returnList := models.GetMenus(list, 0)

	if page == -1 {
		ctx.JSON(
			http.StatusOK, gin.H{
				"status": 200,
				"data": map[string]interface{}{
					"list": returnList,
				},
				"message": "请求成功",
			},
		)
		return
	}

	start := (page - 1) * size
	stop := 0
	if page*size > len(returnList) {
		stop = len(returnList)
	} else {
		stop = page * size
	}
	ctx.JSON(
		http.StatusOK, gin.H{
			"status": 200,
			"data": map[string]interface{}{
				"page":  page,
				"size":  size,
				"list":  returnList[start:stop],
				"total": len(returnList),
			},
			"message": "请求成功",
		},
	)
}
