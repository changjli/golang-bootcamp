package repositories

import (
	"bitly/domains/urlmapping/entities"
	"context"
)

type UrlMappingRepository interface {
	Save(ctx context.Context, data *entities.URLMapping) (*entities.URLMapping, error)
	FindByShortCode(ctx context.Context, shortCode string) (*entities.URLMapping, error)
	FindByLongUrl(ctx context.Context, longUrl string) (*entities.URLMapping, error)
	Update(ctx context.Context, data *entities.URLMapping) (*entities.URLMapping, error)
}
