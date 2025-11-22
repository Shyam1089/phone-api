package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	validator *PhoneNumberValidator
}

func NewHandler() *Handler {
	return &Handler{
		validator: NewPhoneNumberValidator(),
	}
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"service": "phone-number-lookup",
	})
}

func (h *Handler) PhoneNumberLookup(c *gin.Context) {
	var req PhoneValidationRequest
	
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			PhoneNumber: req.PhoneNumber,
			Error: map[string]string{
				"validation": "invalid request parameters",
			},
		})
		return
	}

	response, err := h.validator.ValidatePhoneNumber(req.PhoneNumber, req.CountryCode)
	if err != nil {
		errorMsg := h.mapValidationError(err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			PhoneNumber: req.PhoneNumber,
			Error:       errorMsg,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) mapValidationError(errMsg string) map[string]string {
	switch {
	case errMsg == "phoneNumber is required":
		return map[string]string{
			"phoneNumber": "required value is missing",
		}
	case errMsg == "countryCode is required for numbers without country code":
		return map[string]string{
			"countryCode": "required value is missing",
		}
	case errMsg == "country code must be 2 characters (ISO 3166-1 alpha-2)":
		return map[string]string{
			"countryCode": "invalid format (must be ISO 3166-1 alpha-2)",
		}
	case errMsg == "unsupported country code":
		return map[string]string{
			"countryCode": "unsupported country code",
		}
	case errMsg == "phone number contains invalid characters":
		return map[string]string{
			"phoneNumber": "contains invalid characters",
		}
	case errMsg == "invalid spacing pattern":
		return map[string]string{
			"phoneNumber": "invalid spacing pattern",
		}
	case errMsg == "unsupported country dialing code":
		return map[string]string{
			"phoneNumber": "unsupported country dialing code",
		}
	default:
		if len(errMsg) > 30 && errMsg[:30] == "phone number length is invalid" {
			return map[string]string{
				"phoneNumber": "length is invalid for country",
			}
		}
		return map[string]string{
			"phoneNumber": "invalid format",
		}
	}
}

func (h *Handler) SetupRoutes(router *gin.Engine) {
	router.GET("/health", h.HealthCheck)
	
	v1 := router.Group("/v1")
	{
		v1.GET("/phone-numbers", h.PhoneNumberLookup)
	}
}
