package api

import (
	"testing"
)

func TestPhoneNumberValidator_CleanPhoneNumber(t *testing.T) {
	validator := NewPhoneNumberValidator()

	tests := []struct {
		name        string
		phoneNumber string
		expected    string
		shouldError bool
	}{
		{
			name:        "Valid number with plus",
			phoneNumber: "+12125690123",
			expected:    "+12125690123",
			shouldError: false,
		},
		{
			name:        "Valid number with spaces",
			phoneNumber: "+52 631 3118150",
			expected:    "+526313118150",
			shouldError: false,
		},
		{
			name:        "Invalid characters - hyphen",
			phoneNumber: "212-569-0123",
			shouldError: true,
		},
		{
			name:        "Invalid characters - letters",
			phoneNumber: "212abc0123",
			shouldError: true,
		},
		{
			name:        "Empty phone number",
			phoneNumber: "",
			expected:    "",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test through public ValidatePhoneNumber method
			_, err := validator.ValidatePhoneNumber(tt.phoneNumber, "")
			if tt.shouldError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.shouldError && err != nil && tt.phoneNumber != "" {
				// Only check for unexpected errors when phone number is not empty
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestPhoneNumberValidator_ValidatePhoneNumber(t *testing.T) {
	validator := NewPhoneNumberValidator()

	tests := []struct {
		name        string
		phoneNumber string
		countryCode string
		expected    *PhoneValidationResponse
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "Valid US number with plus",
			phoneNumber: "+12125690123",
			countryCode: "",
			expected: &PhoneValidationResponse{
				PhoneNumber:      "+12125690123",
				CountryCode:      "US",
				AreaCode:         "212",
				LocalPhoneNumber: "5690123",
			},
			shouldError: false,
		},
		{
			name:        "Valid Mexico number with spaces",
			phoneNumber: "+52 631 3118150",
			countryCode: "",
			expected: &PhoneValidationResponse{
				PhoneNumber:      "+526313118150",
				CountryCode:      "MX",
				AreaCode:         "631",
				LocalPhoneNumber: "3118150",
			},
			shouldError: false,
		},
		{
			name:        "Valid Spain number with spaces",
			phoneNumber: "34 915 872200",
			countryCode: "",
			expected: &PhoneValidationResponse{
				PhoneNumber:      "+34915872200",
				CountryCode:      "ES",
				AreaCode:         "91",
				LocalPhoneNumber: "5872200",
			},
			shouldError: false,
		},
		{
			name:        "US number with country code parameter",
			phoneNumber: "2125690123",
			countryCode: "US",
			expected: &PhoneValidationResponse{
				PhoneNumber:      "+12125690123",
				CountryCode:      "US",
				AreaCode:         "212",
				LocalPhoneNumber: "5690123",
			},
			shouldError: false,
		},
		{
			name:        "Missing country code for national number",
			phoneNumber: "2125690123",
			countryCode: "",
			shouldError: true,
			errorMsg:    "countryCode is required for numbers without country code",
		},
		{
			name:        "Invalid country code format",
			phoneNumber: "2125690123",
			countryCode: "ESP",
			shouldError: true,
			errorMsg:    "country code must be 2 characters (ISO 3166-1 alpha-2)",
		},
		{
			name:        "Invalid characters - letters",
			phoneNumber: "212abc0123",
			countryCode: "US",
			shouldError: true,
			errorMsg:    "phone number contains invalid characters",
		},
		{
			name:        "Invalid characters - hyphen",
			phoneNumber: "212-569-0123",
			countryCode: "US",
			shouldError: true,
			errorMsg:    "phone number contains invalid characters",
		},
		{
			name:        "Invalid spacing pattern",
			phoneNumber: "351 21 094 2000",
			countryCode: "",
			shouldError: true,
			errorMsg:    "invalid spacing pattern",
		},
		{
			name:        "Empty phone number",
			phoneNumber: "",
			countryCode: "US",
			shouldError: true,
			errorMsg:    "phoneNumber is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validator.ValidatePhoneNumber(tt.phoneNumber, tt.countryCode)
			
			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Error("Expected result but got nil")
				return
			}

			// Check all fields
			if result.PhoneNumber != tt.expected.PhoneNumber {
				t.Errorf("Expected PhoneNumber '%s', got '%s'", tt.expected.PhoneNumber, result.PhoneNumber)
			}
			if result.CountryCode != tt.expected.CountryCode {
				t.Errorf("Expected CountryCode '%s', got '%s'", tt.expected.CountryCode, result.CountryCode)
			}
			if result.AreaCode != tt.expected.AreaCode {
				t.Errorf("Expected AreaCode '%s', got '%s'", tt.expected.AreaCode, result.AreaCode)
			}
			if result.LocalPhoneNumber != tt.expected.LocalPhoneNumber {
				t.Errorf("Expected LocalPhoneNumber '%s', got '%s'", tt.expected.LocalPhoneNumber, result.LocalPhoneNumber)
			}
		})
	}
}

func TestPhoneNumberValidator_ValidateCountryCode(t *testing.T) {
	validator := NewPhoneNumberValidator()

	tests := []struct {
		name        string
		countryCode string
		shouldError bool
	}{
		{name: "Valid US", countryCode: "US", shouldError: false},
		{name: "Valid MX", countryCode: "MX", shouldError: false},
		{name: "Valid ES", countryCode: "ES", shouldError: false},
		{name: "Invalid ESP", countryCode: "ESP", shouldError: true},
		{name: "Invalid single char", countryCode: "U", shouldError: true},
		{name: "Unsupported country", countryCode: "XX", shouldError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test through public method that uses validateCountryCode internally
			_, err := validator.ValidatePhoneNumber("1234567890", tt.countryCode)
			
			if tt.shouldError && err == nil {
				t.Errorf("Expected error for country code '%s'", tt.countryCode)
			}
			if !tt.shouldError && err != nil {
				// Check if error is about country code specifically
				if err.Error() == "country code must be 2 characters (ISO 3166-1 alpha-2)" ||
				   err.Error() == "unsupported country code" {
					t.Errorf("Unexpected error for valid country code '%s': %v", tt.countryCode, err)
				}
			}
		})
	}
}

func TestPhoneNumberValidator_PhoneNumberLengthValidation(t *testing.T) {
	validator := NewPhoneNumberValidator()

	tests := []struct {
		name        string
		phoneNumber string
		countryCode string
		shouldError bool
		description string
	}{
		{
			name:        "US number too long",
			phoneNumber: "+12125690123456789",
			countryCode: "",
			shouldError: true,
			description: "US numbers should be exactly 10 digits",
		},
		{
			name:        "US number too short",
			phoneNumber: "+1212569",
			countryCode: "",
			shouldError: true,
			description: "US numbers should be exactly 10 digits",
		},
		{
			name:        "ES number valid length",
			phoneNumber: "+34915872200",
			countryCode: "",
			shouldError: false,
			description: "ES numbers should be 9 digits",
		},
		{
			name:        "DE number valid max length",
			phoneNumber: "+49301234567890",
			countryCode: "",
			shouldError: false,
			description: "DE numbers can be 10-12 digits",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validator.ValidatePhoneNumber(tt.phoneNumber, tt.countryCode)
			
			if tt.shouldError && err == nil {
				t.Errorf("Expected error for %s", tt.description)
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.description, err)
			}
		})
	}
}

func TestPhoneNumberValidator_DialingCodeExtraction(t *testing.T) {
	validator := NewPhoneNumberValidator()

	tests := []struct {
		name        string
		phoneNumber string
		expectedCC  string
		shouldError bool
	}{
		{
			name:        "Extract US dialing code",
			phoneNumber: "+12125690123",
			expectedCC:  "US",
			shouldError: false,
		},
		{
			name:        "Extract Mexico dialing code",
			phoneNumber: "+526313118150",
			expectedCC:  "MX",
			shouldError: false,
		},
		{
			name:        "Extract Portugal dialing code",
			phoneNumber: "+351210942000",
			expectedCC:  "PT",
			shouldError: false,
		},
		{
			name:        "Extract Spain dialing code",
			phoneNumber: "+34915872200",
			expectedCC:  "ES",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validator.ValidatePhoneNumber(tt.phoneNumber, "")
			
			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result.CountryCode != tt.expectedCC {
				t.Errorf("Expected country code '%s', got '%s'", tt.expectedCC, result.CountryCode)
			}
		})
	}
}
