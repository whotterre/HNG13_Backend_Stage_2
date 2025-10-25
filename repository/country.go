package repository

import (
	"task_2/models"

	"gorm.io/gorm"
)

type countryRepository struct {
	db *gorm.DB
}

type CountryRepository interface {
	CreateNewCountry(country models.Country) (*models.Country, error)
	GetCountryByName(countryName string) (*models.Country, error)
	DeleteCountryByName(countryName string) error
}

func NewCountryRepository(db *gorm.DB) CountryRepository {
	return &countryRepository{
		db:db,
	}
}

func (r countryRepository) CreateNewCountry(country models.Country) (*models.Country, error){
	if err := r.db.Create(&country).Error; err != nil {
		return nil, err
	}
	return &country, nil 
}


func (r countryRepository) GetCountryByName(countryName string) (*models.Country, error){
	var country models.Country
	if err := r.db.Where("name = ?", countryName).First(&models.Country{}).Error; err != nil {
		return nil, err
	}
	return &country, nil
}

func (r countryRepository) DeleteCountryByName(countryName string) error {
	if err := r.db.Where("name = ?").Delete(&models.Country{}).Error; err != nil {
		return err
	}
	return nil
}

func (r countryRepository) GetStats(){
	if err := r.db.Exec("SELECT COUNT(*), last_refreshed_at FROM countries").Error; err != nil {

	}
	// return 
}

