package model

type TeamStatusEnum int

const (
	PUBLIC TeamStatusEnum = iota
	SECRET
)

var TeamStatusMap = map[TeamStatusEnum]string{
	PUBLIC: "公开",
	SECRET: "加密",
}
