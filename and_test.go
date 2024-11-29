package comver

import (
	"errors"
	"math/rand/v2"
	"reflect"
	"slices"
	"testing"
)

func TestAnd(t *testing.T) {
	t.Parallel()

	var nilC CeilingFloorConstrainter

	tests := []struct {
		name    string
		es      []Endless
		want    CeilingFloorConstrainter
		wantErr error
	}{
		{
			name:    "empty",
			es:      []Endless{},
			want:    nilC,
			wantErr: errNoEndlessGiven,
		},
		{
			name:    "single_lessThan",
			es:      []Endless{NewLessThan(MustParse("1"))},
			want:    NewLessThan(MustParse("1")),
			wantErr: nil,
		},
		{
			name:    "single_lessThanOrEqualTo",
			es:      []Endless{NewLessThanOrEqualTo(MustParse("2"))},
			want:    NewLessThanOrEqualTo(MustParse("2")),
			wantErr: nil,
		},
		{
			name:    "single_greaterThan",
			es:      []Endless{NewGreaterThan(MustParse("3"))},
			want:    NewGreaterThan(MustParse("3")),
			wantErr: nil,
		},
		{
			name:    "single_greaterThanOrEqualTo",
			es:      []Endless{NewGreaterThanOrEqualTo(MustParse("4"))},
			want:    NewGreaterThanOrEqualTo(MustParse("4")),
			wantErr: nil,
		},
		{
			name:    "single_wildcard",
			es:      []Endless{NewWildcard()},
			want:    NewWildcard(),
			wantErr: nil,
		},
		{
			name:    "multiple_wildcards",
			es:      []Endless{NewWildcard(), NewWildcard()},
			want:    NewWildcard(),
			wantErr: nil,
		},
		{
			name: "multiple_ceiling_inclusive",
			es: []Endless{
				NewLessThan(MustParse("1")),
				NewLessThanOrEqualTo(MustParse("1")),
				NewLessThan(MustParse("2")),
				NewLessThanOrEqualTo(MustParse("2")),
			},
			want:    NewLessThan(MustParse("1")),
			wantErr: nil,
		},
		{
			name: "multiple_ceiling_non_inclusive",
			es: []Endless{
				NewLessThan(MustParse("1")),
				NewLessThanOrEqualTo(MustParse("1")),
				NewLessThan(MustParse("2")),
				NewLessThanOrEqualTo(MustParse("2")),
			},
			want:    NewLessThan(MustParse("1")),
			wantErr: nil,
		},
		{
			name: "multiple_floors_inclusive",
			es: []Endless{
				NewGreaterThanOrEqualTo(MustParse("3")),
				NewGreaterThan(MustParse("3")),
				NewGreaterThanOrEqualTo(MustParse("4")),
			},
			want:    NewGreaterThanOrEqualTo(MustParse("4")),
			wantErr: nil,
		},
		{
			name: "multiple_floors_non_inclusive",
			es: []Endless{
				NewGreaterThanOrEqualTo(MustParse("3")),
				NewGreaterThan(MustParse("3")),
				NewGreaterThanOrEqualTo(MustParse("4")),
				NewGreaterThan(MustParse("4")),
			},
			want:    NewGreaterThan(MustParse("4")),
			wantErr: nil,
		},
		{
			name: "impossible_interval",
			es: []Endless{
				NewLessThan(MustParse("1")),
				NewGreaterThan(MustParse("4")),
			},
			want:    nilC,
			wantErr: errImpossibleInterval,
		},
		{
			name: "impossible_interval_wildcard",
			es: []Endless{
				NewLessThan(MustParse("1")),
				NewGreaterThan(MustParse("4")),
				NewWildcard(),
			},
			want:    nilC,
			wantErr: errImpossibleInterval,
		},
		{
			name: "multiple_wildcards_multiple_endlesses",
			es: []Endless{
				NewLessThan(MustParse("1")),
				NewLessThanOrEqualTo(MustParse("1")),
				NewLessThan(MustParse("2")),
				NewLessThanOrEqualTo(MustParse("2")),
				NewGreaterThanOrEqualTo(MustParse("3")),
				NewGreaterThan(MustParse("3")),
				NewGreaterThanOrEqualTo(MustParse("4")),
				NewGreaterThan(MustParse("4")),
				NewWildcard(),
				NewWildcard(),
			},
			want:    nilC,
			wantErr: errImpossibleInterval,
		},
		{
			name: "impossible_interval_same_version",
			es: []Endless{
				NewLessThan(MustParse("1")),
				NewGreaterThan(MustParse("1")),
			},
			want:    nilC,
			wantErr: errImpossibleInterval,
		},

		{
			name: "exact_version_ceiling_inclusive",
			es: []Endless{
				NewLessThanOrEqualTo(MustParse("1")),
				NewGreaterThan(MustParse("1")),
			},
			want:    nilC,
			wantErr: errImpossibleInterval,
		},
		{
			name: "exact_version_floor_inclusive",
			es: []Endless{
				NewLessThan(MustParse("1")),
				NewGreaterThanOrEqualTo(MustParse("1")),
			},
			want:    nilC,
			wantErr: errImpossibleInterval,
		},
		{
			name: "exact_version_ceiling_inclusive_floor_inclusive",
			es: []Endless{
				NewLessThanOrEqualTo(MustParse("1")),
				NewGreaterThanOrEqualTo(MustParse("1")),
			},
			want:    NewExactConstraint(MustParse("1")),
			wantErr: nil,
		},
		{
			name: "exact_version_wildcard",
			es: []Endless{
				NewLessThanOrEqualTo(MustParse("1")),
				NewGreaterThanOrEqualTo(MustParse("1")),
				NewWildcard(),
			},
			want:    NewExactConstraint(MustParse("1")),
			wantErr: nil,
		},
		{
			name: "impossible_interval_same_version",
			es: []Endless{
				NewLessThanOrEqualTo(MustParse("1")),
				NewLessThanOrEqualTo(MustParse("1")),
				NewLessThan(MustParse("1")),
				NewLessThan(MustParse("1")),
				NewGreaterThanOrEqualTo(MustParse("1")),
				NewGreaterThanOrEqualTo(MustParse("1")),
				NewGreaterThan(MustParse("1")),
				NewGreaterThan(MustParse("1")),
				NewWildcard(),
				NewWildcard(),
			},
			want:    nilC,
			wantErr: errImpossibleInterval,
		},
		{
			name: "impossible_interval_same_version",
			es: []Endless{
				NewLessThan(MustParse("1")),
				NewLessThan(MustParse("1")),
				NewGreaterThanOrEqualTo(MustParse("1")),
				NewGreaterThanOrEqualTo(MustParse("1")),
				NewGreaterThan(MustParse("1")),
				NewGreaterThan(MustParse("1")),
				NewWildcard(),
				NewWildcard(),
			},
			want:    nilC,
			wantErr: errImpossibleInterval,
		},
		{
			name: "exact_version_multiple",
			es: []Endless{
				NewLessThanOrEqualTo(MustParse("1")),
				NewLessThanOrEqualTo(MustParse("1")),
				NewGreaterThanOrEqualTo(MustParse("1")),
				NewGreaterThanOrEqualTo(MustParse("1")),
				NewWildcard(),
				NewWildcard(),
			},
			want:    NewExactConstraint(MustParse("1")),
			wantErr: nil,
		},
		{
			name: "impossible_interval_same_version",
			es: []Endless{
				NewLessThanOrEqualTo(MustParse("1")),
				NewLessThanOrEqualTo(MustParse("1")),
				NewLessThan(MustParse("1")),
				NewLessThan(MustParse("1")),
				NewWildcard(),
				NewWildcard(),
				NewGreaterThan(MustParse("1")),
			},
			want:    nilC,
			wantErr: errImpossibleInterval,
		},
		{
			name: "impossible_interval_same_version",
			es: []Endless{
				NewGreaterThanOrEqualTo(MustParse("1")),
				NewGreaterThanOrEqualTo(MustParse("1")),
				NewGreaterThan(MustParse("1")),
				NewGreaterThan(MustParse("1")),
				NewWildcard(),
				NewWildcard(),
				NewLessThan(MustParse("1")),
			},
			want:    nilC,
			wantErr: errImpossibleInterval,
		},
		{
			name: "impossible_interval_same_version",
			es: []Endless{
				NewLessThanOrEqualTo(MustParse("1")),
				NewLessThanOrEqualTo(MustParse("1")),
				NewLessThan(MustParse("1")),
				NewLessThan(MustParse("1")),
				NewWildcard(),
				NewWildcard(),
				NewGreaterThanOrEqualTo(MustParse("1")),
				NewGreaterThanOrEqualTo(MustParse("1")),
				NewGreaterThan(MustParse("1")),
				NewGreaterThan(MustParse("1")),
			},
			want:    nilC,
			wantErr: errImpossibleInterval,
		},
		{
			name: "impossible_interval_same_version",
			es: []Endless{
				NewLessThanOrEqualTo(MustParse("1")),
				NewLessThanOrEqualTo(MustParse("1")),
				NewLessThan(MustParse("1")),
				NewLessThan(MustParse("1")),
				NewGreaterThanOrEqualTo(MustParse("1")),
				NewGreaterThanOrEqualTo(MustParse("1")),
				NewGreaterThan(MustParse("1")),
				NewGreaterThan(MustParse("1")),
			},
			want:    nilC,
			wantErr: errImpossibleInterval,
		},
		{
			name: "simple",
			es: []Endless{
				NewLessThan(MustParse("4")),
				NewGreaterThan(MustParse("1")),
			},
			want: interval{
				upper: NewLessThan(MustParse("4")),
				lower: NewGreaterThan(MustParse("1")),
			},
			wantErr: nil,
		},
		{
			name: "compact_ceiling",
			es: []Endless{
				NewLessThanOrEqualTo(MustParse("4")),
				NewLessThan(MustParse("4")),
				NewLessThanOrEqualTo(MustParse("3")),
				NewLessThan(MustParse("3")),
				NewGreaterThan(MustParse("1")),
			},
			want: interval{
				upper: NewLessThan(MustParse("3")),
				lower: NewGreaterThan(MustParse("1")),
			},
			wantErr: nil,
		},
		{
			name: "compact_ceiling_inclusive",
			es: []Endless{
				NewLessThanOrEqualTo(MustParse("4")),
				NewLessThan(MustParse("4")),
				NewLessThanOrEqualTo(MustParse("3")),
				NewGreaterThan(MustParse("1")),
			},
			want: interval{
				upper: NewLessThanOrEqualTo(MustParse("3")),
				lower: NewGreaterThan(MustParse("1")),
			},
			wantErr: nil,
		},
		{
			name: "compact_floor",
			es: []Endless{
				NewLessThan(MustParse("4")),
				NewGreaterThan(MustParse("2")),
				NewGreaterThanOrEqualTo(MustParse("2")),
				NewGreaterThan(MustParse("1")),
				NewGreaterThanOrEqualTo(MustParse("1")),
			},
			want: interval{
				upper: NewLessThan(MustParse("4")),
				lower: NewGreaterThan(MustParse("2")),
			},
			wantErr: nil,
		},
		{
			name: "compact_floor_non_inclusive",
			es: []Endless{
				NewLessThan(MustParse("4")),
				NewGreaterThanOrEqualTo(MustParse("2")),
				NewGreaterThan(MustParse("1")),
				NewGreaterThanOrEqualTo(MustParse("1")),
			},
			want: interval{
				upper: NewLessThan(MustParse("4")),
				lower: NewGreaterThanOrEqualTo(MustParse("2")),
			},
			wantErr: nil,
		},
		{
			name: "compact_ceiling_floor",
			es: []Endless{
				NewLessThanOrEqualTo(MustParse("4")),
				NewLessThan(MustParse("4")),
				NewLessThanOrEqualTo(MustParse("3")),
				NewLessThan(MustParse("3")),
				NewGreaterThan(MustParse("2")),
				NewGreaterThanOrEqualTo(MustParse("2")),
				NewGreaterThan(MustParse("1")),
				NewGreaterThanOrEqualTo(MustParse("1")),
			},
			want: interval{
				upper: NewLessThan(MustParse("3")),
				lower: NewGreaterThan(MustParse("2")),
			},
			wantErr: nil,
		},
		{
			name: "compact_ceiling_floor_wildcard",
			es: []Endless{
				NewLessThanOrEqualTo(MustParse("4")),
				NewLessThan(MustParse("4")),
				NewLessThanOrEqualTo(MustParse("3")),
				NewLessThan(MustParse("3")),
				NewGreaterThan(MustParse("2")),
				NewGreaterThanOrEqualTo(MustParse("2")),
				NewGreaterThan(MustParse("1")),
				NewGreaterThanOrEqualTo(MustParse("1")),
				NewWildcard(),
			},
			want: interval{
				upper: NewLessThan(MustParse("3")),
				lower: NewGreaterThan(MustParse("2")),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rand.Shuffle(len(tt.es), func(i, j int) {
				tt.es[i], tt.es[j] = tt.es[j], tt.es[i]
			})

			got, err := And(tt.es...)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("And() error = %#v, wantErr %#v", err, tt.wantErr)
			}

			if err != nil {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("And() got = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("no_arg", func(t *testing.T) {
		t.Parallel()

		wantErr := errNoEndlessGiven

		_, err := And()

		if !errors.Is(err, wantErr) {
			t.Errorf("And() error = %#v, wantErr %#v", err, wantErr)
		}
	})
}

func Test_minBoundedCeiling(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name   string
		cs     []Endless
		want   Endless
		wantOk bool
	}
	tests := []testCase{
		{
			name:   "empty",
			cs:     []Endless{},
			want:   NewWildcard(),
			wantOk: false,
		},
		{
			name:   "nil",
			cs:     nil,
			want:   NewWildcard(),
			wantOk: false,
		},
		{
			name:   "single_lessThan",
			cs:     []Endless{NewLessThan(MustParse("1"))},
			want:   NewLessThan(MustParse("1")),
			wantOk: true,
		},
		{
			name:   "single_lessThanOrEqualTo",
			cs:     []Endless{NewLessThanOrEqualTo(MustParse("2"))},
			want:   NewLessThanOrEqualTo(MustParse("2")),
			wantOk: true,
		},
		{
			name:   "single_greaterThan",
			cs:     []Endless{NewGreaterThan(MustParse("3"))},
			want:   NewWildcard(),
			wantOk: false,
		},
		{
			name:   "single_greaterThanOrEqualTo",
			cs:     []Endless{NewGreaterThanOrEqualTo(MustParse("4"))},
			want:   NewWildcard(),
			wantOk: false,
		},
		{
			name:   "single_wildcard",
			cs:     []Endless{NewWildcard()},
			want:   NewWildcard(),
			wantOk: false,
		},
		{
			name: "multiple_no_ceiling",
			cs: []Endless{
				NewGreaterThanOrEqualTo(MustParse("3")),
				NewGreaterThan(MustParse("3")),
				NewGreaterThanOrEqualTo(MustParse("4")),
				NewGreaterThan(MustParse("4")),
			},
			want:   NewWildcard(),
			wantOk: false,
		},
		{
			name: "multiple_no_ceiling_wildcard",
			cs: []Endless{
				NewGreaterThanOrEqualTo(MustParse("3")),
				NewGreaterThan(MustParse("3")),
				NewGreaterThanOrEqualTo(MustParse("4")),
				NewGreaterThan(MustParse("4")),
				NewWildcard(),
			},
			want:   NewWildcard(),
			wantOk: false,
		},
		{
			name: "multiple_ceilings_non_inclusive",
			cs: []Endless{
				NewLessThanOrEqualTo(MustParse("1")),
				NewLessThan(MustParse("2")),
				NewLessThanOrEqualTo(MustParse("2")),
				NewGreaterThanOrEqualTo(MustParse("3")),
				NewGreaterThan(MustParse("3")),
				NewGreaterThanOrEqualTo(MustParse("4")),
				NewGreaterThan(MustParse("4")),
			},
			want:   NewLessThanOrEqualTo(MustParse("1")),
			wantOk: true,
		},
		{
			name: "multiple_ceilings_inclusive",
			cs: []Endless{
				NewLessThan(MustParse("1")),
				NewLessThanOrEqualTo(MustParse("1")),
				NewLessThan(MustParse("2")),
				NewLessThanOrEqualTo(MustParse("2")),
				NewGreaterThanOrEqualTo(MustParse("3")),
				NewGreaterThan(MustParse("3")),
				NewGreaterThanOrEqualTo(MustParse("4")),
				NewGreaterThan(MustParse("4")),
			},
			want:   NewLessThan(MustParse("1")),
			wantOk: true,
		},
		{
			name: "multiple_ceilings_wildcard",
			cs: []Endless{
				NewLessThan(MustParse("1")),
				NewLessThanOrEqualTo(MustParse("1")),
				NewLessThan(MustParse("2")),
				NewLessThanOrEqualTo(MustParse("2")),
				NewGreaterThanOrEqualTo(MustParse("3")),
				NewGreaterThan(MustParse("3")),
				NewGreaterThanOrEqualTo(MustParse("4")),
				NewGreaterThan(MustParse("4")),
				NewWildcard(),
			},
			want:   NewLessThan(MustParse("1")),
			wantOk: true,
		},
		{
			name: "single_ceiling",
			cs: []Endless{
				NewLessThanOrEqualTo(MustParse("2")),
				NewGreaterThanOrEqualTo(MustParse("3")),
				NewGreaterThan(MustParse("3")),
				NewGreaterThanOrEqualTo(MustParse("4")),
				NewGreaterThan(MustParse("4")),
			},
			want:   NewLessThanOrEqualTo(MustParse("2")),
			wantOk: true,
		},
		{
			name: "single_ceiling_wildcard",
			cs: []Endless{
				NewLessThanOrEqualTo(MustParse("2")),
				NewGreaterThanOrEqualTo(MustParse("3")),
				NewGreaterThan(MustParse("3")),
				NewGreaterThanOrEqualTo(MustParse("4")),
				NewGreaterThan(MustParse("4")),
				NewWildcard(),
			},
			want:   NewLessThanOrEqualTo(MustParse("2")),
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rand.Shuffle(len(tt.cs), func(i, j int) {
				tt.cs[i], tt.cs[j] = tt.cs[j], tt.cs[i]
			})

			original := slices.Clone(tt.cs)

			got, gotOk := minBoundedCeiling(tt.cs...)

			if gotOk != tt.wantOk {
				t.Fatalf("minBoundedCeiling() gotOk = %v, want %v", gotOk, tt.wantOk)
			}

			if !tt.wantOk {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("minBoundedCeiling() got = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(tt.cs, original) {
				t.Errorf("minBoundedCeiling() changed the original slice got = %v, want %v", tt.cs, original)
			}
		})
	}
}

func Test_maxBoundedFloor(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name   string
		fs     []Endless
		want   Endless
		wantOk bool
	}
	tests := []testCase{
		{
			name:   "empty",
			fs:     []Endless{},
			want:   NewWildcard(),
			wantOk: false,
		},
		{
			name:   "nil",
			fs:     nil,
			want:   NewWildcard(),
			wantOk: false,
		},
		{
			name:   "single_lessThan",
			fs:     []Endless{NewLessThan(MustParse("1"))},
			want:   NewWildcard(),
			wantOk: false,
		},
		{
			name:   "single_lessThanOrEqualTo",
			fs:     []Endless{NewLessThanOrEqualTo(MustParse("2"))},
			want:   NewWildcard(),
			wantOk: false,
		},
		{
			name:   "single_greaterThan",
			fs:     []Endless{NewGreaterThan(MustParse("3"))},
			want:   NewGreaterThan(MustParse("3")),
			wantOk: true,
		},
		{
			name:   "single_greaterThanOrEqualTo",
			fs:     []Endless{NewGreaterThanOrEqualTo(MustParse("4"))},
			want:   NewGreaterThanOrEqualTo(MustParse("4")),
			wantOk: true,
		},
		{
			name:   "single_wildcard",
			fs:     []Endless{NewWildcard()},
			want:   NewWildcard(),
			wantOk: false,
		},
		{
			name: "multiple_no_floor",
			fs: []Endless{
				NewLessThan(MustParse("1")),
				NewLessThanOrEqualTo(MustParse("1")),
				NewLessThan(MustParse("2")),
				NewLessThanOrEqualTo(MustParse("2")),
			},
			want:   NewWildcard(),
			wantOk: false,
		},
		{
			name: "multiple_no_floor_wildcard",
			fs: []Endless{
				NewLessThan(MustParse("1")),
				NewLessThanOrEqualTo(MustParse("1")),
				NewLessThan(MustParse("2")),
				NewLessThanOrEqualTo(MustParse("2")),
				NewWildcard(),
			},
			want:   NewWildcard(),
			wantOk: false,
		},
		{
			name: "multiple_floors_non_inclusive",
			fs: []Endless{
				NewLessThan(MustParse("1")),
				NewLessThanOrEqualTo(MustParse("1")),
				NewLessThan(MustParse("2")),
				NewLessThanOrEqualTo(MustParse("2")),
				NewGreaterThanOrEqualTo(MustParse("3")),
				NewGreaterThan(MustParse("3")),
				NewGreaterThan(MustParse("4")),
			},
			want:   NewGreaterThan(MustParse("4")),
			wantOk: true,
		},
		{
			name: "multiple_floors_inclusive",
			fs: []Endless{
				NewLessThan(MustParse("1")),
				NewLessThanOrEqualTo(MustParse("1")),
				NewLessThan(MustParse("2")),
				NewLessThanOrEqualTo(MustParse("2")),
				NewGreaterThanOrEqualTo(MustParse("3")),
				NewGreaterThan(MustParse("3")),
				NewGreaterThanOrEqualTo(MustParse("4")),
				NewGreaterThan(MustParse("4")),
			},
			want:   NewGreaterThan(MustParse("4")),
			wantOk: true,
		},
		{
			name: "multiple_floors_wildcard",
			fs: []Endless{
				NewLessThan(MustParse("1")),
				NewLessThanOrEqualTo(MustParse("1")),
				NewLessThan(MustParse("2")),
				NewLessThanOrEqualTo(MustParse("2")),
				NewGreaterThanOrEqualTo(MustParse("3")),
				NewGreaterThan(MustParse("3")),
				NewGreaterThanOrEqualTo(MustParse("4")),
				NewGreaterThan(MustParse("4")),
				NewWildcard(),
			},
			want:   NewGreaterThan(MustParse("4")),
			wantOk: true,
		},
		{
			name: "single_floor",
			fs: []Endless{
				NewLessThan(MustParse("1")),
				NewLessThanOrEqualTo(MustParse("1")),
				NewLessThan(MustParse("2")),
				NewLessThanOrEqualTo(MustParse("2")),
				NewGreaterThan(MustParse("3")),
			},
			want:   NewGreaterThan(MustParse("3")),
			wantOk: true,
		},
		{
			name: "single_floor_wildcard",
			fs: []Endless{
				NewLessThan(MustParse("1")),
				NewLessThanOrEqualTo(MustParse("1")),
				NewLessThan(MustParse("2")),
				NewLessThanOrEqualTo(MustParse("2")),
				NewGreaterThan(MustParse("3")),
				NewWildcard(),
			},
			want:   NewGreaterThan(MustParse("3")),
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rand.Shuffle(len(tt.fs), func(i, j int) {
				tt.fs[i], tt.fs[j] = tt.fs[j], tt.fs[i]
			})

			original := slices.Clone(tt.fs)

			got, gotOk := maxBoundedFloor(tt.fs...)

			if gotOk != tt.wantOk {
				t.Fatalf("maxBoundedFloor() gotOk = %v, want %v", gotOk, tt.wantOk)
			}

			if !tt.wantOk {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("maxBoundedFloor() got = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(tt.fs, original) {
				t.Errorf("maxBoundedFloor() changed the original slice got = %v, want %v", tt.fs, original)
			}
		})
	}
}
