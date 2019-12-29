package engines

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	vcs "github.com/alranel/go-vcsurl"
	parser "github.com/sebbalex/issue-opener/analyzer"
	httpclient "github.com/sebbalex/issue-opener/httpclient"
	. "github.com/sebbalex/issue-opener/model"
	log "github.com/sirupsen/logrus"
)

// Comments type
type Comments []Comment

// Issues type
type Issues []GHIssue

// SingleRepoHandler returns the client handler for an a
// single repository (every domain has a different handler implementation).
type SingleRepoHandler func(domain Domain, event *Event) error

// CommentsHandler ...
type CommentsHandler func(domain Domain, url *url.URL) error

var ghUsername string = "sebbalex" //developers-italia-bot

// Ex:
// time="2019-11-18T01:05:25Z" level=error msg="Error parsing publiccode.yml for https://raw.githubusercontent.com/AgID/pat/master/publiccode.yml."
// time="2019-11-18T01:05:25Z" level=error msg="[AgID/pat] invalid publiccode.yml: logo: invalid image size of 63 (min 120px of width): src/app/grafica/pat_semplice.png"
// time="2019-11-18T01:05:25Z" level=error msg="Appending the bad file URL to the list: https://raw.githubusercontent.com/AgID/pat/master/publiccode.yml"
// time="2019-11-18T01:05:26Z" level=error msg="Error parsing publiccode.yml for https://raw.githubusercontent.com/AgID/pat/master/publiccode.yml."
// time="2019-11-18T01:05:26Z" level=error msg="[AgID/pat] invalid publiccode.yml: logo: invalid image size of 63 (min 120px of width): src/app/grafica/pat_semplice.png"
// time="2019-11-18T01:05:26Z" level=error msg="Appending the bad file URL to the list: https://raw.githubusercontent.com/AgID/pat/master/publiccode.yml"

func getAPIUrlAndHeaders(domain Domain, u url.URL) (*url.URL, map[string]string) {
	// Set BasicAuth header.
	headers := make(map[string]string)
	headers["Authorization"] = githubBasicAuth(domain)

	// Set domain host to new host.
	domain.Host = u.Hostname()

	u.Path = path.Join("repos", u.Path, "issues")
	u.Path = strings.Trim(u.Path, "/")
	u.Host = "api." + u.Host
	return &u, headers
}

// RegisterAppendIssueGithubAPI register API that will post an issue for given validation
// errors using ghUsername user on GitHub
func RegisterAppendIssueGithubAPI() SingleRepoHandler {
	return func(domain Domain, event *Event) error {
		log.Infof("call RegisterAppendIssueGithubAPI()")
		u, headers := getAPIUrlAndHeaders(domain, *event.URL)
		log.Debugf("headers: %s", headers)
		if len(event.Message) == 0 {
			log.Errorf("message is not present %v", event)
			return errors.New("message not present")
		}
		message := event.Message[0]
		messageJSON, err := message.MessageToJSON()
		u.Path = path.Join(u.Path, strconv.Itoa(message.IssueID), "comments")
		if err != nil {
			log.Errorf("error converting to JSON: %s", err)
		}
		if event.DryRun {
			log.Infof("dry-run mode appending: %s to: %s", messageJSON, u)
		} else {
			log.Debugf("appending: %s to %s", messageJSON, u)
			resp, err := httpclient.PostURL(u.String(), headers, bytes.NewReader(messageJSON))
			if err != nil {
				log.Errorf("error posting issues api: %v", err)
				return err
			}
			if resp.Status.Code != http.StatusOK {
				log.Warnf("Request returned: %s", string(resp.Body))
				return errors.New("request returned an incorrect http.Status: " + resp.Status.Text)
			}
		}
		return nil
	}
}

// RegisterPostIssueGithubAPI register API that will post an issue for given validation
// errors using ghUsername user on GitHub
func RegisterPostIssueGithubAPI() SingleRepoHandler {
	return func(domain Domain, event *Event) error {
		log.Infof("call RegisterPostIssueGithubAPI()")
		u, headers := getAPIUrlAndHeaders(domain, *event.URL)
		log.Debugf("headers: %s", headers)
		if len(event.Message) == 0 {
			log.Errorf("message is not present %v", event)
			return errors.New("message not present")
		}
		message := event.Message[0]
		messageJSON, err := message.MessageToJSON()
		if err != nil {
			log.Errorf("error converting to JSON: %s", err)

		}
		if event.DryRun {
			log.Infof("dry-run mode posting: %s", messageJSON)
		} else {
			log.Debugf("posting: %s", messageJSON)
			resp, err := httpclient.PostURL(u.String(), headers, bytes.NewReader(messageJSON))
			if err != nil {
				log.Errorf("error posting issues api: %v", err)
				return err
			}
			if resp.Status.Code != http.StatusOK {
				log.Warnf("Request returned: %s", string(resp.Body))
				return errors.New("request returned an incorrect http.Status: " + resp.Status.Text)
			}
		}
		return nil
	}
}

//https://developer.github.com/v3/issues/#list-issues-for-a-repository GET /repos/:owner/:repo/issues

// RegisterSingleGithubAPI register API to get all issues
// which have been created from ghUsername user on GitHub
func RegisterSingleGithubAPI() SingleRepoHandler {
	return func(domain Domain, event *Event) error {
		urlBase := event.URL
		u, headers := getAPIUrlAndHeaders(domain, *urlBase)

		// filtering for created by me
		q := u.Query()
		q.Set("creator", ghUsername)
		u.RawQuery = q.Encode()

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

		// filtering mine
		v, err = filterMyIssues(v)
		if err != nil {
			log.Errorf("error filtering issues %v", err)
			return err
		}
		log.Tracef("filtered issues: %v", v)

		// Get comments, parse 'em and populate Message chan in Event
		for _, issue := range v {
			comments, errs := enrichWithComments(domain, urlBase, issue.Number)
			if errs != nil {
				return errs
			}
			comments = append(comments, createCommentFromIssueBody(v)...)
			return parser.ParseGHComments(event, comments, issue.Number)
		}

		return nil
	}
}

func createCommentFromIssueBody(issues Issues) (comments Comments) {
	for _, issue := range issues {
		comment := Comment{}
		comment.Body = issue.Body
		comments = append(comments, comment)
		// ... other fields
	}
	return comments
}

func enrichWithComments(domain Domain, urlBase *url.URL, issueID int) (Comments, error) {
	if vcs.IsRawFile(urlBase) {
		urlBase = vcs.GetRepo(urlBase)
	}

	u, headers := getAPIUrlAndHeaders(domain, *urlBase)
	u.Path = path.Join(u.Path, strconv.Itoa(issueID), "comments")

	// Get List of issues for repository.
	log.Debugf("calling API %s", u)
	resp, err := httpclient.GetURL(u.String(), headers)
	if err != nil {
		log.Errorf("error getting issues comments api: %v", err)
		return nil, err
	}
	if resp.Status.Code != http.StatusOK {
		log.Warnf("Request returned: %s", string(resp.Body))
		return nil, errors.New("request returned an incorrect http.Status: " + resp.Status.Text)
	}

	var v Comments
	err = json.Unmarshal(resp.Body, &v)
	if err != nil {
		log.Errorf("error unmarshalling response from GH comments issues API: %v", err)
		return nil, err
	}

	// filtering mine
	v, err = filterMyComments(v)
	if err != nil {
		log.Errorf("error filtering comments %v", err)
		return nil, err
	}
	log.Tracef("filtered comments: %v", v)

	return v, nil
}

func filterMyComments(ghis Comments) (Comments, error) {
	b := ghis[:0]
	for _, x := range ghis {
		log.Debugf("filterMyComments() comment ID: %v", x.ID)
		if x.User.Login == ghUsername {
			log.Debugf("filterMyComments() comment belongs to me %v", x.ID)
			b = append(b, x)
		}
	}
	return b, nil
}

func filterMyIssues(ghis Issues) (Issues, error) {
	b := ghis[:0]
	for _, x := range ghis {
		log.Debugf("filterMyIssues() issues ID: %v", x.Number)
		if x.User.Login == ghUsername && x.Title == "Validation errors on publiccode.yml" {
			log.Debugf("filterMyIssues() issue belongs to me: %v", x.Number)
			b = append(b, x)
		}
	}
	return b, nil
}

func githubBasicAuth(domain Domain) string {
	if len(domain.BasicAuth) > 0 {
		auth := domain.BasicAuth[rand.Intn(len(domain.BasicAuth))]
		return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	}
	return ""
}
