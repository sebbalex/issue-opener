package engines

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ghodss/yaml"
	. "github.com/sebbalex/issue-opener/model"
	log "github.com/sirupsen/logrus"
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
