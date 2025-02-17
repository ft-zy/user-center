package service

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"my-go-user-center/src/config"
	"my-go-user-center/src/constant"
	"my-go-user-center/src/model"
	"my-go-user-center/src/utils"
	"regexp"
	"strconv"
	"time"
)

// RegisterRequest 注册请求结构体
type RegisterRequest struct {
	UserAccount   string `json:"userAccount"`   // 用户账号
	UserPassword  string `json:"userPassword"`  // 用户密码
	CheckPassword string `json:"checkPassword"` // 确认密码
}

// LoginRequest 登录请求结构体
type LoginRequest struct {
	UserAccount  string `json:"userAccount"`  // 用户账号
	UserPassword string `json:"userPassword"` // 用户密码
}

// Login godoc
// @Summary		  登录
// @Schemes		  https
// @Description	          登录
// @Tags	          account
// @accept	          json
// @Produce		  json
// @Param	          account	body param.LoginReq	true	"login"
// @Success		  200		{object}	param.JSONResult{data=param.LoginRes}
// @Router	          /user/login [post]

func Login(req *LoginRequest) *model.User {
	if req.UserAccount == "" || req.UserPassword == "" {
		log.Info("account and password is null")
		panic("输入为空")
	}
	if len(req.UserAccount) < 6 || len(req.UserPassword) < 6 {
		log.Info("account and password is too short")
		panic("密码太短")
	}

	// 账号不能包含特殊字符
	_, err := regexp.MatchString("/[`~!@#$%^&*()_\\-+=<>?:\"{}|,.\\/;'\\\\[\\]·~！@#￥%……&*（）——\\-+={}|《》？：“”【】、；‘'，。、]/", req.UserAccount)
	if err != nil {
		log.Infof("UserAccount have not available char:%v", err)
		panic("账号不能包含特殊字符")
	}
	encryPassWord := utils.EncryptMd5(req.UserPassword)
	user := &model.User{}
	result := config.Db.Model(&model.User{}).Where("isDelete=0 and userAccount=? and userPassword=?", req.UserAccount, encryPassWord).First(&user)
	if result.Error != nil || result.RowsAffected == 0 {
		log.Warningf("查无此人")
		panic("查无此人")
	}
	return user
}

func SearchUser(username string, c *gin.Context) []model.User {
	var user []model.User
	if err := config.Db.Where("isDelete=0 and username like ? and username is not null", "%"+username+"%").Find(&user).Error; err != nil {
		log.Warningf("Serch user information fail:%v", err)
	}
	return user
}

// 根据 id 删除用户
func DeleteById(id int64) {
	config.Db.Model(&model.User{}).Where("id=?", id).Update("isDelete", 1)
}

// 获取当前用户信息
func GetCurrentUser(c *gin.Context) *model.User {

	sessionId, err := c.Cookie(constant.SessionKey)
	if err != nil || sessionId == "" {
		log.Warningf("Get session fail: %v", err)
		panic("当前用户信息cookie拿不到")
	}
	r, _ := config.Red.Get(context.Background(), sessionId).Result()
	var currentUser model.User
	if err := json.Unmarshal([]byte(r), &currentUser); err != nil {
		log.Errorf("Json Unmarshal error:%v", err)
		panic("当前用户信息session拿不到")
	}

	user := model.User{}
	result := config.Db.Model(&model.User{}).Where("isDelete=0 and id = ?", currentUser.Id).Find(&user)
	if result.Error != nil || result.RowsAffected == 0 {
		config.Red.Del(context.Background(), sessionId)
		panic("用户不存在")
	}
	// 设置session
	SetSession(c, &user)
	return &user
}

func UpdateUser(user *model.User, loginUser *model.User) bool {

	if !isAdmin(loginUser) && !(loginUser.Id == user.Id) {
		panic("无权限")
	}
	var olderUser model.User
	result := config.Db.Where("id = ?", user.Id).Find(&olderUser)
	if result.Error != nil || result.RowsAffected == 0 {
		panic("查询失败,没查到数据")
	}
	result = config.Db.Model(&user).Updates(user)
	if result.Error != nil {
		panic("更新用户失败")
	}
	return true
}

func isAdmin(loginUser *model.User) bool {
	return loginUser.UserRole == constant.ADMIN_ROLE
}

func SearchUsersByTags(tagList []string) []model.UserVO {
	var userList []model.User
	config.Db.Model(&model.User{}).Find(&userList)
	filterUserList := make([]model.User, 0)
	for _, user := range userList {
		tags := user.Tags
		if tags == "" {
			continue
		}
		tagNameSet := getTagMapFromString(tags)
		if allContain(tagNameSet, tagList) {
			filterUserList = append(filterUserList, user)
		}
	}
	var userVoList []model.UserVO
	for _, user := range filterUserList {
		userVo := model.UserVO{}
		err := utils.CopyStructFields(&user, &userVo)
		if err != nil {
			panic("copy error")
		}

		var tagsList []string
		err = json.Unmarshal([]byte(user.Tags), &tagsList)
		userVo.TagsList = tagsList
		userVoList = append(userVoList, userVo)
	}
	return userVoList

}

// 判断是否包含所有标签
func allContain(tagNameSet map[string]bool, tagList []string) bool {
	for _, tag := range tagList {
		if _, ok := tagNameSet[tag]; ok == false {
			return false
		}
	}
	return true
}

func getTagMapFromString(tags string) map[string]bool {
	var result []string
	m := make(map[string]bool)
	err := json.Unmarshal([]byte(tags), &result)
	if err != nil {
		log.Errorf("error:%v", err)
		panic(err.Error())
	}
	for _, v := range result {
		m[v] = true
	}
	return m
}

func LoginOut(c *gin.Context) bool {
	idStr, err := c.Request.Cookie(constant.SessionKey)
	if err != nil {
		panic("cookie信息不存在")
	}
	// 删除session信息
	c.SetCookie(constant.SessionKey, "", 0, "/", "", false, true)

	err = config.Red.Del(context.Background(), idStr.Value).Err()

	return err == nil
}

// 设置session信息
func SetSession(c *gin.Context, user *model.User) {
	sessionJson, err := json.Marshal(user)
	if err != nil {
		log.Errorf("Json marshal error:%v", err)
		panic("系统错误")
	}
	// sessionId
	sessionId := constant.SessionId + strconv.FormatInt(user.Id, 10)
	err = config.Red.Set(context.Background(), sessionId, sessionJson, 60*time.Minute).Err()

	if err != nil {
		log.Errorf("Redis set error:%v", err)
		panic("系统错误")
	}
	c.SetCookie(constant.SessionKey, sessionId, constant.CookieExpire, "/", "", false, true)
}

func IsAdmin(c *gin.Context) bool {
	user := GetCurrentUser(c)
	return user.UserRole == constant.ADMIN_ROLE
}
