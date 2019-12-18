package main

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/sebbalex/issue-opener/model"
	"github.com/stretchr/testify/assert"
)

var urlString string = "https://raw.githubusercontent.com/sebbalex/issue-opener/master/publiccode.yml"
var valid bool = false
var valErrors string = `[
		{"key": "name", "reason": "missing mandatory key"}, 
		{"key": "localisation_ready", "reason": "missing mandatory key"}
	]`

// TODO disable logging in stdout while testing

func TestStart(t *testing.T) {
	urlParsed, err := url.Parse(urlString)
	if err != nil {
		t.Errorf("error on parsing url %s", err)
	}

	var verr []model.Error
	// deserialize valErrors
	err = json.Unmarshal([]byte(valErrors), &verr)
	if err != nil {
		t.Errorf("error on unmarsalling validation errors %s", err)
	}

	assert.Nil(t, Start(urlParsed, valid, verr))
}
