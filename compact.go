package comver

import (
	"slices"
)

// Compact returns a new [Constrainter] that is logically equivalent to the input [Or].
// The returned [Constrainter] may or may be not be an [Or] instance.
// When it is, Compact tries to return the shortest [Or] possible.
func Compact(o Or) Constrainter { //nolint:cyclop,ireturn
	if len(o) == 0 {
		return Or{}
	}
	if len(o) == 1 {
		return o[0]
	}
	if slices.ContainsFunc(o, matchAll) {
		return NewMatchAll()
	}

	o = slices.Clone(o)

	ceiling, ceilingOk := maxFloorlessCeiling(o...)
	floor, floorOk := minCeilinglessFloor(o...)

	// short circuit if we have a match all
	if ceilingOk && floorOk && disjunctivelyCombineToMatchAll(ceiling, floor) {
		return NewMatchAll()
	}

	o = slices.DeleteFunc(o, func(c CeilingFloorConstrainter) bool {
		return c.ceiling().matchAll() || c.floor().matchAll() ||
			(ceilingOk && ceiling.compare(c.ceiling()) >= 0) ||
			(floorOk && floor.compare(c.floor()) <= 0)
	})

	// important to sort before compacting
	slices.SortFunc(o, compare)
	o = slices.CompactFunc(o, func(a, b CeilingFloorConstrainter) bool {
		return compare(a, b) == 0
	})

	r := compactMultiple(o)

	if floorOk {
		r = append(r, floor)
	}
	if ceilingOk {
		r = append(r, ceiling)
	}

	if len(r) == 1 {
		return r[0]
	}

	slices.SortFunc(r, compare)

	return slices.Clip(r)
}

func matchAll(c CeilingFloorConstrainter) bool {
	return c.floor().matchAll() && c.ceiling().matchAll()
}

func minCeilinglessFloor(cs ...CeilingFloorConstrainter) (Endless, bool) {
	cs = slices.Clone(cs)
	o := slices.Clone(cs)

	cs = slices.DeleteFunc(cs, func(c CeilingFloorConstrainter) bool {
		return !c.ceiling().matchAll()
	})

	if len(cs) == 0 {
		var nilE Endless

		return nilE, false
	}

	m := slices.MinFunc(cs, func(a, b CeilingFloorConstrainter) int {
		return a.floor().compare(b.floor())
	}).floor()

	o = slices.DeleteFunc(o, func(c CeilingFloorConstrainter) bool {
		return c.ceiling().matchAll() ||
			c.floor().matchAll()
	})

	for i := range o {
		if o[i].floor().compare(m) < 0 {
			if o[i].Check(*m.version) || disjunctivelyCombineToMatchAll(o[i].ceiling(), m) {
				m = o[i].floor()
			}

			continue
		}

		if m.versionCompare(o[i].floor().version) == 0 {
			if o[i].floor().inclusive() {
				m = o[i].floor()
			}

			continue
		}
	}

	return m.floor(), true
}

func maxFloorlessCeiling(cs ...CeilingFloorConstrainter) (Endless, bool) {
	cs = slices.Clone(cs)
	o := slices.Clone(cs)

	cs = slices.DeleteFunc(cs, func(c CeilingFloorConstrainter) bool {
		return !c.floor().matchAll()
	})

	if len(cs) == 0 {
		var nilE Endless

		return nilE, false
	}

	m := slices.MaxFunc(cs, func(a, b CeilingFloorConstrainter) int {
		return a.ceiling().compare(b.ceiling())
	}).ceiling()

	o = slices.DeleteFunc(o, func(c CeilingFloorConstrainter) bool {
		return c.ceiling().matchAll() || c.floor().matchAll()
	})

	for i := range o {
		if o[i].ceiling().compare(m) > 0 {
			if o[i].Check(*m.version) || disjunctivelyCombineToMatchAll(o[i].floor(), m) {
				m = o[i].ceiling()
			}

			continue
		}

		if m.versionCompare(o[i].ceiling().version) == 0 {
			if o[i].ceiling().inclusive() {
				m = o[i].ceiling()
			}

			continue
		}
	}

	return m, true
}

func disjunctivelyCombineToMatchAll(e, f Endless) bool {
	if e.ceilingBounded() && f.ceilingBounded() {
		return false
	}

	if e.floorBounded() && f.floorBounded() {
		return false
	}

	cmp := e.compare(f)

	if cmp == 0 {
		return false
	}

	if cmp > 0 {
		e, f = f, e
	}

	if e.versionCompare(f.version) == 0 {
		return e.inclusive() || f.inclusive()
	}

	return !e.ceilingBounded() && !f.floorBounded()
}

func compare(a, b CeilingFloorConstrainter) int {
	cmp := a.floor().compare(b.floor())

	if cmp != 0 {
		return cmp
	}

	return a.ceiling().compare(b.ceiling())
}

func compactTwo(a, b CeilingFloorConstrainter) (CeilingFloorConstrainter, bool) { //nolint:ireturn
	cmp := compare(a, b)
	if cmp == 0 {
		return a, true
	}

	if cmp > 0 {
		a, b = b, a
	}

	if !overlap(a, b) && !continuous(a, b) {
		return a, false
	}

	if a.ceiling().compare(b.ceiling()) > 0 {
		return a, true
	}

	return interval{
		upper: b.ceiling(),
		lower: a.floor(),
	}, true
}

func compactMultiple(o []CeilingFloorConstrainter) Or {
	r := make(Or, 0, len(o)+2) //nolint:mnd

	if len(o) != 0 {
		p := o[0]
		for i := range o {
			q, ok := compactTwo(p, o[i])

			if ok {
				p = q
			} else {
				r = append(r, p)
				p = o[i]
			}

			// always append the last r
			if i == len(o)-1 {
				r = append(r, p)
			}
		}
	}

	return r
}

func overlap(a, b CeilingFloorConstrainter) bool {
	return a.Check(*b.floor().version) ||
		b.Check(*a.floor().version)
}

func continuous(a, b CeilingFloorConstrainter) bool {
	f := func(a, b CeilingFloorConstrainter) bool {
		return (a.ceiling().inclusive() || b.floor().inclusive()) &&
			a.ceiling().version.Compare(*b.floor().version) == 0
	}

	return f(a, b) || f(b, a)
}
