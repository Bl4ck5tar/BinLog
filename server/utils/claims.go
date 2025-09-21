package utils

import (
	"BinLog/server/global"
	"BinLog/server/model/appTypes"
	"BinLog/server/model/request"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

//setCookie 设置指定名称和值的 cookie
func setCookie(c *gin.Context,name, value string, maxAge int, host string) {
	//判断 host 是否是 IP 地址
	if net.ParseIP(host) != nil {
		//如果是 IP 地址，设置 cookie 的 domain 为“/”
		c.SetCookie(name, value, maxAge, "/", "", false, true)
	}else {
		//如果是域名，设置 cookie 的 domain 为域名
		c.SetCookie(name, value, maxAge, "/", host, false, true)
	}
}

//SetRefreshToken 设置 Refresh Token 的 cookie
func SetRefreshToken(c *gin.Context, token string, maxAge int) {
	//获取请求的 host，如果失败则取原始请求 host
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}
	//调用 setCookie设置 refresh-token
	setCookie(c, "x-refresh-token", token, maxAge, host)
}

//ClearRefreshToken 清除 Refresh Token 的 cookie
func ClearRefreshToken(c *gin.Context) {
	//获取请求的 host，如果失败则取原始请求 host
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}
	//调用 setCookie 设置 cookie 值为空并过期，删除 refresh-token
	setCookie(c, "x-refresh-token", "", -1, host)
}

//GetAccessToken 从请求头获取 Access Token
func GetAccessToken(c *gin.Context) string {
	//获取 x-access-token 头部值
	token := c.Request.Header.Get("x-access-token")
	return token
}

//GetRefreshToken 从 cookie 获取 Refresh Token
func GetRefreshToken(c *gin.Context) string {
	//尝试从 cookie 中获取 refresh-token
	token, _ := c.Cookie("x-refresh-token")
	return token
}

//GetClaims 从 Gin 的 Context 中解析并获取 JWT 的 Claims
func GetClaims(c *gin.Context) (*request.JwtCustomClaims, error) {
	//获取请求头中的 Access Token
	token := GetAccessToken(c)
	j := NewJWT()
	//解析 Access Token
	claims, err := j.ParseAccessToken(token)
	if err != nil {
		//如果解析失败，记录错误日志
		global.Log.Error("Failed to retrieve JWT parsing information from Gin's Context.")
	}
	return claims, err
}

//GetRefreshClaims 从 Gin 的 Context 中解析并获取 Refresh Token 的 Claims
func GetRefreshClaims(c *gin.Context) (*request.JwtCustomRefreshClaims, error) {
	//获取 Refresh Token
	token := GetRefreshToken(c)
	//创建 JWT 实例
	j := NewJWT()
	
	//解析 Refresh Token
	claims, err := j.ParseRefreshToken(token)
	if err != nil {
		//如果解析失败，记录错误日志
		global.Log.Error("Failed to retrieve JWT parsing information from Gin's Context.")
	}
	return claims, err
}

//GetUserInfo 从 Gin 的 Context 中获取 JWT 解析出来的用户信息（Claims）
func GetUserInfo(c *gin.Context) *request.JwtCustomClaims {
	//首先尝试从 Context 中获取 “claims”
	if claims, exists := c.Get("claims"); !exists {
		//如果不存在，则重新解析 Access Token
		if cl, err := GetClaims(c); err != nil {
			//如果解析失败，返回nil
			return nil
		}else {
			//返回解析出来的用户信息
			return cl
		}
	}else {
		//如果已存在 claims，则直接返回
		waitUse := claims.(*request.JwtCustomClaims)
		return waitUse
	}
}

//GetUserID 从 Gin 的 Context 中获取 JWT 解析出来的用户 ID
func GetUserID(c *gin.Context) uint {
	//如果不存在，则重新解析 Access Token
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			//如果解析失败，返回0
			return 0
		}else {
			//返回解析出的用户ID
			return cl.UserID
		}
	}else {
		//如果已存在 claims，则直接返回用户ID
		waitUse := claims.(*request.JwtCustomClaims)
		return waitUse.UserID
	}
}

//GetUUID 从 Gin 的 Context 中获取 JWT 解析出来的用户 UUID
func GetUUID(c *gin.Context) uuid.UUID {
	if claims, exists := c.Get("claims"); !exists {
		//如果不存在，则重新解析 Access Token
		if cl, err := GetClaims(c); err != nil {
			//如果解析失败，返回一个空 UUID
			return uuid.UUID{}
		}else {
			//返回解析出来的 UUID
			return cl.UUID
		}
	}else {
		//如果已存在 claims，则直接返回 UUID
		waitUse := claims.(*request.JwtCustomClaims)
		return waitUse.UUID
	}
}

//GetRoldID 从 Gin 的 Context 中获取 JWT 解析出来的用户角色ID
func GetRoleID(c *gin.Context) appTypes.RoleID {
	//首先尝试从 Context 中获取“claims”
	if claims, exists := c.Get("claims"); !exists {
		//如果不存在，则重新解析 Access Token
		if cl, err := GetClaims(c); err != nil {
			//如果解析失败，返回0
			return 0
		}else {
			//返回解析出的角色ID
			return cl.RoleID
		}
	}else {
		//如果已存在 claims，则直接返回角色ID
		waitUse := claims.(*request.JwtCustomClaims)
		return waitUse.RoleID
	}
}