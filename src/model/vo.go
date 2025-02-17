package model

import (
	"time"
)

type UserVO struct {
	Id           int64     `json:"id"`
	Username     string    `json:"username"`
	UserAccount  string    `json:"userAccount"`
	AvatarUrl    string    `json:"avatarUrl"`
	Gender       int8      `json:"gender"`
	UserPassword string    `json:"userPassword"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	UserStatus   int       `json:"userStatus"`
	UserRole     int       `json:"userRole"`
	CreateTime   time.Time `json:"createTime"`
	UpdateTime   time.Time `json:"updateTime"`
	IsDelete     int8      `json:"isDelete"`
	Tags         string    `json:"tags"`
	TagsList     []string  `json:"tagsList"`
}

// TeamUserVO 结构体定义
type TeamUserVO struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	MaxNum      int       `json:"maxNum"`
	ExpireTime  time.Time `json:"expireTime"`
	UserId      int64     `json:"userId"`
	Status      int       `json:"status"`
	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
	CreateUser  UserVO    `json:"createUser"`
	HasJoinNum  int       `json:"hasJoinNum"`
	HasJoin     bool      `json:"hasJoin"`
}
