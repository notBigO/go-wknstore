package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

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
			if dbExists(".wkn") {
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

			replLoop()
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

func saveToFile() error {
	data, err := json.Marshal(arrays)
	if err != nil {
		return err
	}
	return os.WriteFile(".wkn", data, 0644)
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
			if len(args) < 1 {
				fmt.Println("Usage: new <array_name> [num1 num2 ...]")
				break
			}

			name := args[0]
			if _, exists := arrays[name]; exists {
				fmt.Printf("Array '%s' already exists\n", name)
				break
			}

			nums := []int{}
			for _, arg := range args[1:] {
				num, err := strconv.Atoi(arg)
				if err != nil {
					fmt.Printf("Error: Invalid number: %s\n", arg)
					continue
				}
				nums = append(nums, num)
			}

			arrays[name] = nums

			if err := saveToFile(); err != nil {
				fmt.Println("Error: Failed to save to .wkn:", err)
				break
			}

			fmt.Printf("CREATED (%d)\n", len(nums))

		case "show":
			if len(arrays) == 0 {
				fmt.Println("No arrays present")
				break
			}

			if len(args) == 1 {
				name := args[0]
				if arr, ok := arrays[name]; ok {
					fmt.Printf("%s: %v\n", name, arr)
				} else {
					fmt.Printf("Error: Array '%s' not found\n", name)
				}
				break
			}

			for name, data := range arrays {
				fmt.Printf("%s: %v\n", name, data)
			}

		case "del":
			if len(args) < 1 {
				fmt.Println("Usage: del <array_name>")
				break
			}

			_, exists := arrays[args[0]]
			if !exists {
				fmt.Printf("Error: Array '%s' does not exist\n", args[0])
				break
			}

			delete(arrays, args[0])

			if err := saveToFile(); err != nil {
				fmt.Println("Error: Failed to save to .wkn:", err)
				break
			}

			fmt.Println("DELETED")

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
