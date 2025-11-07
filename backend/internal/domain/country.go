package domain

// Country represents the countries table
type Country struct {
	ID           uint                 `gorm:"primaryKey"`
	Code         string               `gorm:"uniqueIndex;size:2"` // ISO 3166-1 alpha-2 code
	Translations []CountryTranslation `gorm:"foreignKey:CountryCode;references:Code"`
}

// CountryTranslation represents the country_translations table
type CountryTranslation struct {
	ID           uint   `gorm:"primaryKey"`
	CountryCode  string `gorm:"index;size:2"`
	LanguageCode string `gorm:"index;size:2"` // ISO 639-1 code
	Name         string
}
