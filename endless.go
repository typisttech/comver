package comver

// Endless represents a constraint that is either floor bounded, ceiling bounded,
// or [match all].
// The zero value for Endless is a [match all] constraint which satisfied by any [Version].
//
// [match all]: https://github.com/composer/semver/blob/main/src/Constraint/MatchAllConstraint.php
type Endless struct {
	// The version used in the constraint check,
	// e.g.: the version representing 1.2.3 in '<=1.2.3'.
	// If nil, the Endless is a match all constraint which satisfied by any version.
	version *Version
	op      op
}

func NewLessThanOrEqualTo(v Version) Endless {
	return Endless{
		version: &v,
		op:      lessThanOrEqualTo,
	}
}

func NewLessThan(v Version) Endless {
	return Endless{
		version: &v,
		op:      lessThan,
	}
}

func NewGreaterThan(v Version) Endless {
	return Endless{
		version: &v,
		op:      greaterThan,
	}
}

func NewGreaterThanOrEqualTo(v Version) Endless {
	return Endless{
		version: &v,
		op:      greaterThanOrEqualTo,
	}
}

func NewMatchAll() Endless {
	return Endless{ //nolint:exhaustruct
		version: nil,
	}
}

// Check reports whether a [Version] satisfies the constraint.
func (b Endless) Check(v Version) bool {
	if b.version == nil {
		// this is a match all
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
		// this should never happen
		panic("unexpected constraint operator")
	}
}

func (b Endless) String() string {
	if b.matchAll() {
		return "*"
	}

	return b.op.String() + b.version.Short()
}

func (b Endless) ceiling() Endless {
	if !b.ceilingBounded() {
		return NewMatchAll()
	}

	return b
}

func (b Endless) floor() Endless {
	if !b.floorBounded() {
		return NewMatchAll()
	}

	return b
}

func (b Endless) matchAll() bool {
	return b.version == nil
}

// compare returns an integer comparing two [Endless] instances.
//
// The comparison is done by comparing the version first, then the operator.
//   - Versions are compared according to their semantic precedence
//   - Operators are compared in the following order (lowest to highest): >=, >, <, <=
//   - Match all [Endless] is considered to be higher than ceiling bounded [Endless] while
//     lower than floor bounded [Endless]
//
// The result is 0 when b == d, -1 when b < d, or +1 when b > d.
func (b Endless) compare(d Endless) int {
	switch {
	case b.matchAll() && d.matchAll():
		return 0
	case b.matchAll():
		if d.floorBounded() {
			return -1
		}

		return +1
	case d.matchAll():
		if b.floorBounded() {
			return +1
		}

		return -1
	}

	if cmp := b.versionCompare(d.version); cmp != 0 {
		return cmp
	}

	return b.op.compare(d.op)
}

func (b Endless) ceilingBounded() bool {
	return b.op.ceilingBounded()
}

func (b Endless) floorBounded() bool {
	return b.op.floorBounded()
}

func (b Endless) inclusive() bool {
	return b.op.inclusive()
}

func (b Endless) versionCompare(v *Version) int {
	return b.version.Compare(*v)
}
