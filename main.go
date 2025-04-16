package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/notBigO/wkn/repl"
	"github.com/notBigO/wkn/utils"
	"github.com/spf13/cobra"
)

var arrays = make(map[string][]int)

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
			if utils.DbExists(".wkn") {
				data, err := os.ReadFile(".wkn")
				if err == nil {
					fmt.Println(".wkn file already exists. Loading...")
					json.Unmarshal(data, &arrays)
				} else {
					fmt.Println("Failed to read .wkn file:", err)
				}
			} else {
				file, err := os.Create(".wkn")
				if err != nil {
					fmt.Println("Error creating a file: ", err)
					return
				}
				defer file.Close()
				fmt.Println("Created .wkn file")
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
