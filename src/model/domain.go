package model

import (
	"time"
)

type User struct {
	Id           int64     `gorm:"column:id" json:"id"`
	Username     string    `gorm:"column:username" json:"username"`
	UserAccount  string    `gorm:"column:userAccount" json:"userAccount"`
	AvatarUrl    string    `gorm:"column:avatarUrl" json:"avatarUrl"`
	Gender       int8      `gorm:"column:gender" json:"gender"`
	UserPassword string    `gorm:"column:userPassword" json:"userPassword"`
	Phone        string    `gorm:"column:phone" json:"phone"`
	Email        string    `gorm:"column:email" json:"email"`
	UserStatus   int       `gorm:"column:userStatus" json:"userStatus"`
	UserRole     int       `gorm:"column:userRole" json:"userRole"`
	CreateTime   time.Time `gorm:"column:createTime" json:"createTime"`
	UpdateTime   time.Time `gorm:"column:updateTime" json:"updateTime"`
	IsDelete     int8      `gorm:"column:isDelete" json:"isDelete"`
	Tags         string    `gorm:"column:tags" json:"tags"`
}

func (User) TableName() string {
	return "yp_user"
}

// Team 对应Java中的Team类
type Team struct {
	Id          int64     `gorm:"column:id" json:"id"`                   // 使用int64代替Long，因为Go中没有Long类型
	Name        string    `gorm:"column:name" json:"name"`               // 队伍名称
	Description string    `gorm:"column:description" json:"description"` // 描述
	MaxNum      int       `gorm:"column:maxNum" json:"maxNum"`           // 最大人数
	ExpireTime  time.Time `gorm:"column:expireTime" json:"expireTime"`   // 过期时间
	UserID      int64     `gorm:"column:userId" json:"userId"`           // 用户id
	Status      int       `gorm:"column:status" json:"status"`           // 队伍状态
	Password    string    `gorm:"column:password" json:"password"`       // 密码
	CreateTime  time.Time `gorm:"column:createTime" json:"createTime"`   // 创建时间
	UpdateTime  time.Time `gorm:"column:updateTime" json:"updateTime"`   // 更新时间
	IsDelete    int       `gorm:"column:isDelete" json:"isDelete"`       // 是否删除
}

func (Team) TableName() string {
	return "team"
}

type YpTags struct {
	ID         int64     `json:"id" gorm:"column:id"`
	TagsName   string    `json:"tagsName" gorm:"column:tagsName"`
	ParentID   int64     `json:"parentId,omitempty" gorm:"column:parentId"`
	IsParent   bool      `json:"isParent,omitempty" gorm:"column:isParent"`
	CreateTime time.Time `json:"createTime,omitempty" gorm:"column:createTime"`
	UpdateTime time.Time `json:"updateTime,omitempty" gorm:"column:updateTime"`
	UserID     int64     `json:"userId,omitempty" gorm:"column:userId"`
	IsDelete   bool      `json:"isDelete,omitempty" gorm:"column:isDelete"`
}

func (YpTags) TableName() string {
	return "yp_tags"
}

// UserTeam 对应数据库中的 user_team 表
type UserTeam struct {
	ID         int64     `gorm:"primaryKey;autoIncrement;column:column:id;comment:'id'"`
	UserID     int64     `gorm:"column:userId;comment:'用户id';index"`
	TeamID     int64     `gorm:"column:teamId;comment:'队伍id';index"`
	JoinTime   time.Time `gorm:"column:joinTime;comment:'加入时间'"`
	CreateTime time.Time `gorm:"column:createTime;default:CURRENT_TIMESTAMP;comment:'创建时间'"`
	UpdateTime time.Time `gorm:"column:updateTime;default:CURRENT_TIMESTAMP;type:timestamp;autoUpdateTime;comment:'更新时间'"`
	IsDelete   bool      `gorm:"column:isDelete;default:0;not null;comment:'是否删除'"`
}

func (UserTeam) TableName() string {
	return "user_team"
}
