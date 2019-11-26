package engines

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"testing"

	log "github.com/sirupsen/logrus"
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
	// not yet implemented
}

func TestFilterValidIssues(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	f := []string{"../tests/issues.json"}
	for _, file := range f {
		var is Issues

		data, err := ioutil.ReadFile(file)
		if err != nil {
			t.Errorf("error in reading %s file: %v", f, err)
		}
		err = json.Unmarshal(data, &is)
		if err != nil {
			t.Errorf("error unmarshalling response from GH issues API: %v", err)
		}

		out, err := filterMyIssues(is)
		if err != nil {
			t.Errorf("error filtering GH issues %v", err)
		}
		assert.Equal(t, len(out), 1)

		for _, o := range out {
			assert.Equal(t, o.User.Login, ghUsername)
		}
	}
}
func TestFilterInvalidIssues(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	f := []string{"../tests/issues_min.json"}
	for _, file := range f {
		var is Issues

		data, err := ioutil.ReadFile(file)
		if err != nil {
			t.Errorf("error in reading %s file: %v", f, err)
		}
		err = json.Unmarshal(data, &is)
		if err != nil {
			t.Errorf("error unmarshalling response from GH issues API: %v", err)
		}

		out, err := filterMyIssues(is)
		if err != nil {
			t.Errorf("error filtering GH issues %v", err)
		}
		assert.Equal(t, len(out), 0)
		assert.Empty(t, out)
	}
}
