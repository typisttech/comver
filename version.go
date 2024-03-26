package comver

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	classicalVersioningRegex               = `(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:\.(\d+))?`
	dateOnlyVersioningRegex                = `(\d{4})(?:[.:-]?(\d{2}))(?:[.:-]?(\d{2}))?(?:\.(\d+))?`
	modifierRegex                          = `[._-]?(?:(stable|beta|b|rc|alpha|a|patch|pl|p)((?:[.-]?\d+)+)?)?`
	ErrEmptyString             stringError = "version string is empty"
	ErrInvalidVersionString    stringError = "invalid version string"
	ErrNotFixedVersion         stringError = "not a fixed version"
	ErrDateVersionWithFourBits stringError = "date versions with 4 bits"
)

var classicalVersioningRegexp = regexp.MustCompile("^" + classicalVersioningRegex + modifierRegex + "$")
var dateOnlyVersioningRegexp = regexp.MustCompile("^" + dateOnlyVersioningRegex + modifierRegex + "$")

// Version represents a single composer version.
type Version struct {
	major, minor, patch, tweak uint64
	modifier                   modifier
	preRelease                 string
	original                   string
}

// NewVersion parses a given version string and returns an instance of [Version] or
// an error if unable to parse the version.
func NewVersion(v string) (Version, error) {
	original := v

	// normalize to lowercase for easier pattern matching
	v = strings.ToLower(v)

	v = strings.TrimSpace(v)
	if v == "" {
		return Version{}, &ParseError{original, ErrEmptyString}
	}

	v = strings.TrimPrefix(v, "v")
	if v == "" {
		return Version{}, &ParseError{original, ErrInvalidVersionString}
	}

	if strings.Contains(v, " as ") {
		return Version{}, &ParseError{original, ErrNotFixedVersion}
	}

	if hasSuffixAnyOf(v, "@stable", "@rc", "@beta", "@alpha", "@dev") {
		return Version{}, &ParseError{original, ErrNotFixedVersion}
	}

	if containsAnyOf(v, "master", "trunk", "default") {
		return Version{}, &ParseError{original, ErrNotFixedVersion}
	}

	if strings.HasPrefix(v, "dev-") {
		return Version{}, &ParseError{original, ErrNotFixedVersion}
	}

	// strip off build metadata
	v, metadata, _ := strings.Cut(v, "+")
	if v == "" || strings.Contains(metadata, " ") {
		return Version{}, &ParseError{original, ErrInvalidVersionString}
	}

	if strings.HasSuffix(v, "dev") {
		return Version{}, &ParseError{original, ErrNotFixedVersion}
	}

	cv := Version{
		original: original,
	}

	var match []string

	if cm := classicalVersioningRegexp.FindStringSubmatch(v); cm != nil {
		match = cm
	} else if dm := dateOnlyVersioningRegexp.FindStringSubmatch(v); dm != nil {
		match = dm
	}
	if match == nil || len(match) != 7 {
		return Version{}, &ParseError{original, ErrInvalidVersionString}
	}

	var err error

	if cv.major, err = strconv.ParseUint(match[1], 10, 64); err != nil {
		return Version{}, &ParseError{original, err}
	}
	// CalVer (as MAJOR) must be in YYYYMMDDhhmm or YYYYMMDD formats
	if s := strconv.FormatUint(cv.major, 10); len(s) > 12 || len(s) == 11 || len(s) == 9 || len(s) == 7 {
		return Version{}, &ParseError{original, ErrInvalidVersionString}
	}

	if cv.minor, err = strconv.ParseUint(match[2], 10, 64); match[2] != "" && err != nil {
		return Version{}, &ParseError{original, err}
	}

	if cv.patch, err = strconv.ParseUint(match[3], 10, 64); match[3] != "" && err != nil {
		return Version{}, &ParseError{original, err}
	}

	if cv.major >= 1000_00 && match[4] != "" {
		return Version{}, &ParseError{original, ErrDateVersionWithFourBits}
	}
	if cv.tweak, err = strconv.ParseUint(match[4], 10, 64); match[4] != "" && err != nil {
		return Version{}, &ParseError{original, err}
	}

	if cv.modifier, err = newStability(match[5]); err != nil {
		return Version{}, &ParseError{original, err}
	}

	cv.preRelease = strings.TrimPrefix(strings.TrimPrefix(match[6], "-"), ".")

	return cv, nil
}

func hasSuffixAnyOf(s string, suffixes ...string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}
	return false
}

func containsAnyOf(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// String returns the normalized string representation of the version.
func (v Version) String() string {
	s := fmt.Sprintf("%d.%d.%d.%d", v.major, v.minor, v.patch, v.tweak)

	if v.modifier != modifierStable {
		s += "-" + v.modifier.String() + v.preRelease
	}

	return s
}

func (v Version) Original() string {
	return v.original
}
