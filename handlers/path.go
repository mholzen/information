package handlers

import (
	"os"
	"path/filepath"
)

func splitAndCheck(p string) (string, []string, error) {
	// Normalize the path to ensure it's clean
	p = filepath.Clean(p)

	// Split the path into components
	components := filepath.SplitList(p)

	// Start with the root or current directory
	var currentPath string
	if filepath.IsAbs(p) {
		currentPath = string(filepath.Separator)
	}

	// Iterate over the components and build the path
	for i, component := range components {
		nextPath := filepath.Join(currentPath, component)
		if _, err := os.Stat(nextPath); os.IsNotExist(err) {
			// If path doesn't exist, return the current path and the remaining path
			return currentPath, components[i:], nil
		}
		currentPath = nextPath
	}

	// If we reach here, the full path exists
	return p, nil, nil
}
