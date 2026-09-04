/*
Copyright © 2026 ANTONIO RODRIGUEZ <kontakt@antoniorodriguez.no>
*/
package cmd

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
)

var VERSION = "unknown"
var ORIGINAL_DATABASE = "resources/testbrukere.txt"
var DEFAULT_DATABASE = getDefaultDatabase()

func getDefaultDatabase() string {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		fmt.Println("Error resolving resource path:", err)
		return ""
	}

	defaultDatabase := filepath.Join(homeDir, ".config/folkctl/databases/testbrukere.txt")
	return defaultDatabase
}

type Config struct {
	ActiveDatabase string `toml:"active_database"`
}

var rootCmd = &cobra.Command{
	Use:   "folkctl",
	Short: "CLI tool for Folkomaten",
	Long: `Folkctl is a CLI tool for finding and copying national identity numbers
(fødselsnummer) for BankID test users that also exist in the Norwegian
National Population Register (DSF) test database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		databasePath, err := cmd.Flags().GetString("database")

		if err != nil {
			fmt.Println("Error reading flag:", err)
			return err
		}

		if cmd.Flags().Changed("database") {
			if databasePath == "" {
				return fmt.Errorf("database path cannot be empty")
			}
			createDefaultConfig(databasePath)
			fmt.Printf("Config updated: active database set to %s\n", databasePath)
			return nil
		}

		cmd.Help()
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func loadConfig() string {
	configDir, err := getConfigPath()

	if err != nil {
		fmt.Println("Error resolving resource path:", err)
		return ""
	}

	configFile := filepath.Join(configDir, "config.toml")

	var config Config

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		fmt.Println("Error reading config file:", err)
		return ""
	}

	return config.ActiveDatabase
}

func createDefaultConfig(activeDatabase string) {
	configDir, err := getConfigPath()
	if err != nil {
		fmt.Println("Error resolving resource path:", err)
		return
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Println("Error creating config directory:", err)
		return
	}

	configFile := filepath.Join(configDir, "config.toml")

	err = os.WriteFile(
		configFile,
		[]byte(fmt.Sprintf("active_database = %q\n", activeDatabase)),
		0644,
	)

	if err != nil {
		fmt.Println("Error creating config file:", err)
	}
}

func createDefaultDatabase() {
	configDir, err := getConfigPath()

	if err != nil {
		fmt.Println("Error resolving resource path:", err)
		return
	}

	databaseDir := filepath.Join(configDir, "databases")
	databaseFile := filepath.Join(databaseDir, "testbrukere.txt")

	if err := os.MkdirAll(databaseDir, 0755); err != nil {
		fmt.Println("Error creating database directory:", err)
		return
	}

	if _, err := os.Stat(databaseFile); err == nil {
		return
	} else if !os.IsNotExist(err) {
		fmt.Println("Error checking database:", err)
		return
	}

	source, err := os.Open(ORIGINAL_DATABASE)

	if err != nil {
		fmt.Println("Error opening default database:", err)
		return
	}

	defer source.Close()

	destination, err := os.Create(databaseFile)

	if err != nil {
		fmt.Println("Error creating database:", err)
		return
	}

	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		fmt.Println("Error copying default database:", err)
		return
	}
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	configDir, err := filepath.Abs(homeDir + "/.config/folkctl/")

	if err != nil {
		fmt.Println("Error resolving resource path:", err)
		return "", err
	}

	return configDir, err
}

func init() {
	rootCmd.Flags().StringP("database", "d", "", "change the 'active database' from the config to the path specified")
	rootCmd.Version = VERSION
	rootCmd.SetVersionTemplate("folkctl version {{.Version}}\n")

	configPath, _ := getConfigPath()
	configFile := filepath.Join(configPath, "config.toml")

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		createDefaultConfig(DEFAULT_DATABASE)
	}

	createDefaultDatabase()
}
