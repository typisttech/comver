package comver

type ExactConstraint struct {
	version Version
}

func NewExactConstraint(v Version) ExactConstraint {
	return ExactConstraint{
		version: v,
	}
}

// Check reports whether a [Version] satisfies the constraint.
func (e ExactConstraint) Check(v Version) bool {
	return e.version.Compare(v) == 0
}

func (e ExactConstraint) String() string {
	return e.version.Short()
}

func (e ExactConstraint) ceiling() Endless {
	return NewLessThanOrEqualTo(e.version)
}

func (e ExactConstraint) floor() Endless {
	return NewGreaterThanOrEqualTo(e.version)
}
