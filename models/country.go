package models

import "time"

// Country represents a country and some economic metadata.
type Country struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string    `gorm:"size:255;not null" json:"name"`
	Capital         string    `gorm:"size:255" json:"capital"`
	Region          string    `gorm:"size:255" json:"region"`
	Population      int64     `gorm:"not null" json:"population"`
	CurrencyCode    *string   `gorm:"size:10" json:"country_code,omitempty"`
	ExchangeRate    *float64  `json:"exchange_rate,omitempty"`
	EstimatedGDP    *float64  `json:"estimated_gdp,omitempty"`
	FlagURL         string    `gorm:"size:512" json:"flag_url"`
	LastRefreshedAt time.Time `gorm:"autoUpdateTime" json:"last_refreshed_at"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
