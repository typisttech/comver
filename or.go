package comver

import "slices"

type Or []SimpleConstrainter

func (o Or) Check(v Version) bool {
	return true // TODO: Implement!
}

func CompactOr(cs Or) Constrainter {
	if len(cs) == 0 {
		return NewWildcard()
	}
	if len(cs) == 1 {
		return cs[0]
	}

	cs = slices.Clone(cs)

	// if there cs a wildcard, return only the wildcard
	w := slices.ContainsFunc(cs, func(c SimpleConstrainter) bool {
		return c.Floor().version == nil && c.Ceiling().version == nil
	})
	if w {
		return NewWildcard()
	}

	head, headOk := minBoundedFloor(cs)
	cs = slices.DeleteFunc(cs, func(c SimpleConstrainter) bool {
		return c.Floor().version == nil
	})
	if headOk {
		cs = append(cs, head)
	}

	tail, tailOk := maxBoundedCelling(cs)
	cs = slices.DeleteFunc(cs, func(c SimpleConstrainter) bool {
		return c.Ceiling().version == nil
	})
	if tailOk {
		cs = append(cs, tail)
	}

	// TODO: Check whether And{head, tail} are wildcard.
	// If yes, early return wildcard

	// important!
	slices.SortFunc(cs, compareSimpleConstrainters)

	// remove duplicates
	cs = slices.CompactFunc(cs, func(a, b SimpleConstrainter) bool {
		return compareSimpleConstrainters(a, b) == 0
	})

	vals := make(Or, 0, len(cs))
	pendingI := cs[0]
	for index := range cs {
		i, ok := compactTwoSCs(pendingI, cs[index])
		if ok {
			pendingI = i
		} else {
			vals = append(vals, pendingI)
			pendingI = cs[index]
		}

		if index == len(cs)-1 {
			vals = append(vals, pendingI)
		}
	}

	return slices.Clip(vals)
}

func overlap(a, b SimpleConstrainter) bool {
	// TODO: Assume version != nil?
	return a.Check(*b.Floor().version) ||
		b.Check(*a.Floor().version)
}

// TODO: Does overlap check for this?
func seamless(a, b SimpleConstrainter) bool {
	// TODO: Assume version != nil?
	ab := (a.Ceiling().inclusive() || b.Floor().inclusive()) &&
		a.Ceiling().version.Compare(*b.Floor().version) == 0

	if ab {
		return true
	}

	return (b.Ceiling().inclusive() || a.Floor().inclusive()) &&
		b.Ceiling().version.Compare(*a.Floor().version) == 0
}

func compactTwoSCs(a, b SimpleConstrainter) (SimpleConstrainter, bool) {
	cmp := compareSimpleConstrainters(a, b)
	if cmp > 0 {
		a, b = b, a
	}

	if cmp == 0 {
		return a, true
	}

	if !overlap(a, b) && !seamless(a, b) {
		return a, false
	}

	if cmp := a.Ceiling().compare(b.Ceiling()); cmp > 0 {
		return a, true
	}

	return And{a.Floor(), b.Ceiling()}, true
}

func compareSimpleConstrainters(a, b SimpleConstrainter) int {
	cmp := a.Floor().compare(b.Floor())
	if cmp != 0 {
		return cmp
	}
	return a.Ceiling().compare(b.Ceiling())
}
