package wizards

import "github.com/gin-gonic/gin"

func Routes(router *gin.Engine) {
	api := router.Group("api")
	short := api.Group("short")
	{
		short.POST("", UrlMappingController.ShortenUrl)
	}
	router.GET(":code", UrlMappingController.RedirectUrl)
}
