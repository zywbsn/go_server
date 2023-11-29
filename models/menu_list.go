package models

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goinggo/mapstructure"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type MenuList struct {
	Id         int       `json:"id" gorm:"column:id;type:int(11);"`
	Path       string    `json:"key" gorm:"column:path;type:varchar(255);"`
	Component  string    `json:"component" gorm:"column:component;type:varchar(255);"`
	Label      string    `json:"label" gorm:"column:label;type:varchar(255);" `
	ParentId   int       `json:"parent_id" gorm:"column:parent_id;type:int(11);"`
	Sort       int       `json:"sort" gorm:"column:sort;type:int(11);" `
	CreateTime time.Time `json:"create_time" gorm:"column:create_time;type:varchar(255);"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time;type:varchar(255);"`
	DeleteTime time.Time `json:"delete_time" gorm:"column:delete_time;type:varchar(255);"`
}

type Menus struct {
	MenuList
	Children []Menus `json:"children"`
}

//获取侧边菜单
func GetMenus(List []MenuList, PId int) (treeList []Menus) {
	for _, v := range List {
		if v.ParentId == PId {
			child := GetMenus(List, v.Id)
			node := Menus{
				MenuList: MenuList{
					Id:         v.Id,
					Path:       v.Path,
					Component:  v.Component,
					Label:      v.Label,
					ParentId:   PId,
					Sort:       v.Sort,
					CreateTime: v.CreateTime,
					UpdateTime: v.UpdateTime,
					DeleteTime: v.DeleteTime,
				},
			}
			node.Children = child

			treeList = append(treeList, node)
			continue
		}
	}
	return treeList
}

// 菜单参数
func FormatQuery(query map[string]interface{}) MenuList {
	q := map[string]interface{}{
		"Id":        query["id"],
		"Path":      query["key"],
		"Component": query["component"],
		"Label":     query["label"],
		"ParentId":  query["parent_id"],
		"Sort":      query["sort"],
	}
	var data MenuList
	//将 map 转换为指定的结构体
	if err := mapstructure.Decode(q, &data); err != nil {
		fmt.Println(err)
	}
	return data
}

// GetMenuInfo 获取菜单详情
func GetMenuInfo(ctx *gin.Context, id int) (info *MenuList) {
	fmt.Println("id", id)
	info = new(MenuList)
	err := DB.Where("id = ?", id).First(&info).Error
	//if err != nil {
	if err != nil {
		ctx.JSON(
			http.StatusOK, gin.H{
				"status":  -1,
				"error":   "User Info err:" + err.Error(),
				"message": "菜单不存在",
			},
		)
		return nil
	}
	//return nil, err
	//}
	return info
}

// 菜单列表
func GetMenuList() *gorm.DB {
	return DB.Model(new(MenuList))
}

func (table *MenuList) TableName() string {
	return "menu_list"
}
