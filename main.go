package main

import (
	"net/url"

	"github.com/sebbalex/issue-opener/cmd"
	"github.com/sebbalex/issue-opener/engines"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	//init API engines
	engines.RegisterClientAPIs()

	// starting CLI
	cmd.Execute()
}

// Start will get go API request and populate Event struct
// - urlString is a string representing URL pointing to publiccode.yml
//   but will accept also repo url
// - valid is a bool representing publiccode validation status
// - valErrors is a string in JSON format that will be deserialized
//   it contains all validation errors
func Start(url *url.URL, valid bool, valErrors interface{}) error {
	e := engines.NewEngine()
	return e.Start(url, valid, valErrors)
}
