package main

import (
	"fmt"
	"os"

	"github.com/notBigO/wkn/repl"
	"github.com/notBigO/wkn/utils"
	"github.com/spf13/cobra"
)

func main() {
	var dbPath string

	var rootCmd = &cobra.Command{
		Use:   "wkn",
		Short: "Webknot Numbers (WKN) - a simple number-only database",
		Long:  `WKN is a minimal REPL-based database that works with integer arrays.`,
		Run: func(cmd *cobra.Command, args []string) {
			if dbPath != "" {
				arrays, err := utils.LoadFromSpecificFile(dbPath)
				if err != nil {
					fmt.Println(err)
					return
				}
				repl.ReplLoop(arrays)
			} else {
				fmt.Println("Use `wkn new` to create a new database or `wkn --db-path <path>` to load an existing one.")
			}
		},
	}

	rootCmd.Flags().StringVar(&dbPath, "db-path", "", "Path to an existing .wkn database file")

	var newCmd = &cobra.Command{
		Use:   "new",
		Short: "Create a new .wkn database file and start the REPL",
		Run: func(cmd *cobra.Command, args []string) {
			arrays, err := utils.LoadFromFile()
			if err != nil {
				fmt.Println("Error loading database:", err)
				return
			}

			repl.ReplLoop(arrays)
		},
	}

	rootCmd.AddCommand(newCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
