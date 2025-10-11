package router

import (
	"BinLog/server/api"

	"github.com/gin-gonic/gin"
)

type CommentRouter struct {

}

func (c *CommentRouter) InitCommentRouter(Router *gin.RouterGroup, PubilcRouter *gin.RouterGroup, AdminRouter *gin.RouterGroup) {
	commentRouter := Router.Group("comment")
	commentPublicRouter := PubilcRouter.Group("comment")
	commentAdminRouter := AdminRouter.Group("comment")

	commentApi := api.ApiGroupApp.CommentApi
	{
		commentRouter.POST("create", commentApi.CommentCreate)
		commentRouter.DELETE("delete", commentApi.CommentDelete)
		commentRouter.GET("info", commentApi.CommentInfo)
	}
	{
		commentPublicRouter.GET(":article_id", commentApi.CommentInfoByArticleID)
		commentPublicRouter.GET("new", commentApi.CommentNew)
	}
	{
		commentAdminRouter.GET("list", commentApi.CommentList)
	}
}