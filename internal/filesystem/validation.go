package filesystem

import (
	"fmt"
	"path/filepath"
	"strings"
)

const BaseDir = "/Users/xu/mv"

// ValidatePath ensures the path is within the allowed directory and returns the clean path
func ValidatePath(path string) (string, error) {
	cleanPath := filepath.Clean(path)

	if filepath.IsAbs(cleanPath) {
		if !strings.HasPrefix(cleanPath, BaseDir) {
			return "", fmt.Errorf("access denied: path outside allowed directory")
		}
		return cleanPath, nil
	}

	fullPath := filepath.Join(BaseDir, cleanPath)
	cleanFullPath := filepath.Clean(fullPath)

	if !strings.HasPrefix(cleanFullPath, BaseDir) {
		return "", fmt.Errorf("access denied: path outside allowed directory")
	}

	return cleanFullPath, nil
}