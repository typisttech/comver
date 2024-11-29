package comver

import "slices"

const (
	errNoEndlessGiven     stringError = "no endless given"
	errUnexpectedAndLogic stringError = "unexpected and logic"
	errImpossibleInterval stringError = "impossible interval"
)

// And returns a [CeilingFloorConstrainter] instance representing the logical AND of
// the given [Endless] instances; or return an error if the given [Endless] instances
// could never be satisfied at the same time.
func And(es ...Endless) (CeilingFloorConstrainter, error) { //nolint:cyclop,ireturn
	var nilC CeilingFloorConstrainter

	if len(es) == 0 {
		return nilC, errNoEndlessGiven
	}

	es = slices.Clone(es)
	es = slices.DeleteFunc(es, Endless.wildcard)

	if len(es) == 0 {
		return NewWildcard(), nil
	}
	if len(es) == 1 {
		return es[0], nil
	}

	ceiling, ceilingOk := minBoundedCeiling(es...)
	floor, floorOk := maxBoundedFloor(es...)

	if !ceilingOk && !floorOk {
		// logic error! This should never happen
		return nilC, errUnexpectedAndLogic
	}
	if ceilingOk && !floorOk {
		return ceiling, nil
	}
	if !ceilingOk { // floorOk is always true here
		return floor, nil
	}

	vCmp := floor.floor().versionCompare(ceiling.ceiling().version)

	if vCmp > 0 {
		return nilC, errImpossibleInterval
	}

	if vCmp == 0 {
		if !floor.floor().inclusive() || !ceiling.ceiling().inclusive() {
			return nilC, errImpossibleInterval
		}

		return NewExactConstraint(*floor.floor().version), nil
	}

	return interval{
		upper: ceiling,
		lower: floor,
	}, nil
}

// MustAnd is like [And] but panics if an error occurs.
func MustAnd(es ...Endless) CeilingFloorConstrainter { //nolint:ireturn
	c, err := And(es...)
	if err != nil {
		panic(err)
	}

	return c
}

func minBoundedCeiling(es ...Endless) (Endless, bool) {
	es = slices.Clone(es)

	bcs := slices.DeleteFunc(es, func(b Endless) bool {
		return b.ceiling().version == nil
	})

	if len(bcs) == 0 {
		var nilF Endless

		return nilF, false
	}

	return slices.MinFunc(bcs, func(a, b Endless) int {
		return a.ceiling().compare(b.ceiling())
	}), true
}

func maxBoundedFloor(es ...Endless) (Endless, bool) {
	es = slices.Clone(es)

	bfs := slices.DeleteFunc(es, func(c Endless) bool {
		return c.floor().wildcard()
	})

	if len(bfs) == 0 {
		var nilF Endless

		return nilF, false
	}

	return slices.MaxFunc(bfs, func(a, b Endless) int {
		return a.floor().compare(b.floor())
	}), true
}
