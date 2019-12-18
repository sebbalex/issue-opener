package engines

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadAndParseDomains(t *testing.T) {
	domainsFile := "domains.yml"
	domains, err := ReadAndParseDomains(domainsFile)
	if err != nil {
		t.Fatalf("error on %s", err)
	}
	assert.Len(t, domains, 3)

	for _, d := range domains {
		assert.NotEmpty(t, d.API())
	}
	dom, err := Domain{Host: "gitlab"}.mapDomainForAuth(domains)
	assert.NotEmpty(t, dom.Host)
	assert.NoError(t, err)

	dom, err = Domain{Host: "gisdtlab"}.mapDomainForAuth(domains)
	assert.Empty(t, dom.Host)
	assert.Error(t, err)
}
