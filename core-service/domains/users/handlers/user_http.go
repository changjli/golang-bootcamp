package handlers

import (
	"core-service/domains/users"
	"core-service/domains/users/models"
	"core-service/shared/models/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHttp struct {
	UserUseCase users.UserUseCase
}

func NewUserHttp(userUseCase users.UserUseCase) *UserHttp {
	return &UserHttp{
		UserUseCase: userUseCase,
	}
}

func (h *UserHttp) Login(c *gin.Context) {
	var input models.LoginReq
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, err := h.UserUseCase.Login(c, input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, responses.BasicResponse{
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, responses.BasicResponse{
		Data: res,
	})
}

func (h *UserHttp) Register(c *gin.Context) {
	var input models.RegisterReq
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, err := h.UserUseCase.Register(c, input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, responses.BasicResponse{
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, responses.BasicResponse{
		Data: res,
	})
}

func (h *UserHttp) Me(c *gin.Context) {
	res, err := h.UserUseCase.Me(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.BasicResponse{
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, responses.BasicResponse{
		Data: res,
	})
}

func (h *UserHttp) Logout(c *gin.Context) {
	err := h.UserUseCase.Logout(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.BasicResponse{
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, responses.BasicResponse{
		Message: "Logout Successful",
	})
}
