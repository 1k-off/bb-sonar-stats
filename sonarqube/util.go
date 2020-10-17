package sonarqube

import (
	b64 "encoding/base64"
	"log"
)

// Base64Encode is an util function to get base64 encoded strung. Need for auth to sonarqube server.
func Base64Encode(s string) string {
	return b64.StdEncoding.EncodeToString([]byte(s))
}

// HandleError is a wrapper for al error handling.
// TODO: error handling.
func HandleError(err error) {
	if err != nil {
		log.Println(err)
	}
}

// IsBetween is a function for check if number is between two numbers
func IsBetween(num, min, max float64) bool {
	if (num >= min) && (num <= max) {
		return true
	} else {
		return false
	}
}
