package handlers

import (
	"login-api/domains/users/models/requests"
	"login-api/domains/users/usecases"
	"login-api/shared/models/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHttp struct {
	UserUseCase usecases.UserUseCaseInterface
}

func NewUserHttp(userUseCase usecases.UserUseCaseInterface) *UserHttp {
	return &UserHttp{
		UserUseCase: userUseCase,
	}
}

func (h *UserHttp) Login(c *gin.Context) {
	var input requests.UserLoginRequest
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, err := h.UserUseCase.Login(c, input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *UserHttp) Me(c *gin.Context) {
	res, err := h.UserUseCase.Me(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *UserHttp) Logout(c *gin.Context) {
	err := h.UserUseCase.Logout(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, responses.BasicResponse{
		Data: "Logout success",
	})
}
