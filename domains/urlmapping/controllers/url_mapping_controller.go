package controllers

import (
	"bitly/domains/urlmapping"
	requests "bitly/domains/urlmapping/models/requests"
	"bitly/domains/urlmapping/services"
	"bitly/shared/models/responses"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UrlMappingControllerImpl struct {
	service services.UrlMappingService
}

func NewUrlMappingController(service services.UrlMappingService) *UrlMappingControllerImpl {
	return &UrlMappingControllerImpl{
		service: service,
	}
}

func (controller *UrlMappingControllerImpl) ShortenUrl(c *gin.Context) {

	request := &requests.ShortenUrlRequest{}

	if err := c.ShouldBindJSON(request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, responses.BasicResponse{
			Error: err.Error(),
		})
		return
	}

	res, err := controller.service.Create(c, request)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, responses.BasicResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (controller *UrlMappingControllerImpl) RedirectUrl(c *gin.Context) {
	ctx := c.Request.Context()

	code := c.Param("code")

	res, err := controller.service.GetByShortcode(ctx, code)

	if err != nil {
		if errors.Is(err, urlmapping.ErrNotFound) {
			c.JSON(http.StatusNotFound, responses.BasicResponse{
				Error: err.Error(),
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, responses.BasicResponse{
			Error: err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, res)
}
