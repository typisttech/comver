package comver_test

import (
	"fmt"

	"github.com/typisttech/comver"
)

func ExampleCompact() {
	o := comver.Or{
		comver.MustAnd(
			comver.NewLessThan(comver.MustParse("2")),
			comver.NewGreaterThan(comver.MustParse("1")),
		),
		comver.MustAnd(
			comver.NewLessThan(comver.MustParse("5")),
			comver.NewGreaterThan(comver.MustParse("3")),
		),
		comver.MustAnd(
			comver.NewLessThan(comver.MustParse("6")),
			comver.NewGreaterThan(comver.MustParse("4")),
		),
	}

	c := comver.Compact(o)

	fmt.Println("Before:", o)
	fmt.Println("After:", c)

	// Output:
	// Before: >1 <2 || >3 <5 || >4 <6
	// After: >1 <2 || >3 <6
}

func ExampleCompact_endless() {
	o := comver.Or{
		comver.MustAnd(
			comver.NewLessThan(comver.MustParse("5")),
			comver.NewGreaterThan(comver.MustParse("3")),
		),
		comver.NewGreaterThan(comver.MustParse("4")),
	}

	c := comver.Compact(o)

	fmt.Println("Before:", o)
	fmt.Println("After:", c)

	// Output:
	// Before: >3 <5 || >4
	// After: >3
}

func ExampleCompact_matchAll() {
	o := comver.Or{
		comver.NewLessThan(comver.MustParse("3")),
		comver.NewGreaterThan(comver.MustParse("2")),
	}

	c := comver.Compact(o)

	fmt.Println("Before:", o)
	fmt.Println("After:", c)

	// Output:
	// Before: <3 || >2
	// After: *
}

func ExampleCompact_matchAllTrumps() {
	o := comver.Or{
		comver.MustAnd(
			comver.NewLessThan(comver.MustParse("2")),
			comver.NewGreaterThan(comver.MustParse("1")),
		),
		comver.NewMatchAll(),
	}

	c := comver.Compact(o)

	fmt.Println("Before:", o)
	fmt.Println("After:", c)

	// Output:
	// Before: >1 <2 || *
	// After: *
}
