package services

import (
	"errors"
	"strings"
	"task_2/clients"
	"task_2/dto"
	"task_2/models"
	"task_2/repository"
	"task_2/utils"
	"time"

	"gorm.io/gorm"
)

type ValidationError struct {
	Message string
	Details map[string]string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type CountryService interface {
	RefreshCountries() (dto.RefreshCountriesResponse, error)
	GetStats() (*dto.GetCountryStatsResponse, error)
	GetCountryByName(name string) (*models.Country, error)
}	

type countryService struct {
	countryRepository repository.CountryRepository
	db                *gorm.DB
}

func NewCountryService(countryRepo repository.CountryRepository, db *gorm.DB) CountryService {
	return &countryService{
		countryRepository: countryRepo,
		db:                db,
	}
}

// Call the client to get the list of countries
func (s countryService) RefreshCountries() (dto.RefreshCountriesResponse, error) {
	countries, err := clients.GetCountries()
	if err != nil {
		return dto.RefreshCountriesResponse{}, errors.New("failed to fetch country data from external API")
	}

	rates, err := clients.GetExchangeRates()
	if err != nil {
		return dto.RefreshCountriesResponse{}, errors.New("failed to fetch exchange rates from external API")
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		for _, country := range *countries {
			normalizedName := strings.ToLower(country.Name)

			var ct models.Country
			findErr := tx.Session(&gorm.Session{Logger: tx.Logger.LogMode(4)}).Where("LOWER(name) = ?", normalizedName).First(&ct).Error

			// Prepare values to persist
			var currencyPtr *string
			var ratePtr *float64
			var estimatedPtr *float64

			if len(country.Currencies) == 0 {
				currencyPtr = nil
				ratePtr = nil
				zero := float64(0)
				estimatedPtr = &zero
			} else {
				currencyCode := country.Currencies[0].Code
				currencyPtr = &currencyCode
				if r, ok := rates.Rates[currencyCode]; ok && r > 0 {
					rateValue := r
					ratePtr = &rateValue
					ev := utils.ComputeEstimatedGDP(country.Population, rateValue)
					estimatedPtr = &ev
				} else {
					// rate not found
					ratePtr = nil
					estimatedPtr = nil
				}
			}

			// Build record to insert/update
			record := models.Country{
				Name:            country.Name,
				Capital:         country.Capital,
				Region:          country.Region,
				Population:      country.Population,
				CurrencyCode:    currencyPtr,
				ExchangeRate:    ratePtr,
				EstimatedGDP:    estimatedPtr,
				FlagURL:         country.FlagURL,
				LastRefreshedAt: now,
			}

			if findErr != nil {
				if errors.Is(findErr, gorm.ErrRecordNotFound) {
					validationDetails := make(map[string]string)
					if strings.TrimSpace(country.Name) == "" {
						validationDetails["name"] = "is required"
					}
					if country.Population < 0 {
						validationDetails["population"] = "must be non-negative"
					}
					if len(country.Currencies) > 0 && (currencyPtr == nil || strings.TrimSpace(*currencyPtr) == "") {
						validationDetails["currency_code"] = "is required"
					}

					if len(validationDetails) > 0 {
						return &ValidationError{
							Message: "Validation failed",
							Details: validationDetails,
						}
					}

					if err := tx.Create(&record).Error; err != nil {
						return err
					}
					continue
				}
				return findErr
			}

			if err := tx.Model(&ct).Updates(record).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return dto.RefreshCountriesResponse{}, err
	}

	response := dto.RefreshCountriesResponse{
		Status: "Successfully refreshed countries",
	}
	return response, nil
}


func (s countryService) GetStats() (*dto.GetCountryStatsResponse, error){
	countriesCount, lastRefreshedTime, err := s.countryRepository.GetStats()
	if err != nil {
		return nil, err
	}

	statistics := dto.GetCountryStatsResponse{
		TotalCountries: int(countriesCount),
		LastRefreshedAt: lastRefreshedTime,
	}

	return &statistics, nil
}

func (s countryService) GetCountryByName(name string) (*models.Country, error){
	// Normalize the name
	normalizedName := strings.ToLower(name)
	// Call the repo method
	country, err := s.countryRepository.GetCountryByName(normalizedName)
	if err != nil {
		return nil, errors.New("Country not found") 
	}

	return country, nil
}