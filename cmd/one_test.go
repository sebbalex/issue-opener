package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var urlString string = "https://raw.githubusercontent.com/AgID/pat/master/publiccode.yml"
var valid bool = false
var valErrors string = `[
		{"key": "name", "reason": "missing mandatory key"}, 
		{"key": "localisation_ready", "reason": "missing mandatory key"}
	]`

func TestStartf(t *testing.T) {
	assert.Nil(t, StartCLI(urlString, valid, valErrors))
}
