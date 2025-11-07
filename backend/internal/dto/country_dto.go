package dto

// CountryResponse defines the structure for country data sent to the client.
type CountryResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// CreateCountryRequest defines the structure for creating a new country.
type CreateCountryRequest struct {
	Code string `json:"code" validate:"required,min=2,max=2"`
	Name string `json:"name" validate:"required,min=3,max=100"`
}

// UpdateCountryRequest defines the structure for updating an existing country.
type UpdateCountryRequest struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
}
