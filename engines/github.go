package engines

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/italia/developers-italia-backend/crawler/httpclient"
	"github.com/sebbalex/issue-opener/model"
	log "github.com/sirupsen/logrus"
)

// Comment GH comment struct
type Comment model.Comment

// Issues type
type Issues []model.GHIssue

// SingleRepoHandler returns the client handler for an a
// single repository (every domain has a different handler implementation).
type SingleRepoHandler func(domain Domain, url *url.URL, comments chan Comment) error

// CommentsHandler ...
type CommentsHandler func(domain Domain, url *url.URL) error

var ghUsername string = "developers-italia-bot"

// Ex:
// time="2019-11-18T01:05:25Z" level=error msg="Error parsing publiccode.yml for https://raw.githubusercontent.com/AgID/pat/master/publiccode.yml."
// time="2019-11-18T01:05:25Z" level=error msg="[AgID/pat] invalid publiccode.yml: logo: invalid image size of 63 (min 120px of width): src/app/grafica/pat_semplice.png"
// time="2019-11-18T01:05:25Z" level=error msg="Appending the bad file URL to the list: https://raw.githubusercontent.com/AgID/pat/master/publiccode.yml"
// time="2019-11-18T01:05:26Z" level=error msg="Error parsing publiccode.yml for https://raw.githubusercontent.com/AgID/pat/master/publiccode.yml."
// time="2019-11-18T01:05:26Z" level=error msg="[AgID/pat] invalid publiccode.yml: logo: invalid image size of 63 (min 120px of width): src/app/grafica/pat_semplice.png"
// time="2019-11-18T01:05:26Z" level=error msg="Appending the bad file URL to the list: https://raw.githubusercontent.com/AgID/pat/master/publiccode.yml"

//https://developer.github.com/v3/issues/#list-issues-for-a-repository GET /repos/:owner/:repo/issues

// RegisterSingleGithubAPI ....
func RegisterSingleGithubAPI() SingleRepoHandler {
	return func(domain Domain, u *url.URL, comments chan Comment) error {
		// Set BasicAuth header.
		headers := make(map[string]string)
		headers["Authorization"] = githubBasicAuth(domain)

		// Set domain host to new host.
		domain.Host = u.Hostname()

		u.Path = path.Join("repos", u.Path, "issues")
		u.Path = strings.Trim(u.Path, "/")
		u.Host = "api." + u.Host

		// Get List of issues for repository.
		log.Debugf("calling API %s", u)
		resp, err := httpclient.GetURL(u.String(), headers)
		if err != nil {
			log.Errorf("error getting issues api: %v", err)
			return err
		}
		if resp.Status.Code != http.StatusOK {
			log.Warnf("Request returned: %s", string(resp.Body))
			return errors.New("request returned an incorrect http.Status: " + resp.Status.Text)
		}

		var v Issues
		err = json.Unmarshal(resp.Body, &v)
		if err != nil {
			log.Errorf("error unmarshalling response from GH issues API: %v", err)
			return err
		}
		log.Debugf("issues: %v", v)

		// filtering mine
		v, err = filterMyIssues(v)
		if err != nil {
			log.Errorf("error filtering issues %v", err)
			return err
		}

		return nil
	}
}

func filterMyIssues(ghis Issues) (Issues, error) {
	log.Debugf("filterMyIssues() issues: %v", ghis)
	b := ghis[:0]
	for _, x := range ghis {
		if x.User.Login == ghUsername {
			b = append(b, x)
		}
	}
	log.Debugf("filtered issues: %v", b)
	return b, nil
}

func githubBasicAuth(domain Domain) string {
	if len(domain.BasicAuth) > 0 {
		auth := domain.BasicAuth[rand.Intn(len(domain.BasicAuth))]
		return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	}
	return ""
}
