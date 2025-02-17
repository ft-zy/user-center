package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"my-go-user-center/src/common"
	"my-go-user-center/src/config"
	"my-go-user-center/src/model"
	"my-go-user-center/src/service"
	"my-go-user-center/src/utils"
	"regexp"
	"strconv"
	"sync"
	"time"
)

// 更新用户
func UpdateUser(c *gin.Context) {
	user := &model.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Errorf("updateUser json err %v", err)
		common.RespFail(c.Writer, "请求参数错误")
		return
	}
	loginUser := service.GetCurrentUser(c)
	result := service.UpdateUser(user, loginUser)
	common.RespOK(c.Writer, result, "更新成功")
}

// 查标签
func SearchUsersByTags(c *gin.Context) {
	var tagNameList []string
	err := c.ShouldBindJSON(&tagNameList)
	if err != nil {
		log.Errorf(" json err %v", err)
		common.RespFail(c.Writer, "updateUser json err")
	}
	if len(tagNameList) == 0 {
		common.RespFail(c.Writer, "len is 0")
	}
	var userVoList = service.SearchUsersByTags(tagNameList)
	common.RespOK(c.Writer, userVoList, "查询成功")
}

// 推荐的朋友
func Recommend(c *gin.Context) {
	var userList []model.User
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	loginUser := service.GetCurrentUser(c)
	redisKey := fmt.Sprintf("user:recommend:%v", loginUser.Id)
	r, _ := config.Red.Get(context.Background(), redisKey).Result()
	if r != "" {
		json.Unmarshal([]byte(r), &userList)
		common.RespOK(c.Writer, userList, "查询成功")
		return
	}

	// 计算 Offset
	offset := (pageNum - 1) * pageSize
	// 使用 Offset 和 Limit 进行分页
	config.Db.Offset(offset).Limit(pageSize).Find(&userList)
	marshal, err := json.Marshal(userList)
	if err != nil {
		log.Errorf("json.Marshal error:%v", err)
	}
	stat := config.Red.Set(context.Background(), redisKey, marshal, time.Minute)
	if stat.Err() != nil {
		log.Errorf("Redis set error:%v", stat.Err())
	}
	common.RespOK(c.Writer, userList, "查询成功")
}

// 注册
func Register(c *gin.Context) {
	req := &service.RegisterRequest{}
	// 从前端接收请求
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("Login request json err %v", err)
		common.RespFail(c.Writer, err.Error())
		return
	}
	if req.UserAccount == "" || req.UserPassword == "" || req.CheckPassword == "" {
		common.RespFail(c.Writer, "必填项为空")
		return
	}
	if len(req.UserAccount) < 6 || len(req.UserPassword) < 6 || len(req.CheckPassword) < 6 {
		log.Info("UserAccount or UserPassword or CheckPassword is too short")
		common.RespFail(c.Writer, "输入不能小于六位")
		return
	}

	// 账号不能包含特殊字符
	_, err := regexp.MatchString("/[`~!@#$%^&*()_\\-+=<>?:\"{}|,.\\/;'\\\\[\\]·~！@#￥%……&*（）——\\-+={}|《》？：“”【】、；‘'，。、]/", req.UserAccount)
	if err != nil {
		log.Infof("UserAccount have not available char:%v", err)
		common.RespFail(c.Writer, "账号不能包含特殊字符")
		return
	}

	// 密码与校验密码
	if req.UserPassword != req.CheckPassword {
		log.Info("UserPassword != CheckPassword")
		common.RespFail(c.Writer, "密码与校验密码不相同")
		return
	}

	if row := config.Db.Where("userAccount=?", req.UserAccount).First(&model.User{}).RowsAffected; row > 0 {
		log.Info("userAccount is already setup")
		common.RespFail(c.Writer, "账号已被注册")

		return
	}
	// 2. 加密
	encryptPassword := utils.EncryptMd5(req.UserPassword)

	// 3. 插入数据
	user := &model.User{
		UserAccount:  req.UserAccount,
		UserPassword: encryptPassword,
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}
	log.Infof("user ====== %+v", user)
	if err = config.Db.Model(&model.User{}).Create(user).Error; err != nil {
		log.Warningf("Create user fail:%v", err)
		common.RespFail(c.Writer, "数据库创建用户失败")
		return
	}

	common.RespOK(c.Writer, user.Id, "成功")
}

// Login 登录
// @Summary		  登录
// @Description	          登录
// @Tags	          user
// @accept	          json
// @Produce		  json
// @Router	          /api/user/login [post]
func Login(c *gin.Context) {
	req := &service.LoginRequest{}
	// 前端请求
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("Login request json err %v", err)
		common.RespFail(c.Writer, err.Error())
		return
	}
	// 处理请求
	user := service.Login(req)
	// 设置session
	service.SetSession(c, user)
	// 脱敏
	common.RespOK(c.Writer, utils.GetSafetyUser(user), "成功")
}

// LoginOut 注销
// @Summary		  注销
// @Description	          注销
// @Tags	          user
// @Produce		  json
// @Router	          /api/user/loginout [post]
func LoginOut(c *gin.Context) {
	result := service.LoginOut(c)
	common.RespOK(c.Writer, result, "注销成功")
}

// 获取用户消息
func SearchUser(c *gin.Context) {

	username := c.Query("username")

	userList := service.SearchUser(username, c)
	var safetyUserList []model.User
	for _, v := range userList {
		safetyUser := utils.GetSafetyUser(&v)
		safetyUserList = append(safetyUserList, safetyUser)
	}
	common.RespOK(c.Writer, safetyUserList, "查询成功")
}

// 删除用户信息
func DeleteUser(c *gin.Context) {

	// 从前端接收请求
	sid := c.Query("id")
	id, _ := strconv.ParseInt(sid, len(sid), 64)

	// 根据 id 删除用户
	service.DeleteById(id)

	common.RespOK(c.Writer, "删除成功", "删除成功")
}

// 获取当前用户
func Current(c *gin.Context) {
	common.RespOK(c.Writer, service.GetCurrentUser(c), "查询成功")
}
func TestUnLock(c *gin.Context) {

}

func Test(c *gin.Context) {
	lock := utils.NewRedisLock(config.Red, "join_teamfdsfdsfds", utils.NextSnowflakeID(), 10*time.Minute)
	if lock.Lock() {
		common.RespOK(c.Writer, "suces", "ss")
	} else {
		common.RespFail(c.Writer, "fail")

	}

}

// 模拟用户数据导入
func UserSaveBatch(c *gin.Context) {

	errorChan := make(chan error, 10)
	batchSize := 5000
	// 控制go程的并发数量为5
	goChan := make(chan int, 10)
	var wg sync.WaitGroup
	start := time.Now()
	// 模拟用户数据
	for i := 0; i < 20; i++ {
		userList := make([]model.User, 0, batchSize)
		for j := 0; j < batchSize; j++ {
			userList = append(userList, model.User{
				UserAccount:  "test" + time.Now().Format("20060102150405"),
				UserPassword: "123456",
				Username:     "test",
				Gender:       1,
				UserRole:     0,
				UserStatus:   0,
				Email:        "",
				Phone:        "",
				Tags:         "",
				AvatarUrl:    "",
				IsDelete:     0,
				CreateTime:   time.Now(),
				UpdateTime:   time.Now(),
			})
		}
		wg.Add(1)
		go func(user []model.User, index int) {
			defer wg.Done()
			goChan <- 1
			fmt.Printf("Started saving batch in goroutine %d\n", index)
			// 模拟耗时2s
			time.Sleep(1 * time.Second)
			if i%2 == 0 {
				errorChan <- error(fmt.Errorf("error occurred in goroutine %d", index))
			}
			<-goChan
		}(userList, i)
	}
	wg.Wait()
	close(errorChan) // 所有goroutine完成后关闭channel
	// 从channel中读取错误
	var errors []error
	for err := range errorChan {
		errors = append(errors, err)
	}

	// 检查是否有错误发生
	if len(errors) > 0 {
		fmt.Println("导入过程中发生错误：")
		for _, err := range errors {
			fmt.Println(err)
		}
	} else {
		fmt.Println("所有数据导入成功")
	}
	elapsed := time.Since(start)
	fmt.Printf("Total time taken: %v\n", elapsed)
}
