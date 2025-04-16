package main

import (
	"bufio"
	"encoding/json"
	"errors"
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
					fmt.Printf("Error: '%s' does not exist\n", name)
				}
				break
			}

			for name, data := range arrays {
				fmt.Printf("%s: %v\n", name, data)
			}

		case "merge":
			if len(args) < 2 {
				fmt.Println("Usage: merge <target_array> <source_array>")
				break
			}
			target, source := args[0], args[1]
			tArr, ok1 := arrays[target]
			sArr, ok2 := arrays[source]

			if !ok1 {
				fmt.Printf("Error: “%s” does not exist\n", target)
				break
			}
			if !ok2 {
				fmt.Printf("Error: “%s” does not exist\n", source)
				break
			}

			arrays[target] = append(tArr, sArr...)
			if err := saveToFile(); err != nil {
				fmt.Println("Error: Failed to save to .wkn:", err)
				break
			}
			fmt.Println("MERGED")

		case "pow":
			if len(args) < 2 {
				fmt.Println("Usage: pow <arrayA.indexA> <arrayB.indexB>")
				break
			}
			base, err1 := getValueFromReference(args[0])
			exp, err2 := getValueFromReference(args[1])

			if err1 != nil {
				fmt.Println("Error:", err1)
				break
			}
			if err2 != nil {
				fmt.Println("Error:", err2)
				break
			}

			fmt.Println(binaryExponentiation(base, exp))

		case "del":
			if len(args) < 1 {
				fmt.Println("Usage: del <array_name>")
				break
			}

			_, exists := arrays[args[0]]
			if !exists {
				fmt.Printf("Error:  “%s” does not exist\n", args[0])
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

func getValueFromReference(ref string) (int, error) {
	parts := strings.Split(ref, ".")
	if len(parts) != 2 {
		return 0, errors.New("Invalid format, expected array.index")
	}
	name := parts[0]
	idx, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, errors.New("Invalid index")
	}
	arr, ok := arrays[name]
	if !ok {
		return 0, fmt.Errorf("“%s” does not exist", name)
	}
	if idx < 0 || idx >= len(arr) {
		return 0, fmt.Errorf("“%s” index out of bounds for %d", name, idx)
	}
	return arr[idx], nil
}

func binaryExponentiation(base, exp int) int {
	result := 1
	mod := int(1e9 + 7)
	base %= mod
	for exp > 0 {
		if exp%2 == 1 {
			result = (result * base) % mod
		}
		base = (base * base) % mod
		exp /= 2
	}
	return result
}
