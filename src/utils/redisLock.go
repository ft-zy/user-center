package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisLock struct {
	rdb     *redis.Client
	key     string
	value   string // 保存一个随机数，用于lua脚本解锁
	timeout time.Duration
}

func NewRedisLock(rdb *redis.Client, key, value string, timeout time.Duration) *RedisLock {
	return &RedisLock{
		rdb:     rdb,
		key:     key,
		value:   value,
		timeout: timeout,
	}
}

// Lock 尝试获取锁
func (l *RedisLock) Lock() bool {
	// 如果键不存在，则设置成功并设置过期时间
	result, err := l.rdb.SetNX(context.Background(), l.key, l.value, l.timeout).Result()
	if err != nil || result == false {
		fmt.Println("Failed to set lock:", err)
		return false
	}
	return true
}

// Unlock 释放锁
// 使用Lua脚本来确保只有锁的持有者才能释放锁
func (l *RedisLock) Unlock() bool {
	// 使用Lua脚本来安全地释放锁
	// Lua脚本检查键是否存在并且值是否与预期值匹配，然后删除键
	script := `  
	if redis.call("get", KEYS[1]) == ARGV[1] then  
		return redis.call("del", KEYS[1])  
	else  
		return 0  
	end  
	`
	result, err := l.rdb.Eval(context.Background(), script, []string{l.key}, l.value).Result()
	if err != nil || result.(int64) == 0 {
		fmt.Println("Failed to unlock:", err)
		return false
	}
	return true
}
