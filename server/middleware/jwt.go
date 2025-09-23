package middleware

import (
	"BinLog/server/global"
	"BinLog/server/model/database"
	"BinLog/server/model/request"
	"BinLog/server/model/response"
	"BinLog/server/service"
	"BinLog/server/utils"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

var jwtService = service.ServiceGroupApp.JwtService

//JWTAuth 是一个中间件函数，验证请求中的 JWT token 是否合法
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取请求中的 Access Token 和 Refresh Token
		accessToken := utils.GetAccessToken(c)
		refreshToken := utils.GetRefreshToken(c)

		//检查 Refresh Token 是否在黑名单，如果是，则清除 Refresh Token 并返回未授权错误
		if jwtService.IsInBlacklist(refreshToken) {
			utils.ClearRefreshToken(c)
			response.NoAuth("Account logged in from another location or token is invalid", c)
			c.Abort()	//终止请求的后续处理
			return 
		}

		//创建一个 JWT 实例，用于后续的 token 解析与验证
		j := utils.NewJWT()
		
		//解析 Access Token
		claims, err := j.ParseAccessToken(accessToken)
		if err != nil {
			//如果解析失败并且 Access Token 为空或过期
			if accessToken == "" || errors.Is(err, utils.TokenExpired) {
				//尝试解析 Refresh Token
				refreshClaims, err := j.ParseRefreshToken(refreshToken)
				if err != nil {
					//如果 Refresh Token 也无法解析，清除 Refresh Token 并返回未授权错误
					utils.ClearRefreshToken(c)
					response.NoAuth("Refresh token expired or invalid", c)
					c.Abort()
					return 
				}

				//如果 Refresh Token 有效，通过其 UserID 获取用户信息
				var user database.User
				if err := global.DB.Select("uuid", "role_id").Take(&user, refreshClaims.UserID).Error; err != nil {
					//如果没有找到该用户，清除 Refresh Token并返回未授权错误
					utils.ClearRefreshToken(c)
					response.NoAuth("The user does not exist", c)
					c.Abort()
					return 
				}

				//使用 Refresh Token 的用户信息创建一个新的 Access Token的 Claims
				newAccessClaims := j.CreateAccessClaims(request.BaseClaims{
					UserID: 	refreshClaims.UserID,
					UUID: 		user.UUID,
					RoleID: 	user.RoleID,
				})

				//创建新的 Access Token
				newAccessToken, err := j.CreateAccessToken(newAccessClaims)
				if err != nil {
					//如果生成新的 Access Token 失败，清除 Refresh Token 并返回未授权错误
					utils.ClearRefreshToken(c)
					response.NoAuth("Failed to create new access token", c)
					c.Abort()
					return 
				}

				//将新的 Access Token 和过期时间添加到响应头中
				c.Header("new-access-token", newAccessToken)
				c.Header("new-access-expires-at", strconv.FormatInt(newAccessClaims.ExpiresAt.Unix(), 10))

				//将新的 claims 信息存入 Context，供后续启用
				c.Set("claims", &newAccessClaims)
				c.Next()		//继续后续的处理
				return 
			}

			//如果 Access Token 无效且不满足刷新条件，清除 Refresh Token 并返回未授权错误
			utils.ClearRefreshToken(c)
			response.NoAuth("Invalid access token", c)
			c.Abort()
			return 
		}
		//如果 Access Token 合法，将其 Claims 信息存入 Context
		c.Set("claims", claims)
		c.Next()	//继续后续的处理

	}
}