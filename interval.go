package comver

// interval represents the intersection of two [constraint].
type interval [2]*constraint

const (
	ErrImpossibleInterval stringError = "impossible interval"
)

// NewInterval creates a new interval representing the intersection (logical AND) of two [constraint].
//
// If either c1 or c2 is nil, it represents a boundless range.
// If both c1 and c2 are nil, it matches all versions (wildcard).
func NewInterval(c1, c2 *constraint) (interval, error) {
	cmp := c1.compare(c2)
	// ensure c1 is the lower than c2
	if cmp > 0 {
		c1, c2 = c2, c1
	}

	switch {
	case c1 == nil && c2 == nil:
		return interval{}, nil
	case c1 == nil:
		return interval{c2}, nil
	case c2 == nil:
		return interval{c1}, nil
	case cmp == 0: // exactly the same
		return interval{c1}, nil
	case c1.op.ceillingless() && c2.op.ceillingless():
		// same direction
		return interval{c2}, nil
	case c1.op.floorless() && c2.op.floorless():
		// same direction
		return interval{c1}, nil
	case c1.version.Compare(c2.version) == 0 && c1.Check(c1.version) && c2.Check(c2.version):
		// same version & different directions & overlapping
		return interval{c1, c2}, nil
	case c1.Check(c2.version) && c2.Check(c1.version):
		return interval{c1, c2}, nil
	default:
		// different directions & no overlap
		return interval{}, ErrImpossibleInterval
	}
}

// Check tests if a [Version] lies within the interval.
func (i interval) Check(v Version) bool {
	for _, c := range i {
		if c != nil && !c.Check(v) {
			return false
		}
	}
	return true
}

func (i interval) String() string {
	switch {
	case i[0] == nil && i[1] == nil:
		return "*"
	case i[0] == nil:
		return i[1].String()
	case i[1] == nil:
		return i[0].String()
	}

	if i.exactVersionOnly() {
		return i[0].version.Short()
	}

	cmp := i[0].compare(i[1])
	switch {
	case cmp < 0:
		return i[0].String() + " " + i[1].String()
	case cmp > 0:
		return i[1].String() + " " + i[0].String()
	default:
		return i[0].String()
	}
}

func (i interval) wildcard() bool {
	return i[0] == nil && i[1] == nil
}

func (i interval) floorless() bool {
	return i.floor() == nil
}

func (i interval) floor() *constraint {
	switch {
	case i.wildcard():
		return nil
	// i[0] only
	case i[0] != nil && i[1] == nil && i[0].ceillingless():
		return i[0]
	case i[0] != nil && i[1] == nil && !i[0].ceillingless():
		return nil
	// i[1] only
	case i[0] == nil && i[1] != nil && i[1].ceillingless():
		return i[1]
	case i[0] == nil && i[1] != nil && !i[1].ceillingless():
		return nil
	}

	// both i[0] and i[1] are not nil
	if !i[0].ceillingless() && !i[1].ceillingless() {
		return nil
	}

	if i[0].ceillingless() && i[1].ceillingless() {
		cmp := i[0].compare(i[1])
		switch {
		case cmp < 0:
			return i[0]
		case cmp > 0:
			return i[1]
		default:
			return i[0]
		}
	}

	// exactly one of them is lower bounded
	if i[0].ceillingless() {
		return i[0]
	}
	return i[1]
}

func (i interval) ceilingless() bool {
	return i.ceiling() == nil
}

func (i interval) ceiling() *constraint {
	switch {
	case i.wildcard():
		return nil
	// i[0] only
	case i[0] != nil && i[1] == nil && i[0].floorless():
		return i[0]
	case i[0] != nil && i[1] == nil && !i[0].floorless():
		return nil
	// i[1] only
	case i[0] == nil && i[1] != nil && i[1].floorless():
		return i[1]
	case i[0] == nil && i[1] != nil && !i[1].floorless():
		return nil
	}

	// both i[0] and i[1] are not nil

	if !i[0].floorless() && !i[1].floorless() {
		return nil
	}

	if i[0].floorless() && i[1].floorless() {
		cmp := i[0].compare(i[1])
		switch {
		case cmp < 0:
			return i[1]
		case cmp > 0:
			return i[0]
		default:
			return i[0]
		}
	}

	// exactly one of them is upper bounded
	if i[0].floorless() {
		return i[0]
	}
	return i[1]
}

func (i interval) exactVersionOnly() bool {
	if i[0] == nil || i[1] == nil {
		return false
	}

	if i[0].version.Compare(i[1].version) != 0 {
		return false
	}

	return (i[0].ceillingless() && i[1].floorless()) || (i[0].floorless() && i[1].ceillingless())
}

func (i interval) compare(j interval) int {
	cmp := i.floor().compare(j.floor())
	if cmp != 0 {
		return cmp
	}
	return j.ceiling().compare(j.ceiling())
}