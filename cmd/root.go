package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "issue-opener",
	Short: "Issue Opener bot for malformed publiccode.yml",
	Long: `Issue Opener bot for malformed publiccode.yml.
Complete documentation is available at https://github.com/sebbalex/issue-opener`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			log.Fatal(err)
		}
	},
}

// Execute is the entrypoint for cmd package Cobra.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
