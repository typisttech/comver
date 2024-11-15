package comver_test

import (
	"fmt"

	"github.com/typisttech/comver"
)

func Example_version() {
	ss := []string{
		"1.2.3",
		"v1.2.p5+foo",
		"v1.2.3.4.p5+foo",
		"2010-01-02",
		"2010-01-02.5",
		"not a version",
		"1.0.0-meh",
		"20100102.0.3.4",
		"1.0.0-alpha.beta",
	}

	for _, s := range ss {
		v, err := comver.NewVersion(s)
		if err != nil {
			fmt.Println(s, " => ", err)
			continue
		}
		fmt.Println(s, " => ", v)
	}

	// Output:
	// 1.2.3  =>  1.2.3.0
	// v1.2.p5+foo  =>  1.2.0.0-patch5
	// v1.2.3.4.p5+foo  =>  1.2.3.4-patch5
	// 2010-01-02  =>  2010.1.2.0
	// 2010-01-02.5  =>  2010.1.2.5
	// not a version  =>  error parsing version string "not a version"
	// 1.0.0-meh  =>  error parsing version string "1.0.0-meh"
	// 20100102.0.3.4  =>  error parsing version string "20100102.0.3.4"
	// 1.0.0-alpha.beta  =>  error parsing version string "1.0.0-alpha.beta"
}

func Example_constraint() {
	v1, _ := comver.NewVersion("1")
	v2, _ := comver.NewVersion("2")
	v3, _ := comver.NewVersion("3")
	v4, _ := comver.NewVersion("4")

	cs := []any{
		comver.NewGreaterThanConstraint(v1),
		comver.NewGreaterThanOrEqualToConstraint(v2),
		comver.NewLessThanOrEqualToConstraint(v3),
		comver.NewLessThanConstraint(v4),
	}

	for _, c := range cs {
		fmt.Println(c)
	}

	// Output:
	// >1
	// >=2
	// <=3
	// <4
}

func Example_interval() {
	v1, _ := comver.NewVersion("1")
	v2, _ := comver.NewVersion("2")
	v3, _ := comver.NewVersion("3")

	g1l3, _ := comver.NewInterval(
		comver.NewGreaterThanConstraint(v1),
		comver.NewLessThanConstraint(v3),
	)

	if g1l3.Check(v2) {
		fmt.Println(v2.Short(), "satisfies", g1l3)
	}

	if !g1l3.Check(v3) {
		fmt.Println(v2.Short(), "doesn't satisfy", g1l3)
	}

	// Output:
	// 2 satisfies >1 <3
	// 2 doesn't satisfy >1 <3
}

func Example_intervals() {
	v1, _ := comver.NewVersion("1")
	v2, _ := comver.NewVersion("2")
	v3, _ := comver.NewVersion("3")
	v4, _ := comver.NewVersion("4")

	g1l3, _ := comver.NewInterval(
		comver.NewGreaterThanConstraint(v1),
		comver.NewLessThanConstraint(v3),
	)

	ge2le4, _ := comver.NewInterval(
		comver.NewGreaterThanOrEqualToConstraint(v2),
		comver.NewLessThanOrEqualToConstraint(v4),
	)

	is := comver.Intervals{g1l3, ge2le4}
	fmt.Println(is)

	is = comver.Compact(is)
	fmt.Println(is)

	// Output:
	// >1 <3 || >=2 <=4
	// >1 <=4
}
