package comver

const (
	greaterThanOrEqualTo constraintOp = iota
	greaterThan
	lessThan
	lessThanOrEqualTo

	ErrUnexpectedConstraintOp stringError = "unexpected constraintOp"
)

type constraintOp int8

func (op constraintOp) compare(other constraintOp) int {
	return int(op) - int(other)
}

func (op constraintOp) ceillingless() bool {
	return op == greaterThan || op == greaterThanOrEqualTo
}

func (op constraintOp) floorless() bool {
	return op == lessThanOrEqualTo || op == lessThan
}

func (op constraintOp) String() string {
	switch op {
	case lessThanOrEqualTo:
		return "<="
	case lessThan:
		return "<"
	case greaterThan:
		return ">"
	case greaterThanOrEqualTo:
		return ">="
	default:
		return ErrUnexpectedConstraintOp.Error()
	}
}

type constraint struct {
	// The Version used in the constraint check, e.g.: the Version representing 1.2.3 in '<=1.2.3'.
	version Version
	op      constraintOp
}

func NewLessThanOrEqualToConstraint(v Version) *constraint {
	return &constraint{
		version: v,
		op:      lessThanOrEqualTo,
	}
}

func NewLessThanConstraint(v Version) *constraint {
	return &constraint{
		version: v,
		op:      lessThan,
	}
}

func NewGreaterThanConstraint(v Version) *constraint {
	return &constraint{
		version: v,
		op:      greaterThan,
	}
}

func NewGreaterThanOrEqualToConstraint(v Version) *constraint {
	return &constraint{
		version: v,
		op:      greaterThanOrEqualTo,
	}
}

func (c *constraint) lowerbounded() bool {
	return c.op == greaterThan || c.op == greaterThanOrEqualTo
}

func (c *constraint) upperbounded() bool {
	return c.op == lessThanOrEqualTo || c.op == lessThan
}

// Check tests if a [Version] satisfies the constraints.
func (c *constraint) Check(v Version) bool {
	if c == nil {
		// this should never happen
		return true
	}

	cmp := v.Compare(c.version)

	switch c.op {
	case lessThan:
		return cmp < 0
	case lessThanOrEqualTo:
		return cmp <= 0
	case greaterThanOrEqualTo:
		return cmp >= 0
	case greaterThan:
		return cmp > 0
	default:
		// this should never happen
		return true
	}
}

func (c *constraint) String() string {
	return c.op.String() + c.version.Short()
}

// compare returns an integer comparing two constraints.
//
// The comparison is done by comparing the version first, then the operator.
//   - Versions are compared according to their semantic precedence
//   - Operators are compared in the following order (lowest to highest): >=, >, <, <=
//   - nil is considered to be the higher than upperbounded constraints
//   - nil is considered to be the lower than lowerbounded constraints
//
// The result is 0 when c == d, -1 when c < d, or +1 when c > d.
func (c *constraint) compare(d *constraint) int {
	switch {
	case c == nil && d == nil:
		return 0
	case c == nil:
		if d.lowerbounded() {
			return -1
		}
		return +1
	case d == nil:
		if c.lowerbounded() {
			return +1
		}
		return -1
	}

	if cmp := c.version.Compare(d.version); cmp != 0 {
		return cmp
	}
	return c.op.compare(d.op)
}
