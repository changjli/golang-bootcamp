package services

import (
	"bitly/domains/urlclick"
	"bitly/domains/urlclick/entities"

	"github.com/gin-gonic/gin"
)

var (
	baseDomain string
)

type UrlClickServiceImpl struct {
	repo urlclick.UrlClickRepository
}

func NewUrlMappingServiceImpl(repo urlclick.UrlClickRepository) urlclick.UrlClickService {
	return &UrlClickServiceImpl{
		repo: repo,
	}
}

func (srvc *UrlClickServiceImpl) Create(ctx *gin.Context, mappingId uint) (*entities.URLClick, error) {
	mapped := entities.URLClick{
		MappingID: mappingId,
		IPAddress: ctx.ClientIP(),
		UserAgent: ctx.Request.UserAgent(),
	}
	_, err := srvc.repo.Save(&mapped)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
