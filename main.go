package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "wkn",
		Short: "Webknot Numbers (WKN) - a simple number-only database",
		Long:  `WKN is a minimal REPL-based database that works with integer arrays.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Use `wkn new` to create a new database or `wkn --db-path <path>` to load an existing one.")
		},
	}

	var newCmd = &cobra.Command{
		Use:   "new",
		Short: "Create a new .wkn database file and start the REPL",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Creating new .wkn database...")
		},
	}

	rootCmd.AddCommand(newCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
