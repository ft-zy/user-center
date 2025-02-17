package router

import (
	"github.com/gin-gonic/gin"
	"my-go-user-center/src/api"
	"my-go-user-center/src/middleware"
)

func InitRouterAndServe() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	// 全局异常处理
	r.Use(middleware.ErrorHandler())
	userGroup := r.Group("/api/user")
	{
		authGroup := userGroup.Group("/auth")
		authGroup.Use(middleware.AuthMiddleWare())
		{
			authGroup.GET("/search", api.SearchUser)
			authGroup.GET("/delete", api.DeleteUser)
			// 标签
			authGroup.POST("/createTag", api.CreateTag)
			authGroup.POST("/updateTag", api.UpdateTag)
			authGroup.POST("/deleteTag", api.DeleteTag)

		}
		userGroup.POST("/listTag", api.ListTag)
		userGroup.POST("/register", api.Register)
		userGroup.POST("/login", api.Login)
		userGroup.POST("/loginout", api.LoginOut)
		userGroup.POST("/update", api.UpdateUser)
		userGroup.GET("/current", api.Current)
		userGroup.POST("/searchusersbytags", api.SearchUsersByTags)
		userGroup.GET("/recommend", api.Recommend)
		userGroup.GET("/test", api.Test)
		userGroup.GET("/userSaveBath", api.UserSaveBatch)
	}
	// 队伍接口
	teamGroup := r.Group("/api/team")
	{
		teamGroup.POST("/addTeam", api.AddTeam)
		teamGroup.POST("/listTeams", api.ListTeams)
		teamGroup.POST("/joinTeam", api.JoinTeam)
		teamGroup.POST("/quit", api.QuitTeam)
	}
	return r

}
