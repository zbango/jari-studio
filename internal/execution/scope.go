package execution

import (
	"errors"
	"path/filepath"
	"strings"
)

var ErrScopeViolation = errors.New("changed path is outside the card scope")

func ValidateChangedFiles(changed, allowed, forbidden []string) error {
	for _, path := range changed {
		clean := filepath.ToSlash(filepath.Clean(path))
		if clean == "." || filepath.IsAbs(path) || strings.HasPrefix(clean, "../") {
			return ErrScopeViolation
		}
		if matchesAny(clean, forbidden) || !matchesAny(clean, allowed) {
			return ErrScopeViolation
		}
	}
	return nil
}

func matchesAny(path string, patterns []string) bool {
	for _, pattern := range patterns {
		pattern = filepath.ToSlash(filepath.Clean(pattern))
		if strings.HasSuffix(pattern, "/**") && strings.HasPrefix(path, strings.TrimSuffix(pattern, "**")) {
			return true
		}
		if ok, _ := filepath.Match(pattern, path); ok {
			return true
		}
	}
	return false
}
