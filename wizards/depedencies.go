package wizards

import (
	urlClickRepository "bitly/domains/urlclick/repositories"
	urlClickService "bitly/domains/urlclick/services"
	urlMappingController "bitly/domains/urlmapping/controllers"
	urlMappingRepo "bitly/domains/urlmapping/repositories"
	urlMappingService "bitly/domains/urlmapping/services"
	"bitly/infrastructures"
)

var (
	PostgresDatabase     = infrastructures.NewPostgresDatabase()
	UrlMappingRepo       = urlMappingRepo.NewUrlMappingRepositoryImpl(PostgresDatabase)
	UrlMappingService    = urlMappingService.NewUrlMappingServiceImpl(UrlMappingRepo, UrlClickService)
	UrlMappingController = urlMappingController.NewUrlMappingController(UrlMappingService)
	UrlClickRepo         = urlClickRepository.NewUrlClickRepositoryImpl(PostgresDatabase)
	UrlClickService      = urlClickService.NewUrlMappingServiceImpl(UrlClickRepo)
)
