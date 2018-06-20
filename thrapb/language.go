package thrapb

import (
	"errors"
	"strings"

	version "github.com/hashicorp/go-version"
)

var (
	errLangNotSpecified = errors.New("language not specified")
)

// LanguageID is the programming language for a component actually being
// developed. It contains the name and version
type LanguageID string

// Validate validates the language and its version to ensure both are supplied
func (lang LanguageID) Validate() error {
	if lang == "" {
		return errLangNotSpecified
	}
	vers := lang.Version()
	if vers == "" {
		return nil
	}

	_, err := version.NewVersion(vers)
	return err
}

// Lang returns the language name without the version.  It returns an empty
// string if the language format is invalid
func (lang LanguageID) Lang() string {
	v := strings.Split(string(lang), ":")
	return v[0]
}

// Version returns the language version from the string of format <lang>:<version>
// If the format is invalid or version is not specified, an empty string is
// returned
func (lang LanguageID) Version() string {
	v := strings.Split(string(lang), ":")
	if len(v) == 2 {
		return v[1]
	}

	return ""
}
