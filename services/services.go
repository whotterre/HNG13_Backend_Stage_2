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
	GetCountryByName(name string) (*dto.GetCountryByNameResponse, error)
	GetAllCountries(region string, currency string, sort string) ([]dto.FilterCountriesResponse, error)
	DeleteCountryByName(name string) error
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

	// Generate summary image after successful refresh
	totalCount, _, err := s.countryRepository.GetStats()
	if err == nil {
		topCountries, err := s.countryRepository.GetTopCountriesByGDP(5)
		if err == nil {
			// Generate the image with current timestamp
			_ = utils.GenerateSummaryImage(int(totalCount), topCountries, time.Now(), "cache/summary.png")
		}
	}

	response := dto.RefreshCountriesResponse{
		Status: "Successfully refreshed countries",
	}
	return response, nil
}

func (s countryService) GetStats() (*dto.GetCountryStatsResponse, error) {
	countriesCount, lastRefreshedTime, err := s.countryRepository.GetStats()
	if err != nil {
		return nil, err
	}

	statistics := dto.GetCountryStatsResponse{
		TotalCountries:  int(countriesCount),
		LastRefreshedAt: lastRefreshedTime,
	}

	return &statistics, nil
}

func (s countryService) GetCountryByName(name string) (*dto.GetCountryByNameResponse, error) {
	// Normalize the name
	normalizedName := strings.ToLower(name)
	// Call the repo method
	country, err := s.countryRepository.GetCountryByName(normalizedName)
	if err != nil {
		return nil, errors.New("Country not found")
	}

	// Convert to DTO with ISO 8601 formatted timestamp
	response := &dto.GetCountryByNameResponse{
		ID:              country.ID,
		Name:            country.Name,
		Capital:         country.Capital,
		Region:          country.Region,
		Population:      country.Population,
		CurrencyCode:    country.CurrencyCode,
		ExchangeRate:    country.ExchangeRate,
		EstimatedGDP:    country.EstimatedGDP,
		FlagURL:         country.FlagURL,
		LastRefreshedAt: country.LastRefreshedAt.Format(time.RFC3339),
	}

	return response, nil
}

func (s countryService) DeleteCountryByName(name string) error {
	// Normalize the name
	normalizedName := strings.ToLower(name)
	// Call the repo method
	err := s.countryRepository.DeleteCountryByName(normalizedName)
	if err != nil {
		return errors.New("Failed to delete country")
	}

	return nil
}

func (s countryService) GetAllCountries(region string, currency string, sort string) ([]dto.FilterCountriesResponse, error) {
	countries, err := s.countryRepository.GetAllCountriesWithFilters(region, currency, sort)
	if err != nil {
		return nil, err
	}
	var res []dto.FilterCountriesResponse

	for _, country := range *countries {
		record := dto.FilterCountriesResponse{
			ID:              country.ID,
			Name:            country.Name,
			Capital:         country.Capital,
			Region:          country.Region,
			Population:      country.Population,
			CurrencyCode:    country.CurrencyCode,
			ExchangeRate:    country.ExchangeRate,
			EstimatedGDP:    country.EstimatedGDP,
			FlagURL:         country.FlagURL,
			LastRefreshedAt: country.LastRefreshedAt.Format(time.RFC3339),
		}

		res = append(res, record)
	}

	return res, nil
}
