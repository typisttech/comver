package comver_test

import (
	"fmt"

	"github.com/typisttech/comver"
)

func ExampleAnd() {
	a, _ := comver.And(
		comver.NewGreaterThanOrEqualTo(comver.MustParse("2")),
		comver.NewLessThan(comver.MustParse("3")),
	)

	fmt.Println(a)
	// Output: >=2 <3
}

func ExampleAnd_wildcard() {
	a, _ := comver.And(
		comver.NewGreaterThanOrEqualTo(comver.MustParse("2")),
		comver.NewLessThan(comver.MustParse("3")),
		comver.NewWildcard(),
	)

	fmt.Println(a)
	// Output: >=2 <3
}

func ExampleAnd_exactConstraint() {
	a, _ := comver.And(
		comver.NewGreaterThanOrEqualTo(comver.MustParse("2")),
		comver.NewLessThanOrEqualTo(comver.MustParse("2")),
	)

	fmt.Println(a)
	// Output: 2
}

func ExampleAnd_impossibleInterval() {
	_, err := comver.And(
		comver.NewGreaterThanOrEqualTo(comver.MustParse("3")),
		comver.NewLessThan(comver.MustParse("2")),
	)

	fmt.Println(err)
	// Output: impossible interval
}
