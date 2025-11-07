package dto

// LanguageResponse defines the structure for language data sent to the client.
type LanguageResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
