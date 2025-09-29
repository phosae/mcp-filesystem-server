package filesystem

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Validator provides path validation for a specific base directory
type Validator struct {
	baseDir string
}

// NewValidator creates a new path validator for the given base directory
func NewValidator(baseDir string) *Validator {
	// Clean the base directory path
	cleanBaseDir := filepath.Clean(baseDir)
	return &Validator{
		baseDir: cleanBaseDir,
	}
}

// ValidatePath ensures the path is within the allowed directory and returns the clean path
func (v *Validator) ValidatePath(path string) (string, error) {
	cleanPath := filepath.Clean(path)

	if filepath.IsAbs(cleanPath) {
		if !strings.HasPrefix(cleanPath, v.baseDir) {
			return "", fmt.Errorf("access denied: path outside allowed directory %s", v.baseDir)
		}
		return cleanPath, nil
	}

	fullPath := filepath.Join(v.baseDir, cleanPath)
	cleanFullPath := filepath.Clean(fullPath)

	if !strings.HasPrefix(cleanFullPath, v.baseDir) {
		return "", fmt.Errorf("access denied: path outside allowed directory %s", v.baseDir)
	}

	return cleanFullPath, nil
}

// GetBaseDir returns the base directory for this validator
func (v *Validator) GetBaseDir() string {
	return v.baseDir
}