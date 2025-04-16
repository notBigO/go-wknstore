package repl

import (
	"fmt"
	"strconv"

	"github.com/notBigO/wkn/utils"
)

func ReplLoop(arrays map[string][]int) {
	for {
		print("wkn> ")
		input := utils.ReadLine()
		cmd, args := utils.Parse(input)

		// for commands that need fresh data  reload from file first
		if cmd == "show" || cmd == "new" || cmd == "merge" || cmd == "pow" || cmd == "del" {
			freshArrays, err := utils.LoadFromFile()
			if err == nil {
				// update the arrays map with fresh data
				arrays = freshArrays
			} else {
				fmt.Println("Error refreshing data:", err)
			}
		}

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

			if err := utils.SaveToFile(arrays); err != nil {
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
					fmt.Println(arr)
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
				fmt.Printf("Error: %s does not exist\n", target)
				break
			}
			if !ok2 {
				fmt.Printf("Error: %s does not exist\n", source)
				break
			}

			arrays[target] = append(tArr, sArr...)
			if err := utils.SaveToFile(arrays); err != nil {
				fmt.Println("Error: Failed to save to .wkn:", err)
				break
			}
			fmt.Println("MERGED")

		case "pow":
			if len(args) < 2 {
				fmt.Println("Usage: pow <arrayA.indexA> <arrayB.indexB>")
				break
			}
			base, err1 := utils.GetValueFromReference(args[0], arrays)
			exp, err2 := utils.GetValueFromReference(args[1], arrays)

			if err1 != nil {
				fmt.Println("Error:", err1)
				break
			}
			if err2 != nil {
				fmt.Println("Error:", err2)
				break
			}

			fmt.Println(utils.BinaryExponentiation(base, exp))

		case "del":
			if len(args) < 1 {
				fmt.Println("Usage: del <array_name>")
				break
			}

			_, exists := arrays[args[0]]
			if !exists {
				fmt.Printf("Error: %s does not exist\n", args[0])
				break
			}

			delete(arrays, args[0])

			if err := utils.SaveToFile(arrays); err != nil {
				fmt.Println("Error: Failed to save to .wkn:", err)
				break
			}

			fmt.Println("DELETED")

		default:
			fmt.Printf("Error: %s is not a supported operation\n", cmd)
		}
	}
}
