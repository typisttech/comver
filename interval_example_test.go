package comver_test

import (
	"fmt"
	"github.com/typisttech/comver"
)

func ExampleNewInterval() {
	v1, _ := comver.NewVersion("1")
	v2, _ := comver.NewVersion("2")

	c1 := comver.NewGreaterThanConstraint(v1)
	c2 := comver.NewLessThanConstraint(v2)

	i, _ := comver.NewInterval(c1, c2)

	fmt.Println(i)
	// Output: >1 <2
}

func ExampleNewInterval_error() {
	v1, _ := comver.NewVersion("1")
	v2, _ := comver.NewVersion("2")

	c1 := comver.NewLessThanConstraint(v1)
	c2 := comver.NewGreaterThanConstraint(v2)

	_, err := comver.NewInterval(c1, c2)

	fmt.Println(err)
	// Output: impossible interval
}

func ExampleNewInterval_ceillingless() {
	v1, _ := comver.NewVersion("1")

	c1 := comver.NewGreaterThanOrEqualToConstraint(v1)

	i, _ := comver.NewInterval(c1, nil)

	fmt.Println(i)
	// Output: >=1
}

func ExampleNewInterval_floorless() {
	v1, _ := comver.NewVersion("1")

	c1 := comver.NewLessThanOrEqualToConstraint(v1)

	i, _ := comver.NewInterval(c1, nil)

	fmt.Println(i)
	// Output: <=1
}

func ExampleNewInterval_exactVersion() {
	v1, _ := comver.NewVersion("1")

	c1 := comver.NewLessThanOrEqualToConstraint(v1)
	c2 := comver.NewGreaterThanOrEqualToConstraint(v1)

	i, _ := comver.NewInterval(c1, c2)

	fmt.Println(i)
	// Output: 1
}

func ExampleNewInterval_wildcard() {
	i, _ := comver.NewInterval(nil, nil)

	fmt.Println(i)
	// Output: *
}

func ExampleNewInterval_compact() {
	v1, _ := comver.NewVersion("1")
	v2, _ := comver.NewVersion("2")

	c1 := comver.NewLessThanConstraint(v1)
	c2 := comver.NewLessThanConstraint(v2)

	i, _ := comver.NewInterval(c1, c2)

	fmt.Println(i)
	// Output: <1
}

func ExampleNewInterval_compactWildcard() {
	v, _ := comver.NewVersion("1")

	c1 := comver.NewLessThanConstraint(v)
	c2 := comver.NewGreaterThanOrEqualToConstraint(v)

	i, _ := comver.NewInterval(c1, c2)

	fmt.Println(i)
	// Output: *
}
