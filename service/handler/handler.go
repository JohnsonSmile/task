package handler

import (
	"net/http"
	"twitter_task/infra/config"
	"twitter_task/infra/database"
	"twitter_task/service/client"
	"twitter_task/service/request"
	"twitter_task/service/response"
	"twitter_task/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	BIND_TWITTER_MESSAGE = "bind your twitter to AKMD"
)

func PostGetOAuthToken(c *gin.Context) {
	cli := client.NewTwitterClient()
	authToken, err := cli.GetOAuthToken()
	if err != nil {
		zap.L().Error("get token failed", zap.Error(err))
		c.JSON(http.StatusBadGateway, &response.Response{
			Code: http.StatusBadRequest,
			Msg:  "get token failed",
		})
		return
	}
	c.JSON(http.StatusOK, &response.Response{
		Code: http.StatusOK,
		Msg:  "success",
		Data: authToken,
	})
}

func PostGetAccessToken(c *gin.Context) {
	req := request.PostGetAccessTokenRequest{}
	if err := c.ShouldBind(&req); err != nil {
		zap.L().Error("param error", zap.Error(err))
		HandleValidatorError(c, err)
		return
	}
	// verify signedMessage
	address, err := util.VerifyMessage(BIND_TWITTER_MESSAGE, req.SignedMessage)
	if err != nil {
		zap.L().Error("verify message failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, &response.Response{
			Code: http.StatusBadRequest,
			Msg:  "verify message failed",
		})
		return
	}

	if address != req.Address {
		zap.L().Error("verify address failed")
		c.JSON(http.StatusBadRequest, &response.Response{
			Code: http.StatusBadRequest,
			Msg:  "verify address failed",
		})
		return
	}

	cli := client.NewTwitterClient()
	accessToken, err := cli.GetAccessToken(req.OAuthToken, req.OAuthVerifier)
	if err != nil {
		zap.L().Error("get token failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, &response.Response{
			Code: http.StatusBadRequest,
			Msg:  "get token failed",
		})
		return
	}
	// save access token into database...
	user, err := database.SaveAccessToken(address, accessToken.UserID, accessToken.ScreenName, accessToken.OAuthToken, accessToken.OAuthTokenSecret)
	if err != nil {
		zap.L().Error("save token failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, &response.Response{
			Code: http.StatusBadRequest,
			Msg:  "save token failed",
		})
		return
	}
	c.JSON(http.StatusOK, &response.Response{
		Code: http.StatusOK,
		Msg:  "success",
		Data: gin.H{
			"user": gin.H{
				"username": user.Twitter.TWUsername,
				"user_id":  user.Twitter.TWUserID,
			},
		},
	})

}

func PostGetOAuthTokenV2(c *gin.Context) {
	cli := client.NewTwitterClient()
	authURL, err := cli.GenerateOauthV2AuthorizeURL()
	if err != nil {
		zap.L().Error("get authorize URL failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, &response.Response{
			Code: http.StatusBadRequest,
			Msg:  "get authorize URL failed",
		})
		return
	}
	c.JSON(http.StatusOK, &response.Response{
		Code: http.StatusOK,
		Msg:  "success",
		Data: gin.H{
			"url": authURL,
		},
	})
}

func PostFollowUserV2(c *gin.Context) {
	// get account address
	address := c.Param("address")
	if address == "" {
		zap.L().Error("param error, not enough params")
		c.JSON(http.StatusBadRequest, &response.Response{
			Code: http.StatusBadRequest,
			Msg:  "param error, address not passed",
		})
		return
	}
	// get user with address
	user, err := database.GetUserWithAddress(address)
	if err != nil {
		zap.L().Error("get user failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "get user failed",
		})
		return
	}
	cli := client.NewTwitterClient()
	conf := config.GetServerConfig()
	followingInfo, err := cli.FollowUser(user.Twitter.TWUserID, user.Twitter.OAuthToken, user.Twitter.OAuthSecret, conf.AKMDTwitterId)
	if err != nil {
		zap.L().Error("following failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, &response.Response{
			Code: http.StatusBadRequest,
			Msg:  "following failed",
		})
		return
	}
	c.JSON(http.StatusOK, &response.Response{
		Code: http.StatusOK,
		Msg:  "success",
		Data: gin.H{
			"isFollowing": followingInfo.Following,
		},
	})
}

func GetIsFollowV2(c *gin.Context) {
	// get account address
	address := c.Param("address")
	if address == "" {
		zap.L().Error("param error, not enough params")
		c.JSON(http.StatusBadRequest, &response.Response{
			Code: http.StatusBadRequest,
			Msg:  "param error, address not passed",
		})
		return
	}
	// get user with address
	user, err := database.GetUserWithAddress(address)
	if err != nil {
		zap.L().Error("get user failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "get user failed",
		})
		return
	}
	cli := client.NewTwitterClient()
	followingInfo, err := cli.IsFollowUserV2(user.Twitter.TWUserID, user.Twitter.OAuthToken, user.Twitter.OAuthSecret)
	if err != nil {
		zap.L().Error("get user following failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, &response.Response{
			Code: http.StatusBadRequest,
			Msg:  "get user following failed",
		})
		return
	}
	c.JSON(http.StatusOK, &response.Response{
		Code: http.StatusOK,
		Msg:  "success",
		Data: gin.H{
			"isFollowing": followingInfo.Following,
		},
	})
}

func PostTweetV2(c *gin.Context) {
	// get account address
	address := c.Param("address")
	if address == "" {
		zap.L().Error("param error, not enough params")
		c.JSON(http.StatusBadRequest, &response.Response{
			Code: http.StatusBadRequest,
			Msg:  "param error, address not passed",
		})
		return
	}
	// get user with address
	user, err := database.GetUserWithAddress(address)
	if err != nil {
		zap.L().Error("get user failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "get user failed",
		})
		return
	}

	cli := client.NewTwitterClient()
	tweetInfo, err := cli.TweetText(user.Twitter.OAuthToken, user.Twitter.OAuthSecret, "hello")
	if err != nil {
		zap.L().Error("tweet text failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, &response.Response{
			Code: http.StatusBadRequest,
			Msg:  "tweet text failed",
		})
		return
	}
	c.JSON(http.StatusOK, &response.Response{
		Code: http.StatusOK,
		Msg:  "success",
		Data: tweetInfo,
	})
}

func GetIsTweetPostV2(c *gin.Context) {
	// get account address
	address := c.Param("address")
	if address == "" {
		zap.L().Error("param error, not enough params")
		c.JSON(http.StatusBadRequest, &response.Response{
			Code: http.StatusBadRequest,
			Msg:  "param error, address not passed",
		})
		return
	}
	// get user with address
	user, err := database.GetUserWithAddress(address)
	if err != nil {
		zap.L().Error("get user failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "get user failed",
		})
		return
	}

	cli := client.NewTwitterClient()
	tweetInfo, err := cli.IsTweetPost(user.Twitter.OAuthToken, user.Twitter.OAuthSecret, "hello")
	if err != nil {
		zap.L().Error("tweet text failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, &response.Response{
			Code: http.StatusBadRequest,
			Msg:  "tweet text failed",
		})
		return
	}
	c.JSON(http.StatusOK, &response.Response{
		Code: http.StatusOK,
		Msg:  "success",
		Data: tweetInfo,
	})
}

// v1
func GetIsFollow(c *gin.Context) {
	// get account address
	address := c.Param("address")
	if address == "" {
		zap.L().Error("param error, not enough params")
		c.JSON(http.StatusBadRequest, &response.Response{
			Code: http.StatusBadRequest,
			Msg:  "param error, address not passed",
		})
		return
	}
	// get user with address
	user, err := database.GetUserWithAddress(address)
	if err != nil {
		zap.L().Error("get user failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "get user failed",
		})
		return
	}
	cli := client.NewTwitterClient()
	followingInfo, err := cli.IsFollowUser(user.Twitter.OAuthToken, user.Twitter.OAuthSecret, user.Twitter.TWUsername, "akmd-fiance")
	if err != nil {
		zap.L().Error("get user following failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, &response.Response{
			Code: http.StatusBadRequest,
			Msg:  "get user following failed",
		})
		return
	}
	c.JSON(http.StatusOK, &response.Response{
		Code: http.StatusOK,
		Msg:  "success",
		Data: gin.H{
			"isFollowing": followingInfo.Following,
		},
	})
}
