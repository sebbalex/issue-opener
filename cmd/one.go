package cmd

import (
	"encoding/json"
	"net/url"

	"github.com/sebbalex/issue-opener/engines"
	. "github.com/sebbalex/issue-opener/model"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(oneCmd)
}

var oneCmd = &cobra.Command{
	Use:   "one [repo url] [json validation errors]",
	Short: "The URL pointing to repository which contains the publiccode.yml.",
	Long: `The URL pointing to repository which contains the publiccode.yml.
No organizations! Only single repositories!`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		validationErrors := "[]"
		if len(args) > 1 && args[1] != "" {
			validationErrors = args[1]
		}
		StartCLI(url, false, validationErrors)
	},
}

// StartCLI same as above but accepting a more generic
// type for handy usage
// TODO to be renamed
func StartCLI(urlString string, valid bool, valErrors string) error {
	log.Println("Handle event")
	var e = engines.NewEngine()

	urlParsed, err := url.Parse(urlString)
	if err != nil {
		log.Fatal(err)
		return err
	}

	var verr []Error
	// deserialize valErrors
	err = json.Unmarshal([]byte(valErrors), &verr)
	if err != nil {
		log.Fatal(err)
		return err
	}

	if err = e.Start(urlParsed, valid, verr); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
