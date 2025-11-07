package domain

// Language represents the languages table
type Language struct {
	ID           uint                  `gorm:"primaryKey"`
	Code         string                `gorm:"uniqueIndex;size:10"` // e.g., en, en-US, tr
	IsActive     bool                  `gorm:"default:true"`
	Translations []LanguageTranslation `gorm:"foreignKey:LanguageCode;references:Code"`
}

// LanguageTranslation represents the language_translations table
type LanguageTranslation struct {
	ID                      uint   `gorm:"primaryKey"`
	LanguageCode            string `gorm:"index;size:10"`
	TranslationLanguageCode string `gorm:"index;size:10"` // The language of the translation itself
	Name                    string `gorm:"size:50"`        // e.g., English, İngilizce
}
