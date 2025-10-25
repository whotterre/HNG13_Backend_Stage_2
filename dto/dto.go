package dto


type RefreshCountriesResponse struct {
	Status string `json:"status"`
}

type GetCountryStatsResponse struct {
	TotalCountries  int     `json:"total_countries"`
	LastRefreshedAt string `json:"last_refreshed_at"`
}

type GetCountryByNameResponse struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string    `gorm:"size:255;not null" json:"name"`
	Capital         string    `gorm:"size:255" json:"capital"`
	Region          string    `gorm:"size:255" json:"region"`
	Population      int64     `gorm:"not null" json:"population"`
	CurrencyCode    *string   `gorm:"size:10" json:"currency_code"`
	ExchangeRate    *float64  `json:"exchange_rate"`
	EstimatedGDP    *float64  `json:"estimated_gdp"`
	FlagURL         string    `gorm:"size:512" json:"flag_url"`
	LastRefreshedAt string `gorm:"autoUpdateTime" json:"last_refreshed_at"`
}