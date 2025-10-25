package routes

import (
	"task_2/handlers"
	"task_2/repository"
	"task_2/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	countryRepo := repository.NewCountryRepository(db)
	countryServices := services.NewCountryService(countryRepo, db)
	countryHandlers := handlers.NewCountryHandler(countryServices)

	router.POST("/countries/refresh", countryHandlers.RefreshCountries)
}