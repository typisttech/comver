package comver_test

import (
	"fmt"
	"github.com/typisttech/comver"
)

func ExampleNewVersion() {
	v, _ := comver.NewVersion("1.2.3")

	fmt.Println(v)
	// Output: 1.2.3.0
}

func ExampleNewVersion_full() {
	v, _ := comver.NewVersion("1.2.3.4-beta.5+foo")

	fmt.Println(v)
	// Output: 1.2.3.4-beta5
}

func ExampleNewVersion_withLeadingV() {
	v, _ := comver.NewVersion("v1.2.3.4-beta.5+foo")

	fmt.Println(v)
	// Output: 1.2.3.4-beta5
}

func ExampleNewVersion_error() {
	_, err := comver.NewVersion("not a version")

	fmt.Println(err)
	// Output: error parsing version string "not a version"
}

func ExampleNewVersion_dateCal() {
	v, _ := comver.NewVersion("2010-01-02")

	fmt.Println(v)
	// Output: 2010.1.2.0
}

func ExampleNewVersion_spacePadding() {
	v, _ := comver.NewVersion(" 1.0.0")

	fmt.Println(v)
	// Output: 1.0.0.0
}

func ExampleNewVersion_modifierShorthand() {
	v, _ := comver.NewVersion("1.2.3-b5")

	fmt.Println(v)
	// Output: 1.2.3.0-beta5
}

func ExampleNewVersion_zeroPadding() {
	v, _ := comver.NewVersion("00.01.03.04")

	fmt.Println(v)
	// Output: 0.1.3.4
}

func ExampleNewVersion_dateWithFourBits() {
	_, err := comver.NewVersion("20100102.0.3.4")

	fmt.Println(err)
	// Output: error parsing version string "20100102.0.3.4"
}

func ExampleNewVersion_invalidModifier() {
	_, err := comver.NewVersion("1.0.0-meh")

	fmt.Println(err)
	// Output: error parsing version string "1.0.0-meh"
}

func ExampleNewVersion_tooManyBits() {
	_, err := comver.NewVersion("1.0.0.0.0")

	fmt.Println(err)
	// Output: error parsing version string "1.0.0.0.0"
}

func ExampleVersion_Compare() {
	v1, _ := comver.NewVersion("1")
	v2, _ := comver.NewVersion("2")
	v3, _ := comver.NewVersion("3")

	v2v1 := v2.Compare(v1)
	fmt.Println(v2v1)

	v2v2 := v2.Compare(v2)
	fmt.Println(v2v2)

	v2v3 := v2.Compare(v3)
	fmt.Println(v2v3)

	// Output:
	// 1
	// 0
	// -1
}

func ExampleVersion_Compare_patch() {
	v1, _ := comver.NewVersion("1")
	v1p, _ := comver.NewVersion("1.patch")

	v1v1p := v1.Compare(v1p)

	fmt.Println(v1v1p)
	// Output: -1
}

func ExampleVersion_Compare_preRelease() {
	v1b5, _ := comver.NewVersion("1.0.0-beta.5")
	v1b6, _ := comver.NewVersion("1.0.0-beta.6")

	v1b5v1b6 := v1b5.Compare(v1b6)

	fmt.Println(v1b5v1b6)
	// Output: -1
}

func ExampleVersion_Short_major() {
	v, _ := comver.NewVersion("1")

	s := v.Short()

	fmt.Println(s)
	// Output: 1
}

func ExampleVersion_Short_minor() {
	v, _ := comver.NewVersion("1.2")

	s := v.Short()

	fmt.Println(s)
	// Output: 1.2
}

func ExampleVersion_Short_patch() {
	v, _ := comver.NewVersion("1.2.3")

	s := v.Short()

	fmt.Println(s)
	// Output: 1.2.3
}

func ExampleVersion_Short_tweak() {
	v, _ := comver.NewVersion("1.2.3.4")

	s := v.Short()

	fmt.Println(s)
	// Output: 1.2.3.4
}

func ExampleVersion_Short_modifier() {
	v, _ := comver.NewVersion("1.2.3.4.beta")

	s := v.Short()

	fmt.Println(s)
	// Output: 1.2.3.4-beta
}

func ExampleVersion_Short_preRelease() {
	v, _ := comver.NewVersion("1.2.3.4-beta5")

	s := v.Short()

	fmt.Println(s)
	// Output: 1.2.3.4-beta5
}

func ExampleVersion_Short_metadata() {
	v, _ := comver.NewVersion("1.2.3.4-beta5+foo")

	s := v.Short()

	fmt.Println(s)
	// Output: 1.2.3.4-beta5
}

func ExampleVersion_Original() {
	v, _ := comver.NewVersion("1.b5+foo")

	s := v.Original()

	fmt.Println(s)
	// Output: 1.b5+foo
}
