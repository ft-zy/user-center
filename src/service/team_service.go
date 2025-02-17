package service

import (
	log "github.com/sirupsen/logrus"
	"my-go-user-center/src/config"
	"my-go-user-center/src/model"
	"my-go-user-center/src/utils"
	"strconv"
	"time"
)

func AddTeam(teamAddRequest *model.TeamAddRequest, loginUser *model.User) int64 {
	team := model.Team{}
	err := utils.CopyStructFields(&teamAddRequest, &team)
	if err != nil {
		log.Warning("CopyStructFields: %v", err.Error())
		panic(err.Error())
	}
	team.ExpireTime = teamAddRequest.ExpireTime.Time
	team.CreateTime = time.Now()
	team.UpdateTime = time.Now()
	// 3. 校验信息
	//   1. 队伍人数 > 1 且 <= 20
	if team.MaxNum < 1 || team.MaxNum > 20 {
		panic("队伍人数不满足要求")
	}
	//   2. 队伍标题 <= 20
	if len(team.Name) < 1 || len(team.Name) > 20 {
		panic("队伍标题不满足要求")
	}
	//   3. 描述 <= 512
	if len(team.Name) > 512 {
		panic("队伍描述过长")
	}
	//   4. status 是否公开（int）不传默认为 0（公开）

	if _, ok := model.TeamStatusMap[model.TeamStatusEnum(team.Status)]; !ok {
		panic("队伍状态不满足要求")
	}

	//   5. 如果 status 是加密状态，一定要有密码，且密码 <= 32
	if model.SECRET == model.TeamStatusEnum(team.Status) {
		if len(team.Password) == 0 || len(team.Password) > 32 {
			panic("密码设置不正确")
		}
	}
	// 6. 超时时间 早于 当前时间
	if time.Now().After(team.ExpireTime) {
		// 抛出错误
		panic("超时时间 早于 当前时间")
	}

	// 7. 校验用户最多创建 5 个队伍
	// todo 有 bug，可能同时创建 100 个队伍
	var count int64
	config.Db.Model(&model.Team{}).Where("userId = ? and isDelete = 0", loginUser.Id).Count(&count)
	if count >= 5 {
		panic("用户最多创建 5 个队伍")
	}
	// 开启事务
	tx := config.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("回滚")
			tx.Rollback()
			panic(r)
		}
	}()
	// 8. 插入队伍信息到队伍表
	team.UserID = loginUser.Id
	config.Db.Model(&model.Team{}).Create(team)
	// 9. 插入用户  => 队伍关系到关系表
	userTeam := &model.UserTeam{
		UserID:     loginUser.Id,
		TeamID:     team.Id,
		JoinTime:   time.Now(),
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	config.Db.Model(&model.UserTeam{}).Create(userTeam)
	tx.Commit()
	return team.Id
}

func ListTeams(m *model.TeamQuery, admin bool) []model.TeamUserVO {
	var teams []model.Team
	db := config.Db.Model(&model.Team{})
	// 计算 Offset
	offset := (m.PageNum - 1) * m.PageSize
	// 使用 Offset 和 Limit 进行分页
	db.Offset(offset).Limit(m.PageNum)
	// 查询条件
	if m.SearchText != "" {
		db.Where("name like ? or description like ?", "%"+m.SearchText+"%", "%"+m.SearchText+"%")
	}
	if !admin && model.TeamStatusEnum(m.Status) == model.SECRET {
		panic("非管理员不能看私有房间")
	}
	// 不展示已过期的队伍
	// expireTime is null or expireTime > now()
	db.Where("isDelete = 0")
	db.Where("expireTime is null or expireTime > ?", time.Now()).Find(&teams)

	if len(teams) == 0 {
		return []model.TeamUserVO{}
	}
	var teamUserVOList []model.TeamUserVO
	// 关联查询创建人的用户信息
	// 查询用户信息
	userIdsMap := make(map[int64]bool)
	userIds := make([]int64, 0, 10)
	for _, team := range teams {
		if flag, ok := userIdsMap[team.Id]; ok && flag {
			continue
		}
		userIdsMap[team.UserID] = true
		userIds = append(userIds, team.UserID)
	}
	// 查询关联用户
	var users []model.User
	config.Db.Model(&model.User{}).Where("id in (?)", userIds).Find(&users)
	userId2UserMap := make(map[int64]model.User)
	for _, user1 := range users {
		userId2UserMap[user1.Id] = user1
	}
	// 保存用户信息
	for _, team := range teams {
		teamUserVO := model.TeamUserVO{}
		utils.CopyStructFields(&team, &teamUserVO)
		if user, ok := userId2UserMap[team.UserID]; ok {
			teamUserVO.CreateUser = model.UserVO{}
			utils.CopyStructFields(&user, &teamUserVO.CreateUser)
		}
		teamUserVOList = append(teamUserVOList, teamUserVO)
	}

	return teamUserVOList
}

func UpdateTeam(m *model.TeamAddRequest, user *model.User) interface{} {
	// 校验入参合法性
	// 检查队伍是否过期和存在
	// 修改
	return nil
}

// 入队
func JoinTeam(joinRequest *model.TeamJoinRequest, user *model.User) bool {
	team := model.Team{}
	result := config.Db.Model(&model.Team{}).Where("id = ? and isDelete = 0", joinRequest.TeamId).Find(&team)
	if result.Error != nil || result.RowsAffected == 0 {
		panic("队伍不存在")
	}
	if team.ExpireTime.Before(time.Now()) {
		panic("队伍已过期")
	}
	if model.SECRET == model.TeamStatusEnum(team.Status) && team.Password != joinRequest.Password {
		panic("密码错误")
	}
	// 加锁防止重复入队
	lock := utils.NewRedisLock(config.Red, "join_team:"+strconv.FormatInt(team.Id, 10), utils.NextSnowflakeID(), 10*time.Minute)
	if lock.Lock() {
		defer lock.Unlock()
		var count int64
		config.Db.Model(&model.UserTeam{}).Where("userId = ? and isDelete = 0", user.Id).Count(&count)
		if count > 5 {
			panic("用户最多加入 5 个队伍")
		}
		config.Db.Model(&model.UserTeam{}).Where("userId = ? and teamId = ? and isDelete = 0", user.Id, team.Id).Count(&count)
		if count > 0 {
			panic("用户已加入该队伍")
		}
		// 当前队伍人数
		var teamUserCount int64
		config.Db.Model(&model.UserTeam{}).Where("teamId = ? and isDelete = 0", team.Id).Count(&teamUserCount)
		if teamUserCount > int64(team.MaxNum) {
			panic("队伍已满")
		}
		// 添加队伍信息
		userTeam := &model.UserTeam{
			UserID:     user.Id,
			TeamID:     team.Id,
			JoinTime:   time.Now(),
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
		config.Db.Model(&model.UserTeam{}).Create(userTeam)
		return true
	}
	return false
}

// 退出队伍
func QuitTeam(quitRequest *model.TeamQuitRequest, loginUser *model.User) bool {
	tx := config.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("回滚")
			tx.Rollback()
			panic(r)
		}
	}()
	team := &model.Team{}
	result := config.Db.Model(&model.Team{}).Where("id = ? and isDelete = 0", quitRequest.TeamId).Find(team)
	if result.Error != nil || result.RowsAffected == 0 {
		panic("队伍不存在")
	}
	var teamHasJoinNum int64
	config.Db.Model(&model.UserTeam{}).Where("teamId = ? and isDelete = 0", quitRequest.TeamId).Count(&teamHasJoinNum)
	// 队伍只剩一人，解散
	if teamHasJoinNum == 1 {
		// 删除队伍
		config.Db.Model(&model.Team{}).Where("id = ?", quitRequest.TeamId).Update("isDelete", 1)
	} else {
		// 队伍还剩至少两人
		// 是队长
		if team.UserID == loginUser.Id {
			// 把队伍转移给最早加入的用户
			// 1. 查询已加入队伍的所有用户和加入时间
			var userTeams []model.UserTeam
			result := config.Db.Model(&model.UserTeam{}).Where("teamId = ? and isDelete = 0", quitRequest.TeamId).Order("joinTime ASC").Limit(2).Find(&userTeams)
			if result.Error != nil || result.RowsAffected == 0 {
				panic("系统错误")
			}
			nextTeamLeaderId := userTeams[1].UserID
			result = config.Db.Model(&model.Team{}).Where("id = ?", quitRequest.TeamId).Update("userId", nextTeamLeaderId)
			if result.Error != nil || result.RowsAffected == 0 {
				panic("更新队伍队长失败")
			}
		}
	}
	config.Db.Model(&model.UserTeam{}).Where("userId = ? and teamId = ? and isDelete = 0", loginUser.Id, quitRequest.TeamId).Update("isDelete", 1)
	tx.Commit()
	return true
}
