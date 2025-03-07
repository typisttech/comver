package comver

//nolint:godox
// TODO: Make Or to be []Constrainter so that we can nest Or

// Or represents a logical OR operation between multiple [CeilingFloorConstrainter] instances.
// The zero value for Or is a [match none] constraint which could never be satisfied.
//
// [match none]: https://github.com/composer/semver/blob/main/src/Constraint/MatchNoneConstraint.php
type Or []CeilingFloorConstrainter

// Check reports whether a [Version] satisfies the constraint.
func (o Or) Check(v Version) bool {
	for i := range o {
		if o[i].Check(v) {
			return true
		}
	}

	return false
}

func (o Or) String() string {
	s := ""

	for i := range o {
		if i > 0 {
			s += " || "
		}
		s += o[i].String()
	}

	return s
}
