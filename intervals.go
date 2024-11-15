package comver

import (
	"slices"
)

// Intervals represent the union (logical OR) of multiple intervals.
type Intervals []interval

func (is Intervals) String() string {
	s := ""
	for _, i := range is {
		if len(s) > 0 {
			s += " || "
		}
		s += i.String()
	}
	return s
}

// Compact returns a new [Intervals] covering the same version ranges but with the smallest number of intervals.
//
// This function compacts it by looking at the real version ranges covered by all the intervals and then creates
// a new [Intervals] containing only the smallest numbers of intervals to cover the same version ranges.
func Compact(is Intervals) Intervals {
	if len(is) == 1 {
		return is
	}
	if len(is) == 0 {
		return Intervals{}
	}

	is = slices.Clone(is)

	// if there is a wildcard, return only the wildcard
	wIndex := slices.IndexFunc(is, func(i interval) bool {
		return i.wildcard()
	})
	if wIndex >= 0 {
		return Intervals{is[wIndex]}
	}

	head := lowestCeilingless(is)
	tail := highestFloorless(is)

	is = slices.DeleteFunc(is, func(i interval) bool {
		return i.floorless() || i.ceilingless()
	})

	is = slices.Clip(slices.Concat(head, is, tail))

	// sort the intervals
	slices.SortFunc(is, func(a, b interval) int {
		return a.compare(b)
	})
	// remove duplicates
	is = slices.CompactFunc(is, func(a, b interval) bool {
		return a.compare(b) == 0
	})

	vals := make(Intervals, 0, len(is))
	pendingI := is[0]
	for index := range is {
		i, ok := compactTwo(pendingI, is[index])
		if ok {
			pendingI = i
		} else {
			vals = append(vals, pendingI)
			pendingI = is[index]
		}

		if index == len(is)-1 {
			vals = append(vals, pendingI)
		}
	}

	return slices.Clip(vals)
}

func compactTwo(a, b interval) (interval, bool) {
	cmp := a.compare(b)
	if cmp > 0 {
		a, b = b, a
	}

	if a.compare(b) == 0 {
		return a, true
	}

	overlap := a.Check(b.floor().version) ||
		(a.ceiling().version.Compare(b.floor().version) == 0 &&
			(a.ceiling().op == lessThanOrEqualTo || b.floor().op == greaterThanOrEqualTo))

	if !overlap {
		return a, false
	}

	ccmp := a.ceiling().compare(b.ceiling())
	if ccmp > 0 {
		return a, true
	}

	i, err := NewInterval(a.floor(), b.ceiling())
	if err != nil {
		// this should not happen
		return a, false
	}
	return i, true
}

func lowestCeilingless(is Intervals) Intervals {
	ceilingless := slices.DeleteFunc(slices.Clone(is), func(i interval) bool {
		return !i.ceilingless()
	})

	if len(ceilingless) == 0 {
		return Intervals{}
	}

	i := slices.MinFunc(ceilingless, func(a, b interval) int {
		return a.floor().compare(b.floor())
	})

	return Intervals{i}
}

func highestFloorless(is Intervals) Intervals {
	floorlesses := slices.DeleteFunc(slices.Clone(is), func(i interval) bool {
		return !i.floorless()
	})

	if len(floorlesses) == 0 {
		return Intervals{}
	}

	i := slices.MaxFunc(floorlesses, func(a, b interval) int {
		return a.ceiling().compare(b.ceiling())
	})

	return Intervals{i}
}
