package utils

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadLine() string {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func Parse(input string) (string, []string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", nil
	}
	return parts[0], parts[1:]
}

func SaveToFile(arrays map[string][]int) error {
	data, err := json.Marshal(arrays)
	if err != nil {
		return err
	}
	return os.WriteFile(".wkn", data, 0644)
}

func DbExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func GetValueFromReference(ref string, arrays map[string][]int) (int, error) {
	parts := strings.Split(ref, ".")
	if len(parts) != 2 {
		return 0, errors.New("invalid format, expected array.index")
	}
	name := parts[0]
	idx, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, errors.New("invalid index")
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

func BinaryExponentiation(base, exp int) int {
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
