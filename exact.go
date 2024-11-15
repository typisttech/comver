package comver

type ExactConstraint struct {
	version Version
}

func NewExactConstraint(v Version) ExactConstraint {
	return ExactConstraint{
		version: v,
	}
}

func (e ExactConstraint) Ceiling() Boundless {
	return NewLessThanOrEqualTo(e.version)
}

func (e ExactConstraint) Floor() Boundless {
	return NewGreaterThanOrEqualTo(e.version)
}

// Check reports whether a [Version] satisfies the constraint.
func (e ExactConstraint) Check(v Version) bool {
	return e.version.Compare(v) == 0
}

func (e ExactConstraint) String() string {
	return e.version.Short()
}
