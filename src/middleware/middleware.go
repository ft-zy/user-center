package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"my-go-user-center/src/common"
	"my-go-user-center/src/config"
	"my-go-user-center/src/constant"
	"my-go-user-center/src/model"
	"runtime"
	"strings"
)

func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

// 全局异常处理
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				//
				c.Abort()
				common.RespFail(c.Writer, message)
			}
		}()
		c.Next()

	}
}

// 中间件鉴权
func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr, err := c.Cookie(constant.SessionKey)
		if err != nil || idStr == "" {
			panic("cookie信息不存在")
		}
		// 如果没有出现错误且 session 不为空，说明存在有效的 session
		// 则调用 c.Next() 继续处理后续的请求处理函数，即允许通过该中间件
		var user model.User
		ctx := context.Background()
		r, _ := config.Red.Get(ctx, idStr).Result()
		if err := json.Unmarshal([]byte(r), &user); err != nil {
			log.Errorf("Json Unmarshal error:%v", err)
			panic(err)
		}
		// 管理员权限
		if user.UserRole != 1 {
			panic("权限不足")
		}
		c.Next()
	}
}
