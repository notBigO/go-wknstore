package utils

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	dbFile     = ".wkn"
	lockFile   = ".wkn.lock"
	maxRetries = 5
	retryDelay = 200 * time.Millisecond
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

// tryigg to create a lock file
func acquireLock() error {
	for i := 0; i < maxRetries; i++ {
		file, err := os.OpenFile(lockFile, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err == nil {
			file.Close()
			return nil
		}

		// check if file exists error (lock is held by someone else)
		if os.IsExist(err) {
			// wait and retry
			time.Sleep(retryDelay)
			continue
		}

		return err
	}

	return fmt.Errorf("failed to acquire lock after %d attempts", maxRetries)
}

// remove the lock file
func releaseLock() error {
	return os.Remove(lockFile)
}

// load arrays from the database file with locking
func LoadFromFile() (map[string][]int, error) {
	if err := acquireLock(); err != nil {
		return nil, fmt.Errorf("couldn't acquire lock: %v", err)
	}
	defer releaseLock()

	data, err := os.ReadFile(dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			// if the file doesn't exist, return an empty map
			return make(map[string][]int), nil
		}
		return nil, err
	}

	// empty fiile case
	if len(data) == 0 {
		return make(map[string][]int), nil
	}

	var arrays map[string][]int
	if err := json.Unmarshal(data, &arrays); err != nil {
		return make(map[string][]int), nil
	}

	return arrays, nil
}

// save arrays to the database file with locking
func SaveToFile(arrays map[string][]int) error {
	if err := acquireLock(); err != nil {
		return fmt.Errorf("couldn't acquire lock: %v", err)
	}
	defer releaseLock()

	data, err := json.Marshal(arrays)
	if err != nil {
		return err
	}
	return os.WriteFile(dbFile, data, 0644)
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
		return 0, fmt.Errorf("%s does not exist", name)
	}
	if idx < 0 || idx >= len(arr) {
		return 0, fmt.Errorf("%s index out of bounds for %d", name, idx)
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
