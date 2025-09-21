package service

import (
	"BinLog/server/global"
	"BinLog/server/model/database"
	"BinLog/server/utils"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type JwtService struct {

}

//SetRedisJWT 将 JWT 存储到 Redis 中
func (jwtService *JwtService) SetRedisJWT(jwt string, uuid uuid.UUID) error {
	//解析配置中的 JWT 过期时间
	dr, err := utils.ParseDuration(global.Config.Jwt.RefreshTokenExpiryTime)
	if err != nil {
		return err
	}
	//设置 JWT 在 Redis 中过期的时间
	return global.Redis.Set(uuid.String(), jwt, dr).Err()
}

//GetRedisJWT 从 Redis 中获取 JWT
func (jwtService *JwtService) GetRedisJWT(uuid uuid.UUID) (string, error) {
	//从 Redis 获取指定 uuid 对应的 JWT
	return global.Redis.Get(uuid.String()).Result()
}

//JoinInBlacklist 将 JWT 添加到黑名单
func (jwtService *JwtService) JoinInBlacklist(jwtList database.JwtBlacklist) error {
	//将 JWT 记录插入到数据库中的黑名单表
	if err := global.DB.Create(&jwtList).Error; err != nil {
		return err
	}
	//将 JWT 添加到内存中的黑名单缓存
	global.BlackCache.SetDefault(jwtList.Jwt, struct{}{})
	return nil
}

//IsInBlacklist 检查 JWT 是否在黑名单中
func (jwtService *JwtService) IsInBlacklist(jwt string) bool {
	//从黑名单缓存中检查 JWT 是否存在
	_, ok := global.BlackCache.Get(jwt)
	return ok
}

//LoadAll 从数据库加载所有的 JWT 黑名单并加入缓存
func LoadAll() {
	var data []string
	//从数据库中获取所有的黑名单 JWT
	if err := global.DB.Model(&database.JwtBlacklist{}).Pluck("jwt", &data).Error; err != nil {
		//如果获取失败，记录错误日志
		global.Log.Error("Failed to load JWT blacklist from the database", zap.Error(err))
	}
	//将所有 JWT 添加到 BlackCache 缓存中
	for i:=0; i < len(data); i++ {
		global.BlackCache.SetDefault(data[i], struct{}{})
	}
}