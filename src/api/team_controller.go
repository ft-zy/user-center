package api

import (
	"github.com/gin-gonic/gin"
	"my-go-user-center/src/common"
	"my-go-user-center/src/config"
	"my-go-user-center/src/constant"
	"my-go-user-center/src/model"
	"my-go-user-center/src/service"
)

// AddTeam 创建队伍
func AddTeam(c *gin.Context) {
	teamAddRequest := model.TeamAddRequest{}
	if err := c.ShouldBindJSON(&teamAddRequest); err != nil {
		panic(err.Error())
	}
	loginUser := service.GetCurrentUser(c)

	teamId := service.AddTeam(&teamAddRequest, loginUser)
	common.RespOK(c.Writer, teamId, "添加成功")
}

// UpdateTeam 修改队伍
func UpdateTeam(c *gin.Context) {
	teamAddRequest := model.TeamAddRequest{}
	if err := c.ShouldBindJSON(&teamAddRequest); err != nil {
		panic(err.Error())
	}
	loginUser := service.GetCurrentUser(c)

	teamId := service.UpdateTeam(&teamAddRequest, loginUser)
	common.RespOK(c.Writer, teamId, "修改成功")
}

// ListTeams 查询队伍
func ListTeams(c *gin.Context) {
	var teamQuery model.TeamQuery
	if err := c.ShouldBindJSON(&teamQuery); err != nil {
		panic(err.Error())
	}
	loginUser := service.GetCurrentUser(c)
	teamList := service.ListTeams(&teamQuery, loginUser.UserRole == constant.ADMIN_ROLE)
	if len(teamList) == 0 {
		common.RespOK(c.Writer, teamList, "查询成功")
		return
	}
	// 提取队伍ID
	teamIds := make([]int64, len(teamList))
	for i, team := range teamList {
		teamIds[i] = team.Id
	}

	// 登录人加入的队伍
	var loginJoinUserTeamList []model.UserTeam
	config.Db.Model(&model.UserTeam{}).
		Where("userId = ? and teamId IN (?)", loginUser.Id, teamIds).
		Find(&loginJoinUserTeamList)

	// 登录人加入的队伍标记
	hasJoinMap := make(map[int64]bool)
	for _, ut := range loginJoinUserTeamList {
		hasJoinMap[ut.TeamID] = true
	}

	for _, team := range teamList {
		team.HasJoin = hasJoinMap[team.Id]
	}

	// 查询队伍的人数
	var userTeamList []model.UserTeam
	config.Db.Model(&model.UserTeam{}).Where("teamId in (?)", teamIds).Find(&userTeamList)
	teamUserCountMap := make(map[int64]int)
	for _, ut := range userTeamList {
		teamUserCountMap[ut.TeamID]++
	}

	for _, team := range teamList {
		team.HasJoinNum = teamUserCountMap[team.Id]
	}
	common.RespOK(c.Writer, teamList, "查询成功")
}

// JoinTeam 加入队伍
func JoinTeam(c *gin.Context) {
	var teamJoin model.TeamJoinRequest
	if err := c.ShouldBindJSON(&teamJoin); err != nil {
		panic(err.Error())
	}
	loginUser := service.GetCurrentUser(c)
	common.RespOK(c.Writer, service.JoinTeam(&teamJoin, loginUser), "查询成功")

}

// QuitTeam 退出队伍
func QuitTeam(c *gin.Context) {
	request := &model.TeamQuitRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		panic(err.Error())
	}
	loginUser := service.GetCurrentUser(c)
	common.RespOK(c.Writer, service.QuitTeam(request, loginUser), "退出队伍")
}
