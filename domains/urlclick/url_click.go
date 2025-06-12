package urlclick

import (
	"bitly/domains/urlclick/entities"

	"github.com/gin-gonic/gin"
)

type UrlClickRepository interface {
	Save(data *entities.URLClick) (*entities.URLClick, error)
}

type UrlClickService interface {
	Create(ctx *gin.Context, mappingId uint) (*entities.URLClick, error)
}
