
package dto

// Setup2FAResponse represents the data returned when setting up 2FA.
type Setup2FAResponse struct {
	Secret   string   `json:"secret"`
	QRCode   string   `json:"qrCode"`
	Recovery []string `json:"recovery"`
}

// Enable2FARequest represents the request to enable 2FA.
type Enable2FARequest struct {
	Code string `json:"code"`
}

// Login2FARequest represents the request to verify 2FA during login.
type Login2FARequest struct {
	UserID int    `json:"userId"`
	Code   string `json:"code"`
}
