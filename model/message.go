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
	Message []string
	// ValidationErrors
	ValidationErrors []Error
	// Footer
	Footer string
}

/*
### Developers Italia - Issue Opener for publiccode.yml
We have discovered potential issue on validating your `publiccode.yml`.
Here some details:

- url missing mandatory key
- name missing mandatory key
- longDescription too short (2), min 500 chars

Please review your pubbliccode.
## find out more: https://developers.italia.it
*/
