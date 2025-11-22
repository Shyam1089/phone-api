package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"phone-api/api"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := api.NewHandler()
	handler.SetupRoutes(router)
	return router
}

// TestAPIEndpoints focuses purely on API response testing
func TestHealthEndpoint(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Test HTTP response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	// Test JSON response structure
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])
	assert.Equal(t, "phone-number-lookup", response["service"])
}

// Individual test functions have been replaced with comprehensive API test suite above

// Additional comprehensive API endpoint tests
func TestAPIEndpoints_ComprehensiveResponseTesting(t *testing.T) {
	router := setupTestRouter()

	t.Run("Success Responses", func(t *testing.T) {
		testCases := []struct {
			name     string
			url      string
			expected map[string]string
		}{
			{
				name: "US Number with Plus",
				url:  "/v1/phone-numbers?phoneNumber=%2B12125690123",
				expected: map[string]string{
					"phoneNumber":      "+12125690123",
					"countryCode":      "US",
					"areaCode":         "212",
					"localPhoneNumber": "5690123",
				},
			},
			{
				name: "Mexico Number with Spaces",
				url:  "/v1/phone-numbers?phoneNumber=%2B52%20631%203118150",
				expected: map[string]string{
					"phoneNumber":      "+526313118150",
					"countryCode":      "MX",
					"areaCode":         "631",
					"localPhoneNumber": "3118150",
				},
			},
			{
				name: "Spain Number",
				url:  "/v1/phone-numbers?phoneNumber=34%20915%20872200",
				expected: map[string]string{
					"phoneNumber":      "+34915872200",
					"countryCode":      "ES",
					"areaCode":         "91",
					"localPhoneNumber": "5872200",
				},
			},
			{
				name: "US Number with Country Code Parameter",
				url:  "/v1/phone-numbers?phoneNumber=2125690123&countryCode=US",
				expected: map[string]string{
					"phoneNumber":      "+12125690123",
					"countryCode":      "US",
					"areaCode":         "212",
					"localPhoneNumber": "5690123",
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req, _ := http.NewRequest("GET", tc.url, nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				// Verify HTTP response
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

				// Verify JSON structure
				var response api.PhoneValidationResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				// Verify response fields
				assert.Equal(t, tc.expected["phoneNumber"], response.PhoneNumber)
				assert.Equal(t, tc.expected["countryCode"], response.CountryCode)
				assert.Equal(t, tc.expected["areaCode"], response.AreaCode)
				assert.Equal(t, tc.expected["localPhoneNumber"], response.LocalPhoneNumber)
			})
		}
	})

	t.Run("Error Responses", func(t *testing.T) {
		errorTestCases := []struct {
			name               string
			url                string
			expectedStatus     int
			expectedErrorField string
			expectedPhoneNum   string
		}{
			{
				name:               "Missing Phone Number",
				url:                "/v1/phone-numbers",
				expectedStatus:     http.StatusBadRequest,
				expectedErrorField: "phoneNumber",
				expectedPhoneNum:   "",
			},
			{
				name:               "Missing Country Code",
				url:                "/v1/phone-numbers?phoneNumber=2125690123",
				expectedStatus:     http.StatusBadRequest,
				expectedErrorField: "countryCode",
				expectedPhoneNum:   "2125690123",
			},
			{
				name:               "Invalid Country Code Format",
				url:                "/v1/phone-numbers?phoneNumber=2125690123&countryCode=ESP",
				expectedStatus:     http.StatusBadRequest,
				expectedErrorField: "countryCode",
				expectedPhoneNum:   "2125690123",
			},
			{
				name:               "Invalid Characters - Letters",
				url:                "/v1/phone-numbers?phoneNumber=212abc0123&countryCode=US",
				expectedStatus:     http.StatusBadRequest,
				expectedErrorField: "phoneNumber",
				expectedPhoneNum:   "212abc0123",
			},
			{
				name:               "Invalid Characters - Hyphen",
				url:                "/v1/phone-numbers?phoneNumber=212-569-0123&countryCode=US",
				expectedStatus:     http.StatusBadRequest,
				expectedErrorField: "phoneNumber",
				expectedPhoneNum:   "212-569-0123",
			},
			{
				name:               "Invalid Spacing Pattern",
				url:                "/v1/phone-numbers?phoneNumber=351%2021%20094%202000",
				expectedStatus:     http.StatusBadRequest,
				expectedErrorField: "phoneNumber",
				expectedPhoneNum:   "351 21 094 2000",
			},
			{
				name:               "Number Too Long",
				url:                "/v1/phone-numbers?phoneNumber=%2B1212569012398877",
				expectedStatus:     http.StatusBadRequest,
				expectedErrorField: "phoneNumber",
				expectedPhoneNum:   "+1212569012398877",
			},
			{
				name:               "Number Too Short",
				url:                "/v1/phone-numbers?phoneNumber=%2B1212569",
				expectedStatus:     http.StatusBadRequest,
				expectedErrorField: "phoneNumber",
				expectedPhoneNum:   "+1212569",
			},
		}

		for _, tc := range errorTestCases {
			t.Run(tc.name, func(t *testing.T) {
				req, _ := http.NewRequest("GET", tc.url, nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				// Verify HTTP error response
				assert.Equal(t, tc.expectedStatus, w.Code)
				assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

				// Verify error JSON structure
				var response api.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				// Verify error response fields
				assert.Equal(t, tc.expectedPhoneNum, response.PhoneNumber)
				assert.Contains(t, response.Error, tc.expectedErrorField)
				assert.NotEmpty(t, response.Error[tc.expectedErrorField])
			})
		}
	})

	t.Run("URL Encoding Handling", func(t *testing.T) {
		encodingTests := []struct {
			name     string
			url      string
			expected string
		}{
			{
				name:     "Plus Sign Encoded",
				url:      "/v1/phone-numbers?phoneNumber=%2B12125690123",
				expected: "+12125690123",
			},
			{
				name:     "Spaces Encoded",
				url:      "/v1/phone-numbers?phoneNumber=%2B52%20631%203118150",
				expected: "+526313118150",
			},
			{
				name:     "Mixed Encoding",
				url:      "/v1/phone-numbers?phoneNumber=34%20915%20872200",
				expected: "+34915872200",
			},
		}

		for _, tc := range encodingTests {
			t.Run(tc.name, func(t *testing.T) {
				req, _ := http.NewRequest("GET", tc.url, nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				assert.Equal(t, http.StatusOK, w.Code)
				
				var response api.PhoneValidationResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, response.PhoneNumber)
			})
		}
	})

	t.Run("HTTP Methods", func(t *testing.T) {
		// Test that unsupported methods return appropriate responses
		methods := []string{"POST", "PUT", "DELETE", "PATCH"}
		
		for _, method := range methods {
			t.Run("Method_"+method, func(t *testing.T) {
				req, _ := http.NewRequest(method, "/v1/phone-numbers?phoneNumber=%2B12125690123", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				// Since we only define GET routes, other methods should return 404
				assert.Equal(t, http.StatusNotFound, w.Code)
			})
		}
	})

	t.Run("Response Headers", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// Test that we get JSON content type
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
	})
}
