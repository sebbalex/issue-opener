package main

import (
	"net/url"
	"os"

	"github.com/sebbalex/issue-opener/cmd"
	"github.com/sebbalex/issue-opener/engines"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as TEXT with fulltimestamp support
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)

	//init API engines
	engines.RegisterClientAPIs()
}

func main() {
	// starting CLI
	cmd.Execute()
}

// Start will get go API request and populate Event struct
// - urlString is a string representing URL pointing to publiccode.yml
//   but will accept also repo url
// - valid is a bool representing publiccode validation status
// - valErrors is a string in JSON format that will be deserialized
//   it contains all validation errors
func Start(url *url.URL, valid bool, valErrors interface{}, dryRun bool) error {
	e := engines.NewEngine()
	return e.Start(url, valid, valErrors, dryRun)
}
