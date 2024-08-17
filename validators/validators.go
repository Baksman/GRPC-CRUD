package validators

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

func IsValidBirthDay(fl validator.FieldLevel) bool {
	layout := "2006-01-02"
	t, err := time.Parse(layout, fl.Field().String())
	if err != nil {
		return false
	}
	// Check if the parsing resulted in a valid time object (not location specific)
	now := time.Now()
	// Adjust for potential rounding errors during parsing (#42795)
	return t.Year() > now.Year()-5 && t.After(time.Time{})
}

func IsURL(fl validator.FieldLevel) bool {
	// Use a regular expression to validate basic URL format
	urlRegex := regexp.MustCompile(`^(http|https)://([\w.]+)+(:\d+)?(/[-._/~:@!$&'\(\)\*+,;=?:%]*)?(\?\S+)?$`)
	return urlRegex.MatchString(fl.Field().String()) && isValidDomain(fl.Field().String())
}
func isValidDomain(domain string) bool {
	// You can customize this regex for stricter domain name validation
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9](?:\.[a-zA-Z0-9]{1,61})*$`)
	return domainRegex.MatchString(domain)
}

func IsEmail(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`).MatchString(fl.Field().String())
}

func ValidatorErrorFormater(err error) error {
	customMessages := map[string]string{
		"Username.required":    "username field is required",
		"Username.min":         "username should be atleast 3 characters",
		"Username.max":         "Username should be at most 12 characters",
		"Email.required":       "email field is required",
		"Email.IsEmail":        "email must be a valid email address",
		"Country.required":     "The country must be a valid email address",
		"Password.required":    "The password field is required",
		"Password.min":         "The password requires minimum of 8 characters",
		"Password.max":         "The password requires maximum of 40 characters",
		"OTP.required":         "OTP field is required",
		"OTP.min":              "OTP field should be 5 characters",
		"OTP.max":              "OTP field should be 5 characters",
		"ID.required":          "ID is required",
		"DOB.IsValidBirthDay":  "Invalid birthday",
		"Dp.isValidDomain":     "Invalid image url",
		"NewPassword.required": "New password field is required",
		"NewPassword.min":      "New password  should be minimum of 8 characters",
		"NewPassword.max":      "New password  should be maximum of 40 characters",
		"FcmToken.required":    "Fcm token is required",
		"FcmToken.min":         "Fcm token should be minimum of 24 characters",
		"PhoneNumber.min":      "PhoneNumber should be minimum of 10 characters",
		"PhoneNumber.max":      "PhoneNumber should be maximum of 14 characters",
	}

	if err != nil {
		// Handle validation errors (e.g., return specific error messages)
		validationErrors := err.(validator.ValidationErrors)
		errorMessage := ""
		for _, fieldError := range validationErrors {
			fieldName := fmt.Sprintf("%s.%s", fieldError.StructField(), fieldError.Tag())
			if customMessage, ok := customMessages[fieldName]; ok {
				errorMessage += customMessage + ", "
			} else {
				errorMessage += fmt.Sprintf(fieldError.Field() + " is invalid, ")
			}

		}
		errorMessage = strings.Trim(errorMessage, ", ")
		return fmt.Errorf(errorMessage)
	}
	return nil
}
