package initialize

import (
	"github.com/songzhibin97/gkit/cache/local_cache"
	"BinLog/server/global"
	"BinLog/server/utils"
	"os"
	"go.uber.org/zap"
)

func OtherInit() {
	refreshTokenExpiry, err := utils.ParseDuration(global.Config.Jwt.RefreshTokenExpiryTime)
	if err != nil {
		global.Log.Error("Failed to parse access token expiry time configuration:", zap.Error(err))
		os.Exit(1)
	}
	global.BlackCache = local_cache.NewCache(
		local_cache.SetDefaultExpire(refreshTokenExpiry),
	)
}