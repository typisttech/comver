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
	errEmptyString             stringError = "version string is empty"
	errInvalidVersionString    stringError = "invalid version string"
	errNotFixedVersion         stringError = "not a fixed version"
	errDateVersionWithFourBits stringError = "date versions with 4 bits"
)

var (
	classicalVersioningRegexp = regexp.MustCompile(
		"^" + classicalVersioningRegex + modifierRegex + "$",
	)
	dateOnlyVersioningRegexp = regexp.MustCompile(
		"^" + dateOnlyVersioningRegex + modifierRegex + "$",
	)
)

// Version represents a single composer version.
// The zero value for Version is v0.0.0.0 with empty original string.
type Version struct {
	major, minor, patch, tweak uint64   `exhaustruct:"optional"`
	modifier                   modifier `exhaustruct:"optional"`
	preRelease                 string   `exhaustruct:"optional"`
	original                   string   `exhaustruct:"optional"`
}

// Parse parses a given version string, attempts to coerce a version string into
// a [Version] object or return an error if unable to parse the version string.
//
// If there is a leading v or a version listed without all parts (e.g. v1.2.p5+foo) it
// attempt to coerce it into a valid composer version (e.g. 1.2.0.0-patch5). In both cases
// a [Version] object is returned that can be sorted, compared, and used in constraints.
//
// Due to implementation complexity, it only supports a subset of [composer versioning].
// Refer to the [version_test.go] for examples.
//
// [composer versioning]: https://github.com/composer/semver/
// [version_test.go]: https://github.com/typisttech/comver/blob/main/version_test.go
func Parse(v string) (Version, error) { //nolint:cyclop,funlen
	original := v

	// normalize to lowercase for easier pattern matching
	v = strings.ToLower(v)

	v = strings.TrimSpace(v)
	if v == "" {
		return Version{}, &ParseError{original, errEmptyString}
	}

	v = strings.TrimPrefix(v, "v")
	if v == "" {
		return Version{}, &ParseError{original, errInvalidVersionString}
	}

	if strings.Contains(v, " as ") {
		return Version{}, &ParseError{original, errNotFixedVersion}
	}

	if hasSuffixAnyOf(v, "@stable", "@rc", "@beta", "@alpha", "@dev") {
		return Version{}, &ParseError{original, errNotFixedVersion}
	}

	if containsAnyOf(v, "master", "trunk", "default") {
		return Version{}, &ParseError{original, errNotFixedVersion}
	}

	if strings.HasPrefix(v, "dev-") {
		return Version{}, &ParseError{original, errNotFixedVersion}
	}

	// strip off build metadata
	v, metadata, _ := strings.Cut(v, "+")
	if v == "" || strings.Contains(metadata, " ") {
		return Version{}, &ParseError{original, errInvalidVersionString}
	}

	if strings.HasSuffix(v, "dev") {
		return Version{}, &ParseError{original, errNotFixedVersion}
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
		return Version{}, &ParseError{original, errInvalidVersionString}
	}

	var err error
	if cv.major, err = strconv.ParseUint(match[1], 10, 64); err != nil {
		return Version{}, &ParseError{original, err}
	}
	// CalVer (as MAJOR) must be in YYYYMMDDhhmm or YYYYMMDD formats
	if s := strconv.FormatUint(cv.major, 10); len(s) > 12 || len(s) == 11 || len(s) == 9 ||
		len(s) == 7 {
		return Version{}, &ParseError{original, errInvalidVersionString}
	}

	if cv.minor, err = strconv.ParseUint(match[2], 10, 64); match[2] != "" && err != nil {
		return Version{}, &ParseError{original, err}
	}

	if cv.patch, err = strconv.ParseUint(match[3], 10, 64); match[3] != "" && err != nil {
		return Version{}, &ParseError{original, err}
	}

	if cv.major >= 1000_00 && match[4] != "" {
		return Version{}, &ParseError{original, errDateVersionWithFourBits}
	}

	if cv.tweak, err = strconv.ParseUint(match[4], 10, 64); match[4] != "" && err != nil {
		return Version{}, &ParseError{original, err}
	}

	if cv.modifier, err = newModifier(match[5]); err != nil {
		return Version{}, &ParseError{original, err}
	}

	cv.preRelease = strings.TrimPrefix(strings.TrimPrefix(match[6], "-"), ".")

	return cv, nil
}

// MustParse is like [Parse] but panics if the version string cannot be parsed.
func MustParse(v string) Version {
	cv, err := Parse(v)
	if err != nil {
		panic(err)
	}

	return cv
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

// Short returns the shortest string representation of the version.
func (v Version) Short() string {
	s := fmt.Sprintf("%d.%d.%d.%d", v.major, v.minor, v.patch, v.tweak)

	s = strings.TrimSuffix(s, ".0")
	s = strings.TrimSuffix(s, ".0")
	s = strings.TrimSuffix(s, ".0")

	if v.modifier != modifierStable {
		s += "-" + v.modifier.String() + v.preRelease
	}

	return s
}

// Original returns the original version string passed into [Parse].
// Empty string is returned when [Version] is the zero value.
func (v Version) Original() string {
	return v.original
}

// Compare returns an integer comparing two [Version] instances.
//
// Pre-release versions are compared according to [semantic version precedence].
// The result is 0 when v == w, -1 when v < w, or +1 when v > w.
//
// [semantic version precedence]: https://semver.org/#spec-item-11
func (v Version) Compare(w Version) int { //nolint:cyclop,funlen
	switch {
	case v.String() == w.String():
		return 0
	case v.major > w.major:
		return +1
	case v.major < w.major:
		return -1
	case v.minor > w.minor:
		return +1
	case v.minor < w.minor:
		return -1
	case v.patch > w.patch:
		return +1
	case v.patch < w.patch:
		return -1
	case v.tweak > w.tweak:
		return +1
	case v.tweak < w.tweak:
		return -1
	case v.modifier > w.modifier:
		return +1
	case v.modifier < w.modifier:
		return -1
	case v.preRelease != "" && w.preRelease == "":
		return +1
	case v.preRelease == "" && w.preRelease != "":
		return -1
	}

	vPres := strings.Split(v.preRelease, ".")
	wPres := strings.Split(w.preRelease, ".")

	// comparing each dot separated identifier from ceiling to floor
	for i := range vPres {
		// a larger set of pre-release fields has a higher precedence than a smaller set
		if i >= len(wPres) {
			return +1
		}

		vi, wi := vPres[i], wPres[i]
		if vi == wi {
			continue
		}

		vid := isDigits(vi)
		wid := isDigits(wi)

		// identifiers consisting of only digits are compared numerically
		if vid && wid {
			vii, _ := strconv.ParseUint(vi, 10, 64)
			wii, _ := strconv.ParseUint(wi, 10, 64)

			if vii > wii {
				return +1
			}

			return -1
		}

		//nolint:godox
		// TODO: Find out whether composer/semver supports this
		//
		// identifiers with letters or hyphens are compared lexically in ASCII sort order
		if !vid && !wid {
			if vi > wi {
				return +1
			}

			return -1
		}

		//nolint:godox
		// TODO: Find out whether composer/semver supports this
		//
		// numeric identifiers always have floor precedence than non-numeric identifiers
		if !vid && wid {
			return +1
		}

		return -1
	}

	// a larger set of pre-release fields has a higher precedence than a smaller set
	return -1
}

func isDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}
