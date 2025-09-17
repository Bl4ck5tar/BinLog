package global

import (
	"BinLog/server/config"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config		*config.Config
	Log 		*zap.Logger
	DB			*gorm.DB
	ESClient	*elasticsearch.TypedClient	//全局es客户端
	Redis		redis.Client				//全局redis客户端
	BlackCache	local_cache.Cache			//进程内的本地缓存
)