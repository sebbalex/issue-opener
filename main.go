package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/sebbalex/issue-opener/engines"
	"github.com/sebbalex/issue-opener/model"
	log "github.com/sirupsen/logrus"
)

var e = engines.NewEngine()

func main() {
	log.SetLevel(log.DebugLevel)

	//init API engines
	engines.RegisterClientAPIs()

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [ OPTIONS ] URL\n", os.Args[0])
		flag.PrintDefaults()
	}
	repoURL := flag.String("repo-url", "", "The URL pointing to repository which contains the publiccode.yml.")
	validErrors := flag.String("validation", "[]", "JSON representing validation errors array.")
	helpPtr := flag.Bool("help", false, "Display command line usage.")
	flag.Parse()

	if *helpPtr || len(flag.Args()) < 1 {
		flag.Usage()
		return
	}
	if *repoURL != "" {
		StartCLI(*repoURL, true, *validErrors)
	}
}

// Start will get go API request and populate Event struct
// - urlString is a string representing URL pointing to publiccode.yml
//   but will accept also repo url
// - valid is a bool representing publiccode validation status
// - valErrors is a string in JSON format that will be deserialized
//   it contains all validation errors
func Start(url *url.URL, valid bool, valErrors interface{}) error {
	log.Debug("starting...")
	event := model.Event{}
	event.URL = url
	event.Valid = valid
	event.ValidationError = valErrors.([]model.Error)
	log.Debugf("on: %s", event)

	d, err := e.IdentifyVCS(url)
	e.StartFlow(url, d)
	return err
}

// StartCLI same as above but accepting a more generic
// type for handy usage
// TODO to be renamed
func StartCLI(urlString string, valid bool, valErrors string) error {
	log.Println("Handle event")

	urlParsed, err := url.Parse(urlString)
	if err != nil {
		log.Fatal(err)
		return err
	}

	var verr []model.Error
	// deserialize valErrors
	err = json.Unmarshal([]byte(valErrors), &verr)
	if err != nil {
		log.Fatal(err)
		return err
	}

	if err = Start(urlParsed, valid, verr); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
