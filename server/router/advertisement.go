package router

import (
	"BinLog/server/api"

	"github.com/gin-gonic/gin"
)

type AdvertisementRouter struct {

}

func (a *AdvertisementRouter) InitAdvertisementRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	AdvertisementRouter := Router.Group("advertisement")
	advertisementPublicRouter := PublicRouter.Group("advertisement")

	advetisementApi := api.ApiGroupApp.AdvertisementApi
	{
		AdvertisementRouter.POST("create", advetisementApi.AdvertisementCreate)
		AdvertisementRouter.DELETE("delete", advetisementApi.AdvertisementDelete)
		AdvertisementRouter.PUT("update", advetisementApi.AdvertisementUpdate)
		AdvertisementRouter.GET("list", advetisementApi.AdvertisementList)
	}
	{
		advertisementPublicRouter.GET("info", advetisementApi.AdvertisementInfo)
	}
}