package comver

// interval represents a constraint that is both floor bounded and ceiling bounded.
// It must be initialized via [And].
// The zero value for interval is a constraint could never be satisfied.
type interval struct {
	upper Endless
	lower Endless
}

// Check reports whether a [Version] satisfies the constraint.
func (i interval) Check(v Version) bool {
	return i.ceiling().Check(v) && i.floor().Check(v)
}

func (i interval) ceiling() Endless {
	return i.upper
}

func (i interval) floor() Endless {
	return i.lower
}

func (i interval) String() string {
	return i.floor().String() + " " + i.ceiling().String()
}
