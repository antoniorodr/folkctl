/*
Copyright © 2026 ANTONIO RODRIGUEZ <kontakt@antoniorodriguez.no>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "folkctl",
	Short: "CLI tool for Folkomaten",
	Long: `Folkctl is a CLI tool for finding and copying national identity numbers
(fødselsnummer) for BankID test users that also exist in the Norwegian
National Population Register (DSF) test database.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
