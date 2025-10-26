package comver_test

import (
	"fmt"

	"github.com/typisttech/comver"
)

func Example_version() {
	ss := []string{
		"1.2.3",
		"1.2",
		"1",

		"   1.0.0",
		"00.01.03.04",

		"2010-01-02.5",
		"2010-01-02",

		"v1.2.3.4-beta.5+foo",
		"v1.2.3.4.p5+foo",
		"v1.2.3",
		"v1.2.p5+foo",

		"not a version",
		"1.0.0-alpha.beta",
		"1.0.0-meh",
		"1.0.0.0.0",
		"20100102.0.3.4",
	}

	for _, s := range ss {
		v, err := comver.Parse(s)
		if err != nil {
			fmt.Printf("%-21q => %v\n", s, err)

			continue
		}

		fmt.Printf("%-21q => %v\n", s, v)
	}

	// Output:
	// "1.2.3"               => 1.2.3.0
	// "1.2"                 => 1.2.0.0
	// "1"                   => 1.0.0.0
	// "   1.0.0"            => 1.0.0.0
	// "00.01.03.04"         => 0.1.3.4
	// "2010-01-02.5"        => 2010.1.2.5
	// "2010-01-02"          => 2010.1.2.0
	// "v1.2.3.4-beta.5+foo" => 1.2.3.4-beta5
	// "v1.2.3.4.p5+foo"     => 1.2.3.4-patch5
	// "v1.2.3"              => 1.2.3.0
	// "v1.2.p5+foo"         => 1.2.0.0-patch5
	// "not a version"       => error parsing version string "not a version"
	// "1.0.0-alpha.beta"    => error parsing version string "1.0.0-alpha.beta"
	// "1.0.0-meh"           => error parsing version string "1.0.0-meh"
	// "1.0.0.0.0"           => error parsing version string "1.0.0.0.0"
	// "20100102.0.3.4"      => error parsing version string "20100102.0.3.4"
}
