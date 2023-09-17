package router

import (
	"net/http"
	"twitter_task/service/handler"
	"twitter_task/service/middleware"

	"github.com/gin-gonic/gin"
)

func Initialize(isDebug bool) *gin.Engine {
	// gin debug mode
	if isDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	eg := gin.Default()
	
	// file server
	eg.Static("/static", "./service/static")
	eg.LoadHTMLGlob("service/templates/**/*")

	eg.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "登录twitter",
		})
	})

	eg.GET("twitter/callback", func(c *gin.Context) {
		c.HTML(http.StatusOK, "twitter/callback.html", gin.H{
			"title": "获得twitter授权",
		})
	})

	eg.GET("twitter/test", func(c *gin.Context) {
		c.HTML(http.StatusOK, "twitter/test.html", gin.H{
			"title": "测试关注，发送tweet",
		})
	})

	// cors
	eg.Use(middleware.Cors())

	{
		eg.POST("/twitter/oauth_token", handler.PostGetOAuthToken)
		eg.POST("/twitter/access_token", handler.PostGetAccessToken)
		eg.POST("/twitter/v2/oauth_token", handler.PostGetOAuthTokenV2)
		// eg.POST("/twitter/v2/access_token", handler.PostGetOAuthTokenV2)
		eg.POST("/twitter/v2/follow/:address", handler.PostFollowUserV2)
		eg.POST("/twitter/v2/tweet/:address", handler.PostTweetV2)

		eg.GET("/twitter/v2/isFollow/:address", handler.GetIsFollowV2)
		eg.GET("/twitter/v2/isTweetPost/:address", handler.GetIsTweetPostV2)

		// v1.1 not available
		// eg.GET("/twitter/isFollow/:address", handler.GetIsFollow)
	}

	return eg
}
