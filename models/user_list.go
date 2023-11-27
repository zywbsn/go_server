package models

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type UserList struct {
	Id         int       `gorm:"column:id;type:int(11);" json:"id"`
	Identity   string    `gorm:"column:identity;type:varchar(255);" json:"identity"` // 用户的唯一标识
	NickName   string    `gorm:"column:nickname;type:varchar(255);" json:"nickname"`
	UserName   string    `gorm:"column:username;type:varchar(255);" json:"username"`
	Password   string    `gorm:"column:password;type:varchar(255);" json:"password"`
	Phone      string    `gorm:"column:phone;type:varchar(255);" json:"phone"`
	Rule       string    `gorm:"column:rule;type:varchar(255);" json:"rule"`
	CreateTime time.Time `gorm:"column:create_time;type:varchar(255);" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;type:varchar(255);" json:"update_time"`
	DeleteTime time.Time `gorm:"column:delete_time;type:varchar(255);" json:"delete_time"`
}

// GetName 获取用户姓名
func GetName(identity string) string {
	if identity != "" {
		info, _ := GetUserInfo(identity)
		return info.NickName
	}
	return ""
}

//// 获取用户头像
//func GetImage(identity string) string {
//	info, _ := GetUserInfo(identity)
//	return info.AvatarUrl
//}

// GetUserInfo 获取个人信息
func GetUserInfo(identity string) (info *UserList, err error) {
	fmt.Println("identity", identity)
	info = new(UserList)
	err = DB.Where("identity = ?", identity).First(&info).Error
	if err != nil {
		return nil, err
	}
	return
}

// 用户列表
func GetUserList() *gorm.DB {
	return DB.Model(new(UserList))
}

func (table *UserList) TableName() string {
	return "user_list"
}
