package request

import (
	"BinLog/server/model/appTypes"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
)

//jwtCustomClaims 结构体用于存储 JWT 的自定义 Claims，继承自 BaseClaims，并包含标准的 JWT 注册信息
type JwtCustomClaims struct {
	BaseClaims				//基础 Claims，包含用户ID、UUID和角色ID
	jwt.RegisteredClaims	//标准 JWT 声明，例如过期时间、发行者等
}
// JwtCustomRefreshClaims 结构体用于存储刷新 Token 的自定义 Claims，包含用户ID和标准的 JWT 注册信息
type JwtCustomRefreshClaims struct {
	UserID 	uint			//用户ID，用于刷新Token相关的身份验证
	jwt.RegisteredClaims	//标准 JWT 声明
}

//BaseClaims 结构用于存储基本的用户信息，作为 JWT 的 Claim 部分
type BaseClaims struct {
	UserID 	uint				//用户ID，表示用户唯一性
	UUID 	uuid.UUID			//用户的UUID，唯一标识用户
	RoleID	appTypes.RoleID		//用户角色ID，表示用户的权限级别
}