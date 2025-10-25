package dto


type RefreshCountriesResponse struct {
	Status string `json:"status"`
}

type GetCountryStatsResponse struct {
	TotalCountries  int     `json:"total_countries"`
	LastRefreshedAt string `json:"last_refreshed_at"`
}
