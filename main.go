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

func init() {
	log.SetLevel(log.DebugLevel)

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [ OPTIONS ] URL\n", os.Args[0])
		flag.PrintDefaults()
	}
	// remoteBaseURLPtr := flag.String("remote-url", "", "The URL pointing to the directory where the publiccode.yml file is located.")
	helpPtr := flag.Bool("help", false, "Display command line usage.")

	if *helpPtr || len(flag.Args()) < 1 {
		flag.Usage()
		return
	}
}

var e = engines.NewEngine()

func main() {
	flag.Parse()

	//init API engines
	engines.RegisterClientAPIs()
}

// Start will get go API request and populate Event struct
// - urlString is a string representing URL pointing to publiccode.yml
//   but will accept also repo url
// - valid is a bool representing publiccode validation status
// - valErrors is a string in JSON format that will be deserialized
//   it contains all validation errors
func Start(url *url.URL, valid bool, valErrors interface{}) error {
	event := model.Event{}
	event.URL = url
	event.Valid = valid
	event.ValidationError = valErrors.([]model.Error)

	e.IdentifyVCS(url)
	return nil
}

// Startf same as above but accepting a more generic
// type for handy usage
// TODO to be renamed
func Startf(urlString string, valid bool, valErrors string) error {
	log.Println("Handle event")

	urlParsed, err := url.Parse(urlString)
	if err != nil {
		return err
	}

	var verr []model.Error
	// deserialize valErrors
	err = json.Unmarshal([]byte(valErrors), &verr)
	if err != nil {
		return err
	}

	if err = Start(urlParsed, valid, verr); err != nil {
		return err
	}

	return nil
}
