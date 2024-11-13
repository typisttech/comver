package comver

import (
	"math/rand/v2"
	"reflect"
	"slices"
	"testing"
)

func compactTestCases() []struct { //nolint:maintidx
	name string
	is   Intervals
	want Intervals
} {
	return []struct {
		name string
		is   Intervals
		want Intervals
	}{
		{
			name: "nil",
			is:   nil,
			want: Intervals{},
		},
		{
			name: "empty",
			is:   Intervals{},
			want: Intervals{},
		},
		{
			name: "single_wildcard",
			is:   Intervals{{}},
			want: Intervals{{}},
		},
		{
			name: "single_boundless",
			is:   Intervals{{&constraint{Version{major: 9}, lessThanOrEqualTo}}},
			want: Intervals{{&constraint{Version{major: 9}, lessThanOrEqualTo}}},
		},
		{
			name: "single_bounded",
			is: Intervals{
				{&constraint{Version{major: 7}, greaterThanOrEqualTo}, &constraint{Version{major: 8}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 7}, greaterThanOrEqualTo}, &constraint{Version{major: 8}, lessThanOrEqualTo}},
			},
		},
		{
			name: "wildcard_with_<=",
			is:   Intervals{{}, {&constraint{Version{major: 9}, lessThanOrEqualTo}}},
			want: Intervals{{}},
		},
		{
			name: "wildcard_with_<",
			is:   Intervals{{}, {&constraint{Version{major: 9}, lessThan}}},
			want: Intervals{{}},
		},
		{
			name: "wildcard_with_=",
			is: Intervals{
				{},
				{&constraint{Version{major: 9}, greaterThanOrEqualTo}, &constraint{Version{major: 9}, lessThanOrEqualTo}},
			},
			want: Intervals{{}},
		},
		{
			name: "wildcard_with_>",
			is:   Intervals{{}, {&constraint{Version{major: 9}, greaterThan}}},
			want: Intervals{{}},
		},
		{
			name: "wildcard_with_>=",
			is:   Intervals{{}, {&constraint{Version{major: 9}, greaterThanOrEqualTo}}},
			want: Intervals{{}},
		},
		{
			name: "wildcard_with_single_bounded",
			is: Intervals{
				{},
				{&constraint{Version{major: 7}, greaterThanOrEqualTo}, &constraint{Version{major: 8}, lessThanOrEqualTo}},
			},
			want: Intervals{{}},
		},
		{
			name: "wildcard_with_wildcard",
			is:   Intervals{{}, {}},
			want: Intervals{{}},
		},
		{
			name: "<=_<=_same_version",
			is: Intervals{
				{&constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 9}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 9}, lessThanOrEqualTo}}},
		},
		{
			name: "<=_<=_different_versions",
			is: Intervals{
				{&constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 10}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 10}, lessThanOrEqualTo}}},
		},
		{
			name: "<=_<_same_version",
			is: Intervals{
				{&constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 9}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 9}, lessThanOrEqualTo}}},
		},
		{
			name: "<=_<_lower_version",
			is: Intervals{
				{&constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 8}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 9}, lessThanOrEqualTo}}},
		},
		{
			name: "<=_<_higher_version",
			is: Intervals{
				{&constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 10}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 10}, lessThan}}},
		},
		{
			name: "<=_=_same_version",
			is: Intervals{
				{&constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 9}, greaterThanOrEqualTo}, &constraint{Version{major: 9}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 9}, lessThanOrEqualTo}}},
		},
		{
			name: "<=_=_lower_version",
			is: Intervals{
				{&constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 8}, greaterThanOrEqualTo}, &constraint{Version{major: 8}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 9}, lessThanOrEqualTo}}},
		},
		{
			name: "<=_=_higher_version",
			is: Intervals{
				{&constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 10}, greaterThanOrEqualTo}, &constraint{Version{major: 10}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 10}, greaterThanOrEqualTo}, &constraint{Version{major: 10}, lessThanOrEqualTo}},
			},
		},

		{
			name: "<=_>_same_version",
			is: Intervals{
				{&constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 9}, greaterThan}},
			},
			want: Intervals{{}},
		},
		{
			name: "<=_>_lower_version",
			is: Intervals{
				{&constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 8}, greaterThan}},
			},
			want: Intervals{{}},
		},
		{
			name: "<=_>_higher_version",
			is: Intervals{
				{&constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 10}, greaterThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 10}, greaterThan}},
			},
		},

		{
			name: "<=_>=_same_version",
			is: Intervals{
				{&constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 9}, greaterThanOrEqualTo}},
			},
			want: Intervals{{}},
		},
		{
			name: "<=_>=_lower_version",
			is: Intervals{
				{&constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 8}, greaterThanOrEqualTo}},
			},
			want: Intervals{{}},
		},
		{
			name: "<=_>=_higher_version",
			is: Intervals{
				{&constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 10}, greaterThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 10}, greaterThanOrEqualTo}},
			},
		},

		{
			name: "<_<_same_version",
			is: Intervals{
				{&constraint{Version{major: 9}, lessThan}},
				{&constraint{Version{major: 9}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 9}, lessThan}}},
		},
		{
			name: "<_<_different_versions",
			is:   Intervals{{&constraint{Version{major: 9}, lessThan}}, {&constraint{Version{major: 10}, lessThan}}},
			want: Intervals{{&constraint{Version{major: 10}, lessThan}}},
		},

		{
			name: "<_=_same_version",
			is: Intervals{
				{&constraint{Version{major: 9}, lessThan}},
				{&constraint{Version{major: 9}, greaterThanOrEqualTo}, &constraint{Version{major: 9}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 9}, lessThanOrEqualTo}}},
		},
		{
			name: "<_=_lower_version",
			is: Intervals{
				{&constraint{Version{major: 9}, lessThan}},
				{&constraint{Version{major: 8}, greaterThanOrEqualTo}, &constraint{Version{major: 8}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 9}, lessThan}}},
		},
		{
			name: "<_=_higher_version",
			is: Intervals{
				{&constraint{Version{major: 9}, lessThan}},
				{&constraint{Version{major: 10}, greaterThanOrEqualTo}, &constraint{Version{major: 10}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 9}, lessThan}},
				{&constraint{Version{major: 10}, greaterThanOrEqualTo}, &constraint{Version{major: 10}, lessThanOrEqualTo}},
			},
		},

		{
			name: "<_>_same_version",
			is: Intervals{
				{&constraint{Version{major: 9}, lessThan}},
				{&constraint{Version{major: 9}, greaterThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 9}, lessThan}},
				{&constraint{Version{major: 9}, greaterThan}},
			},
		},
		{
			name: "<_>_lower_version",
			is: Intervals{
				{&constraint{Version{major: 9}, lessThan}},
				{&constraint{Version{major: 8}, greaterThan}},
			},
			want: Intervals{{}},
		},
		{
			name: "<_>_higher_version",
			is: Intervals{
				{&constraint{Version{major: 9}, lessThan}},
				{&constraint{Version{major: 10}, greaterThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 9}, lessThan}},
				{&constraint{Version{major: 10}, greaterThan}},
			},
		},

		{
			name: "<_>=_same_version",
			is: Intervals{
				{&constraint{Version{major: 9}, lessThan}},
				{&constraint{Version{major: 9}, greaterThanOrEqualTo}},
			},
			want: Intervals{{}},
		},
		{
			name: "<_>=_lower_version",
			is: Intervals{
				{&constraint{Version{major: 9}, lessThan}},
				{&constraint{Version{major: 8}, greaterThanOrEqualTo}},
			},
			want: Intervals{{}},
		},
		{
			name: "<_>=_higher_version",
			is: Intervals{
				{&constraint{Version{major: 9}, lessThan}},
				{&constraint{Version{major: 10}, greaterThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 9}, lessThan}},
				{&constraint{Version{major: 10}, greaterThanOrEqualTo}},
			},
		},

		{
			name: "=_=_same_version",
			is: Intervals{
				{&constraint{Version{major: 9}, greaterThanOrEqualTo}, &constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 9}, greaterThanOrEqualTo}, &constraint{Version{major: 9}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 9}, greaterThanOrEqualTo}, &constraint{Version{major: 9}, lessThanOrEqualTo}},
			},
		},
		{
			name: "=_=_different_versions",
			is: Intervals{
				{&constraint{Version{major: 9}, greaterThanOrEqualTo}, &constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 10}, greaterThanOrEqualTo}, &constraint{Version{major: 10}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 9}, greaterThanOrEqualTo}, &constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 10}, greaterThanOrEqualTo}, &constraint{Version{major: 10}, lessThanOrEqualTo}},
			},
		},

		{
			name: "=_>_same_version",
			is: Intervals{
				{&constraint{Version{major: 9}, greaterThanOrEqualTo}, &constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 9}, greaterThan}},
			},
			want: Intervals{{&constraint{Version{major: 9}, greaterThanOrEqualTo}}},
		},
		{
			name: "=_>_lower_version",
			is: Intervals{
				{&constraint{Version{major: 9}, greaterThanOrEqualTo}, &constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 8}, greaterThan}},
			},
			want: Intervals{{&constraint{Version{major: 8}, greaterThan}}},
		},
		{
			name: "=_>_higher_version",
			is: Intervals{
				{&constraint{Version{major: 9}, greaterThanOrEqualTo}, &constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 10}, greaterThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 9}, greaterThanOrEqualTo}, &constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 10}, greaterThan}},
			},
		},

		{
			name: "=_>=_same_version",
			is: Intervals{
				{&constraint{Version{major: 9}, greaterThanOrEqualTo}, &constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 9}, greaterThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 9}, greaterThanOrEqualTo}}},
		},
		{
			name: "=_>=_lower_version",
			is: Intervals{
				{&constraint{Version{major: 9}, greaterThanOrEqualTo}, &constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 8}, greaterThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 8}, greaterThanOrEqualTo}}},
		},
		{
			name: "=_>=_higher_version",
			is: Intervals{
				{&constraint{Version{major: 9}, greaterThanOrEqualTo}, &constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 10}, greaterThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 9}, greaterThanOrEqualTo}, &constraint{Version{major: 9}, lessThanOrEqualTo}},
				{&constraint{Version{major: 10}, greaterThanOrEqualTo}},
			},
		},

		{
			name: ">_>_same_version",
			is: Intervals{
				{&constraint{Version{major: 9}, greaterThan}},
				{&constraint{Version{major: 9}, greaterThan}},
			},
			want: Intervals{{&constraint{Version{major: 9}, greaterThan}}},
		},
		{
			name: ">_>_different_versions",
			is: Intervals{
				{&constraint{Version{major: 9}, greaterThan}},
				{&constraint{Version{major: 10}, greaterThan}},
			},
			want: Intervals{{&constraint{Version{major: 9}, greaterThan}}},
		},

		{
			name: ">_>=_same_version",
			is: Intervals{
				{&constraint{Version{major: 9}, greaterThan}},
				{&constraint{Version{major: 9}, greaterThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 9}, greaterThanOrEqualTo}}},
		},
		{
			name: ">_>=_lower_version",
			is: Intervals{
				{&constraint{Version{major: 9}, greaterThan}},
				{&constraint{Version{major: 8}, greaterThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 8}, greaterThanOrEqualTo}}},
		},
		{
			name: ">_>=_higher_version",
			is: Intervals{
				{&constraint{Version{major: 9}, greaterThan}},
				{&constraint{Version{major: 10}, greaterThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 9}, greaterThan}},
			},
		},

		{
			name: ">=_>=_same_version",
			is: Intervals{
				{&constraint{Version{major: 9}, greaterThanOrEqualTo}},
				{&constraint{Version{major: 9}, greaterThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 9}, greaterThanOrEqualTo}}},
		},
		{
			name: ">=_>=_different_versions",
			is: Intervals{
				{&constraint{Version{major: 9}, greaterThanOrEqualTo}},
				{&constraint{Version{major: 10}, greaterThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 9}, greaterThanOrEqualTo}}},
		},

		{
			name: "<=_>=_gap_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 2}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 2}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: "<=_>=_gap_<=_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 2}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 2}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
		},
		{
			name: "<=_>=_gap_<_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 2}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 2}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: "<=_>=_gap_<_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 2}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 2}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
		},

		{
			name: "<=_>_gap_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: "<=_>_gap_<=_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
		},
		{
			name: "<=_>_gap_<_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: "<=_>_gap_<_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
		},

		{
			name: "<_>=_gap_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 2}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 2}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: "<_>=_gap_<=_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 2}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 2}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
		},
		{
			name: "<_>=_gap_<_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 2}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 2}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: "<_>=_gap_<_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 2}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 2}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
		},

		{
			name: "<_>_gap_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: "<_>_gap_<=_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
		},
		{
			name: "<_>_gap_<_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: "<_>_gap_<_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
		},

		{
			name: "<=_>=_seamless_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 2}, lessThanOrEqualTo}},
				{&constraint{Version{major: 2}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: "<=_>=_seamless_<=_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 2}, lessThanOrEqualTo}},
				{&constraint{Version{major: 2}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
		},
		{
			name: "<=_>=_seamless_<_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 2}, lessThanOrEqualTo}},
				{&constraint{Version{major: 2}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: "<=_>=_seamless_<_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 2}, lessThanOrEqualTo}},
				{&constraint{Version{major: 2}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
		},

		{
			name: "<=_>_seamless_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 2}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: "<=_>_seamless_<=_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 2}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
		},
		{
			name: "<=_>_seamless_<_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 2}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 2}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: "<=_>_seamless_<_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 2}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 2}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
		},

		{
			name: "<_>=_seamless_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 2}, lessThanOrEqualTo}},
				{&constraint{Version{major: 2}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: "<_>=_seamless_<=_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 2}, lessThanOrEqualTo}},
				{&constraint{Version{major: 2}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
		},
		{
			name: "<_>=_seamless_<_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 2}, lessThanOrEqualTo}},
				{&constraint{Version{major: 2}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: "<_>=_seamless_<_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 2}, lessThanOrEqualTo}},
				{&constraint{Version{major: 2}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
		},

		{
			name: "<_>_seamless_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 2}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: "<_>_seamless_<=_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 2}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
		},
		{
			name: "<_>_seamless_<_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 2}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 2}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: "<_>_seamless_<_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 2}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 2}, lessThan}},
				{&constraint{Version{major: 2}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
		},

		{
			name: "<=_>=_overlap_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 3}, lessThanOrEqualTo}},
				{&constraint{Version{major: 2}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: "<=_>=_overlap_<=_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 3}, lessThanOrEqualTo}},
				{&constraint{Version{major: 2}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
		},
		{
			name: "<=_>=_overlap_<_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 3}, lessThanOrEqualTo}},
				{&constraint{Version{major: 2}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: "<=_>=_overlap_<_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 3}, lessThanOrEqualTo}},
				{&constraint{Version{major: 2}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
		},

		{
			name: "<=_>_overlap_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 3}, lessThan}},
				{&constraint{Version{major: 2}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: "<=_>_overlap_<=_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 3}, lessThan}},
				{&constraint{Version{major: 2}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
		},
		{
			name: "<=_>_overlap_<_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 3}, lessThan}},
				{&constraint{Version{major: 2}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: "<=_>_overlap_<_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 3}, lessThan}},
				{&constraint{Version{major: 2}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
		},

		{
			name: "<_>=_overlap_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 3}, lessThanOrEqualTo}},
				{&constraint{Version{major: 2}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: "<_>=_overlap_<=_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 3}, lessThanOrEqualTo}},
				{&constraint{Version{major: 2}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
		},
		{
			name: "<_>=_overlap_<_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 3}, lessThanOrEqualTo}},
				{&constraint{Version{major: 2}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: "<_>=_overlap_<_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 3}, lessThanOrEqualTo}},
				{&constraint{Version{major: 2}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
		},

		{
			name: "<_>_overlap_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 3}, lessThan}},
				{&constraint{Version{major: 2}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: "<_>_overlap_<=_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 3}, lessThan}},
				{&constraint{Version{major: 2}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
		},
		{
			name: "<_>_overlap_<_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 3}, lessThan}},
				{&constraint{Version{major: 2}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: "<_>_overlap_<_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 3}, lessThan}},
				{&constraint{Version{major: 2}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
		},

		{
			name: "<=_gap_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 5}, greaterThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
				{&constraint{Version{major: 5}, greaterThanOrEqualTo}},
			},
		},
		{
			name: "<=_gap_<=_>",
			is: Intervals{
				{&constraint{Version{major: 5}, greaterThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
				{&constraint{Version{major: 5}, greaterThanOrEqualTo}},
			},
		},
		{
			name: "<=_gap_<_>=",
			is: Intervals{
				{&constraint{Version{major: 5}, greaterThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
				{&constraint{Version{major: 5}, greaterThanOrEqualTo}},
			},
		},
		{
			name: "<=_gap_<_>",
			is: Intervals{
				{&constraint{Version{major: 5}, greaterThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
				{&constraint{Version{major: 5}, greaterThanOrEqualTo}},
			},
		},

		{
			name: "<_gap_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 5}, greaterThan}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
				{&constraint{Version{major: 5}, greaterThan}},
			},
		},
		{
			name: "<_gap_<=_>",
			is: Intervals{
				{&constraint{Version{major: 5}, greaterThan}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
				{&constraint{Version{major: 5}, greaterThan}},
			},
		},
		{
			name: "<_gap_<_>=",
			is: Intervals{
				{&constraint{Version{major: 5}, greaterThan}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
				{&constraint{Version{major: 5}, greaterThan}},
			},
		},
		{
			name: "<_gap_<_>",
			is: Intervals{
				{&constraint{Version{major: 5}, greaterThan}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
				{&constraint{Version{major: 5}, greaterThan}},
			},
		},

		{
			name: ">=_gap_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: ">=_gap_<=_>",
			is: Intervals{
				{&constraint{Version{major: 1}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
		},
		{
			name: ">=_gap_<_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: ">=_gap_<_>",
			is: Intervals{
				{&constraint{Version{major: 1}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
		},

		{
			name: ">_gap_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, lessThan}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, lessThan}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: ">_gap_<=_>",
			is: Intervals{
				{&constraint{Version{major: 1}, lessThan}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, lessThan}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
		},
		{
			name: ">_gap_<_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, lessThan}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, lessThan}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: ">_gap_<_>",
			is: Intervals{
				{&constraint{Version{major: 1}, lessThan}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, lessThan}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
		},

		{
			name: "<=_cover_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 2}, greaterThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 2}, greaterThanOrEqualTo}}},
		},
		{
			name: "<=_cover_<=_>",
			is: Intervals{
				{&constraint{Version{major: 2}, greaterThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 2}, greaterThanOrEqualTo}}},
		},
		{
			name: "<=_cover_<_>=",
			is: Intervals{
				{&constraint{Version{major: 2}, greaterThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 2}, greaterThanOrEqualTo}}},
		},
		{
			name: "<=_cover_<_>",
			is: Intervals{
				{&constraint{Version{major: 2}, greaterThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 2}, greaterThanOrEqualTo}}},
		},

		{
			name: "<_cover_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 2}, greaterThan}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 2}, greaterThan}}},
		},
		{
			name: "<_cover_<=_>",
			is: Intervals{
				{&constraint{Version{major: 2}, greaterThan}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 2}, greaterThan}}},
		},
		{
			name: "<_cover_<_>=",
			is: Intervals{
				{&constraint{Version{major: 2}, greaterThan}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 2}, greaterThan}}},
		},
		{
			name: "<_cover_<_>",
			is: Intervals{
				{&constraint{Version{major: 2}, greaterThan}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 2}, greaterThan}}},
		},

		{
			name: ">=_cover_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 5}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 5}, lessThanOrEqualTo}}},
		},
		{
			name: ">=_cover_<=_>",
			is: Intervals{
				{&constraint{Version{major: 5}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 5}, lessThanOrEqualTo}}},
		},
		{
			name: ">=_cover_<_>=",
			is: Intervals{
				{&constraint{Version{major: 5}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 5}, lessThanOrEqualTo}}},
		},
		{
			name: ">=_cover_<_>",
			is: Intervals{
				{&constraint{Version{major: 5}, lessThanOrEqualTo}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 5}, lessThanOrEqualTo}}},
		},

		{
			name: ">_cover_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 5}, lessThan}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 5}, lessThan}}},
		},
		{
			name: ">_cover_<=_>",
			is: Intervals{
				{&constraint{Version{major: 5}, lessThan}},
				{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 5}, lessThan}}},
		},
		{
			name: ">_cover_<_>=",
			is: Intervals{
				{&constraint{Version{major: 5}, lessThan}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 5}, lessThan}}},
		},
		{
			name: ">_cover_<_>",
			is: Intervals{
				{&constraint{Version{major: 5}, lessThan}},
				{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 5}, lessThan}}},
		},

		{
			name: "<=_within_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 2}, greaterThanOrEqualTo}},
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 1}, greaterThanOrEqualTo}}},
		},
		{
			name: "<=_within_<=_>",
			is: Intervals{
				{&constraint{Version{major: 2}, greaterThanOrEqualTo}},
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 1}, greaterThanOrEqualTo}}},
		},
		{
			name: "<=_within_<_>=",
			is: Intervals{
				{&constraint{Version{major: 2}, greaterThanOrEqualTo}},
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 1}, greaterThan}}},
		},
		{
			name: "<=_within_<_>",
			is: Intervals{
				{&constraint{Version{major: 2}, greaterThanOrEqualTo}},
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 1}, greaterThan}}},
		},

		{
			name: "<_within_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 2}, greaterThan}},
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 1}, greaterThanOrEqualTo}}},
		},
		{
			name: "<_within_<=_>",
			is: Intervals{
				{&constraint{Version{major: 2}, greaterThan}},
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 1}, greaterThanOrEqualTo}}},
		},
		{
			name: "<_within_<_>=",
			is: Intervals{
				{&constraint{Version{major: 2}, greaterThan}},
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 1}, greaterThan}}},
		},
		{
			name: "<_within_<_>",
			is: Intervals{
				{&constraint{Version{major: 2}, greaterThan}},
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 1}, greaterThan}}},
		},

		{
			name: ">=_within_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 3}, lessThanOrEqualTo}},
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 4}, lessThanOrEqualTo}}},
		},
		{
			name: ">=_within_<=_>",
			is: Intervals{
				{&constraint{Version{major: 3}, lessThanOrEqualTo}},
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 4}, lessThan}}},
		},
		{
			name: ">=_within_<_>=",
			is: Intervals{
				{&constraint{Version{major: 3}, lessThanOrEqualTo}},
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 4}, lessThanOrEqualTo}}},
		},
		{
			name: ">=_within_<_>",
			is: Intervals{
				{&constraint{Version{major: 3}, lessThanOrEqualTo}},
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 4}, lessThan}}},
		},

		{
			name: ">_within_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 3}, lessThan}},
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 4}, lessThanOrEqualTo}}},
		},
		{
			name: ">_within_<=_>",
			is: Intervals{
				{&constraint{Version{major: 3}, lessThan}},
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 4}, lessThan}}},
		},
		{
			name: ">_within_<_>=",
			is: Intervals{
				{&constraint{Version{major: 3}, lessThan}},
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 4}, lessThanOrEqualTo}}},
		},
		{
			name: ">_within_<_>",
			is: Intervals{
				{&constraint{Version{major: 3}, lessThan}},
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 4}, lessThan}}},
		},

		{
			name: "<=_seamless_floor_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}},
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 1}, greaterThanOrEqualTo}}},
		},
		{
			name: "<=_seamless_floor_<=_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}},
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 1}, greaterThanOrEqualTo}}},
		},
		{
			name: "<=_seamless_floor_<_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}},
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 1}, greaterThanOrEqualTo}}},
		},
		{
			name: "<=_seamless_floor_<_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}},
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 1}, greaterThanOrEqualTo}}},
		},

		{
			name: "<_seamless_floor_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}},
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 1}, greaterThanOrEqualTo}}},
		},
		{
			name: "<_seamless_floor_<=_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}},
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 1}, greaterThanOrEqualTo}}},
		},
		{
			name: "<_seamless_floor_<_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}},
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 1}, greaterThan}}},
		},
		{
			name: "<_seamless_floor_<_>",
			is: Intervals{
				{&constraint{Version{major: 1}, greaterThan}},
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 1}, greaterThan}}},
		},

		{
			name: ">=_seamless_floor_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, lessThanOrEqualTo}},
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 4}, lessThanOrEqualTo}}},
		},
		{
			name: ">=_seamless_floor_<=_>",
			is: Intervals{
				{&constraint{Version{major: 1}, lessThanOrEqualTo}},
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 4}, lessThan}}},
		},
		{
			name: ">=_seamless_floor_<_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, lessThanOrEqualTo}},
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 4}, lessThanOrEqualTo}}},
		},
		{
			name: ">=_seamless_floor_<_>",
			is: Intervals{
				{&constraint{Version{major: 1}, lessThanOrEqualTo}},
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 4}, lessThan}}},
		},

		{
			name: ">_seamless_floor_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, lessThan}},
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 4}, lessThanOrEqualTo}}},
		},
		{
			name: ">_seamless_floor_<=_>",
			is: Intervals{
				{&constraint{Version{major: 1}, lessThan}},
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 4}, lessThan}}},
		},
		{
			name: ">_seamless_floor_<_>=",
			is: Intervals{
				{&constraint{Version{major: 1}, lessThan}},
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, lessThan}},
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
		},
		{
			name: ">_seamless_floor_<_>",
			is: Intervals{
				{&constraint{Version{major: 1}, lessThan}},
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, lessThan}},
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
		},

		{
			name: "<=_seamless_ceiling_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 4}, greaterThanOrEqualTo}},
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 1}, greaterThanOrEqualTo}}},
		},
		{
			name: "<=_seamless_ceiling_<=_>",
			is: Intervals{
				{&constraint{Version{major: 4}, greaterThanOrEqualTo}},
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 1}, greaterThanOrEqualTo}}},
		},
		{
			name: "<=_seamless_ceiling_<_>=",
			is: Intervals{
				{&constraint{Version{major: 4}, greaterThanOrEqualTo}},
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 1}, greaterThan}}},
		},
		{
			name: "<=_seamless_ceiling_<_>",
			is: Intervals{
				{&constraint{Version{major: 4}, greaterThanOrEqualTo}},
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 1}, greaterThan}}},
		},

		{
			name: "<_seamless_ceiling_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 4}, greaterThan}},
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 1}, greaterThanOrEqualTo}}},
		},
		{
			name: "<_seamless_ceiling_<=_>",
			is: Intervals{
				{&constraint{Version{major: 4}, greaterThan}},
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
				{&constraint{Version{major: 4}, greaterThan}},
			},
		},
		{
			name: "<_seamless_ceiling_<_>=",
			is: Intervals{
				{&constraint{Version{major: 4}, greaterThan}},
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 1}, greaterThan}}},
		},
		{
			name: "<_seamless_ceiling_<_>",
			is: Intervals{
				{&constraint{Version{major: 4}, greaterThan}},
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
				{&constraint{Version{major: 4}, greaterThan}},
			},
		},

		{
			name: ">=_seamless_ceiling_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 4}, lessThanOrEqualTo}},
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 4}, lessThanOrEqualTo}}},
		},
		{
			name: ">=_seamless_ceiling_<=_>",
			is: Intervals{
				{&constraint{Version{major: 4}, lessThanOrEqualTo}},
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 4}, lessThanOrEqualTo}}},
		},
		{
			name: ">=_seamless_ceiling_<_>=",
			is: Intervals{
				{&constraint{Version{major: 4}, lessThanOrEqualTo}},
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 4}, lessThanOrEqualTo}}},
		},
		{
			name: ">=_seamless_ceiling_<_>",
			is: Intervals{
				{&constraint{Version{major: 4}, lessThanOrEqualTo}},
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 4}, lessThanOrEqualTo}}},
		},

		{
			name: ">_seamless_ceiling_<=_>=",
			is: Intervals{
				{&constraint{Version{major: 4}, lessThan}},
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 4}, lessThanOrEqualTo}}},
		},
		{
			name: ">_seamless_ceiling_<=_>",
			is: Intervals{
				{&constraint{Version{major: 4}, lessThan}},
				{&constraint{Version{major: 1}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 4}, lessThan}}},
		},
		{
			name: ">_seamless_ceiling_<_>=",
			is: Intervals{
				{&constraint{Version{major: 4}, lessThan}},
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			},
			want: Intervals{{&constraint{Version{major: 4}, lessThanOrEqualTo}}},
		},
		{
			name: ">_seamless_ceiling_<_>",
			is: Intervals{
				{&constraint{Version{major: 4}, lessThan}},
				{&constraint{Version{major: 1}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			},
			want: Intervals{{&constraint{Version{major: 4}, lessThan}}},
		},
	}
}

func TestCompact(t *testing.T) {
	t.Parallel()

	for _, tt := range compactTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := Compact(tt.is); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Compact(%q) = %q, want %q", tt.is, got, tt.want)
			}
		})
	}
}

func TestCompact_reverse(t *testing.T) {
	t.Parallel()

	for _, tt := range compactTestCases() {
		t.Run(tt.name+"_reverse", func(t *testing.T) {
			t.Parallel()

			slices.Reverse(tt.is)

			if got := Compact(tt.is); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Compact(%q) = %q, want %q", tt.is, got, tt.want)
			}
		})
	}
}

func TestCompact_shuffle(t *testing.T) {
	t.Parallel()

	for _, tt := range compactTestCases() {
		t.Run(tt.name+"_shuffle", func(t *testing.T) {
			t.Parallel()

			rand.Shuffle(len(tt.is), func(i, j int) {
				tt.is[i], tt.is[j] = tt.is[j], tt.is[i]
			})

			if got := Compact(tt.is); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Compact(%q) = %q, want %q", tt.is, got, tt.want)
			}
		})
	}
}
