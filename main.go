package main

import (
	urlClick "bitly/domains/urlclick/entities"
	urlMapping "bitly/domains/urlmapping/entities"
	"bitly/wizards"

	"github.com/gin-gonic/gin"
)

func main() {
	wizards.PostgresDatabase.GetInstance().AutoMigrate(
		&urlClick.URLClick{},
		&urlMapping.URLMapping{},
	)

	router := gin.Default()

	wizards.Routes(router)

	router.Run(":8080")
}
