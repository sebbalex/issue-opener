package engines

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"

	. "github.com/sebbalex/issue-opener/model"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// Domain is a single code hosting service.
type Domain struct {
	// Domains.yml data
	Host      string   `yaml:"host"`
	BasicAuth []string `yaml:"basic-auth"`
}

// API returns a Domain without tld.
func (domain Domain) API() string {
	truncateIndex := strings.LastIndexAny(domain.Host, ".")
	// It is already an API without tld.
	if truncateIndex == -1 {
		return domain.Host
	}
	return domain.Host[:truncateIndex]
}

func (domain Domain) mapDomainForAuth(dom []Domain) (Domain, error) {
	for _, d := range dom {
		if domain.Host == d.API() {
			d.Host = d.API() // hack to have API endpoint (withoud tld) in domain
			return d, nil
		}
	}
	return Domain{}, errors.New("No domain registered in domains.yml for " + domain.Host)
}

func (domain Domain) processPostOrAppendIssue(event *Event) error {
	if len(event.Message) == 0 {
		log.Errorf("message is not present %v", event)
		return errors.New("message not present")
	}
	message := event.Message[0]
	var engine SingleRepoHandler
	var err error
	if message.Append {
		engine, err = GetAppendIssueClientAPIEngine(domain.API())
	} else {
		engine, err = GetPostIssueClientAPIEngine(domain.API())
	}

	if err != nil {
		return err
	}
	return engine(domain, event)
}

func (domain Domain) processSingleRepo(event *Event) error {
	engine, err := GetSingleClientAPIEngine(domain.API())
	if err != nil {
		return err
	}
	return engine(domain, event)
}

// ReadAndParseDomains read domainsFile and return the parsed content in a Domain slice.
func ReadAndParseDomains(domainsFile string) ([]Domain, error) {
	// Getting absolute path to be called from different packages/test files
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	// Open and read domains file list.
	data, err := ioutil.ReadFile(basepath + "/../" + domainsFile)
	if err != nil {
		return nil, fmt.Errorf("error in reading %s file: %v", domainsFile, err)
	}
	// Parse domains file list.
	domains, err := parseDomainsFile(data)
	if err != nil {
		return nil, fmt.Errorf("error in parsing %s file: %v", domainsFile, err)
	}
	log.Infof("Loaded and parsed %s", domainsFile)

	return domains, err
}

// parseDomainsFile parses the domains file to build a slice of Domain.
func parseDomainsFile(data []byte) ([]Domain, error) {
	domains := []Domain{}

	// Unmarshal the yml in domains list.
	err := yaml.Unmarshal(data, &domains)
	if err != nil {
		log.Errorf("error %v", err)
		return nil, err
	}
	return domains, err
}
