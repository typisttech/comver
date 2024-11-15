package comver_test

import (
	"fmt"

	"github.com/typisttech/comver"
)

func ExampleParse() {
	v, _ := comver.Parse("1.2.3")

	fmt.Println(v)
	// Output: 1.2.3.0
}

func ExampleParse_full() {
	v, _ := comver.Parse("1.2.3.4-beta.5+foo")

	fmt.Println(v)
	// Output: 1.2.3.4-beta5
}

func ExampleParse_error() {
	_, err := comver.Parse("not a version")

	fmt.Println(err)
	// Output: error parsing version string "not a version"
}

func ExampleVersion_Compare() {
	v1 := comver.MustParse("1")
	v2 := comver.MustParse("2")
	w2 := comver.MustParse("2")
	v3 := comver.MustParse("3")

	v2v1 := v2.Compare(v1)
	fmt.Println(v2v1)

	v2w2 := v2.Compare(w2)
	fmt.Println(v2w2)

	v2v3 := v2.Compare(v3)
	fmt.Println(v2v3)

	// Output:
	// 1
	// 0
	// -1
}

func ExampleVersion_Compare_patch() {
	v1 := comver.MustParse("1")
	v1p := comver.MustParse("1.patch")

	got := v1.Compare(v1p)

	fmt.Println(got)
	// Output: -1
}

func ExampleVersion_Compare_preRelease() {
	v1b5 := comver.MustParse("1.0.0-beta.5")
	v1b6 := comver.MustParse("1.0.0-beta.6")

	got := v1b5.Compare(v1b6)

	fmt.Println(got)
	// Output: -1
}

func ExampleVersion_Compare_metadata() {
	foo := comver.MustParse("1.0.0+foo")
	bar := comver.MustParse("1.0.0+bar")

	got := foo.Compare(bar)

	fmt.Println(got)
	// Output: 0
}

func ExampleVersion_Short() {
	ss := []string{
		"1",
		"1.2",
		"1.2.3",
		"1.2.3+foo",
		"1.2.3.4",
		"1.2.3.4.beta",
		"1.2.3.4-beta5",
		"1.2.3.4-beta5+foo",
		"1.b5+foo",
	}

	for _, s := range ss {
		v := comver.MustParse(s)
		got := v.Short()

		fmt.Printf("%-19q => %v\n", s, got)
	}

	// Output:
	// "1"                 => 1
	// "1.2"               => 1.2
	// "1.2.3"             => 1.2.3
	// "1.2.3+foo"         => 1.2.3
	// "1.2.3.4"           => 1.2.3.4
	// "1.2.3.4.beta"      => 1.2.3.4-beta
	// "1.2.3.4-beta5"     => 1.2.3.4-beta5
	// "1.2.3.4-beta5+foo" => 1.2.3.4-beta5
	// "1.b5+foo"          => 1-beta5
}

func ExampleVersion_Original() {
	ss := []string{
		"1",
		"1.2",
		"1.2.3",
		"1.2.3+foo",
		"1.2.3.4",
		"1.2.3.4.beta",
		"1.2.3.4-beta5",
		"1.2.3.4-beta5+foo",
		"1.b5+foo",
	}

	for _, s := range ss {
		v := comver.MustParse(s)
		got := v.Original()

		fmt.Printf("%-19q => %v\n", s, got)
	}

	// Output:
	// "1"                 => 1
	// "1.2"               => 1.2
	// "1.2.3"             => 1.2.3
	// "1.2.3+foo"         => 1.2.3+foo
	// "1.2.3.4"           => 1.2.3.4
	// "1.2.3.4.beta"      => 1.2.3.4.beta
	// "1.2.3.4-beta5"     => 1.2.3.4-beta5
	// "1.2.3.4-beta5+foo" => 1.2.3.4-beta5+foo
	// "1.b5+foo"          => 1.b5+foo
}
