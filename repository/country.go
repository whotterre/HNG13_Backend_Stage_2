package repository

import (
	"strings"
	"task_2/models"
	"time"

	"gorm.io/gorm"
)

type countryRepository struct {
	db *gorm.DB
}

type CountryRepository interface {
	CreateNewCountry(country *models.Country) (*models.Country, error)
	GetCountryByName(countryName string) (*models.Country, error)
	UpdateCountry(countryId uint, updateData *models.Country) error
	DeleteCountryByName(countryName string) error
	GetAllCountries() (*[]models.Country, error)
	GetAllCountriesWithFilters(region string, currency string, sort string) (*[]models.Country, error)
	GetStats() (int64, string, error)
	GetTopCountriesByGDP(limit int) ([]models.Country, error)
}

func NewCountryRepository(db *gorm.DB) CountryRepository {
	return &countryRepository{
		db: db,
	}
}

func (r countryRepository) CreateNewCountry(country *models.Country) (*models.Country, error) {
	if err := r.db.Create(country).Error; err != nil {
		return nil, err
	}
	return country, nil
}

func (r countryRepository) GetCountryByName(countryName string) (*models.Country, error) {
	var country models.Country
	if err := r.db.Where("LOWER(name) = ?", strings.ToLower(countryName)).First(&country).Error; err != nil {
		return nil, err
	}
	return &country, nil
}

func (r countryRepository) GetAllCountries() (*[]models.Country, error) {
	var countries []models.Country
	if err := r.db.Find(&countries).Error; err != nil {
		return nil, err
	}
	return &countries, nil
}

func (r countryRepository) GetAllCountriesWithFilters(region string, currency string, sort string) (*[]models.Country, error) {
	var countries []models.Country

	q := r.db.Model(&models.Country{})

	if strings.TrimSpace(region) != "" {
		q = q.Where("region = ?", region)
	}

	if strings.TrimSpace(currency) != "" {
		q = q.Where("currency_code = ?", currency)
	}

	switch sort {
	case "gdp_desc":
		q = q.Order("estimated_gdp DESC")
	case "gdp_asc":
		q = q.Order("estimated_gdp ASC")
	default:
		q = q.Order("name ASC")
	}

	if err := q.Find(&countries).Error; err != nil {
		return nil, err
	}
	return &countries, nil
}

func (r countryRepository) UpdateCountry(countryId uint, updateData *models.Country) error {
	res := r.db.Model(&models.Country{}).Where("id = ?", countryId).Updates(updateData)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r countryRepository) DeleteCountryByName(countryName string) error {
	if err := r.db.Where("LOWER(name) = ?", strings.ToLower(countryName)).Delete(&models.Country{}).Error; err != nil {
		return err
	}
	return nil
}

func (r countryRepository) GetStats() (int64, string, error) {
	var count int64

	if err := r.db.Model(&models.Country{}).Count(&count).Error; err != nil {
		return 0, "", err
	}

	var result struct {
		LastRefreshedAt *time.Time
	}

	err := r.db.Model(&models.Country{}).
		Select("MAX(last_refreshed_at) as last_refreshed_at").
		Scan(&result).Error

	if err != nil {
		return count, "", err
	}

	// Format to ISO 8601 (RFC3339)
	var lastRefreshedStr string
	if result.LastRefreshedAt != nil {
		lastRefreshedStr = result.LastRefreshedAt.Format(time.RFC3339)
	} else {
		lastRefreshedStr = ""
	}

	return count, lastRefreshedStr, nil
}

func (r countryRepository) GetTopCountriesByGDP(limit int) ([]models.Country, error) {
	var countries []models.Country
	if err := r.db.Where("estimated_gdp IS NOT NULL").
		Order("estimated_gdp DESC").
		Limit(limit).
		Find(&countries).Error; err != nil {
		return nil, err
	}
	return countries, nil
}
