package services

import (
	"bitly/domains/urlclick"
	"bitly/domains/urlmapping/entities"
	"bitly/domains/urlmapping/models/requests"
	"bitly/domains/urlmapping/models/responses"
	"bitly/domains/urlmapping/repositories"
	"context"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	baseDomain string
)

type UrlMappingServiceImpl struct {
	repo            repositories.UrlMappingRepository
	urlClickService urlclick.UrlClickService
}

func encodeID(id uint) string {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(id))

	return base64.RawURLEncoding.EncodeToString(buf)
}

func NewUrlMappingServiceImpl(repo repositories.UrlMappingRepository, urlClickService urlclick.UrlClickService) UrlMappingService {
	err := godotenv.Load()
	if err != nil {
		panic(".env file not found")
	}

	baseDomain = os.Getenv("BASE_DOMAIN")
	return &UrlMappingServiceImpl{
		repo:            repo,
		urlClickService: urlClickService,
	}
}

func (s *UrlMappingServiceImpl) Create(ctx *gin.Context, request *requests.ShortenUrlRequest) (*responses.ShortenUrlResponse, error) {
	existing, err := s.repo.FindByLongUrl(ctx, request.LongUrl)
	if err == nil {
		shortURL := baseDomain + "/" + existing.ShortCode
		return &responses.ShortenUrlResponse{ShortUrl: shortURL, ShortCode: existing.ShortCode}, nil
	}

	urlMapping := &entities.URLMapping{LongURL: request.LongUrl}
	urlMapping, err = s.repo.Save(ctx, urlMapping)
	if err != nil {
		return nil, err
	}

	shortCode := encodeID(urlMapping.ID)
	urlMapping.ShortCode = shortCode

	_, err = s.repo.Update(ctx, urlMapping)
	if err != nil {
		return nil, err
	}

	// Logging
	_, err = s.urlClickService.Create(ctx, urlMapping.ID)
	if err != nil {
		return nil, err
	}

	shortURL := baseDomain + "/" + shortCode
	return &responses.ShortenUrlResponse{ShortUrl: shortURL, ShortCode: shortCode}, nil
}

func (s *UrlMappingServiceImpl) GetByShortcode(ctx context.Context, shortCode string) (string, error) {
	url, err := s.repo.FindByShortCode(ctx, shortCode)
	if err != nil {
		return "", errors.New("URL not found")
	}

	return url.LongURL, nil
}
