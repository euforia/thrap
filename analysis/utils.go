package analysis

import (
	"strings"
)

// EstimateLanguage estimates the programming language by percent of file types
// present in dir.  If less than 50% an empty string is returned
func EstimateLanguage(dir string) string {
	fts := BuildFileTypeSpread(dir)
	highest := fts.Highest()
	if highest != nil && highest.Percent > 50 {
		return strings.ToLower(highest.Language)
	}
	return ""
}
