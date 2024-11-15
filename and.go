package comver

import (
	"slices"
)

type And struct {
	upper Boundless
	lower Boundless
}

func NewAnd(b0 Boundless, bs ...Boundless) SimpleConstrainter { // TODO return err
	if len(bs) == 0 {
		return b0
	}

	bs = slices.Clone(bs)
	bs = append(bs, b0)

	if containsWildcard(bs...) {
		return NewWildcard()
	}

	// TODO: Rethink! Should this be a max bounded floor?
	minBF, minBFOk := minBoundedFloor(bs)
	// TODO: Rethink! Should this be a min bounded floor?
	maxBC, maxBFOk := maxBoundedCelling(bs)

	// TODO: Check for And{minBF, maxBC} for wildcard.

	if !minBFOk && !maxBFOk {
		panic("impossible constrainter") // TODO!
	}

	if minBFOk && !maxBFOk {
		return minBF
	}

	if !minBFOk && maxBFOk {
		return maxBC
	}

	// different directions & overlapping
	if overlap(minBF, maxBC) {
		return NewWildcard()
	}

	if seamless(minBF, maxBC) {
		ev := minBF.Floor().version
		if ev == nil {
			ev = minBF.Ceiling().version
		}

		return NewExactConstraint(*ev)
	}

	//panic("impossible constrainter") // TODO! e.g: >2 <1

	floor, floorOk := minBF.(Boundless)
	if !floorOk {
		panic("impossible constrainter") // TODO!
	}
	ceiling, ceilingOk := maxBC.(Boundless)
	if !ceilingOk {
		panic("impossible constrainter") // TODO!
	}

	return And{
		upper: ceiling,
		lower: floor,
	}
}

func (a And) Ceiling() Boundless {
	return a.upper.Ceiling()
}

func (a And) Floor() Boundless {
	return a.lower.Floor()
}

// Check reports whether a [Version] satisfies the constraint.
func (a And) Check(v Version) bool {
	return a.upper.Check(v) && a.lower.Check(v)
}

func (a And) String() string {
	return a.Floor().String() + " " + a.Ceiling().String()
}
