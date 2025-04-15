package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
			if dbExists(".wkn") {
				fmt.Printf(".wkn file already exists")
				return
			} else {
				file, err := os.Create(".wkn")
				if err != nil {
					fmt.Println("Error creating a file: ", err)
					return
				}
				defer file.Close()

				fmt.Println("Created .wkn file")
				replLoop()
			}
		},
	}

	rootCmd.AddCommand(newCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func parse(input string) (string, []string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", nil
	}
	return parts[0], parts[1:]
}

func replLoop() {
	for {
		print("wkn> ")
		input := readLine()
		cmd, args := parse(input)

		switch cmd {
		case "exit":
			fmt.Println("Bye!")
			return
		case "new":
			fmt.Println("Creating a new array.....")
		case "show":
			fmt.Println("Showing all arrays....")
		default:
			fmt.Printf("Error: “%s” is not a supported operation\n", cmd)
		}
	}
}

func dbExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
