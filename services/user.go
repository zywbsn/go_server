package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-server/define"
	"go-server/helper"
	"go-server/models"
	"net/http"
	"strconv"
	"time"
)

// @Tags 用户
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
func UpdateUserInfo(c *gin.Context) {
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
		c.JSON(
			http.StatusOK, gin.H{
				"status":  -1,
				"error":   "User Info err:" + err.Error(),
				"message": "用户不存在",
			},
		)
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
		c.JSON(
			http.StatusOK, gin.H{
				"status":  -1,
				"error":   "User Info err:" + err.Error(),
				"message": "请求失败",
			},
		)
		return
	}
	c.JSON(
		http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"data":    info,
			"message": "请求成功",
		},
	)
}

// @Tags 用户
// @Summary 新增用户
// @Description 这是一个新增用户接口
// @Router /user/create [post]
// @Param nickname formData string true "昵称"
// @Param username formData string true "账号"
// @Param password formData string true "密码"
// @Param phone formData string true "手机号"
// @Param rule formData string true "权限"
// @Produce application/json
// @Success 200 {string} string
func CreateUser(c *gin.Context) {
	NickName := c.PostForm("nickname")
	UserName := c.PostForm("username")
	Password := c.PostForm("password")
	Phone := c.PostForm("phone")
	Rule := c.PostForm("rule")
	Identity := helper.GetUUID()
	currentTime := time.Now()
	if NickName == "" || UserName == "" || Password == "" || Phone == "" || Rule == "" {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  -1,
				"error":   "参数为空",
				"message": "新增失败",
			},
		)
		return
	}
	data := &models.UserList{
		NickName:   NickName,
		UserName:   UserName,
		Password:   Password,
		Phone:      Phone,
		Rule:       Rule,
		Identity:   Identity,
		CreateTime: currentTime,
		UpdateTime: currentTime,
	}
	err := models.DB.Create(data).Error
	if err != nil {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  -1,
				"error":   err.Error(),
				"message": "新增失败",
			},
		)
		return
	}
	c.JSON(
		http.StatusOK, gin.H{
			"status":  200,
			"data":    data,
			"message": "创建成功",
		},
	)
}

// @Tags 用户
// @Summary 用户详情
// @Description 用户详情接口
// @Router /user/info [get]
// @Param identity query  string true "identity"
// @Produce application/json
// @Success 200 {string} string
func GetUserInfo(c *gin.Context) {
	identity := c.Query("identity")
	info, err := models.GetUserInfo(identity)
	if err != nil {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  -1,
				"message": "User Info err:" + err.Error(),
			},
		)
		return
	}
	c.JSON(
		http.StatusOK, gin.H{
			"status": http.StatusOK,
			"info":   info,
		},
	)
}

// @Tags 用户
// @Summary 导出
// @Description 用户列表接口
// @Router /user/export [get]
// @Param page query  string true "page"
// @Param size query  string true "size"
// @Produce application/json
// @Success 200 {string} string
func ExportUserList(c *gin.Context) {
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, _ := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if size == 0 {
		size = 10
	}
	if page == 0 {
		page = 1
	}

	list := make([]*models.UserList, 0)
	tx := models.GetUserList()
	err := tx.Omit("content").Offset((page - 1) * size).Limit(size).Find(&list).Error
	if err != nil {
		fmt.Printf("Get Express Error:", err)
		return
	}

	var res []interface{}
	for _, v := range list {
		res = append(
			res, &models.UserList{
				Id:         v.Id,
				Identity:   v.Identity,
				NickName:   v.NickName,
				UserName:   v.UserName,
				Password:   v.Password,
				Phone:      v.Phone,
				Rule:       v.Rule,
				CreateTime: v.CreateTime,
				UpdateTime: v.UpdateTime,
				DeleteTime: v.DeleteTime,
			},
		)
	}

	excelData := ToExcel(
		[]string{
			"id",
			"唯一标识",
			"昵称",
			"用户名",
			"密码",
			"手机号",
			"权限",
			"创建时间",
			"更新时间",
			"删除时间",
		}, res,
	)

	helper.ReturnExcel(c, excelData, "用户列表")
}

// @Tags 用户
// @Summary 用户列表
// @Description 用户列表接口
// @Router /user/list [get]
// @Param page query  string true "page"
// @Param size query  string true "size"
// @Produce application/json
// @Success 200 {string} string
func GetUserList(c *gin.Context) {
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, _ := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if size == 0 {
		size = 10
	}
	if page == 0 {
		page = 1
	}
	var count int64

	list := make([]*models.UserList, 0)
	tx := models.GetUserList()
	err := tx.Debug().Omit("content").Offset((page - 1) * size).Limit(size).Find(&list).Count(&count).Error
	if err != nil {
		fmt.Printf("Get Express Error:", err)
		return
	}

	c.JSON(
		http.StatusOK, gin.H{
			"status": 200,
			"data": map[string]interface{}{
				"list":  list,
				"page":  page,
				"size":  size,
				"count": count,
			},
			"message": "请求成功",
		},
	)
}

//  ------------------------------------------------------------------   //

type WXLoginResp struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// @Tags 用户
// @Summary 用户登录
// @Description 用户登录接口
// @Router /login [post]
// @Param username formData  string true "username"
// @Param password formData  string true "password"
// @Produce application/json
// @Success 200 {string} string
func AdminLogin(c *gin.Context) {
	query := helper.PostJson(c)
	username := query["username"]
	password := query["password"]
	fmt.Println(username)
	fmt.Println(password)

	data := new(models.UserList)
	err := models.DB.Where("username = ?", username).First(&data).Error
	if err != nil {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  500,
				"error":   err.Error(),
				"message": "该用户不存在",
			},
		)
		return
	}
	if password != data.Password {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  500,
				"message": "账号密码错误",
			},
		)
		return
	}
	token, _ := helper.GenerateToken(data.Identity)

	c.JSON(
		http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "登陆成功",
			"data": map[string]interface{}{
				"info":  data,
				"token": token,
			},
		},
	)
}

// Login @Tags 用户
// @Summary 用户登录
// @Description 用户登录接口
// @Router /user/login [post]
// @Param code formData  string true "code"
// @Param name formData  string true "名字"
// @Param avatarUrl formData  string true "头像"
// @Produce application/json
// @Success 200 {string} string
func Login(c *gin.Context) {
	code := c.PostForm("code") //  获取 code

	// 根据code获取 openID 和 session_key
	wxLoginResp, err := WXLogin(code)
	fmt.Println("wxLoginResp:%v", wxLoginResp)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}
	token, err := helper.GenerateToken(wxLoginResp.OpenId)
	if err != nil {
		c.JSON(
			http.StatusOK, gin.H{
				"status": -1,
				"msg":    "GenerateToken Error:" + err.Error(),
			},
		)
	}
	data := new(models.UserList)
	err = models.DB.Where("identity = ?", wxLoginResp.OpenId).First(&data).Error
	if err != nil {
		info := &models.UserList{
			Identity: wxLoginResp.OpenId,
		}
		err = models.DB.Create(&info).Error
		if err != nil {
			c.JSON(
				http.StatusOK, gin.H{
					"status":  -1,
					"message": "User Create Error:" + err.Error(),
				},
			)
			return
		}
		c.JSON(
			http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"message": "登陆成功",
				"info":    info,
				"token":   token,
			},
		)
		return
	}
	c.JSON(
		http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "登陆成功",
			"info":    data,
			"token":   token,
		},
	)
}

// 这个函数以 code 作为输入, 返回调用微信接口得到的对象指针和异常情况
func WXLogin(code string) (*WXLoginResp, error) {
	fmt.Printf("code:%v", code)
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"

	// 合成url, 这里的 appId 和 secret 是在微信公众平台上获取的
	url = fmt.Sprintf(url, define.AppId, define.Secret, code)

	// 创建http get请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 解析http请求中body 数据到我们定义的结构体中
	wxResp := WXLoginResp{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&wxResp); err != nil {
		return nil, err
	}

	// 判断微信接口返回的是否是一个异常情况
	if wxResp.ErrCode != 0 {
		return nil, errors.New(fmt.Sprintf("ErrCode:%s  ErrMsg:%s", wxResp.ErrCode, wxResp.ErrMsg))
	}

	return &wxResp, nil
}
