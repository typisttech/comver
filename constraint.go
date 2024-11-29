package comver

type Constrainter interface {
	// Check reports whether a [Version] satisfies the constraint.
	Check(v Version) bool
	String() string
}

type CeilingFloorConstrainter interface {
	ceiling() Endless
	floor() Endless

	Constrainter
}
