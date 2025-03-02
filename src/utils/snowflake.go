package utils

import (
	"github.com/bwmarrin/snowflake"
	"strconv"
)

var node *snowflake.Node

// 初始化雪花算法节点
func InitSnowflake() error {
	var err error
	node, err = snowflake.NewNode(1) // 1 是机器ID，可以根据实际情况进行配置
	return err
}

// 获取一个新的雪花ID
func NextSnowflakeID() string {
	if node == nil {
		if err := InitSnowflake(); err != nil {
			return "0"
		}
	}
	return strconv.FormatInt(node.Generate().Int64(), 10)
}
