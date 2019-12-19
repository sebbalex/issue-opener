package engines

import (
	"errors"
	"fmt"
	"github.com/sebbalex/issue-opener/analyzer"
	"net/url"

	vcs "github.com/alranel/go-vcsurl"
	validator "github.com/asaskevich/govalidator"
	. "github.com/sebbalex/issue-opener/model"
	log "github.com/sirupsen/logrus"
)

// Engine is a helper class representing an engine.
type Engine struct {
	domains  []Domain
	comments Comments
}

// Repository is a single code repository. FileRawURL contains the direct url to the raw file.
type Repository struct {
	Name        string
	Hostname    string
	FileRawURL  string
	GitCloneURL string
	GitBranch   string
	Domain      Domain
	Headers     map[string]string
	Metadata    []byte
}

// ClientAPI contains all the API function in a single Client.
type ClientAPI struct {
	Single      SingleRepoHandler
	PostIssue   SingleRepoHandler
	AppendIssue SingleRepoHandler
}

var clientAPIs map[string]ClientAPI

// RegisterClientAPIs register all the client APIs for all the clients.
func RegisterClientAPIs() {

	clientAPIs = make(map[string]ClientAPI)

	clientAPIs["bitbucket"] = ClientAPI{
		// Single: RegisterSingleBitbucketAPI(),
	}

	clientAPIs["github"] = ClientAPI{
		Single:      RegisterSingleGithubAPI(),
		PostIssue:   RegisterPostIssueGithubAPI(),
		AppendIssue: RegisterAppendIssueGithubAPI(),
	}

	clientAPIs["gitlab"] = ClientAPI{
		// Single: RegisterSingleGitlabAPI(),
	}

}

// NewEngine istance
func NewEngine() *Engine {
	var e Engine
	var err error
	// Read and parse list of domains.
	e.domains, err = ReadAndParseDomains("domains.yml")
	if err != nil {
		log.Fatal(err)
	}

	return &e
}

// Start will get go API request and populate Event struct
// - urlString is a string representing URL pointing to publiccode.yml
//   but will accept also repo url
// - valid is a bool representing publiccode validation status
// - valErrors is a string in JSON format that will be deserialized
//   it contains all validation errors
func (e *Engine) Start(url *url.URL, valid bool, valErrors interface{}, dryRun bool) error {
	if !validator.IsURL(url.String()) {
		return errors.New("Error parsing url, please specify a good one")
	}
	log.Debug("starting...")
	event := Event{}
	event.URL = url
	event.Valid = valid
	event.DryRun = dryRun
	event.ValidationError = valErrors.([]Error)

	log.Debugf("on: %v", event)

	d, err := e.IdentifyVCS(url)
	if err != nil {
		return err
	}

	*d, err = d.mapDomainForAuth(e.domains)
	if err != nil {
		return err
	}
	err = e.StartFlow(&event, d)
	if err != nil {
		return err
	}

	err = analyzer.CompareMessages(&event)
	if err != nil {
		return err
	}

	err = d.processPostOrAppendIssue(&event)
	if err != nil {
		return err
	}

	return nil
}

// StartFlow ..
func (e *Engine) StartFlow(event *Event, d *Domain) error {
	url := event.URL
	if vcs.IsRawFile(url) {
		event.URL = vcs.GetRepo(url)
	}
	return d.processSingleRepo(event)
}

// IdentifyVCS Will identify which VCS platform come
// the request and address it through correct engine
// this will handle vcs recognition and initiate with correct
// engine
func (e *Engine) IdentifyVCS(url *url.URL) (*Domain, error) {
	if !validator.IsURL(url.String()) {
		return &Domain{Host: "none"}, errors.New("Error parsing url, please specify a good one")
	}
	if vcs.IsBitBucket(url) {
		return &Domain{Host: "bitbucket"}, errors.New("Not yet implemented")
	} else if vcs.IsGitLab(url) {
		return &Domain{Host: "gitlab"}, errors.New("Not yet implemented")
	} else if vcs.IsGitHub(url) {
		return &Domain{Host: "github"}, nil
	} else {
		return &Domain{Host: "none"}, errors.New("Not yet implemented")
	}
}

// GetSingleClientAPIEngine checks if the API client for the requested
// single repository clientAPI exists and return its handler.
func GetSingleClientAPIEngine(clientAPI string) (SingleRepoHandler, error) {
	if clientAPIs[clientAPI].Single != nil {
		return clientAPIs[clientAPI].Single, nil
	}
	return nil, fmt.Errorf("no single client found for %s", clientAPI)
}

// GetPostIssueClientAPIEngine checks if the API client
// for the requested Post Issue clientAPI exists and return its handler.
func GetPostIssueClientAPIEngine(clientAPI string) (SingleRepoHandler, error) {
	if clientAPIs[clientAPI].PostIssue != nil {
		return clientAPIs[clientAPI].PostIssue, nil
	}
	return nil, fmt.Errorf("no PostIssue client found for %s", clientAPI)
}

// GetAppendIssueClientAPIEngine checks if the API client
// for the requested Post Issue clientAPI exists and return its handler.
func GetAppendIssueClientAPIEngine(clientAPI string) (SingleRepoHandler, error) {
	if clientAPIs[clientAPI].AppendIssue != nil {
		return clientAPIs[clientAPI].AppendIssue, nil
	}
	return nil, fmt.Errorf("no AppendIssue client found for %s", clientAPI)
}
