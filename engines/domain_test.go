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
}
