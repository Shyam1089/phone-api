package api

import (
	"errors"
	"regexp"
	"strings"
)

var CountryPhoneLengths = map[string][2]int{
	"US": {10, 10},	
	"CA": {10, 10},
	"MX": {10, 10},
	"ES": {9, 9},
	"PT": {9, 9},
	"GB": {10, 11},
	"FR": {10, 10},
	"DE": {10, 12},
	"IT": {9, 11},
	"BR": {10, 11},
}

var CountryDialingCodes = map[string]string{
	"US": "1",
	"CA": "1",
	"MX": "52",
	"ES": "34",
	"PT": "351",
	"GB": "44",
	"FR": "33",
	"DE": "49",
	"IT": "39",
	"BR": "55",
}

var DialingCodeToCountry = map[string]string{
	"1":   "US",
	"52":  "MX",
	"34":  "ES",
	"351": "PT",
	"44":  "GB",
	"33":  "FR",
	"49":  "DE",
	"39":  "IT",
	"55":  "BR",
}

type PhoneValidationRequest struct {
	PhoneNumber string `form:"phoneNumber" json:"phoneNumber"`
	CountryCode string `form:"countryCode" json:"countryCode"`
}

type PhoneValidationResponse struct {
	PhoneNumber      string `json:"phoneNumber"`
	CountryCode      string `json:"countryCode"`
	AreaCode         string `json:"areaCode"`
	LocalPhoneNumber string `json:"localPhoneNumber"`
}

type ErrorResponse struct {
	PhoneNumber string            `json:"phoneNumber"`
	Error       map[string]string `json:"error"`
}

type PhoneNumberValidator struct{}

func NewPhoneNumberValidator() *PhoneNumberValidator {
	return &PhoneNumberValidator{}
}

func (v *PhoneNumberValidator) ValidatePhoneNumber(phoneNumber, countryCode string) (*PhoneValidationResponse, error) {
	if phoneNumber == "" {
		return nil, errors.New("phoneNumber is required")
	}

	if err := v.validateSpacing(phoneNumber); err != nil {
		return nil, err
	}

	cleanedNumber, err := v.cleanPhoneNumber(phoneNumber)
	if err != nil {
		return nil, err
	}

	extractedCountryCode, areaCode, localNumber, err := v.parsePhoneNumber(cleanedNumber, countryCode)
	if err != nil {
		return nil, err
	}

	if err := v.validateCountryCode(extractedCountryCode); err != nil {
		return nil, err
	}

	if err := v.validatePhoneLength(areaCode+localNumber, extractedCountryCode); err != nil {
		return nil, err
	}

	response := &PhoneValidationResponse{
		PhoneNumber:      v.formatPhoneNumber(extractedCountryCode, areaCode, localNumber),
		CountryCode:      extractedCountryCode,
		AreaCode:         areaCode,
		LocalPhoneNumber: localNumber,
	}

	return response, nil
}

func (v *PhoneNumberValidator) cleanPhoneNumber(phoneNumber string) (string, error) {
	validChars := regexp.MustCompile(`^[\d\s+]+$`)
	if !validChars.MatchString(phoneNumber) {
		return "", errors.New("phone number contains invalid characters")
	}

	cleaned := strings.ReplaceAll(phoneNumber, " ", "")
	
	return cleaned, nil
}

func (v *PhoneNumberValidator) parsePhoneNumber(phoneNumber, providedCountryCode string) (string, string, string, error) {
	hasPlus := strings.HasPrefix(phoneNumber, "+")
	if hasPlus {
		phoneNumber = phoneNumber[1:]
	}

	var countryCode string
	var nationalNumber string

	if hasPlus || v.hasDialingCode(phoneNumber) {
		dialingCode, remaining, err := v.extractDialingCode(phoneNumber)
		if err != nil {
			return "", "", "", err
		}
		
		country, exists := DialingCodeToCountry[dialingCode]
		if !exists {
			return "", "", "", errors.New("unsupported country dialing code")
		}
		
		countryCode = country
		nationalNumber = remaining
	} else {
		if providedCountryCode == "" {
			return "", "", "", errors.New("countryCode is required for numbers without country code")
		}
		countryCode = providedCountryCode
		nationalNumber = phoneNumber
	}

	areaCode, localNumber := v.splitNationalNumber(nationalNumber, countryCode)
	
	return countryCode, areaCode, localNumber, nil
}

func (v *PhoneNumberValidator) validateSpacing(originalPhoneNumber string) error {
	if strings.Contains(originalPhoneNumber, " ") {
		parts := strings.Split(originalPhoneNumber, " ")
		if len(parts) == 4 {
			return errors.New("invalid spacing pattern")
		}
	}
	
	return nil
}

func (v *PhoneNumberValidator) hasDialingCode(phoneNumber string) bool {
	for dialingCode := range DialingCodeToCountry {
		if strings.HasPrefix(phoneNumber, dialingCode) {
			return true
		}
	}
	return false
}

func (v *PhoneNumberValidator) extractDialingCode(phoneNumber string) (string, string, error) {
	if len(phoneNumber) >= 3 {
		threeDigit := phoneNumber[:3]
		if _, exists := DialingCodeToCountry[threeDigit]; exists {
			return threeDigit, phoneNumber[3:], nil
		}
	}

	if len(phoneNumber) >= 2 {
		twoDigit := phoneNumber[:2]
		if _, exists := DialingCodeToCountry[twoDigit]; exists {
			return twoDigit, phoneNumber[2:], nil
		}
	}

	if len(phoneNumber) >= 1 {
		oneDigit := phoneNumber[:1]
		if _, exists := DialingCodeToCountry[oneDigit]; exists {
			return oneDigit, phoneNumber[1:], nil
		}
	}

	return "", "", errors.New("unable to extract dialing code")
}

func (v *PhoneNumberValidator) splitNationalNumber(nationalNumber, countryCode string) (string, string) {
	if len(nationalNumber) >= 3 {
		switch countryCode {
		case "US", "CA", "MX":
			return nationalNumber[:3], nationalNumber[3:]
		case "ES", "PT", "FR", "IT", "BR":
			return nationalNumber[:2], nationalNumber[2:]
		case "GB":
			if len(nationalNumber) >= 4 {
				return nationalNumber[:4], nationalNumber[4:]
			}
			return nationalNumber[:3], nationalNumber[3:]
		case "DE":
			return nationalNumber[:3], nationalNumber[3:]
		}
	}
	
	return "", nationalNumber
}

func (v *PhoneNumberValidator) validateCountryCode(countryCode string) error {
	if len(countryCode) != 2 {
		return errors.New("country code must be 2 characters (ISO 3166-1 alpha-2)")
	}

	if _, exists := CountryPhoneLengths[countryCode]; !exists {
		return errors.New("unsupported country code")
	}

	return nil
}

func (v *PhoneNumberValidator) validatePhoneLength(nationalNumber, countryCode string) error {
	lengths, exists := CountryPhoneLengths[countryCode]
	if !exists {
		return errors.New("unsupported country code")
	}

	minLength, maxLength := lengths[0], lengths[1]
	actualLength := len(nationalNumber)

	if actualLength < minLength || actualLength > maxLength {
		return errors.New("phone number length is invalid for country " + countryCode)
	}

	return nil
}

func (v *PhoneNumberValidator) formatPhoneNumber(countryCode, areaCode, localNumber string) string {
	dialingCode := CountryDialingCodes[countryCode]
	return "+" + dialingCode + areaCode + localNumber
}
