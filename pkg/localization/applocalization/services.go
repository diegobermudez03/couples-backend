package applocalization

import (
	"errors"

	"github.com/diegobermudez03/couples-backend/pkg/localization"
	"github.com/pariz/gountries"
	"golang.org/x/text/language"
)

type LocalizationServiceImpl struct {
}

func NewLocalizationServiceImpl() localization.LocalizationService{
	return &LocalizationServiceImpl{}
}


func (s *LocalizationServiceImpl) ValidateCountry(code string) error{
	query := gountries.New()
	if _, err := query.FindCountryByAlpha(code); err != nil{
		return errors.New("invalid country code")
	}
	return nil
}

func (s *LocalizationServiceImpl) ValidateLanguage(code string) error{
	_, err := language.Parse(code)
	if err != nil {
		return errors.New("invalid language code")
	}
	return nil
}