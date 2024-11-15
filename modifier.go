package comver

type modifier int8

const (
	modifierPatch         modifier    = 10
	modifierStable        modifier    = 0
	modifierRC            modifier    = -10
	modifierBeta          modifier    = -20
	modifierAlpha         modifier    = -30
	errUnexpectedModifier stringError = "unexpected modifier"
)

func newModifier(s string) (modifier, error) {
	switch s {
	case "":
		return modifierStable, nil
	case "patch", "pl", "p":
		return modifierPatch, nil
	case "rc":
		return modifierRC, nil
	case "beta", "b":
		return modifierBeta, nil
	case "alpha", "a":
		return modifierAlpha, nil
	}

	return modifierStable, errUnexpectedModifier
}

func (s modifier) String() string {
	switch s {
	case modifierPatch:
		return "patch"
	case modifierRC:
		return "RC"
	case modifierBeta:
		return "beta"
	case modifierAlpha:
		return "alpha"
	case modifierStable:
		return ""
	default:
		return ""
	}
}
