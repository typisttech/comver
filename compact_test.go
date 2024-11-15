package comver

import (
	"math/rand/v2"
	"reflect"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCompact(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		o    Or
		want Constrainter
	}{
		{
			name: "empty",
			o:    Or{},
			want: Or{},
		},
		{
			name: "nil",
			o:    nil,
			want: Or{},
		},
		{
			name: "single_lessThan",
			o:    Or{NewLessThan(MustParse("1"))},
			want: NewLessThan(MustParse("1")),
		},
		{
			name: "single_lessThanOrEqualTo",
			o:    Or{NewLessThanOrEqualTo(MustParse("2"))},
			want: NewLessThanOrEqualTo(MustParse("2")),
		},
		{
			name: "single_greaterThan",
			o:    Or{NewGreaterThan(MustParse("3"))},
			want: NewGreaterThan(MustParse("3")),
		},
		{
			name: "single_greaterThanOrEqualTo",
			o:    Or{NewGreaterThanOrEqualTo(MustParse("4"))},
			want: NewGreaterThanOrEqualTo(MustParse("4")),
		},
		{
			name: "single_exact",
			o:    Or{NewExactConstraint(MustParse("5"))},
			want: NewExactConstraint(MustParse("5")),
		},
		{
			name: "single_interval",
			o: Or{
				interval{
					upper: NewLessThan(MustParse("7")),
					lower: NewGreaterThan(MustParse("6")),
				},
			},
			want: interval{
				upper: NewLessThan(MustParse("7")),
				lower: NewGreaterThan(MustParse("6")),
			},
		},
		{
			name: "single_wildcard",
			o:    Or{NewWildcard()},
			want: NewWildcard(),
		},
		{
			name: "multiple_wildcards",
			o:    Or{NewWildcard(), NewWildcard()},
			want: NewWildcard(),
		},
		{
			name: "wildcard_trumps_everything_else",
			o: Or{
				NewLessThan(MustParse("1")),
				NewLessThanOrEqualTo(MustParse("1")),
				NewLessThan(MustParse("2")),
				NewLessThanOrEqualTo(MustParse("2")),
				NewGreaterThanOrEqualTo(MustParse("3")),
				NewGreaterThan(MustParse("3")),
				NewGreaterThanOrEqualTo(MustParse("4")),
				NewGreaterThan(MustParse("4")),
				NewExactConstraint(MustParse("5")),
				interval{
					upper: NewLessThan(MustParse("7")),
					lower: NewGreaterThan(MustParse("6")),
				},
				NewWildcard(),
			},
			want: NewWildcard(),
		},
		{
			name: "match_all",
			o: Or{
				NewLessThan(MustParse("10")),
				NewGreaterThan(MustParse("9")),
			},
			want: NewWildcard(),
		},
		{
			name: "match_all_same_version_ceiling_inclusive_floor_inclusive",
			o: Or{
				NewLessThanOrEqualTo(MustParse("10")),
				NewGreaterThanOrEqualTo(MustParse("10")),
			},
			want: NewWildcard(),
		},
		{
			name: "match_all_same_version_ceiling_non_inclusive_floor_inclusive",
			o: Or{
				NewLessThan(MustParse("10")),
				NewGreaterThanOrEqualTo(MustParse("10")),
			},
			want: NewWildcard(),
		},
		{
			name: "match_all_same_version_ceiling_inclusive_floor_non_inclusive",
			o: Or{
				NewLessThanOrEqualTo(MustParse("10")),
				NewGreaterThan(MustParse("10")),
			},
			want: NewWildcard(),
		},
		{
			name: "same_version_not_match_all",
			o: Or{
				NewGreaterThan(MustParse("10")),
				NewLessThan(MustParse("10")),
			},
			want: Or{
				NewLessThan(MustParse("10")),
				NewGreaterThan(MustParse("10")),
			},
		},
		{
			name: "unrelated_intervals",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("14")),
				},
				interval{
					lower: NewGreaterThan(MustParse("17")),
					upper: NewLessThan(MustParse("20")),
				},
				interval{
					lower: NewGreaterThan(MustParse("23")),
					upper: NewLessThan(MustParse("26")),
				},
			},
			want: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("14")),
				},
				interval{
					lower: NewGreaterThan(MustParse("17")),
					upper: NewLessThan(MustParse("20")),
				},
				interval{
					lower: NewGreaterThan(MustParse("23")),
					upper: NewLessThan(MustParse("26")),
				},
			},
		},
		{
			name: "overlapping_intervals",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("18")),
				},
				interval{
					lower: NewGreaterThan(MustParse("17")),
					upper: NewLessThan(MustParse("20")),
				},
				interval{
					lower: NewGreaterThan(MustParse("23")),
					upper: NewLessThan(MustParse("26")),
				},
			},
			want: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("20")),
				},
				interval{
					lower: NewGreaterThan(MustParse("23")),
					upper: NewLessThan(MustParse("26")),
				},
			},
		},
		{
			name: "overlapping_intervals",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("14")),
				},
				interval{
					lower: NewGreaterThan(MustParse("13")),
					upper: NewLessThan(MustParse("20")),
				},
				interval{
					lower: NewGreaterThan(MustParse("23")),
					upper: NewLessThan(MustParse("26")),
				},
			},
			want: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("20")),
				},
				interval{
					lower: NewGreaterThan(MustParse("23")),
					upper: NewLessThan(MustParse("26")),
				},
			},
		},
		{
			name: "overlapping_intervals",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("14")),
				},
				interval{
					lower: NewGreaterThan(MustParse("17")),
					upper: NewLessThan(MustParse("24")),
				},
				interval{
					lower: NewGreaterThan(MustParse("23")),
					upper: NewLessThan(MustParse("26")),
				},
			},
			want: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("14")),
				},
				interval{
					lower: NewGreaterThan(MustParse("17")),
					upper: NewLessThan(MustParse("26")),
				},
			},
		},
		{
			name: "overlapping_intervals",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("14")),
				},
				interval{
					lower: NewGreaterThan(MustParse("17")),
					upper: NewLessThan(MustParse("20")),
				},
				interval{
					lower: NewGreaterThan(MustParse("19")),
					upper: NewLessThan(MustParse("26")),
				},
			},
			want: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("14")),
				},
				interval{
					lower: NewGreaterThan(MustParse("17")),
					upper: NewLessThan(MustParse("26")),
				},
			},
		},
		{
			name: "overlapping_intervals",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("25")),
				},
				interval{
					lower: NewGreaterThan(MustParse("17")),
					upper: NewLessThan(MustParse("20")),
				},
				interval{
					lower: NewGreaterThan(MustParse("19")),
					upper: NewLessThan(MustParse("26")),
				},
			},
			want: interval{
				lower: NewGreaterThan(MustParse("11")),
				upper: NewLessThan(MustParse("26")),
			},
		},
		{
			name: "overlapping_intervals",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("14")),
				},
				interval{
					lower: NewGreaterThan(MustParse("17")),
					upper: NewLessThan(MustParse("20")),
				},
				interval{
					lower: NewGreaterThan(MustParse("12")),
					upper: NewLessThan(MustParse("26")),
				},
			},
			want: interval{
				lower: NewGreaterThan(MustParse("11")),
				upper: NewLessThan(MustParse("26")),
			},
		},
		{
			name: "overlapping_intervals",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("20")),
				},
				interval{
					lower: NewGreaterThan(MustParse("17")),
					upper: NewLessThan(MustParse("20")),
				},
				interval{
					lower: NewGreaterThan(MustParse("19")),
					upper: NewLessThan(MustParse("26")),
				},
			},
			want: interval{
				lower: NewGreaterThan(MustParse("11")),
				upper: NewLessThan(MustParse("26")),
			},
		},
		{
			name: "continuous_intervals",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThanOrEqualTo(MustParse("17")),
				},
				interval{
					lower: NewGreaterThan(MustParse("17")),
					upper: NewLessThan(MustParse("20")),
				},
				interval{
					lower: NewGreaterThan(MustParse("23")),
					upper: NewLessThan(MustParse("26")),
				},
			},
			want: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("20")),
				},
				interval{
					lower: NewGreaterThan(MustParse("23")),
					upper: NewLessThan(MustParse("26")),
				},
			},
		},
		{
			name: "continuous_intervals",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("17")),
				},
				interval{
					lower: NewGreaterThanOrEqualTo(MustParse("17")),
					upper: NewLessThan(MustParse("20")),
				},
				interval{
					lower: NewGreaterThan(MustParse("23")),
					upper: NewLessThan(MustParse("26")),
				},
			},
			want: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("20")),
				},
				interval{
					lower: NewGreaterThan(MustParse("23")),
					upper: NewLessThan(MustParse("26")),
				},
			},
		},
		{
			name: "continuous_intervals",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("14")),
				},
				interval{
					lower: NewGreaterThan(MustParse("17")),
					upper: NewLessThanOrEqualTo(MustParse("23")),
				},
				interval{
					lower: NewGreaterThan(MustParse("23")),
					upper: NewLessThan(MustParse("26")),
				},
			},
			want: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("14")),
				},
				interval{
					lower: NewGreaterThan(MustParse("17")),
					upper: NewLessThan(MustParse("26")),
				},
			},
		},
		{
			name: "continuous_intervals",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("14")),
				},
				interval{
					lower: NewGreaterThan(MustParse("17")),
					upper: NewLessThanOrEqualTo(MustParse("23")),
				},
				interval{
					lower: NewGreaterThan(MustParse("23")),
					upper: NewLessThan(MustParse("26")),
				},
			},
			want: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("14")),
				},
				interval{
					lower: NewGreaterThan(MustParse("17")),
					upper: NewLessThan(MustParse("26")),
				},
			},
		},
		{
			name: "continuous_intervals",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("14")),
				},
				interval{
					lower: NewGreaterThanOrEqualTo(MustParse("14")),
					upper: NewLessThanOrEqualTo(MustParse("23")),
				},
				interval{
					lower: NewGreaterThan(MustParse("23")),
					upper: NewLessThan(MustParse("26")),
				},
			},
			want: interval{
				lower: NewGreaterThan(MustParse("11")),
				upper: NewLessThan(MustParse("26")),
			},
		},
		{
			name: "tail_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThanOrEqualTo(MustParse("18")),
			},
			want: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThanOrEqualTo(MustParse("18")),
			},
		},
		{
			name: "tail_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThanOrEqualTo(MustParse("17")),
			},
			want: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				NewGreaterThan(MustParse("15")),
			},
		},
		{
			name: "tail_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThanOrEqualTo(MustParse("17")),
				},
				NewGreaterThanOrEqualTo(MustParse("17")),
			},
			want: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				NewGreaterThan(MustParse("15")),
			},
		},
		{
			name: "tail_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThanOrEqualTo(MustParse("16")),
			},
			want: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				NewGreaterThan(MustParse("15")),
			},
		},
		{
			name: "tail_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThanOrEqualTo(MustParse("15")),
			},
			want: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				NewGreaterThanOrEqualTo(MustParse("15")),
			},
		},
		{
			name: "tail_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThanOrEqualTo(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThanOrEqualTo(MustParse("15")),
			},
			want: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				NewGreaterThanOrEqualTo(MustParse("15")),
			},
		},
		{
			name: "tail_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThanOrEqualTo(MustParse("14")),
			},
			want: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				NewGreaterThanOrEqualTo(MustParse("14")),
			},
		},
		{
			name: "tail_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThanOrEqualTo(MustParse("13")),
			},
			want: NewGreaterThan(MustParse("11")),
		},
		{
			name: "tail_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThanOrEqualTo(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThanOrEqualTo(MustParse("13")),
			},
			want: NewGreaterThan(MustParse("11")),
		},
		{
			name: "tail_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThanOrEqualTo(MustParse("12")),
			},
			want: NewGreaterThan(MustParse("11")),
		},
		{
			name: "tail_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThanOrEqualTo(MustParse("11")),
			},
			want: NewGreaterThanOrEqualTo(MustParse("11")),
		},
		{
			name: "tail_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThanOrEqualTo(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThanOrEqualTo(MustParse("11")),
			},
			want: NewGreaterThanOrEqualTo(MustParse("11")),
		},
		{
			name: "tail_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThanOrEqualTo(MustParse("10")),
			},
			want: NewGreaterThanOrEqualTo(MustParse("10")),
		},

		{
			name: "tail_non_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThan(MustParse("18")),
			},
			want: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThan(MustParse("18")),
			},
		},
		{
			name: "tail_non_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThan(MustParse("17")),
			},
			want: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThan(MustParse("17")),
			},
		},
		{
			name: "tail_non_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThanOrEqualTo(MustParse("17")),
				},
				NewGreaterThan(MustParse("17")),
			},
			want: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				NewGreaterThan(MustParse("15")),
			},
		},
		{
			name: "tail_non_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThan(MustParse("16")),
			},
			want: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				NewGreaterThan(MustParse("15")),
			},
		},
		{
			name: "tail_non_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThan(MustParse("15")),
			},
			want: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				NewGreaterThan(MustParse("15")),
			},
		},
		{
			name: "tail_non_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThanOrEqualTo(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThan(MustParse("15")),
			},
			want: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				NewGreaterThanOrEqualTo(MustParse("15")),
			},
		},
		{
			name: "tail_non_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThan(MustParse("14")),
			},
			want: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				NewGreaterThan(MustParse("14")),
			},
		},
		{
			name: "tail_non_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThan(MustParse("13")),
			},
			want: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				NewGreaterThan(MustParse("13")),
			},
		},
		{
			name: "tail_non_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThanOrEqualTo(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThan(MustParse("13")),
			},
			want: NewGreaterThan(MustParse("11")),
		},
		{
			name: "tail_non_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThan(MustParse("12")),
			},
			want: NewGreaterThan(MustParse("11")),
		},
		{
			name: "tail_non_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThan(MustParse("11")),
			},
			want: NewGreaterThan(MustParse("11")),
		},
		{
			name: "tail_non_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThanOrEqualTo(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThan(MustParse("11")),
			},
			want: NewGreaterThanOrEqualTo(MustParse("11")),
		},
		{
			name: "tail_non_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewGreaterThan(MustParse("10")),
			},
			want: NewGreaterThan(MustParse("10")),
		},
		{
			name: "head_non_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewLessThan(MustParse("18")),
			},
			want: NewLessThan(MustParse("18")),
		},
		{
			name: "head_non_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewLessThan(MustParse("17")),
			},
			want: NewLessThan(MustParse("17")),
		},
		{
			name: "head_non_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThanOrEqualTo(MustParse("17")),
				},
				NewLessThan(MustParse("17")),
			},
			want: NewLessThanOrEqualTo(MustParse("17")),
		},
		{
			name: "head_non_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewLessThan(MustParse("16")),
			},
			want: NewLessThan(MustParse("17")),
		},
		{
			name: "head_non_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewLessThan(MustParse("15")),
			},
			want: Or{
				NewLessThan(MustParse("15")),
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
			},
		},
		{
			name: "head_non_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThanOrEqualTo(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewLessThan(MustParse("15")),
			},
			want: NewLessThan(MustParse("17")),
		},
		{
			name: "head_non_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewLessThan(MustParse("14")),
			},
			want: Or{
				NewLessThan(MustParse("14")),
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
			},
		},
		{
			name: "head_non_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewLessThan(MustParse("13")),
			},
			want: Or{
				NewLessThan(MustParse("13")),
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
			},
		},
		{
			name: "head_non_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThanOrEqualTo(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewLessThan(MustParse("13")),
			},
			want: Or{
				NewLessThanOrEqualTo(MustParse("13")),
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
			},
		},
		{
			name: "head_non_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewLessThan(MustParse("12")),
			},
			want: Or{
				NewLessThan(MustParse("13")),
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
			},
		},
		{
			name: "head_non_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewLessThan(MustParse("11")),
			},
			want: Or{
				NewLessThan(MustParse("11")),
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
			},
		},
		{
			name: "head_non_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThanOrEqualTo(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewLessThan(MustParse("11")),
			},
			want: Or{
				NewLessThan(MustParse("13")),
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
			},
		},
		{
			name: "head_non_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewLessThan(MustParse("10")),
			},
			want: Or{
				NewLessThan(MustParse("10")),
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
			},
		},
		{
			name: "head_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewLessThanOrEqualTo(MustParse("18")),
			},
			want: NewLessThanOrEqualTo(MustParse("18")),
		},
		{
			name: "head_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewLessThanOrEqualTo(MustParse("17")),
			},
			want: NewLessThanOrEqualTo(MustParse("17")),
		},
		{
			name: "head_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThanOrEqualTo(MustParse("17")),
				},
				NewLessThanOrEqualTo(MustParse("17")),
			},
			want: NewLessThanOrEqualTo(MustParse("17")),
		},
		{
			name: "head_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewLessThanOrEqualTo(MustParse("16")),
			},
			want: NewLessThan(MustParse("17")),
		},
		{
			name: "head_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewLessThanOrEqualTo(MustParse("15")),
			},
			want: NewLessThan(MustParse("17")),
		},
		{
			name: "head_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThanOrEqualTo(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewLessThanOrEqualTo(MustParse("15")),
			},
			want: NewLessThan(MustParse("17")),
		},
		{
			name: "head_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewLessThanOrEqualTo(MustParse("14")),
			},
			want: Or{
				NewLessThanOrEqualTo(MustParse("14")),
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
			},
		},
		{
			name: "head_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewLessThanOrEqualTo(MustParse("13")),
			},
			want: Or{
				NewLessThanOrEqualTo(MustParse("13")),
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
			},
		},
		{
			name: "head_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThanOrEqualTo(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewLessThanOrEqualTo(MustParse("13")),
			},
			want: Or{
				NewLessThanOrEqualTo(MustParse("13")),
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
			},
		},
		{
			name: "head_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewLessThanOrEqualTo(MustParse("12")),
			},
			want: Or{
				NewLessThan(MustParse("13")),
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
			},
		},
		{
			name: "head_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewLessThanOrEqualTo(MustParse("11")),
			},
			want: Or{
				NewLessThan(MustParse("13")),
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
			},
		},
		{
			name: "head_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThanOrEqualTo(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewLessThanOrEqualTo(MustParse("11")),
			},
			want: Or{
				NewLessThan(MustParse("13")),
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
			},
		},
		{
			name: "head_inclusive",
			o: Or{
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
				NewLessThanOrEqualTo(MustParse("10")),
			},
			want: Or{
				NewLessThanOrEqualTo(MustParse("10")),
				interval{
					lower: NewGreaterThan(MustParse("11")),
					upper: NewLessThan(MustParse("13")),
				},
				interval{
					lower: NewGreaterThan(MustParse("15")),
					upper: NewLessThan(MustParse("17")),
				},
			},
		},
		{
			name: "match_all_with_interval",
			o: Or{
				NewGreaterThan(MustParse("1")),
				interval{
					lower: NewGreaterThan(MustParse("2")),
					upper: NewLessThan(MustParse("3")),
				},
				NewLessThan(MustParse("4")),
			},
			want: NewWildcard(),
		},
		{
			name: "match_all_within_interval",
			o: Or{
				NewGreaterThan(MustParse("2")),
				interval{
					lower: NewGreaterThan(MustParse("1")),
					upper: NewLessThan(MustParse("4")),
				},
				NewLessThan(MustParse("3")),
			},
			want: NewWildcard(),
		},
		{
			name: "match_all_within_intervals",
			o: Or{
				NewLessThan(MustParse("3")),
				interval{
					lower: NewGreaterThan(MustParse("2")),
					upper: NewLessThan(MustParse("6")),
				},
				interval{
					lower: NewGreaterThan(MustParse("5")),
					upper: NewLessThan(MustParse("8")),
				},
				NewGreaterThan(MustParse("7")),
			},
			want: NewWildcard(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rand.Shuffle(len(tt.o), func(i, j int) {
				tt.o[i], tt.o[j] = tt.o[j], tt.o[i]
			})

			original := slices.Clone(tt.o)

			got := Compact(tt.o)

			opts := cmp.Options{
				cmp.Comparer(func(a, b Constrainter) bool {
					return a.String() == b.String()
				}),
			}

			if diff := cmp.Diff(tt.want, got, opts); diff != "" {
				t.Errorf("Compact(%q) mismatch (-want +got):\n%s", original, diff)
			}

			if !reflect.DeepEqual(tt.o, original) {
				t.Errorf("Compact() changed the original slice got = %v, want %v", tt.o, original)
			}
		})
	}
}

func Test_wildcard(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		c    CeilingFloorConstrainter
		want bool
	}{
		{
			name: "lessThan",
			c:    NewLessThan(MustParse("1")),
			want: false,
		},
		{
			name: "lessThanOrEqualTo",
			c:    NewLessThanOrEqualTo(MustParse("2")),
			want: false,
		},
		{
			name: "greaterThan",
			c:    NewGreaterThan(MustParse("3")),
			want: false,
		},
		{
			name: "greaterThanOrEqualTo",
			c:    NewGreaterThanOrEqualTo(MustParse("4")),
			want: false,
		},
		{
			name: "exact",
			c:    NewExactConstraint(MustParse("5")),
			want: false,
		},
		{
			name: "interval",
			c: interval{
				upper: NewLessThan(MustParse("7")),
				lower: NewGreaterThan(MustParse("6")),
			},
			want: false,
		},
		{
			name: "wildcard",
			c:    NewWildcard(),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := wildcard(tt.c); got != tt.want {
				t.Errorf("wildcard() = %v, want %v", got, tt.want)
			}
		})
	}
}
