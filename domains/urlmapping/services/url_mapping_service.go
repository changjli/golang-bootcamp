package services

import (
	"bitly/domains/urlmapping/models/requests"
	"bitly/domains/urlmapping/models/responses"
	"context"

	"github.com/gin-gonic/gin"
)

type UrlMappingService interface {
	Create(ctx *gin.Context, request *requests.ShortenUrlRequest) (*responses.ShortenUrlResponse, error)
	GetByShortcode(ctx context.Context, shortCode string) (string, error)
}
