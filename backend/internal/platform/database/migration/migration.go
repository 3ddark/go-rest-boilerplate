package migration

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"ths-erp.com/internal/domain"
)

func Migrate(db *gorm.DB) {
	dropTables(db)
	autoMigrate(db)
	seedData(db)
}

func dropTables(db *gorm.DB) {
	if err := db.Migrator().DropTable(
		&domain.User{},
		&domain.UserPermission{},
		&domain.Report{},
		&domain.Country{},
		&domain.CountryTranslation{},
		&domain.Language{},
		&domain.LanguageTranslation{},
	); err != nil {
		log.Fatalf("Could not drop tables: %v", err)
	}
	log.Println("✓ Tables dropped")
}

func autoMigrate(db *gorm.DB) {
	if err := db.AutoMigrate(
		&domain.User{},
		&domain.UserPermission{},
		&domain.Report{},
		&domain.Country{},
		&domain.CountryTranslation{},
		&domain.Language{},
		&domain.LanguageTranslation{},
		&domain.Unit{},
		&domain.UnitTranslation{},
	); err != nil {
		log.Fatalf("Could not migrate database: %v", err)
	}
	log.Println("✓ Database migrated")
}

func seedData(db *gorm.DB) {
	seedLanguages(db)
	seedCountries(db)
	seedUnits(db)
	seedUsers(db)
	log.Println("✓ Data seeded")
}

func seedLanguages(db *gorm.DB) {
	languages := []domain.Language{
		{Code: "en", IsActive: true},
		{Code: "tr", IsActive: true},
	}
	db.Create(&languages)
}

func seedCountries(db *gorm.DB) {
	countries := []domain.Country{
		{Code: "AD", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Andorra"}, {LanguageCode: "tr", Name: "Andorra"}}},
		{Code: "AE", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "United Arab Emirates"}, {LanguageCode: "tr", Name: "Birleşik Arap Emirlikleri"}}},
		{Code: "AF", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Afghanistan"}, {LanguageCode: "tr", Name: "Afganistan"}}},
		{Code: "AG", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Antigua and Barbuda"}, {LanguageCode: "tr", Name: "Antigua ve Barbuda"}}},
		{Code: "AI", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Anguilla"}, {LanguageCode: "tr", Name: "Anguilla"}}},
		{Code: "AL", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Albania"}, {LanguageCode: "tr", Name: "Arnavutluk"}}},
		{Code: "AM", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Armenia"}, {LanguageCode: "tr", Name: "Ermenistan"}}},
		{Code: "AO", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Angola"}, {LanguageCode: "tr", Name: "Angola"}}},
		{Code: "AQ", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Antarctica"}, {LanguageCode: "tr", Name: "Antarktika"}}},
		{Code: "AR", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Argentina"}, {LanguageCode: "tr", Name: "Arjantin"}}},
		{Code: "AS", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "American Samoa"}, {LanguageCode: "tr", Name: "Amerikan Samoası"}}},
		{Code: "AT", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Austria"}, {LanguageCode: "tr", Name: "Avusturya"}}},
		{Code: "AU", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Australia"}, {LanguageCode: "tr", Name: "Avustralya"}}},
		{Code: "AW", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Aruba"}, {LanguageCode: "tr", Name: "Aruba"}}},
		{Code: "AZ", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Azerbaijan"}, {LanguageCode: "tr", Name: "Azerbaycan"}}},
		{Code: "BA", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Bosnia and Herzegovina"}, {LanguageCode: "tr", Name: "Bosna-Hersek"}}},
		{Code: "BB", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Barbados"}, {LanguageCode: "tr", Name: "Barbados"}}},
		{Code: "BD", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Bangladesh"}, {LanguageCode: "tr", Name: "Bangladeş"}}},
		{Code: "BE", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Belgium"}, {LanguageCode: "tr", Name: "Belçika"}}},
		{Code: "BF", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Burkina Faso"}, {LanguageCode: "tr", Name: "Burkina Faso"}}},
		{Code: "BG", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Bulgaria"}, {LanguageCode: "tr", Name: "Bulgaristan"}}},
		{Code: "BH", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Bahrain"}, {LanguageCode: "tr", Name: "Bahreyn"}}},
		{Code: "BI", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Burundi"}, {LanguageCode: "tr", Name: "Burundi"}}},
		{Code: "BJ", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Benin"}, {LanguageCode: "tr", Name: "Benin"}}},
		{Code: "BL", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Saint Barthélemy"}, {LanguageCode: "tr", Name: "Saint Barthélemy"}}},
		{Code: "BM", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Bermuda"}, {LanguageCode: "tr", Name: "Bermuda"}}},
		{Code: "BN", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Brunei Darussalam"}, {LanguageCode: "tr", Name: "Brunei Darussalam"}}},
		{Code: "BO", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Bolivia (Plurinational State of)"}, {LanguageCode: "tr", Name: "Bolivya"}}},
		{Code: "BQ", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Bonaire, Sint Eustatius and Saba"}, {LanguageCode: "tr", Name: "Bonaire, Sint Eustatius ve Saba"}}},
		{Code: "BR", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Brazil"}, {LanguageCode: "tr", Name: "Brezilya"}}},
		{Code: "BS", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Bahamas (the)"}, {LanguageCode: "tr", Name: "Bahamalar"}}},
		{Code: "BT", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Bhutan"}, {LanguageCode: "tr", Name: "Butan"}}},
		{Code: "BV", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Bouvet Island"}, {LanguageCode: "tr", Name: "Bouvet Adası"}}},
		{Code: "BW", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Botswana"}, {LanguageCode: "tr", Name: "Botsvana"}}},
		{Code: "BY", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Belarus"}, {LanguageCode: "tr", Name: "Belarus"}}},
		{Code: "BZ", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Belize"}, {LanguageCode: "tr", Name: "Belize"}}},
		{Code: "CA", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Canada"}, {LanguageCode: "tr", Name: "Kanada"}}},
		{Code: "CC", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Cocos (Keeling) Islands"}, {LanguageCode: "tr", Name: "Cocos (Keeling) Adaları"}}},
		{Code: "CD", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Congo (the Democratic Republic of the)"}, {LanguageCode: "tr", Name: "Kongo Demokratik Cumhuriyeti"}}},
		{Code: "CF", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Central African Republic (the)"}, {LanguageCode: "tr", Name: "Orta Afrika Cumhuriyeti"}}},
		{Code: "CG", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Congo (the)"}, {LanguageCode: "tr", Name: "Kongo"}}},
		{Code: "CH", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Switzerland"}, {LanguageCode: "tr", Name: "İsviçre"}}},
		{Code: "CI", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Côte d'Ivoire"}, {LanguageCode: "tr", Name: "Fildişi Sahili"}}},
		{Code: "CK", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Cook Islands (the)"}, {LanguageCode: "tr", Name: "Cook Adaları"}}},
		{Code: "CL", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Chile"}, {LanguageCode: "tr", Name: "Şili"}}},
		{Code: "CM", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Cameroon"}, {LanguageCode: "tr", Name: "Kamerun"}}},
		{Code: "CN", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "China"}, {LanguageCode: "tr", Name: "Çin"}}},
		{Code: "CO", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Colombia"}, {LanguageCode: "tr", Name: "Kolombiya"}}},
		{Code: "CR", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Costa Rica"}, {LanguageCode: "tr", Name: "Kosta Rika"}}},
		{Code: "CU", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Cuba"}, {LanguageCode: "tr", Name: "Küba"}}},
		{Code: "CV", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Cabo Verde"}, {LanguageCode: "tr", Name: "Cabo Verde"}}},
		{Code: "CW", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Curaçao"}, {LanguageCode: "tr", Name: "Curaçao"}}},
		{Code: "CX", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Christmas Island"}, {LanguageCode: "tr", Name: "Christmas Adası"}}},
		{Code: "CY", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Cyprus"}, {LanguageCode: "tr", Name: "Kıbrıs"}}},
		{Code: "CZ", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Czechia"}, {LanguageCode: "tr", Name: "Çekya"}}},
		{Code: "DE", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Germany"}, {LanguageCode: "tr", Name: "Almanya"}}},
		{Code: "DJ", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Djibouti"}, {LanguageCode: "tr", Name: "Cibuti"}}},
		{Code: "DK", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Denmark"}, {LanguageCode: "tr", Name: "Danimarka"}}},
		{Code: "DM", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Dominica"}, {LanguageCode: "tr", Name: "Dominika"}}},
		{Code: "DO", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Dominican Republic (the)"}, {LanguageCode: "tr", Name: "Dominik Cumhuriyeti"}}},
		{Code: "DZ", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Algeria"}, {LanguageCode: "tr", Name: "Cezayir"}}},
		{Code: "EC", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Ecuador"}, {LanguageCode: "tr", Name: "Ekvador"}}},
		{Code: "EE", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Estonia"}, {LanguageCode: "tr", Name: "Estonya"}}},
		{Code: "EG", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Egypt"}, {LanguageCode: "tr", Name: "Mısır"}}},
		{Code: "EH", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Western Sahara"}, {LanguageCode: "tr", Name: "Batı Sahra"}}},
		{Code: "ER", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Eritrea"}, {LanguageCode: "tr", Name: "Eritre"}}},
		{Code: "ES", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Spain"}, {LanguageCode: "tr", Name: "İspanya"}}},
		{Code: "ET", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Ethiopia"}, {LanguageCode: "tr", Name: "Etiyopya"}}},
		{Code: "FI", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Finland"}, {LanguageCode: "tr", Name: "Finlandiya"}}},
		{Code: "FJ", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Fiji"}, {LanguageCode: "tr", Name: "Fiji"}}},
		{Code: "FK", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Falkland Islands (the)"}, {LanguageCode: "tr", Name: "Falkland Adaları"}}},
		{Code: "FM", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Micronesia (Federated States of)"}, {LanguageCode: "tr", Name: "Mikronezya"}}},
		{Code: "FO", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Faroe Islands (the)"}, {LanguageCode: "tr", Name: "Faroe Adaları"}}},
		{Code: "FR", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "France"}, {LanguageCode: "tr", Name: "Fransa"}}},
		{Code: "GA", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Gabon"}, {LanguageCode: "tr", Name: "Gabon"}}},
		{Code: "GB", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "United Kingdom of Great Britain and Northern Ireland"}, {LanguageCode: "tr", Name: "Birleşik Krallık"}}},
		{Code: "GD", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Grenada"}, {LanguageCode: "tr", Name: "Grenada"}}},
		{Code: "GE", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Georgia"}, {LanguageCode: "tr", Name: "Gürcistan"}}},
		{Code: "GF", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "French Guiana"}, {LanguageCode: "tr", Name: "Fransız Guyanası"}}},
		{Code: "GG", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Guernsey"}, {LanguageCode: "tr", Name: "Guernsey"}}},
		{Code: "GH", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Ghana"}, {LanguageCode: "tr", Name: "Gana"}}},
		{Code: "GI", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Gibraltar"}, {LanguageCode: "tr", Name: "Cebelitarık"}}},
		{Code: "GL", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Greenland"}, {LanguageCode: "tr", Name: "Grönland"}}},
		{Code: "GM", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Gambia (the)"}, {LanguageCode: "tr", Name: "Gambiya"}}},
		{Code: "GN", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Guinea"}, {LanguageCode: "tr", Name: "Gine"}}},
		{Code: "GP", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Guadeloupe"}, {LanguageCode: "tr", Name: "Guadeloupe"}}},
		{Code: "GQ", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Equatorial Guinea"}, {LanguageCode: "tr", Name: "Ekvator Ginesi"}}},
		{Code: "GR", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Greece"}, {LanguageCode: "tr", Name: "Yunanistan"}}},
		{Code: "GS", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "South Georgia and the South Sandwich Islands"}, {LanguageCode: "tr", Name: "Güney Georgia ve Güney Sandwich Adaları"}}},
		{Code: "GT", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Guatemala"}, {LanguageCode: "tr", Name: "Guatemala"}}},
		{Code: "GU", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Guam"}, {LanguageCode: "tr", Name: "Guam"}}},
		{Code: "GW", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Guinea-Bissau"}, {LanguageCode: "tr", Name: "Gine-Bissau"}}},
		{Code: "GY", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Guyana"}, {LanguageCode: "tr", Name: "Guyana"}}},
		{Code: "HK", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Hong Kong"}, {LanguageCode: "tr", Name: "Hong Kong"}}},
		{Code: "HM", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Heard Island and McDonald Islands"}, {LanguageCode: "tr", Name: "Heard Adası ve McDonald Adaları"}}},
		{Code: "HN", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Honduras"}, {LanguageCode: "tr", Name: "Honduras"}}},
		{Code: "HR", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Croatia"}, {LanguageCode: "tr", Name: "Hırvatistan"}}},
		{Code: "HT", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Haiti"}, {LanguageCode: "tr", Name: "Haiti"}}},
		{Code: "HU", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Hungary"}, {LanguageCode: "tr", Name: "Macaristan"}}},
		{Code: "ID", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Indonesia"}, {LanguageCode: "tr", Name: "Endonezya"}}},
		{Code: "IE", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Ireland"}, {LanguageCode: "tr", Name: "İrlanda"}}},
		{Code: "IL", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Israel"}, {LanguageCode: "tr", Name: "İsrail"}}},
		{Code: "IM", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Isle of Man"}, {LanguageCode: "tr", Name: "Man Adası"}}},
		{Code: "IN", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "India"}, {LanguageCode: "tr", Name: "Hindistan"}}},
		{Code: "IO", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "British Indian Ocean Territory (the)"}, {LanguageCode: "tr", Name: "Britanya Hint Okyanusu Toprakları"}}},
		{Code: "IQ", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Iraq"}, {LanguageCode: "tr", Name: "Irak"}}},
		{Code: "IR", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Iran (Islamic Republic of)"}, {LanguageCode: "tr", Name: "İran"}}},
		{Code: "IS", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Iceland"}, {LanguageCode: "tr", Name: "İzlanda"}}},
		{Code: "IT", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Italy"}, {LanguageCode: "tr", Name: "İtalya"}}},
		{Code: "JE", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Jersey"}, {LanguageCode: "tr", Name: "Jersey"}}},
		{Code: "JM", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Jamaica"}, {LanguageCode: "tr", Name: "Jamaika"}}},
		{Code: "JO", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Jordan"}, {LanguageCode: "tr", Name: "Ürdün"}}},
		{Code: "JP", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Japan"}, {LanguageCode: "tr", Name: "Japonya"}}},
		{Code: "KE", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Kenya"}, {LanguageCode: "tr", Name: "Kenya"}}},
		{Code: "KG", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Kyrgyzstan"}, {LanguageCode: "tr", Name: "Kırgızistan"}}},
		{Code: "KH", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Cambodia"}, {LanguageCode: "tr", Name: "Kamboçya"}}},
		{Code: "KI", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Kiribati"}, {LanguageCode: "tr", Name: "Kiribati"}}},
		{Code: "KM", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Comoros (the)"}, {LanguageCode: "tr", Name: "Komorlar"}}},
		{Code: "KN", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Saint Kitts and Nevis"}, {LanguageCode: "tr", Name: "Saint Kitts ve Nevis"}}},
		{Code: "KP", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Korea (the Democratic People's Republic of)"}, {LanguageCode: "tr", Name: "Kuzey Kore"}}},
		{Code: "KR", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Korea (the Republic of)"}, {LanguageCode: "tr", Name: "Güney Kore"}}},
		{Code: "KW", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Kuwait"}, {LanguageCode: "tr", Name: "Kuveyt"}}},
		{Code: "KY", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Cayman Islands (the)"}, {LanguageCode: "tr", Name: "Cayman Adaları"}}},
		{Code: "KZ", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Kazakhstan"}, {LanguageCode: "tr", Name: "Kazakistan"}}},
		{Code: "LA", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Lao People's Democratic Republic (the)"}, {LanguageCode: "tr", Name: "Laos"}}},
		{Code: "LB", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Lebanon"}, {LanguageCode: "tr", Name: "Lübnan"}}},
		{Code: "LC", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Saint Lucia"}, {LanguageCode: "tr", Name: "Saint Lucia"}}},
		{Code: "LI", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Liechtenstein"}, {LanguageCode: "tr", Name: "Lihtenştayn"}}},
		{Code: "LK", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Sri Lanka"}, {LanguageCode: "tr", Name: "Sri Lanka"}}},
		{Code: "LR", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Liberia"}, {LanguageCode: "tr", Name: "Liberya"}}},
		{Code: "LS", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Lesotho"}, {LanguageCode: "tr", Name: "Lesotho"}}},
		{Code: "LT", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Lithuania"}, {LanguageCode: "tr", Name: "Litvanya"}}},
		{Code: "LU", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Luxembourg"}, {LanguageCode: "tr", Name: "Lüksemburg"}}},
		{Code: "LV", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Latvia"}, {LanguageCode: "tr", Name: "Letonya"}}},
		{Code: "LY", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Libya"}, {LanguageCode: "tr", Name: "Libya"}}},
		{Code: "MA", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Morocco"}, {LanguageCode: "tr", Name: "Fas"}}},
		{Code: "MC", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Monaco"}, {LanguageCode: "tr", Name: "Monako"}}},
		{Code: "MD", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Moldova (the Republic of)"}, {LanguageCode: "tr", Name: "Moldova"}}},
		{Code: "ME", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Montenegro"}, {LanguageCode: "tr", Name: "Karadağ"}}},
		{Code: "MF", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Saint Martin (French part)"}, {LanguageCode: "tr", Name: "Saint Martin (Fransız kısmı)"}}},
		{Code: "MG", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Madagascar"}, {LanguageCode: "tr", Name: "Madagaskar"}}},
		{Code: "MH", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Marshall Islands (the)"}, {LanguageCode: "tr", Name: "Marshall Adaları"}}},
		{Code: "MK", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "North Macedonia"}, {LanguageCode: "tr", Name: "Kuzey Makedonya"}}},
		{Code: "ML", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Mali"}, {LanguageCode: "tr", Name: "Mali"}}},
		{Code: "MM", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Myanmar"}, {LanguageCode: "tr", Name: "Myanmar"}}},
		{Code: "MN", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Mongolia"}, {LanguageCode: "tr", Name: "Moğolistan"}}},
		{Code: "MO", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Macao"}, {LanguageCode: "tr", Name: "Makao"}}},
		{Code: "MP", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Northern Mariana Islands (the)"}, {LanguageCode: "tr", Name: "Kuzey Mariana Adaları"}}},
		{Code: "MQ", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Martinique"}, {LanguageCode: "tr", Name: "Martinik"}}},
		{Code: "MR", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Mauritania"}, {LanguageCode: "tr", Name: "Moritanya"}}},
		{Code: "MS", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Montserrat"}, {LanguageCode: "tr", Name: "Montserrat"}}},
		{Code: "MT", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Malta"}, {LanguageCode: "tr", Name: "Malta"}}},
		{Code: "MU", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Mauritius"}, {LanguageCode: "tr", Name: "Mauritius"}}},
		{Code: "MV", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Maldives"}, {LanguageCode: "tr", Name: "Maldivler"}}},
		{Code: "MW", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Malawi"}, {LanguageCode: "tr", Name: "Malavi"}}},
		{Code: "MX", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Mexico"}, {LanguageCode: "tr", Name: "Meksika"}}},
		{Code: "MY", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Malaysia"}, {LanguageCode: "tr", Name: "Malezya"}}},
		{Code: "MZ", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Mozambique"}, {LanguageCode: "tr", Name: "Mozambik"}}},
		{Code: "NA", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Namibia"}, {LanguageCode: "tr", Name: "Namibya"}}},
		{Code: "NC", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "New Caledonia"}, {LanguageCode: "tr", Name: "Yeni Kaledonya"}}},
		{Code: "NE", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Niger (the)"}, {LanguageCode: "tr", Name: "Nijer"}}},
		{Code: "NF", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Norfolk Island"}, {LanguageCode: "tr", Name: "Norfolk Adası"}}},
		{Code: "NG", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Nigeria"}, {LanguageCode: "tr", Name: "Nijerya"}}},
		{Code: "NI", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Nicaragua"}, {LanguageCode: "tr", Name: "Nikaragua"}}},
		{Code: "NL", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Netherlands (the)"}, {LanguageCode: "tr", Name: "Hollanda"}}},
		{Code: "NO", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Norway"}, {LanguageCode: "tr", Name: "Norveç"}}},
		{Code: "NP", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Nepal"}, {LanguageCode: "tr", Name: "Nepal"}}},
		{Code: "NR", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Nauru"}, {LanguageCode: "tr", Name: "Nauru"}}},
		{Code: "NU", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Niue"}, {LanguageCode: "tr", Name: "Niue"}}},
		{Code: "NZ", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "New Zealand"}, {LanguageCode: "tr", Name: "Yeni Zelanda"}}},
		{Code: "OM", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Oman"}, {LanguageCode: "tr", Name: "Umman"}}},
		{Code: "PA", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Panama"}, {LanguageCode: "tr", Name: "Panama"}}},
		{Code: "PE", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Peru"}, {LanguageCode: "tr", Name: "Peru"}}},
		{Code: "PF", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "French Polynesia"}, {LanguageCode: "tr", Name: "Fransız Polinezyası"}}},
		{Code: "PG", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Papua New Guinea"}, {LanguageCode: "tr", Name: "Papua Yeni Gine"}}},
		{Code: "PH", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Philippines (the)"}, {LanguageCode: "tr", Name: "Filipinler"}}},
		{Code: "PK", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Pakistan"}, {LanguageCode: "tr", Name: "Pakistan"}}},
		{Code: "PL", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Poland"}, {LanguageCode: "tr", Name: "Polonya"}}},
		{Code: "PM", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Saint Pierre and Miquelon"}, {LanguageCode: "tr", Name: "Saint Pierre ve Miquelon"}}},
		{Code: "PN", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Pitcairn"}, {LanguageCode: "tr", Name: "Pitcairn"}}},
		{Code: "PR", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Puerto Rico"}, {LanguageCode: "tr", Name: "Porto Riko"}}},
		{Code: "PS", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Palestine, State of"}, {LanguageCode: "tr", Name: "Filistin"}}},
		{Code: "PT", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Portugal"}, {LanguageCode: "tr", Name: "Portekiz"}}},
		{Code: "PW", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Palau"}, {LanguageCode: "tr", Name: "Palau"}}},
		{Code: "PY", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Paraguay"}, {LanguageCode: "tr", Name: "Paraguay"}}},
		{Code: "QA", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Qatar"}, {LanguageCode: "tr", Name: "Katar"}}},
		{Code: "RE", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Réunion"}, {LanguageCode: "tr", Name: "Réunion"}}},
		{Code: "RO", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Romania"}, {LanguageCode: "tr", Name: "Romanya"}}},
		{Code: "RS", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Serbia"}, {LanguageCode: "tr", Name: "Sırbistan"}}},
		{Code: "RU", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Russian Federation (the)"}, {LanguageCode: "tr", Name: "Rusya"}}},
		{Code: "RW", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Rwanda"}, {LanguageCode: "tr", Name: "Ruanda"}}},
		{Code: "SA", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Saudi Arabia"}, {LanguageCode: "tr", Name: "Suudi Arabistan"}}},
		{Code: "SB", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Solomon Islands"}, {LanguageCode: "tr", Name: "Solomon Adaları"}}},
		{Code: "SC", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Seychelles"}, {LanguageCode: "tr", Name: "Seyşeller"}}},
		{Code: "SD", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Sudan (the)"}, {LanguageCode: "tr", Name: "Sudan"}}},
		{Code: "SE", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Sweden"}, {LanguageCode: "tr", Name: "İsveç"}}},
		{Code: "SG", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Singapore"}, {LanguageCode: "tr", Name: "Singapur"}}},
		{Code: "SH", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Saint Helena, Ascension and Tristan da Cunha"}, {LanguageCode: "tr", Name: "Saint Helena, Ascension ve Tristan da Cunha"}}},
		{Code: "SI", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Slovenia"}, {LanguageCode: "tr", Name: "Slovenya"}}},
		{Code: "SJ", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Svalbard and Jan Mayen"}, {LanguageCode: "tr", Name: "Svalbard ve Jan Mayen"}}},
		{Code: "SK", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Slovakia"}, {LanguageCode: "tr", Name: "Slovakya"}}},
		{Code: "SL", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Sierra Leone"}, {LanguageCode: "tr", Name: "Sierra Leone"}}},
		{Code: "SM", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "San Marino"}, {LanguageCode: "tr", Name: "San Marino"}}},
		{Code: "SN", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Senegal"}, {LanguageCode: "tr", Name: "Senegal"}}},
		{Code: "SO", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Somalia"}, {LanguageCode: "tr", Name: "Somali"}}},
		{Code: "SR", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Suriname"}, {LanguageCode: "tr", Name: "Surinam"}}},
		{Code: "SS", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "South Sudan"}, {LanguageCode: "tr", Name: "Güney Sudan"}}},
		{Code: "ST", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Sao Tome and Principe"}, {LanguageCode: "tr", Name: "Sao Tome ve Principe"}}},
		{Code: "SV", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "El Salvador"}, {LanguageCode: "tr", Name: "El Salvador"}}},
		{Code: "SX", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Sint Maarten (Dutch part)"}, {LanguageCode: "tr", Name: "Sint Maarten (Hollanda kısmı)"}}},
		{Code: "SY", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Syrian Arab Republic"}, {LanguageCode: "tr", Name: "Suriye"}}},
		{Code: "SZ", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Eswatini"}, {LanguageCode: "tr", Name: "Esvatini"}}},
		{Code: "TC", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Turks and Caicos Islands (the)"}, {LanguageCode: "tr", Name: "Turks ve Caicos Adaları"}}},
		{Code: "TD", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Chad"}, {LanguageCode: "tr", Name: "Çad"}}},
		{Code: "TF", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "French Southern Territories (the)"}, {LanguageCode: "tr", Name: "Fransız Güney Toprakları"}}},
		{Code: "TG", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Togo"}, {LanguageCode: "tr", Name: "Togo"}}},
		{Code: "TH", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Thailand"}, {LanguageCode: "tr", Name: "Tayland"}}},
		{Code: "TJ", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Tajikistan"}, {LanguageCode: "tr", Name: "Tacikistan"}}},
		{Code: "TK", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Tokelau"}, {LanguageCode: "tr", Name: "Tokelau"}}},
		{Code: "TL", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Timor-Leste"}, {LanguageCode: "tr", Name: "Timor-Leste"}}},
		{Code: "TM", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Turkmenistan"}, {LanguageCode: "tr", Name: "Türkmenistan"}}},
		{Code: "TN", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Tunisia"}, {LanguageCode: "tr", Name: "Tunus"}}},
		{Code: "TO", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Tonga"}, {LanguageCode: "tr", Name: "Tonga"}}},
		{Code: "TR", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Turkey"}, {LanguageCode: "tr", Name: "Türkiye"}}},
		{Code: "TT", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Trinidad and Tobago"}, {LanguageCode: "tr", Name: "Trinidad ve Tobago"}}},
		{Code: "TV", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Tuvalu"}, {LanguageCode: "tr", Name: "Tuvalu"}}},
		{Code: "TW", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Taiwan (Province of China)"}, {LanguageCode: "tr", Name: "Tayvan"}}},
		{Code: "TZ", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Tanzania, United Republic of"}, {LanguageCode: "tr", Name: "Tanzanya"}}},
		{Code: "UA", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Ukraine"}, {LanguageCode: "tr", Name: "Ukrayna"}}},
		{Code: "UG", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Uganda"}, {LanguageCode: "tr", Name: "Uganda"}}},
		{Code: "UM", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "United States Minor Outlying Islands (the)"}, {LanguageCode: "tr", Name: "Amerika Birleşik Devletleri'nin Küçük Dış Adaları"}}},
		{Code: "US", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "United States of America (the)"}, {LanguageCode: "tr", Name: "Amerika Birleşik Devletleri"}}},
		{Code: "UY", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Uruguay"}, {LanguageCode: "tr", Name: "Uruguay"}}},
		{Code: "UZ", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Uzbekistan"}, {LanguageCode: "tr", Name: "Özbekistan"}}},
		{Code: "VA", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Holy See (the)"}, {LanguageCode: "tr", Name: "Vatikan"}}},
		{Code: "VC", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Saint Vincent and the Grenadines"}, {LanguageCode: "tr", Name: "Saint Vincent ve Grenadinler"}}},
		{Code: "VE", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Venezuela (Bolivarian Republic of)"}, {LanguageCode: "tr", Name: "Venezuela"}}},
		{Code: "VG", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Virgin Islands (British)"}, {LanguageCode: "tr", Name: "Britanya Virjin Adaları"}}},
		{Code: "VI", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Virgin Islands (U.S.)"}, {LanguageCode: "tr", Name: "ABD Virjin Adaları"}}},
		{Code: "VN", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Viet Nam"}, {LanguageCode: "tr", Name: "Vietnam"}}},
		{Code: "VU", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Vanuatu"}, {LanguageCode: "tr", Name: "Vanuatu"}}},
		{Code: "WF", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Wallis and Futuna"}, {LanguageCode: "tr", Name: "Wallis ve Futuna"}}},
		{Code: "WS", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Samoa"}, {LanguageCode: "tr", Name: "Samoa"}}},
		{Code: "YE", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Yemen"}, {LanguageCode: "tr", Name: "Yemen"}}},
		{Code: "YT", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Mayotte"}, {LanguageCode: "tr", Name: "Mayotte"}}},
		{Code: "ZA", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "South Africa"}, {LanguageCode: "tr", Name: "Güney Afrika"}}},
		{Code: "ZM", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Zambia"}, {LanguageCode: "tr", Name: "Zambiya"}}},
		{Code: "ZW", Translations: []domain.CountryTranslation{{LanguageCode: "en", Name: "Zimbabwe"}, {LanguageCode: "tr", Name: "Zimbabve"}}},
	}
	if err := db.Create(&countries).Error; err != nil {
		log.Fatalf("Could not seed countries: %v", err)
	}
}

func seedUnits(db *gorm.DB) {
	units := []domain.Unit{
		{
			Code: "C62",
			Translations: []domain.UnitTranslation{
				{LanguageCode: "en", Name: "Piece"},
				{LanguageCode: "tr", Name: "Adet"},
			},
		},
		{
			Code: "KGM",
			Translations: []domain.UnitTranslation{
				{LanguageCode: "en", Name: "Kilogram"},
				{LanguageCode: "tr", Name: "Kilogram"},
			},
		},
		{
			Code: "LTR",
			Translations: []domain.UnitTranslation{
				{LanguageCode: "en", Name: "Liter"},
				{LanguageCode: "tr", Name: "Litre"},
			},
		},
		{
			Code: "MTR",
			Translations: []domain.UnitTranslation{
				{LanguageCode: "en", Name: "Meter"},
				{LanguageCode: "tr", Name: "Metre"},
			},
		},
		{
			Code: "DAY",
			Translations: []domain.UnitTranslation{
				{LanguageCode: "en", Name: "Day"},
				{LanguageCode: "tr", Name: "Gün"},
			},
		},
		{
			Code: "HUR",
			Translations: []domain.UnitTranslation{
				{LanguageCode: "en", Name: "Hour"},
				{LanguageCode: "tr", Name: "Saat"},
			},
		},
		{
			Code: "MTK",
			Translations: []domain.UnitTranslation{
				{LanguageCode: "en", Name: "Square Meter"},
				{LanguageCode: "tr", Name: "Metrekare"},
			},
		},
	}
	if err := db.Create(&units).Error; err != nil {
		log.Printf("could not seed units: %v", err)
	}
}

func seedUsers(db *gorm.DB) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	users := []domain.User{
		{
			Name:         "Admin User",
			Email:        "admin@example.com",
			PasswordHash: string(hashedPassword),
		},
		{
			Name:         "Test User",
			Email:        "test@example.com",
			PasswordHash: string(hashedPassword),
		},
	}
	db.Create(&users)
}
