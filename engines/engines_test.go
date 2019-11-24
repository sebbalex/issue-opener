package engines

import (
	"net/url"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var goodURLs = []string{
	"https://raw.githubusercontent.com/AgID/pat/master/publiccode.yml",
}

var badURLs = []string{
	"https://bitbucket.org/Comune_Venezia/iris/raw/master/publiccode.yml",
	"https://gitlab.com/fusslab/fuss/raw/master/publiccode.yml",
	"https://gerrit.libreoffice.org/plugins/gitiles/core/+/refs/heads/distro/cib/libreoffice-6-1",
}

var githubURLs = []string{"https://raw.githubusercontent.com/AgID/pat/master/publiccode.yml"}

func TestIdentifyVCS(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	RegisterClientAPIs()
	var e = NewEngine()
	for _, repoURL := range goodURLs {
		urlParsed, err := url.Parse(repoURL)
		if err != nil {
			t.Fail()
		}
		_, err = e.IdentifyVCS(urlParsed)
		assert.Equal(t, err, nil)
	}

	for _, repoURL := range badURLs {
		urlParsed, err := url.Parse(repoURL)
		if err != nil {
			t.Fail()
		}
		_, err = e.IdentifyVCS(urlParsed)
		assert.EqualError(t, err, "Not yet implemented")
	}
}
