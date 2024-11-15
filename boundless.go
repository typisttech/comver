package comver

import "slices"

type Boundless struct {
	// The version used in the constraint check,
	// e.g.: the version representing 1.2.3 in '<=1.2.3'.
	// If nil, the Boundless is a wildcard which match all versions.
	version *Version
	op      ConstraintOp
}

func NewLessThanOrEqualTo(v Version) Boundless {
	return Boundless{
		version: &v,
		op:      lessThanOrEqualTo,
	}
}

func NewLessThan(v Version) Boundless {
	return Boundless{
		version: &v,
		op:      lessThan,
	}
}

func NewGreaterThan(v Version) Boundless {
	return Boundless{
		version: &v,
		op:      greaterThan,
	}
}

func NewGreaterThanOrEqualTo(v Version) Boundless {
	return Boundless{
		version: &v,
		op:      greaterThanOrEqualTo,
	}
}

func NewWildcard() Boundless {
	return Boundless{
		version: nil,
		op:      greaterThanOrEqualTo, // op is unused.
	}
}

func (b Boundless) Ceiling() Boundless {
	if !b.op.upperBounded() {
		return NewWildcard()
	}
	return b
}

func (b Boundless) Floor() Boundless {
	if !b.op.lowerBounded() {
		return NewWildcard()
	}
	return b
}

// Check reports whether a [Version] satisfies the constraint.
func (b Boundless) Check(v Version) bool {
	if b.version == nil {
		// this is wildcard, match all versions
		return true
	}

	cmp := b.version.Compare(v)

	switch b.op {
	case lessThan:
		return cmp > 0
	case lessThanOrEqualTo:
		return cmp >= 0
	case greaterThanOrEqualTo:
		return cmp <= 0
	case greaterThan:
		return cmp < 0
	default:
		// TODO: this should never happen
		panic("invalid constraint operator")
	}
}

func (b Boundless) String() string {
	if b.version == nil {
		return "*"
	}

	return b.op.String() + b.version.Short()
}

// compare returns an integer comparing two [Boundless]es.
//
// The comparison is done by comparing the version first, then the operator.
//   - Versions are compared according to their semantic precedence
//   - Operators are compared in the following order (lowest to highest): >=, >, <, <=
//   - wildcard [Boundless] is considered to be higher than upper bounded [Boundless] while
//     lower than lower bounded [Boundless]
//
// The result is 0 when b == d, -1 when b < d, or +1 when b > d.
func (b Boundless) compare(d Boundless) int {
	switch {
	case b.version == nil && d.version == nil:
		return 0
	case b.version == nil:
		if d.op.lowerBounded() {
			return -1
		}
		return +1
	case d.version == nil:
		if b.op.lowerBounded() {
			return +1
		}
		return -1
	}

	if cmp := b.version.Compare(*d.version); cmp != 0 {
		return cmp
	}
	return b.op.compare(d.op)
}

func (b Boundless) inclusive() bool {
	return b.op.inclusive()
}

func containsWildcard(bs ...Boundless) bool {
	for i := range bs {
		if bs[i].version == nil {
			return true
		}
	}

	return false
}

func minBoundedFloor[S ~[]T, T SimpleConstrainter](bs S) (SimpleConstrainter, bool) {

	//func minBoundedFloor[T SimpleConstrainter](bs ...T) (SimpleConstrainter, bool) {
	bs = slices.Clone(bs)

	bs = slices.DeleteFunc(bs, func(b T) bool {
		return b.Floor().version == nil
	})

	if len(bs) == 0 {
		return NewWildcard(), false
	}

	bs = slices.Clone(bs)

	return slices.MinFunc(bs, func(n, m T) int {
		return n.Floor().compare(m.Floor())
	}), true

	//return minFloor(bs), true
}

func minFloor[S ~[]T, T SimpleConstrainter](bs S) SimpleConstrainter {
	if len(bs) == 0 {
		return NewWildcard()
	}

	bs = slices.Clone(bs)

	return slices.MinFunc(bs, func(n, m T) int {
		return n.Floor().compare(m.Floor())
	})
}

func maxBoundedCelling[S ~[]T, T SimpleConstrainter](bs S) (SimpleConstrainter, bool) {
	bs = slices.Clone(bs)

	bs = slices.DeleteFunc(bs, func(b T) bool {
		return b.Ceiling().version == nil
	})

	if len(bs) == 0 {
		return NewWildcard(), false
	}

	bs = slices.Clone(bs)

	return slices.MaxFunc(bs, func(n, m T) int {
		return n.Ceiling().compare(m.Ceiling())
	}), true

	//return maxCelling(bs), true
}

//func maxCelling[S ~[]T, T SimpleConstrainter](bs S) SimpleConstrainter {
//	if len(bs) == 0 {
//		return NewWildcard()
//	}
//
//	bs = slices.Clone(bs)
//
//	return slices.MaxFunc(bs, func(n, m T) int {
//		return n.Ceiling().compare(m.Ceiling())
//	})
//}
