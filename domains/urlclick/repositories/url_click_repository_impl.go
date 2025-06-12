package repositories

import (
	"bitly/domains/urlclick"
	"bitly/domains/urlclick/entities"
	"bitly/infrastructures"
)

type UrlClickRepositoryImpl struct {
	db infrastructures.Database
}

func NewUrlClickRepositoryImpl(db infrastructures.Database) urlclick.UrlClickRepository {
	return &UrlClickRepositoryImpl{
		db: db,
	}
}

func (repo *UrlClickRepositoryImpl) Save(data *entities.URLClick) (*entities.URLClick, error) {
	result := repo.db.GetInstance().Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}

	return data, nil
}
