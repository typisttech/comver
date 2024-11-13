package comver

import "testing"

func Test_constraint_Check(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		constraint *constraint
		version    Version
		want       bool
	}{
		{
			name:       "lessThan satisfied",
			constraint: &constraint{version: Version{major: 2}, op: lessThan},
			version:    Version{major: 1},
			want:       true,
		},
		{
			name:       "lessThan just satisfied",
			constraint: &constraint{version: Version{major: 2}, op: lessThan},
			version:    Version{major: 2, modifier: modifierRC},
			want:       true,
		},
		{
			name:       "lessThan just not satisfied",
			constraint: &constraint{version: Version{major: 2}, op: lessThan},
			version:    Version{major: 2},
			want:       false,
		},
		{
			name:       "lessThan not satisfied",
			constraint: &constraint{version: Version{major: 2}, op: lessThan},
			version:    Version{major: 3},
			want:       false,
		},

		{
			name:       "lessThanOrEqualTo satisfied",
			constraint: &constraint{version: Version{major: 2}, op: lessThanOrEqualTo},
			version:    Version{major: 1},
			want:       true,
		},
		{
			name:       "lessThanOrEqualTo just satisfied",
			constraint: &constraint{version: Version{major: 2}, op: lessThanOrEqualTo},
			version:    Version{major: 2},
			want:       true,
		},
		{
			name:       "lessThanOrEqualTo just not satisfied",
			constraint: &constraint{version: Version{major: 2}, op: lessThanOrEqualTo},
			version:    Version{major: 2, modifier: modifierPatch},
			want:       false,
		},
		{
			name:       "lessThanOrEqualTo not satisfied",
			constraint: &constraint{version: Version{major: 2}, op: lessThanOrEqualTo},
			version:    Version{major: 3},
			want:       false,
		},

		{
			name:       "greaterThan satisfied",
			constraint: &constraint{version: Version{major: 2}, op: greaterThan},
			version:    Version{major: 3},
			want:       true,
		},
		{
			name:       "greaterThan just satisfied",
			constraint: &constraint{version: Version{major: 2}, op: greaterThan},
			version:    Version{major: 2, modifier: modifierPatch},
			want:       true,
		},
		{
			name:       "greaterThan just not satisfied",
			constraint: &constraint{version: Version{major: 2}, op: greaterThan},
			version:    Version{major: 2},
			want:       false,
		},
		{
			name:       "greaterThan not satisfied",
			constraint: &constraint{version: Version{major: 2}, op: greaterThan},
			version:    Version{major: 1},
			want:       false,
		},

		{
			name:       "greaterThanOrEqualTo satisfied",
			constraint: &constraint{version: Version{major: 2}, op: greaterThanOrEqualTo},
			version:    Version{major: 3},
			want:       true,
		},
		{
			name:       "greaterThanOrEqualTo just satisfied",
			constraint: &constraint{version: Version{major: 2}, op: greaterThanOrEqualTo},
			version:    Version{major: 2},
			want:       true,
		},
		{
			name:       "greaterThanOrEqualTo just not satisfied",
			constraint: &constraint{version: Version{major: 2}, op: greaterThanOrEqualTo},
			version:    Version{major: 2, modifier: modifierRC},
			want:       false,
		},
		{
			name:       "greaterThanOrEqualTo not satisfied",
			constraint: &constraint{version: Version{major: 2}, op: greaterThanOrEqualTo},
			version:    Version{major: 1},
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.constraint.Check(tt.version); got != tt.want {
				t.Errorf("Check() = %v, want %v", got, tt.want)
			}
		})
	}
}
