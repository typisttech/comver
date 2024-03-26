package comver

import (
	"errors"
	"testing"
)

func TestNewVersion(t *testing.T) {
	tests := []struct {
		name string
		v    string
		want string
	}{
		// Taken from composer/semver VersionParserTest::successfulNormalizedVersions()
		// https://github.com/composer/semver/blob/1d09200268e7d1052ded8e5da9c73c96a63d18f5/tests/VersionParserTest.php#L65-L142
		{"none", "1.0.0", "1.0.0.0"},
		{"none/2", "1.2.3.4", "1.2.3.4"},
		{"RC uppercase", "1.0.0-rc1", "1.0.0.0-RC1"},
		{"forces w.x.y.z/2", "0", "0.0.0.0"},
		{"forces w.x.y.z/maximum major", "99999", "99999.0.0.0"}, // https://github.com/composer/semver/pull/158
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
		{"parses CalVer YYYYMMDDhhmm (as MAJOR) versions", "202301310000.0.0", "202301310000.0.0.0"},
		{"strips v/datetime", "v20100102", "20100102.0.0.0"},
		{"parses dates no delimiter", "20100102", "20100102.0.0.0"},
		{"parses dates no delimiter/2", "20100102.0", "20100102.0.0.0"},
		{"parses dates no delimiter/3", "20100102.1.0", "20100102.1.0.0"},
		{"parses dates no delimiter/4", "20100102.0.3", "20100102.0.3.0"},
		{"parses dates no delimiter/earliest year", "100000", "100000.0.0.0"}, // https://github.com/composer/semver/pull/158
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

		// additional tests
		{"parses dates y-m", "2010-01", "2010.1.0.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//t.Parallel()

			got, err := NewVersion(tt.v)
			if err != nil {
				t.Fatalf("NewVersion() error = %v, wantErr %v", err, nil)
			}
			if gotString := got.String(); gotString != tt.want {
				t.Errorf("NewVersion().String() got = %v, want %v", gotString, tt.want)
			}
			if gotOriginal := got.original; gotOriginal != tt.v {
				t.Errorf("NewVersion().original got = %v, want %v", gotOriginal, tt.v)
			}
		})
	}
}

func TestNewVersion_ParseError(t *testing.T) {
	tests := []struct {
		name    string
		v       string
		wantErr error
	}{
		// composer/semver supports a lot of different version formats, but we only support a subset of them.
		// taken from composer/semver VersionParserTest::successfulNormalizedVersions()
		// https://github.com/composer/semver/blob/1d09200268e7d1052ded8e5da9c73c96a63d18f5/tests/VersionParserTest.php#L65-L142
		{"parses state", "1.0.0RC1dev", ErrNotFixedVersion},
		{"CI parsing", "1.0.0-rC15-dev", ErrNotFixedVersion},
		{"delimiters", "1.0.0.RC.15-dev", ErrNotFixedVersion},
		{"patch replace", "1.0.0.pl3-dev", ErrNotFixedVersion},
		{"forces w.x.y.z", "1.0-dev", ErrNotFixedVersion},
		{"parses dates w/ - and .", "2010-01-02-10-20-30.0.3", ErrInvalidVersionString},
		{"parses dates w/ - and ./2", "2010-01-02-10-20-30.5", ErrInvalidVersionString},
		{"parses datetime", "20100102-203040", ErrInvalidVersionString},
		{"parses date dev", "20100102.x-dev", ErrNotFixedVersion},
		{"parses datetime dev", "20100102.203040.x-dev", ErrNotFixedVersion},
		{"parses dt+number", "20100102203040-10", ErrInvalidVersionString},
		{"parses dt+patch", "20100102-203040-p1", ErrInvalidVersionString},
		{"parses dt Ym dev", "201903.x-dev", ErrNotFixedVersion},
		{"parses master", "dev-master", ErrNotFixedVersion},
		{"parses master w/o dev", "master", ErrNotFixedVersion},
		{"parses trunk", "dev-trunk", ErrNotFixedVersion},
		{"parses branches", "1.x-dev", ErrNotFixedVersion},
		{"parses arbitrary", "dev-feature-foo", ErrNotFixedVersion},
		{"parses arbitrary/2", "DEV-FOOBAR", ErrNotFixedVersion},
		{"parses arbitrary/3", "dev-feature/foo", ErrNotFixedVersion},
		{"parses arbitrary/4", "dev-feature+issue-1", ErrNotFixedVersion},
		{"ignores aliases", "dev-master as 1.0.0", ErrNotFixedVersion},
		{"ignores aliases/2", "dev-load-varnish-only-when-used as ^2.0", ErrNotFixedVersion},
		{"ignores aliases/3", "dev-load-varnish-only-when-used@dev as ^2.0@dev", ErrNotFixedVersion},
		{"ignores stability", "1.0.0+foo@dev", ErrNotFixedVersion},
		{"ignores stability/2", "dev-load-varnish-only-when-used@stable", ErrNotFixedVersion},
		{"semver metadata/7", "1.0.0-0.3.7", ErrInvalidVersionString},    // composer/semver doesn't support this
		{"semver metadata/8", "1.0.0-x.7.z.92", ErrInvalidVersionString}, // composer/semver doesn't support this
		{"metadata w/ alias", "1.0.0+foo as 2.0", ErrNotFixedVersion},
		{"keep zero-padding/5", "041.x-dev", ErrNotFixedVersion},
		{"keep zero-padding/6", "dev-041.003", ErrNotFixedVersion},
		{"dev with mad name", "dev-1.0.0-dev<1.0.5-dev", ErrNotFixedVersion},
		{"dev prefix with spaces", "dev-foo bar", ErrNotFixedVersion},

		// composer/semver doesn't support these.
		// taken from composer/semver VersionParserTest::failingNormalizedVersions()
		// https://github.com/composer/semver/blob/1d09200268e7d1052ded8e5da9c73c96a63d18f5/tests/VersionParserTest.php#L158-L183
		{"empty", "", ErrEmptyString},
		{"invalid chars", "a", ErrInvalidVersionString},
		{"invalid type", "1.0.0-meh", ErrInvalidVersionString},
		{"too many bits", "1.0.0.0.0", ErrInvalidVersionString},
		{"non-dev arbitrary", "feature-foo", ErrInvalidVersionString},
		{"metadata w/ space", "1.0.0+foo bar", ErrInvalidVersionString},
		{"maven style release", "1.0.1-SNAPSHOT", ErrInvalidVersionString},
		{"dev with less than", "1.0.0<1.0.5-dev", ErrNotFixedVersion},
		{"dev with less than/2", "1.0.0-dev<1.0.5-dev", ErrNotFixedVersion},
		{"dev suffix with spaces", "foo bar-dev", ErrNotFixedVersion},
		{"any with spaces", "1.0 .2", ErrInvalidVersionString},
		{"no version, no alias", " as ", ErrInvalidVersionString},
		{"no version, only alias", " as 1.2", ErrInvalidVersionString},
		{"just an operator", "^", ErrInvalidVersionString},
		{"just an operator/2", "^8 || ^", ErrInvalidVersionString},
		{"just an operator/3", "~", ErrInvalidVersionString},
		{"just an operator/4", "~1 ~", ErrInvalidVersionString},
		{"constraint", "~1", ErrInvalidVersionString},
		{"constraint/2", "^1", ErrInvalidVersionString},
		{"constraint/3", "1.*", ErrInvalidVersionString},
		{"date versions with 4 bits", "20100102.0.3.4", ErrDateVersionWithFourBits},
		{"date versions with 4 bits/earliest year", "100000.0.0.0", ErrDateVersionWithFourBits},
		{"invalid CalVer (as MAJOR) versions/YYYYMMD", "2023013.0.0", ErrInvalidVersionString},
		{"invalid CalVer (as MAJOR) versions/YYYYMMDDh", "202301311.0.0", ErrInvalidVersionString},
		{"invalid CalVer (as MAJOR) versions/YYYYMMDDhhm", "20230131000.0.0", ErrInvalidVersionString},
		{"invalid CalVer (as MAJOR) versions/YYYYMMDDhhmmX", "2023013100000.0.0", ErrInvalidVersionString},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewVersion(tt.v)

			if err == nil {
				t.Fatalf("NewVersion() got = %s error = %v, wantErr %v", got, err, tt.wantErr)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewVersion() error = %#v, wantErr %#v", err, tt.wantErr)
			}

			var wantParseError *ParseError
			if !errors.As(err, &wantParseError) {
				t.Fatalf("NewVersion() error = %#v, wantErr %#v", err, wantParseError)
			}

			if wantParseError.original != tt.v {
				t.Errorf("NewVersion() error.original = %v, want %v", wantParseError.original, tt.v)
			}
		})
	}
}
