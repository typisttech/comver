package comver

const (
	greaterThanOrEqualTo op = iota
	greaterThan
	lessThan
	lessThanOrEqualTo

	errUnexpectedOp stringError = "unexpected op"
)

type op int8

func (o op) String() string {
	switch o {
	case lessThanOrEqualTo:
		return "<="
	case lessThan:
		return "<"
	case greaterThan:
		return ">"
	case greaterThanOrEqualTo:
		return ">="
	default:
		// logic error! This should never happen
		panic(errUnexpectedOp)
	}
}

func (o op) compare(other op) int {
	i := int(o) - int(other)

	switch {
	case i < 0:
		return -1
	case i > 0:
		return 1
	default:
		return 0
	}
}

func (o op) ceilingBounded() bool {
	return o == lessThan || o == lessThanOrEqualTo
}

func (o op) floorBounded() bool {
	return o == greaterThan || o == greaterThanOrEqualTo
}

func (o op) inclusive() bool {
	return o == lessThanOrEqualTo || o == greaterThanOrEqualTo
}
