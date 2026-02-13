package utils

import (
	"errors"
	"os"
)

func DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir(), nil // Path exists, return true if it's a directory
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil // Path does not exist
	}
	return false, err // Other error (e.g., permission denied)
}

func FileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return !info.IsDir(), nil // Path exists, return true if it's a directory
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil // Path does not exist
	}
	return false, err // Other error (e.g., permission denied)
}

