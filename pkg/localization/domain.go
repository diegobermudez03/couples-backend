package localization

type LocalizationService interface {
	ValidateCountry(code string) error
	ValidateLanguage(code string) error
}