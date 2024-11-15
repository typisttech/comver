package comver_test

import (
	"fmt"

	"github.com/typisttech/comver"
)

func buildIntervals(ks ...string) comver.Intervals {
	v1, _ := comver.NewVersion("1")
	v2, _ := comver.NewVersion("2")
	v3, _ := comver.NewVersion("3")
	v4, _ := comver.NewVersion("4")
	v5, _ := comver.NewVersion("5")
	v6, _ := comver.NewVersion("6")
	v7, _ := comver.NewVersion("7")
	v8, _ := comver.NewVersion("8")
	v9, _ := comver.NewVersion("9")
	v10, _ := comver.NewVersion("10")

	g1 := comver.NewGreaterThanConstraint(v1)
	g2 := comver.NewGreaterThanConstraint(v2)
	g3 := comver.NewGreaterThanConstraint(v3)
	g4 := comver.NewGreaterThanConstraint(v4)
	g7 := comver.NewGreaterThanConstraint(v7)
	g8 := comver.NewGreaterThanConstraint(v8)
	g10 := comver.NewGreaterThanConstraint(v10)

	gt4 := comver.NewGreaterThanOrEqualToConstraint(v4)

	l2 := comver.NewLessThanConstraint(v2)
	l3 := comver.NewLessThanConstraint(v3)
	l4 := comver.NewLessThanConstraint(v4)
	l1 := comver.NewLessThanConstraint(v1)
	l5 := comver.NewLessThanConstraint(v5)
	l6 := comver.NewLessThanConstraint(v6)
	l9 := comver.NewLessThanConstraint(v9)

	g1l2, _ := comver.NewInterval(g1, l2)
	g1l3, _ := comver.NewInterval(g1, l3)
	g2l4, _ := comver.NewInterval(g2, l4)
	g3l4, _ := comver.NewInterval(g3, l4)
	g4l5, _ := comver.NewInterval(g4, l5)
	g8l9, _ := comver.NewInterval(g8, l9)
	gt4l6, _ := comver.NewInterval(gt4, l6)

	ig1, _ := comver.NewInterval(g1, nil)
	ig10, _ := comver.NewInterval(g10, nil)
	ig7, _ := comver.NewInterval(g7, nil)

	il1, _ := comver.NewInterval(l1, nil)
	il3, _ := comver.NewInterval(l3, nil)
	il4, _ := comver.NewInterval(l4, nil)

	wildcard, _ := comver.NewInterval(nil, nil)

	m := map[string]comver.Intervals{
		"*":      {wildcard},
		"<1":     {il1},
		"<3":     {il3},
		"<4":     {il4},
		">=5 <6": {gt4l6},
		">1 <2":  {g1l2},
		">1 <3":  {g1l3},
		">1":     {ig1},
		">10":    {ig10},
		">2 <4":  {g2l4},
		">3 <4":  {g3l4},
		">4 <5":  {g4l5},
		">7":     {ig7},
		">8 <9":  {g8l9},
	}

	is := make(comver.Intervals, 0, len(ks))
	for _, k := range ks {
		is = append(is, m[k]...)
	}

	return is
}

func ExampleCompact() {
	is := buildIntervals(">1 <3", ">2 <4")
	fmt.Println("Original:", is)

	is = comver.Compact(is)
	fmt.Println("Compacted:", is)

	// Output:
	// Original: >1 <3 || >2 <4
	// Compacted: >1 <4
}

func ExampleCompact_wildcard() {
	is := buildIntervals(">1", "<4")
	fmt.Println("Original:", is)

	is = comver.Compact(is)
	fmt.Println("Compacted:", is)

	// Output:
	// Original: >1 || <4
	// Compacted: *
}

func ExampleCompact_wildcard2() {
	is := buildIntervals(">1 <3", "*")
	fmt.Println("Original:", is)

	is = comver.Compact(is)
	fmt.Println("Compacted:", is)

	// Output:
	// Original: >1 <3 || *
	// Compacted: *
}

func ExampleCompact_unrelated() {
	is := buildIntervals(">1 <2", ">3 <4")
	fmt.Println("Original:", is)

	is = comver.Compact(is)
	fmt.Println("Compacted:", is)

	// Output:
	// Original: >1 <2 || >3 <4
	// Compacted: >1 <2 || >3 <4
}

func ExampleCompact_mega() {
	is := buildIntervals(
		"<1",
		">1 <2",
		"<3",
		">4 <5",
		">=5 <6",
		">7",
		">8 <9",
		">10",
	)
	fmt.Println("Original:", is)

	is = comver.Compact(is)
	fmt.Println("Compacted:", is)

	// Output:
	// Original: <1 || >1 <2 || <3 || >4 <5 || >=4 <6 || >7 || >8 <9 || >10
	// Compacted: <3 || >=4 <6 || >7
}
