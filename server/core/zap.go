package core

import (
	"BinLog/server/global"
	"log"
	"os"
	"gopkg.in/natefinch/lumberjack.v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)
//InitLogger 初始化并返回一个基于配置设置的新 zap.Logger实例
func InitLogger() *zap.Logger {
	zapConfig := global.Config.Zap

	//创建一个用于日志输出的 writeSyncer
	writeSyncer := getLogWriter(zapConfig.Filename, zapConfig.MaxSize, zapConfig.MaxBackups, zapConfig.MaxAge)

	//如果配置了控制台输出，则添加控制台输出
	if zapConfig.IsConsolePrint {
		writeSyncer = zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout))
	}
	//创建日志格式化的编码器
	encoder := getEncoder()

	//根据配置确定日志级别
	var logLevel zapcore.Level

	if err := logLevel.UnmarshalText([]byte(zapConfig.Level)); err != nil {
		log.Fatalf("Failed to parse log level: %v", err)
	}

	//创建核心和日志实例
	core := zapcore.NewCore(encoder, writeSyncer, logLevel)
	logger := zap.New(core, zap.AddCaller())
	return logger
}
//getLogWriter 返回一个 zapcore.WriteSyncer,该写入器利用lumberjack包，实现日志的滚动记录
func getLogWriter(filename string, maxSize, maxBackups, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:	filename,		//日志文件的位置
		MaxSize:	maxSize,		//在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups:	maxBackups,		//保留旧文件的最大个数
		MaxAge:		maxAge,			//保留旧文件的最大天数
	}
	return zapcore.AddSync(lumberJackLogger)
}
//getEncoder 返回一个为生产日志配置的JSON解释器
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}