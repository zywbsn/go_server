package controllers

import (
	"github.com/gin-gonic/gin"
	"go-server/define"
	"go-server/helper"
	"go-server/models"
	"go-server/services"
	"strconv"
	"time"
)

// @Tags 菜单
// @Summary 更新菜单信息
// @Description 更新菜单信息接口
// @Router /menu/update [put]
// @Param id query  int true "id"
// @Param key query  string true "key"
// @Param component query string true "component"
// @Param label query string true "label"
// @Param parent_id query int true "parent_id"
// @Param sort query int true "sort"
// @Param icon query string true "icon"
// @Produce application/json
// @Success 200 {string} string
func UpdateMenu(ctx *gin.Context) {
	query := models.FormatQuery(helper.PostJson(ctx))
	currentTime := time.Now()

	query.UpdateTime = currentTime

	info := models.GetMenuInfo(ctx, query.Id)

	if info == nil {
		return
	}

	if query.CreateTime.IsZero() {
		query.CreateTime = info.CreateTime
	}

	if query.DeleteTime.IsZero() {
		query.DeleteTime = info.DeleteTime
	}

	services.UpdateMenu(ctx, query)
}

// @Tags 菜单
// @Summary 菜单详情
// @Description 这是一个新增菜单接口
// @Router /menu/info [get]
// @Param id query int true "id"
// @Produce application/json
// @Success 200 {string} string
func GetMenuInfo(ctx *gin.Context) {
	services.GetMenuInfo(ctx)
}

// @Tags 菜单
// @Summary 新增菜单
// @Description 这是一个新增菜单接口
// @Router /menu/create [post]
// @Param key query string true "key"
// @Param component query string true "component"
// @Param label query string true "label"
// @Param parent_id query int true "parent_id"
// @Param sort query string false "sort"
// @Produce application/json
// @Success 200 {string} string
func CreateMenu(ctx *gin.Context) {
	services.CreateMenu(ctx)
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
func GetMenuList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", define.DefaultPage))
	label := ctx.DefaultQuery("label", "")
	key := ctx.DefaultQuery("key", "")
	component := ctx.DefaultQuery("component", "")

	flag := label == "" && key == "" && component == ""
	if page == -1 || flag {
		services.GetMenu(ctx)
	} else {
		services.GetMenus(ctx)
	}
}
