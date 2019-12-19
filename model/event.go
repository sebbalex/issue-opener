package model

import "net/url"

// Event main class
type Event struct {
	// URL to analyze
	URL *url.URL `json:"url"`
	// is it valid
	Valid bool `json:"valid"`
	// Validation Errors
	ValidationError []Error `json:"validationErrors"`
	// Message slice
	Message []Message `json:"message"`
	// DryRun
	DryRun bool `json:"dryRun"`
}
