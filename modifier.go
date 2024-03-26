package comver

type modifier int8

const (
	modifierPatch         modifier    = 10
	modifierStable        modifier    = 0
	modifierRC            modifier    = -10
	modifierBeta          modifier    = -20
	modifierAlpha         modifier    = -30
	ErrUnexpectedModifier stringError = "unexpected modifier"
)

func newStability(s string) (modifier, error) {
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
	return 0, ErrUnexpectedModifier
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
	}
	return ""
}
