package engines

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"testing"

	. "github.com/sebbalex/issue-opener/model"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGithub(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	RegisterClientAPIs()
	var e = NewEngine()
	event := Event{}
	for _, repoURL := range githubURLs {
		urlParsed, err := url.Parse(repoURL)
		if err != nil {
			t.Errorf("Error parsing url %s err %v", repoURL, err)
		}
		event.URL = urlParsed
		event.Message = make(chan Message)
		d, err := e.IdentifyVCS(urlParsed)
		assert.Equal(t, e.StartFlow(&event, d), nil)
	}
}

func TestEnrichWithComments(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	RegisterClientAPIs()
	var e = NewEngine()
	var comments Comments

	for _, repoURL := range githubURLs {
		urlParsed, err := url.Parse(repoURL)
		if err != nil {
			t.Errorf("Error parsing url %s err %v", repoURL, err)
		}
		d, err := e.IdentifyVCS(urlParsed)

		// should provide a valid id
		issueID := 1
		out, err := enrichWithComments(*d, urlParsed, issueID)
		if err != nil {
			t.Errorf("Error %v", err)
		}
		assert.GreaterOrEqual(t, len(out), 0)
		assert.IsType(t, out, comments)
	}
}

func testGHAuth(t *testing.T) {
	// not yet implemented
}

func TestFilterValidIssue(t *testing.T) {
	ghUsername = "developers-italia-bot"
	log.SetLevel(log.InfoLevel)
	f := []string{"../tests/issues.json", "../tests/comments.json"}
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
		// be carefull, this may fail if username in github.go:34 has changed
		assert.Equal(t, 1, len(out))

		for _, o := range out {
			assert.Equal(t, o.User.Login, ghUsername)
		}
	}
}
func TestFilterInvalidIssues(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	// those won't match our name
	// we expect 0 results
	f := []string{"../tests/issues_min.json", "../tests/comments_min.json"}
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

		assert.Equal(t, 0, len(out))
		assert.Empty(t, out)
	}
}

func TestFilterValidComments(t *testing.T) {
	ghUsername = "developers-italia-bot"
	log.SetLevel(log.InfoLevel)
	f := []string{"../tests/comments.json"}
	for _, file := range f {
		var is Comments

		data, err := ioutil.ReadFile(file)
		if err != nil {
			t.Errorf("error in reading %s file: %v", f, err)
		}
		err = json.Unmarshal(data, &is)
		if err != nil {
			t.Errorf("error unmarshalling response from GH comments API: %v", err)
		}

		out, err := filterMyComments(is)
		if err != nil {
			t.Errorf("error filtering GH comments %v", err)
		}
		// be carefull, this may fail if username in github.go:34 has changed
		assert.Equal(t, 1, len(out))

		for _, o := range out {
			assert.Equal(t, o.User.Login, ghUsername)
		}
	}
}
func TestFilterInvalidComments(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	// those won't match our name
	// we expect 0 results
	f := []string{"../tests/comments_min.json"}
	for _, file := range f {
		var is Comments

		data, err := ioutil.ReadFile(file)
		if err != nil {
			t.Errorf("error in reading %s file: %v", f, err)
		}
		err = json.Unmarshal(data, &is)
		if err != nil {
			t.Errorf("error unmarshalling response from GH comments API: %v", err)
		}

		out, err := filterMyComments(is)
		if err != nil {
			t.Errorf("error filtering GH comments %v", err)
		}

		assert.Equal(t, 0, len(out))
		assert.Empty(t, out)
	}
}
