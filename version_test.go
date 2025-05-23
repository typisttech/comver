package comver

import (
	"errors"
	"testing"
)

func goodVersionTestCases() []struct {
	name string
	v    string
	want string
} {
	return []struct {
		name string
		v    string
		want string
	}{
		// taken from composer/semver VersionParserTest::successfulNormalizedVersions()
		// https://github.com/composer/semver/blob/1d09200268e7d1052ded8e5da9c73c96a63d18f5/tests/VersionParserTest.php#L65-L142
		{"none", "1.0.0", "1.0.0.0"},
		{"none/2", "1.2.3.4", "1.2.3.4"},
		{"RC uppercase", "1.0.0-rc1", "1.0.0.0-RC1"},
		{"forces w.x.y.z/2", "0", "0.0.0.0"},
		{
			"forces w.x.y.z/maximum major",
			"99999",
			"99999.0.0.0",
		}, // https://github.com/composer/semver/pull/158
		{"parses long", "10.4.13-beta", "10.4.13.0-beta"},
		{"parses long/2", "10.4.13beta2", "10.4.13.0-beta2"},
		{"parses long/semver", "10.4.13beta.2", "10.4.13.0-beta2"},
		{"parses long/semver2", "v1.13.11-beta.0", "1.13.11.0-beta0"},
		{"parses long/semver3", "1.13.11.0-beta0", "1.13.11.0-beta0"},
		{"expand shorthand", "10.4.13-b", "10.4.13.0-beta"},
		{"expand shorthand/2", "10.4.13-b5", "10.4.13.0-beta5"},
		{"strips leading v", "v1.0.0", "1.0.0.0"},
		{"parses dates y-m as classical", "2010.01", "2010.1.0.0"},
		{"parses dates w/ . as classical", "2010.01.02", "2010.1.2.0"},
		{"parses dates y.m.Y as classical", "2010.1.555", "2010.1.555.0"},
		{"parses dates y.m.Y/2 as classical", "2010.10.200", "2010.10.200.0"},
		{"parses CalVer YYYYMMDD (as MAJOR) versions", "20230131.0.0", "20230131.0.0.0"},
		{
			"parses CalVer YYYYMMDDhhmm (as MAJOR) versions",
			"202301310000.0.0",
			"202301310000.0.0.0",
		},
		{"strips v/datetime", "v20100102", "20100102.0.0.0"},
		{"parses dates no delimiter", "20100102", "20100102.0.0.0"},
		{"parses dates no delimiter/2", "20100102.0", "20100102.0.0.0"},
		{"parses dates no delimiter/3", "20100102.1.0", "20100102.1.0.0"},
		{"parses dates no delimiter/4", "20100102.0.3", "20100102.0.3.0"},
		{
			"parses dates no delimiter/earliest year",
			"100000",
			"100000.0.0.0",
		}, // https://github.com/composer/semver/pull/158
		{"parses dates w/ -", "2010-01-02", "2010.1.2.0"},
		{"parses dates w/ .", "2012.06.07", "2012.6.7.0"},
		{"parses numbers", "2010-01-02.5", "2010.1.2.5"},
		{"parses dates y.m.Y", "2010.1.555", "2010.1.555.0"},
		{"parses dt Ym", "201903.0", "201903.0.0.0"},
		{"parses dt Ym+patch", "201903.0-p2", "201903.0.0.0-patch2"},
		{"semver metadata/2", "1.0.0-beta.5+foo", "1.0.0.0-beta5"},
		{"semver metadata/3", "1.0.0+foo", "1.0.0.0"},
		{"semver metadata/4", "1.0.0-alpha.3.1+foo", "1.0.0.0-alpha3.1"},
		{"semver metadata/5", "1.0.0-alpha2.1+foo", "1.0.0.0-alpha2.1"},
		{"semver metadata/6", "1.0.0-alpha-2.1-3+foo", "1.0.0.0-alpha2.1-3"},
		{"keep zero-padding", "00.01.03.04", "0.1.3.4"},
		{"keep zero-padding/2", "000.001.003.004", "0.1.3.4"},
		{"keep zero-padding/3", "0.000.103.204", "0.0.103.204"},
		{"keep zero-padding/4", "0700", "700.0.0.0"},
		{"space padding", " 1.0.0", "1.0.0.0"},
		{"space padding/2", "1.0.0 ", "1.0.0.0"},

		// taken from https://semver.org/#spec-item-11
		{"semver pre-release/1", "1.0.0-alpha", "1.0.0.0-alpha"},
		{"semver pre-release/2", "1.0.0-alpha.1", "1.0.0.0-alpha1"},
		{"semver pre-release/3", "1.0.0-beta", "1.0.0.0-beta"},
		{"semver pre-release/4", "1.0.0-beta.2", "1.0.0.0-beta2"},
		{"semver pre-release/5", "1.0.0-beta.11", "1.0.0.0-beta11"},
		{"semver pre-release/6", "1.0.0-rc.1", "1.0.0.0-RC1"},

		// additional tests
		{"parses dates y-m", "2010-01", "2010.1.0.0"},
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	for _, tt := range goodVersionTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Parse(tt.v)
			if err != nil {
				t.Fatalf("Parse() error = %v, wantErr %v", err, nil)
			}
			if gotString := got.String(); gotString != tt.want {
				t.Errorf("Parse().String() got = %q, want %v", gotString, tt.want)
			}
			if gotOriginal := got.Original(); gotOriginal != tt.v {
				t.Errorf("Parse().Original() got = %q, want %v", gotOriginal, tt.v)
			}
		})
	}
}

func TestMustParse(t *testing.T) {
	t.Parallel()

	for _, tt := range goodVersionTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := MustParse(tt.v)

			if gotString := got.String(); gotString != tt.want {
				t.Errorf("MustParse().String() got = %q, want %v", gotString, tt.want)
			}
			if gotOriginal := got.Original(); gotOriginal != tt.v {
				t.Errorf("MustParse().Original() got = %q, want %v", gotOriginal, tt.v)
			}
		})
	}
}

func badVersionTestCases() []struct {
	name    string
	v       string
	wantErr error
} {
	return []struct {
		name    string
		v       string
		wantErr error
	}{
		// composer/semver supports a lot of different version formats, but we only support a subset of them
		// taken from composer/semver VersionParserTest::successfulNormalizedVersions()
		// https://github.com/composer/semver/blob/1d09200268e7d1052ded8e5da9c73c96a63d18f5/tests/VersionParserTest.php#L65-L142
		{"parses state", "1.0.0RC1dev", errNotFixedVersion},
		{"CI parsing", "1.0.0-rC15-dev", errNotFixedVersion},
		{"delimiters", "1.0.0.RC.15-dev", errNotFixedVersion},
		{"patch replace", "1.0.0.pl3-dev", errNotFixedVersion},
		{"forces w.x.y.z", "1.0-dev", errNotFixedVersion},
		{"parses dates w/ - and .", "2010-01-02-10-20-30.0.3", errInvalidVersionString},
		{"parses dates w/ - and ./2", "2010-01-02-10-20-30.5", errInvalidVersionString},
		{"parses datetime", "20100102-203040", errInvalidVersionString},
		{"parses date dev", "20100102.x-dev", errNotFixedVersion},
		{"parses datetime dev", "20100102.203040.x-dev", errNotFixedVersion},
		{"parses dt+number", "20100102203040-10", errInvalidVersionString},
		{"parses dt+patch", "20100102-203040-p1", errInvalidVersionString},
		{"parses dt Ym dev", "201903.x-dev", errNotFixedVersion},
		{"parses master", "dev-master", errNotFixedVersion},
		{"parses master w/o dev", "master", errNotFixedVersion},
		{"parses trunk", "dev-trunk", errNotFixedVersion},
		{"parses branches", "1.x-dev", errNotFixedVersion},
		{"parses arbitrary", "dev-feature-foo", errNotFixedVersion},
		{"parses arbitrary/2", "DEV-FOOBAR", errNotFixedVersion},
		{"parses arbitrary/3", "dev-feature/foo", errNotFixedVersion},
		{"parses arbitrary/4", "dev-feature+issue-1", errNotFixedVersion},
		{"ignores aliases", "dev-master as 1.0.0", errNotFixedVersion},
		{"ignores aliases/2", "dev-load-varnish-only-when-used as ^2.0", errNotFixedVersion},
		{
			"ignores aliases/3",
			"dev-load-varnish-only-when-used@dev as ^2.0@dev",
			errNotFixedVersion,
		},
		{"ignores stability", "1.0.0+foo@dev", errNotFixedVersion},
		{"ignores stability/2", "dev-load-varnish-only-when-used@stable", errNotFixedVersion},
		{
			"semver metadata/7",
			"1.0.0-0.3.7",
			errInvalidVersionString,
		}, // composer/semver doesn't support this
		{
			"semver metadata/8",
			"1.0.0-x.7.z.92",
			errInvalidVersionString,
		}, // composer/semver doesn't support this
		{"metadata w/ alias", "1.0.0+foo as 2.0", errNotFixedVersion},
		{"keep zero-padding/5", "041.x-dev", errNotFixedVersion},
		{"keep zero-padding/6", "dev-041.003", errNotFixedVersion},
		{"dev with mad name", "dev-1.0.0-dev<1.0.5-dev", errNotFixedVersion},
		{"dev prefix with spaces", "dev-foo bar", errNotFixedVersion},

		// composer/semver doesn't support these
		// taken from composer/semver VersionParserTest::failingNormalizedVersions()
		// https://github.com/composer/semver/blob/1d09200268e7d1052ded8e5da9c73c96a63d18f5/tests/VersionParserTest.php#L158-L183
		{"empty", "", errEmptyString},
		{"invalid chars", "a", errInvalidVersionString},
		{"invalid type", "1.0.0-meh", errInvalidVersionString},
		{"too many bits", "1.0.0.0.0", errInvalidVersionString},
		{"non-dev arbitrary", "feature-foo", errInvalidVersionString},
		{"metadata w/ space", "1.0.0+foo bar", errInvalidVersionString},
		{"maven style release", "1.0.1-SNAPSHOT", errInvalidVersionString},
		{"dev with less than", "1.0.0<1.0.5-dev", errNotFixedVersion},
		{"dev with less than/2", "1.0.0-dev<1.0.5-dev", errNotFixedVersion},
		{"dev suffix with spaces", "foo bar-dev", errNotFixedVersion},
		{"any with spaces", "1.0 .2", errInvalidVersionString},
		{"no version, no alias", " as ", errInvalidVersionString},
		{"no version, only alias", " as 1.2", errInvalidVersionString},
		{"just an operator", "^", errInvalidVersionString},
		{"just an operator/2", "^8 || ^", errInvalidVersionString},
		{"just an operator/3", "~", errInvalidVersionString},
		{"just an operator/4", "~1 ~", errInvalidVersionString},
		{"constraint", "~1", errInvalidVersionString},
		{"constraint/2", "^1", errInvalidVersionString},
		{"constraint/3", "1.*", errInvalidVersionString},
		{"date versions with 4 bits", "20100102.0.3.4", errDateVersionWithFourBits},
		{"date versions with 4 bits/earliest year", "100000.0.0.0", errDateVersionWithFourBits},
		{"invalid CalVer (as MAJOR) versions/YYYYMMD", "2023013.0.0", errInvalidVersionString},
		{"invalid CalVer (as MAJOR) versions/YYYYMMDDh", "202301311.0.0", errInvalidVersionString},
		{
			"invalid CalVer (as MAJOR) versions/YYYYMMDDhhm",
			"20230131000.0.0",
			errInvalidVersionString,
		},
		{
			"invalid CalVer (as MAJOR) versions/YYYYMMDDhhmmX",
			"2023013100000.0.0",
			errInvalidVersionString,
		},

		// composer/semver doesn't support these.
		// taken from https://semver.org/#spec-item-11
		{"incompatible semver", "1.0.0-alpha.beta", errInvalidVersionString},
	}
}

func TestParse_ParseError(t *testing.T) {
	t.Parallel()

	for _, tt := range badVersionTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Parse(tt.v)
			if err == nil {
				t.Fatalf("Parse() got = %s error = %v, wantErr %v", got, err, tt.wantErr)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Parse() error = %#v, wantErr %#v", err, tt.wantErr)
			}

			var wantParseError *ParseError
			if !errors.As(err, &wantParseError) {
				t.Fatalf("Parse() error = %#v, wantErr %#v", err, wantParseError)
			}

			if wantParseError.original != tt.v {
				t.Errorf("Parse() error.original = %v, want %v", wantParseError.original, tt.v)
			}
		})
	}
}

func TestMustParse_ParseError(t *testing.T) {
	t.Parallel()

	for _, tt := range badVersionTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got Version

			defer func() {
				err := recover()
				if err == nil {
					t.Fatalf("MustParse() got = %s panic = %v, wantErr %v", got, err, tt.wantErr)
				}

				e, ok := err.(error)
				if !ok {
					t.Fatalf(
						"MustParse() doesn't panic with error got = %s panic = %v, wantErr %v",
						got,
						err,
						tt.wantErr,
					)
				}

				if !errors.Is(e, tt.wantErr) {
					t.Fatalf("MustParse() got = %s panic = %v, wantErr %v", got, err, tt.wantErr)
				}

				var wantParseError *ParseError
				if !errors.As(e, &wantParseError) {
					t.Fatalf("MustParse() panic = %#v, wantErr %#v", err, wantParseError)
				}

				if wantParseError.original != tt.v {
					t.Errorf(
						"MustParse() error.original = %v, want %v",
						wantParseError.original,
						tt.v,
					)
				}
			}()

			got = MustParse(tt.v)
		})
	}
}

func TestVersion_Compare(t *testing.T) {
	t.Parallel()

	tests := []struct {
		v    string
		w    string
		want int
	}{
		{"1", "1", 0},
		{"1", "2", -1},
		{"1.2", "1.2", 0},
		{"1.2", "1.3", -1},
		{"1.2.3", "1.2.3", 0},
		{"1.2.3", "1.2.4", -1},
		{"1.2.3.4", "1.2.3.4", 0},
		{"1.2.3.4", "1.2.3.5", -1},
		{"1.2.3.4-beta1", "1.2.3.4-beta2", -1},
		{"1.2.3.4-beta1.1", "1.2.3.4-beta2", -1},
		{"1.2.3.4-beta2", "1.2.3.4-beta11", -1},
		{"1.2.3.4-beta2.22", "1.2.3.4-beta11", -1},

		{"1-alpha", "1-beta", -1},
		{"1-beta", "1-RC", -1},
		{"1-RC", "1", -1},
		{"1", "1-patch", -1},
		{"1-patch2", "1-patch11", -1},

		// taken from https://semver.org/#spec-item-11
		{"1.0.0-alpha", "1.0.0-alpha.1", -1},
		{"1.0.0-alpha.1", "1.0.0-beta", -1},
		{"1.0.0-beta", "1.0.0-beta.2", -1},
		{"1.0.0-beta.2", "1.0.0-beta.11", -1},
		{"1.0.0-beta.11", "1.0.0-rc.1", -1},
		{"1.0.0-rc.1", "1.0.0", -1},

		// taken from composer/semver ComparatorTest::compareProvider()
		// https://github.com/composer/semver/blob/43f8029888dd52d01df0fc6d0d98d4024ab5bef1/tests/ComparatorTest.php#L213
		{"1.25.0-beta2.1", "1.25.0-b.3", -1},
		{"1.25.0-b2.1", "1.25.0beta.3", -1},
		{"1.25.0-b-2.1", "1.25.0-rc", -1},
	}
	for _, tt := range tests {
		t.Run(tt.v+"<=>"+tt.w, func(t *testing.T) {
			t.Parallel()

			v, err := Parse(tt.v)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v, wantErr %v", tt.v, err, nil)
			}
			w, err := Parse(tt.w)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v, wantErr %v", tt.w, err, nil)
			}

			if got := v.Compare(w); got != tt.want {
				t.Errorf("%q.compare(%q) = %v, want %v", tt.v, tt.w, got, tt.want)
			}

			if got := w.Compare(v); got != -1*tt.want {
				t.Errorf("%q.compare(%q) = %v, want %v", tt.w, tt.v, got, tt.want)
			}
		})
	}
}

func TestVersion_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		v    string
		want string
	}{
		{"1", "1.0.0.0"},
		{"1.0", "1.0.0.0"},
		{"1.0.0", "1.0.0.0"},
		{"1.0.0.0", "1.0.0.0"},
		{"1.2", "1.2.0.0"},
		{"1.2.0", "1.2.0.0"},
		{"1.2.0.0", "1.2.0.0"},
		{"1.2.3", "1.2.3.0"},
		{"1.2.3.0", "1.2.3.0"},
		{"1.2.3.4", "1.2.3.4"},

		{"1-beta", "1.0.0.0-beta"},
		{"1.0-beta", "1.0.0.0-beta"},
		{"1.0.0-beta", "1.0.0.0-beta"},
		{"1.0.0.0-beta", "1.0.0.0-beta"},
		{"1.2-beta", "1.2.0.0-beta"},
		{"1.2.0-beta", "1.2.0.0-beta"},
		{"1.2.0.0-beta", "1.2.0.0-beta"},
		{"1.2.3-beta", "1.2.3.0-beta"},
		{"1.2.3.0-beta", "1.2.3.0-beta"},
		{"1.2.3.4-beta", "1.2.3.4-beta"},

		{"1-beta0", "1.0.0.0-beta0"},
		{"1.0-beta0", "1.0.0.0-beta0"},
		{"1.0.0-beta0", "1.0.0.0-beta0"},
		{"1.0.0.0-beta0", "1.0.0.0-beta0"},
		{"1.2-beta0", "1.2.0.0-beta0"},
		{"1.2.0-beta0", "1.2.0.0-beta0"},
		{"1.2.0.0-beta0", "1.2.0.0-beta0"},
		{"1.2.3-beta0", "1.2.3.0-beta0"},
		{"1.2.3.0-beta0", "1.2.3.0-beta0"},
		{"1.2.3.4-beta0", "1.2.3.4-beta0"},

		{"1-beta999", "1.0.0.0-beta999"},
		{"1.0-beta999", "1.0.0.0-beta999"},
		{"1.0.0-beta999", "1.0.0.0-beta999"},
		{"1.0.0.0-beta999", "1.0.0.0-beta999"},
		{"1.2-beta999", "1.2.0.0-beta999"},
		{"1.2.0-beta999", "1.2.0.0-beta999"},
		{"1.2.0.0-beta999", "1.2.0.0-beta999"},
		{"1.2.3-beta999", "1.2.3.0-beta999"},
		{"1.2.3.0-beta999", "1.2.3.0-beta999"},
		{"1.2.3.4-beta999", "1.2.3.4-beta999"},
	}
	for _, tt := range tests {
		t.Run(tt.v, func(t *testing.T) {
			t.Parallel()

			v, err := Parse(tt.v)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v, wantErr %v", tt.v, err, nil)
			}

			if got := v.String(); got != tt.want {
				t.Errorf("%q.String() = %q, want %q", tt.v, got, tt.want)
			}
		})
	}
}

func TestVersion_Short(t *testing.T) {
	t.Parallel()

	tests := []struct {
		v    string
		want string
	}{
		{"1", "1"},
		{"1.0", "1"},
		{"1.0.0", "1"},
		{"1.0.0.0", "1"},
		{"1.2", "1.2"},
		{"1.2.0", "1.2"},
		{"1.2.0.0", "1.2"},
		{"1.2.3", "1.2.3"},
		{"1.2.3.0", "1.2.3"},
		{"1.2.3.4", "1.2.3.4"},

		{"1-beta", "1-beta"},
		{"1.0-beta", "1-beta"},
		{"1.0.0-beta", "1-beta"},
		{"1.0.0.0-beta", "1-beta"},
		{"1.2-beta", "1.2-beta"},
		{"1.2.0-beta", "1.2-beta"},
		{"1.2.0.0-beta", "1.2-beta"},
		{"1.2.3-beta", "1.2.3-beta"},
		{"1.2.3.0-beta", "1.2.3-beta"},
		{"1.2.3.4-beta", "1.2.3.4-beta"},

		{"1-beta0", "1-beta0"},
		{"1.0-beta0", "1-beta0"},
		{"1.0.0-beta0", "1-beta0"},
		{"1.0.0.0-beta0", "1-beta0"},
		{"1.2-beta0", "1.2-beta0"},
		{"1.2.0-beta0", "1.2-beta0"},
		{"1.2.0.0-beta0", "1.2-beta0"},
		{"1.2.3-beta0", "1.2.3-beta0"},
		{"1.2.3.0-beta0", "1.2.3-beta0"},
		{"1.2.3.4-beta0", "1.2.3.4-beta0"},

		{"1-beta999", "1-beta999"},
		{"1.0-beta999", "1-beta999"},
		{"1.0.0-beta999", "1-beta999"},
		{"1.0.0.0-beta999", "1-beta999"},
		{"1.2-beta999", "1.2-beta999"},
		{"1.2.0-beta999", "1.2-beta999"},
		{"1.2.0.0-beta999", "1.2-beta999"},
		{"1.2.3-beta999", "1.2.3-beta999"},
		{"1.2.3.0-beta999", "1.2.3-beta999"},
		{"1.2.3.4-beta999", "1.2.3.4-beta999"},
	}
	for _, tt := range tests {
		t.Run(tt.v, func(t *testing.T) {
			t.Parallel()

			v, err := Parse(tt.v)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v, wantErr %v", tt.v, err, nil)
			}

			if got := v.Short(); got != tt.want {
				t.Errorf("%q.Short() = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestVersion_zero(t *testing.T) {
	t.Parallel()

	v := Version{}

	if got := v.String(); got != "0.0.0.0" {
		t.Errorf("Version{}.String() = %q, want %q", got, "0.0.0.0")
	}

	if got := v.Short(); got != "0" {
		t.Errorf("Version{}.Short() = %q, want %q", got, "0")
	}

	if got := v.Original(); got != "" {
		t.Errorf("Version{}.Original() = %q, want %q", got, "")
	}
}
