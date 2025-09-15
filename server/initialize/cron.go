package initialize

import (
	"BinLog/server/global"
	"os"
	"BinLog/server/task"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

//ZapLogger 结构体实现了 cron.Logger 接口的Info和Error方法， 这些方法用于接收 cron 包生成的日志并使用 zap 进行记录
type ZapLogger struct {
	logger *zap.Logger
}

func (l *ZapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Info(msg, zap.Any("KeysAndValues", keysAndValues))
}

func (l *ZapLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	l.logger.Error(msg, zap.Error(err), zap.Any("keysAndValues", keysAndValues))
}

func NewZapLogger() *ZapLogger {
	return &ZapLogger{logger: global.Log}
}
func InitCron() {
	c := cron.New(cron.WithLogger(NewZapLogger()))
	err := task.RegisterScheduledTasks(c)
	if err != nil {
		global.Log.Error("Error scheduling cron job:", zap.Error(err))
		os.Exit(1)
	}
	c.Start()
}