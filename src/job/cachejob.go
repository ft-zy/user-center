package job

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"my-go-user-center/src/config"
	"my-go-user-center/src/utils"
)

func InitJob() {
	// 创建一个新的cron实例
	c := cron.New()

	// 添加一个任务，每天1:30执行
	// 注意cron的表达式格式是 "分 时 日 月 星期"，所以这里设置为 "30 1 * * *"
	// 注意：秒字段是必须的，但在每天特定时间执行的场景下，通常设置为0
	_, err := c.AddFunc("40 13 * * *", docache)
	if err != nil {
		fmt.Println("添加任务失败:", err)
		return
	}

	// 启动cron
	c.Start()

	// 为了看到程序不会立即退出，我们可以让它等待一段时间（或者永远等待，比如使用select{}）
	select {}

}

// 每天执行，预热推荐用户
func docache() {
	lock := utils.NewRedisLock(config.Red, "cacheLock", utils.NextSnowflakeID(), 10)
	if lock.Lock() {
		defer lock.Unlock()
		// todo 将用户的推荐好友加入缓存
	}
}
