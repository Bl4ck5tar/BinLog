package initialize

import (
	"BinLog/server/global"
	"os"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

func ConnectRedis() redis.Client {
	redisConfig := global.Config.Redis
	
	//使用单节点配置创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr: 		redisConfig.Address,	//设置Redis 服务器地址
		Password: 	redisConfig.Password,	//设置Redis 密码
		DB: 		redisConfig.DB,			//设置使用的数据库索引
	})

	//Ping Redis服务器以检测连接是否正常
	_, err := client.Ping().Result()
	if err != nil {
		global.Log.Error("Failed to connect to Redis:",zap.Error(err))
		os.Exit(1)
	}


	return *client
}