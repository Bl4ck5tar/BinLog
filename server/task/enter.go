package task

import (
	"BinLog/server/global"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

//每一个小时执行一次更新浏览量操作
func RegisterScheduledTasks(c *cron.Cron) error {
	if _, err := c.AddFunc("@hourly", func() {
		if err := UpdateArticleViewsSyncTask(); err != nil {
			global.Log.Error("Failed to update article views:", zap.Error(err))
		}
	}); err != nil {
		return err
	}
	return nil
}