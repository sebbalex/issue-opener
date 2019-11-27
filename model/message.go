package model

import "net/url"

// Message represent one comment dropped by Issue Opener.
// It contains all vaidation errors reported in related
// repository.
type Message struct {
	// URL Represent comment URL
	URL *url.URL
	// Header
	Header string
	// Message
	Message string
	// Footer
	Footer string
}
