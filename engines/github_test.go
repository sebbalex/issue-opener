package engines

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGithub(t *testing.T) {
	RegisterClientAPIs()
	var e = NewEngine()
	for _, repoURL := range githubURLs {
		urlParsed, err := url.Parse(repoURL)
		if err != nil {
			t.Fail()
		}
		d, err := e.IdentifyVCS(urlParsed)
		assert.Equal(t, e.StartFlow(urlParsed, d), nil)
	}
}

func testGHAuth(t *testing.T) {

}
