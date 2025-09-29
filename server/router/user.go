package router

import (
	"BinLog/server/api"
	"BinLog/server/middleware"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {

}

func (u *UserRouter) InitUserRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup, AdminRouter *gin.RouterGroup) {
	userRouter := Router.Group("user")												//普通用户接口（需登录）
	userPublicRouter :=PublicRouter.Group("user")									//公开接口（无需登录）
	userLoginRouter := PublicRouter.Group("user").Use(middleware.LoginRecord())		//登录/注册接口，每次访问会额外记录日志
	userAdminRouter := AdminRouter.Group("user")									//管理员接口
	userApi := api.ApiGroupApp.UserApi

	{															//普通用户路由
		userRouter.POST("logout", userApi.Logout)									//退出登录
		userRouter.PUT("resetPassword", userApi.UserResetPassword)					//重置密码
		userRouter.GET("info", userApi.UserInfo)									//获取用户信息
		userRouter.PUT("changeInfo", userApi.UserChangeInfo)						//修改用户信息
		userRouter.GET("weather", userApi.UserWeather)								//获取天气
		userRouter.GET("chart", userApi.UserChart)									//获取统计图表数据

	}														
	{															//公共路由
		userPublicRouter.POST("forgotPassword", userApi.ForgotPassword)				//忘记密码
		userPublicRouter.GET("card", userApi.UserCard)								//用户名片
	}
	{															//登录/注册路由
		userLoginRouter.POST("register", userApi.Register)							//注册
		userLoginRouter.POST("login", userApi.Login)								//登录
	}
	{															//管理员路由
		userAdminRouter.GET("list", userApi.UserList)								//获取用户列表
		userAdminRouter.PUT("freeze", userApi.UserFreeze)							//冻结用户
		userAdminRouter.PUT("unfreeze", userApi.UserUnfreeze)						//解冻用户
		userAdminRouter.GET("loginList", userApi.UserLoginList)						//获取用户登录日志
	}
}