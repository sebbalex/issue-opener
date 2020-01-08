package model

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// Message represent one comment dropped by Issue Opener.
// It contains all vaidation errors reported in related
// repository.
type Message struct {
	// URL Represent comment URL
	URL *url.URL
	// IssueID
	IssueID int
	// Header
	Header string
	// Title
	Title string
	// Append indicate if message is in apend to an already
	// exists issue or if it is needed to create a new one
	Append bool
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
	m.Title = "Validation errors on publiccode.yml"
	m.Header = `### Developers Italia - Issue Opener for publiccode.yml
We have discovered potential issue on validating your publiccode.yml.
Here some details:

`
	if len(m.ValidationErrors) > 0 {
		for _, valErr := range m.ValidationErrors {
			m.Message = append(m.Message, fmt.Sprintf("- %s %s", valErr.Key, valErr.Reason))
		}
	}
	m.Footer = `
Please review your publiccode.
## find out more: https://developers.italia.it`
}

func (m *Message) String() string {
	return fmt.Sprintf("%s\n%v\n%s", m.Header, strings.Join(m.Message, "\n"), m.Footer)
}

// MessageToJSON convert actual obj to JSON
// following rules: https://developer.github.com/v3/issues/#create-an-issue
func (m *Message) MessageToJSON() ([]byte, error) {
	obj := map[string]interface{}{"title": m.Title, "body": m.String()}
	return json.Marshal(obj)
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
