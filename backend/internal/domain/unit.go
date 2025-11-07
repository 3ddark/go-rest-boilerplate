package domain

// Unit represents a unit of measurement.
type Unit struct {
	ID           uint              `gorm:"primaryKey"`
	Code         string            `gorm:"uniqueIndex;size:10"` // ISO or other standard code, e.g., "KGM", "C62"
	Translations []UnitTranslation `gorm:"foreignKey:UnitCode;references:Code"`
}

// UnitTranslation stores the language-specific names for a unit of measurement.
type UnitTranslation struct {
	ID           uint   `gorm:"primaryKey"`
	UnitCode     string `gorm:"index;size:10"`
	LanguageCode string `gorm:"index;size:2"` // ISO 639-1 language code
	Name         string `gorm:"size:100"`
}
