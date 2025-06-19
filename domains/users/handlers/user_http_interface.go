package handlers

import "github.com/gin-gonic/gin"

type UserHttpInterface interface {
	Login(c *gin.Context)
	Me(c *gin.Context)
	Logout(c *gin.Context)
}
