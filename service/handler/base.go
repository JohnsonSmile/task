package handler

import (
	"net/http"
	"twitter_task/service/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusBadRequest, response.Response{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"code":    http.StatusBadRequest,
		"message": errs.Error(),
	})
}
