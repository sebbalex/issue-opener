package model

import "net/url"

import "fmt"

import "strings"

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

// Template function will create necessary stuff to publish
// new message on code hosting platform
func (m *Message) Template() {
	m.Header = `
			### Developers Italia - Issue Opener for publiccode.yml
			We have discovered potential issue on validating your publiccode.yml.
			Here some details:

		`
	if len(m.ValidationErrors) > 0 {
		for _, valErr := range m.ValidationErrors {
			m.Message = append(m.Message, fmt.Sprintf("- %s %s", valErr.Key, valErr.Reason))
		}
	}
	m.Footer = `
		Please review your pubbliccode.
		## find out more: https://developers.italia.it
	`
}

func (m *Message) String() string {
	return fmt.Sprintf("%s\n%v\n%s", m.Header, strings.Join(m.Message, "\n"), m.Footer)
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
