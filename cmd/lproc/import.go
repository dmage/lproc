package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/dmage/lproc/pkg/config"
	"github.com/dmage/lproc/pkg/importcsv"
)

var importArgs struct {
	FormatName string
}

var importCmd = &cobra.Command{
	Use:   "import [flags] csv-file...",
	Short: "Import transactions from a CSV file",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rules, err := config.LoadDefaultRules()
		if err != nil {
			log.Fatal(err)
		}

		format, err := rules.GetFormat(importArgs.FormatName)
		if err != nil {
			log.Fatal(err)
		}

		for _, arg := range args {
			if err := importcsv.ImportCSV(arg, format, rules.Classifiers); err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	importCmd.Flags().StringVarP(&importArgs.FormatName, "format-name", "n", "default", "format name to use from the configuration file")
	rootCmd.AddCommand(importCmd)
}
