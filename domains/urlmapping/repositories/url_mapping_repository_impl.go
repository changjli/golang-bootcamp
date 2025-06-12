package repositories

import (
	"bitly/domains/urlmapping"
	"bitly/domains/urlmapping/entities"
	"bitly/infrastructures"
	"context"
	"errors"

	"gorm.io/gorm"
)

type UrlMappingRepositoryImpl struct {
	db infrastructures.Database
}

func NewUrlMappingRepositoryImpl(db infrastructures.Database) UrlMappingRepository {
	return &UrlMappingRepositoryImpl{
		db: db,
	}
}

func (r *UrlMappingRepositoryImpl) Save(ctx context.Context, data *entities.URLMapping) (*entities.URLMapping, error) {
	result := r.db.GetInstance().Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}

	return data, nil
}

func (r *UrlMappingRepositoryImpl) FindByShortCode(ctx context.Context, shortCode string) (*entities.URLMapping, error) {
	var urlMapping *entities.URLMapping

	result := r.db.GetInstance().First(&urlMapping, "short_code = ?", shortCode)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, urlmapping.ErrNotFound
		}
		return nil, result.Error
	}

	return urlMapping, nil
}

func (r *UrlMappingRepositoryImpl) FindByLongUrl(ctx context.Context, longUrl string) (*entities.URLMapping, error) {
	var urlMapping *entities.URLMapping

	result := r.db.GetInstance().First(&urlMapping, "long_url = ?", longUrl)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, urlmapping.ErrNotFound
		}
		return nil, result.Error
	}

	return urlMapping, nil
}

func (r *UrlMappingRepositoryImpl) Update(ctx context.Context, data *entities.URLMapping) (*entities.URLMapping, error) {
	result := r.db.GetInstance().Save(&data)
	if result.Error != nil {
		return nil, result.Error
	}

	return data, nil
}
