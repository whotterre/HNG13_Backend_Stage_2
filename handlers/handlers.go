package handlers

import (
	"log"
	"net/http"
	"strings"
	"task_2/services"

	"github.com/gin-gonic/gin"
)

type ValidationError struct {
	Message string
	Details map[string]string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type CountryHandler struct {
	countryServices services.CountryService
}

func NewCountryHandler(countryServices services.CountryService) *CountryHandler {
	return &CountryHandler{
		countryServices: countryServices,
	}
}

func (h CountryHandler) RefreshCountries(c *gin.Context) {
	// Call the service straight away!!!
	response, err := h.countryServices.RefreshCountries()
	if err != nil {
		handleError(err, c)
		return
	}
	log.Println(response)
	c.JSON(http.StatusOK, response)
}

func (h CountryHandler) GetStatistics(c *gin.Context) {
	stats, err := h.countryServices.GetStats()
	if err != nil {
		handleError(err, c)
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h CountryHandler) GetCountryByName(c *gin.Context) {
	countryName := c.Param("name")

	if countryName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No param passed",
		})
		return
	}

	countryData, err := h.countryServices.GetCountryByName(countryName)
	if err != nil {
		handleError(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"country": countryData,
	})
}

func (h CountryHandler) GetAllCountries(c *gin.Context) {
	region := c.Query("region")
	currency := c.Query("currency")
	sort := c.Query("sort")

	countries, err := h.countryServices.GetAllCountries(region, currency, sort)
	if err != nil {
		handleError(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"countries": countries})
}

func (h CountryHandler) DeleteCountry(c *gin.Context) {
	countryName := c.Param("name")

	if countryName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No param passed",
		})
		return
	}

	err := h.countryServices.DeleteCountryByName(countryName)
	if err != nil {
		handleError(err, c)
		return
	}

	c.Status(http.StatusNoContent)
}

func handleError(err error, c *gin.Context) error {
	errString := err.Error()

	if strings.Contains(errString, "Validation failed") {
		if valErr, ok := err.(*services.ValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": valErr.Details,
			})
			return err
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
		})
		return err
	}

	if strings.Contains(errString, "failed to fetch") {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "External data source unavailable",
			"details": errString,
		})
		return err
	}

	if strings.Contains(errString, "Country not found") {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Country not found",
		})
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "details": errString})
	return nil
}
