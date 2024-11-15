package comver

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func newIntervalTestCases() []struct {
	c1      *constraint
	c2      *constraint
	want    interval
	wantErr bool
} {
	return []struct {
		c1      *constraint
		c2      *constraint
		want    interval
		wantErr bool
	}{
		{
			c1:      nil,
			c2:      nil,
			want:    interval{},
			wantErr: false,
		},
		// With single nil.
		{
			c1:      &constraint{Version{major: 9}, lessThanOrEqualTo},
			c2:      nil,
			want:    interval{&constraint{Version{major: 9}, lessThanOrEqualTo}},
			wantErr: false,
		},
		{
			c1:      &constraint{Version{major: 9}, lessThan},
			c2:      nil,
			want:    interval{&constraint{Version{major: 9}, lessThan}},
			wantErr: false,
		},
		{
			c1:      &constraint{Version{major: 9}, greaterThan},
			c2:      nil,
			want:    interval{&constraint{Version{major: 9}, greaterThan}},
			wantErr: false,
		},
		{
			c1:      &constraint{Version{major: 9}, greaterThanOrEqualTo},
			c2:      nil,
			want:    interval{&constraint{Version{major: 9}, greaterThanOrEqualTo}},
			wantErr: false,
		},

		// Same constraint.
		{
			c1:      &constraint{Version{major: 8}, lessThanOrEqualTo},
			c2:      &constraint{Version{major: 8}, lessThanOrEqualTo},
			want:    interval{&constraint{Version{major: 8}, lessThanOrEqualTo}},
			wantErr: false,
		},
		{
			c1:      &constraint{Version{major: 8}, lessThan},
			c2:      &constraint{Version{major: 8}, lessThan},
			want:    interval{&constraint{Version{major: 8}, lessThan}},
			wantErr: false,
		},
		{
			c1:      &constraint{Version{major: 8}, greaterThan},
			c2:      &constraint{Version{major: 8}, greaterThan},
			want:    interval{&constraint{Version{major: 8}, greaterThan}},
			wantErr: false,
		},
		{
			c1:      &constraint{Version{major: 8}, greaterThanOrEqualTo},
			c2:      &constraint{Version{major: 8}, greaterThanOrEqualTo},
			want:    interval{&constraint{Version{major: 8}, greaterThanOrEqualTo}},
			wantErr: false,
		},

		// Same direction. Different versions. Same op.
		{
			c1:      &constraint{Version{major: 7}, lessThanOrEqualTo},
			c2:      &constraint{Version{major: 6}, lessThanOrEqualTo},
			want:    interval{&constraint{Version{major: 6}, lessThanOrEqualTo}},
			wantErr: false,
		},
		{
			c1:      &constraint{Version{major: 7}, lessThan},
			c2:      &constraint{Version{major: 6}, lessThan},
			want:    interval{&constraint{Version{major: 6}, lessThan}},
			wantErr: false,
		},
		{
			c1:      &constraint{Version{major: 7}, greaterThan},
			c2:      &constraint{Version{major: 6}, greaterThan},
			want:    interval{&constraint{Version{major: 7}, greaterThan}},
			wantErr: false,
		},
		{
			c1:      &constraint{Version{major: 7}, greaterThanOrEqualTo},
			c2:      &constraint{Version{major: 6}, greaterThanOrEqualTo},
			want:    interval{&constraint{Version{major: 7}, greaterThanOrEqualTo}},
			wantErr: false,
		},

		// Same direction. Same version. Different ops.
		{
			c1:      &constraint{Version{major: 5}, lessThanOrEqualTo},
			c2:      &constraint{Version{major: 5}, lessThan},
			want:    interval{&constraint{Version{major: 5}, lessThan}},
			wantErr: false,
		},
		{
			c1:      &constraint{Version{major: 5}, greaterThan},
			c2:      &constraint{Version{major: 5}, greaterThanOrEqualTo},
			want:    interval{&constraint{Version{major: 5}, greaterThan}},
			wantErr: false,
		},

		// Different directions. Same version. Different ops.
		{
			c1:      &constraint{Version{major: 5}, lessThanOrEqualTo},
			c2:      &constraint{Version{major: 5}, greaterThan},
			want:    interval{},
			wantErr: true,
		},
		{
			c1:      &constraint{Version{major: 5}, lessThan},
			c2:      &constraint{Version{major: 5}, greaterThan},
			want:    interval{},
			wantErr: true,
		},
		{
			c1:      &constraint{Version{major: 5}, lessThan},
			c2:      &constraint{Version{major: 5}, greaterThanOrEqualTo},
			want:    interval{},
			wantErr: true,
		},
		{
			c1:      &constraint{Version{major: 5}, lessThanOrEqualTo},
			c2:      &constraint{Version{major: 5}, greaterThanOrEqualTo},
			want:    interval{&constraint{Version{major: 5}, greaterThanOrEqualTo}, &constraint{Version{major: 5}, lessThanOrEqualTo}},
			wantErr: false,
		},

		// Different directions. Different versions. Different ops.
		{
			c1:      &constraint{Version{major: 4}, lessThanOrEqualTo},
			c2:      &constraint{Version{major: 3}, greaterThan},
			want:    interval{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			wantErr: false,
		},
		{
			c1:      &constraint{Version{major: 4}, lessThanOrEqualTo},
			c2:      &constraint{Version{major: 3}, greaterThanOrEqualTo},
			want:    interval{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThanOrEqualTo}},
			wantErr: false,
		},
		{
			c1:      &constraint{Version{major: 4}, lessThan},
			c2:      &constraint{Version{major: 3}, greaterThan},
			want:    interval{&constraint{Version{major: 3}, greaterThan}, &constraint{Version{major: 4}, lessThan}},
			wantErr: false,
		},
		{
			c1:      &constraint{Version{major: 4}, lessThan},
			c2:      &constraint{Version{major: 3}, greaterThanOrEqualTo},
			want:    interval{&constraint{Version{major: 3}, greaterThanOrEqualTo}, &constraint{Version{major: 4}, lessThan}},
			wantErr: false,
		},
		{
			c1:      &constraint{Version{major: 1}, lessThanOrEqualTo},
			c2:      &constraint{Version{major: 2}, greaterThan},
			want:    interval{},
			wantErr: true,
		},
		{
			c1:      &constraint{Version{major: 1}, lessThanOrEqualTo},
			c2:      &constraint{Version{major: 2}, greaterThanOrEqualTo},
			want:    interval{},
			wantErr: true,
		},
		{
			c1:      &constraint{Version{major: 1}, lessThan},
			c2:      &constraint{Version{major: 2}, greaterThan},
			want:    interval{},
			wantErr: true,
		},
		{
			c1:      &constraint{Version{major: 1}, lessThan},
			c2:      &constraint{Version{major: 2}, greaterThanOrEqualTo},
			want:    interval{},
			wantErr: true,
		},
	}
}

func TestNewInterval(t *testing.T) {
	t.Parallel()

	for _, tt := range newIntervalTestCases() {
		t.Run(fmt.Sprintf("%s_%s_reversed", tt.c1, tt.c2), func(t *testing.T) {
			t.Parallel()

			got, err := NewInterval(tt.c1, tt.c2)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewInterval(%q, %q) error = %v, wantErr %v", tt.c1, tt.c2, err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInterval(%q, %q) = %v, want %v", tt.c1, tt.c2, got, tt.want)
			}
		})
	}
}

func TestNewInterval_reversed(t *testing.T) {
	t.Parallel()

	for _, tt := range newIntervalTestCases() {
		t.Run(fmt.Sprintf("%s_%s_reversed", tt.c2, tt.c1), func(t *testing.T) {
			t.Parallel()

			got, err := NewInterval(tt.c2, tt.c1)

			if (err != nil) != tt.wantErr {
				t.Fatalf("NewInterval(%q, %q) error = %v, wantErr %v", tt.c2, tt.c1, err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInterval(%q, %q) = %v, want %v", tt.c2, tt.c1, got, tt.want)
			}
		})
	}
}

func Test_interval_Check(t *testing.T) { //nolint:maintidx
	t.Parallel()

	tests := []struct {
		i    interval
		v    Version
		want bool
	}{
		{
			i:    interval{},
			v:    Version{major: 9},
			want: true,
		},
		{
			i:    interval{&constraint{Version{major: 9}, lessThanOrEqualTo}},
			v:    Version{major: 10},
			want: false,
		},
		{
			i:    interval{&constraint{Version{major: 9}, lessThanOrEqualTo}},
			v:    Version{major: 9},
			want: true,
		},
		{
			i:    interval{&constraint{Version{major: 9}, lessThanOrEqualTo}},
			v:    Version{major: 8},
			want: true,
		},
		{
			i:    interval{&constraint{Version{major: 9}, lessThan}},
			v:    Version{major: 10},
			want: false,
		},
		{
			i:    interval{&constraint{Version{major: 9}, lessThan}},
			v:    Version{major: 9},
			want: false,
		},
		{
			i:    interval{&constraint{Version{major: 9}, lessThan}},
			v:    Version{major: 8},
			want: true,
		},
		{
			i:    interval{&constraint{Version{major: 9}, greaterThan}},
			v:    Version{major: 10},
			want: true,
		},
		{
			i:    interval{&constraint{Version{major: 9}, greaterThan}},
			v:    Version{major: 9},
			want: false,
		},
		{
			i:    interval{&constraint{Version{major: 9}, greaterThan}},
			v:    Version{major: 8},
			want: false,
		},
		{
			i:    interval{&constraint{Version{major: 9}, greaterThanOrEqualTo}},
			v:    Version{major: 10},
			want: true,
		},
		{
			i:    interval{&constraint{Version{major: 9}, greaterThanOrEqualTo}},
			v:    Version{major: 9},
			want: true,
		},
		{
			i:    interval{&constraint{Version{major: 9}, greaterThanOrEqualTo}},
			v:    Version{major: 8},
			want: false,
		},
		{
			i: interval{
				&constraint{Version{major: 9}, greaterThanOrEqualTo},
				&constraint{Version{major: 9}, lessThanOrEqualTo},
			},
			v:    Version{major: 10},
			want: false,
		},
		{
			i: interval{
				&constraint{Version{major: 9}, greaterThanOrEqualTo},
				&constraint{Version{major: 9}, lessThanOrEqualTo},
			},
			v:    Version{major: 9},
			want: true,
		},
		{
			i: interval{
				&constraint{Version{major: 9}, greaterThanOrEqualTo},
				&constraint{Version{major: 9}, lessThanOrEqualTo},
			},
			v:    Version{major: 8},
			want: false,
		},
		{
			i: interval{
				&constraint{Version{major: 4}, greaterThanOrEqualTo},
				&constraint{Version{major: 6}, lessThanOrEqualTo},
			},
			v:    Version{major: 7},
			want: false,
		},
		{
			i: interval{
				&constraint{Version{major: 4}, greaterThanOrEqualTo},
				&constraint{Version{major: 6}, lessThanOrEqualTo},
			},
			v:    Version{major: 6},
			want: true,
		},
		{
			i: interval{
				&constraint{Version{major: 4}, greaterThanOrEqualTo},
				&constraint{Version{major: 6}, lessThanOrEqualTo},
			},
			v:    Version{major: 5},
			want: true,
		},
		{
			i: interval{
				&constraint{Version{major: 4}, greaterThanOrEqualTo},
				&constraint{Version{major: 6}, lessThanOrEqualTo},
			},
			v:    Version{major: 4},
			want: true,
		},
		{
			i: interval{
				&constraint{Version{major: 4}, greaterThanOrEqualTo},
				&constraint{Version{major: 6}, lessThanOrEqualTo},
			},
			v:    Version{major: 3},
			want: false,
		},
		{
			i: interval{
				&constraint{Version{major: 4}, greaterThanOrEqualTo},
				&constraint{Version{major: 6}, lessThan},
			},
			v:    Version{major: 7},
			want: false,
		},
		{
			i: interval{
				&constraint{Version{major: 4}, greaterThanOrEqualTo},
				&constraint{Version{major: 6}, lessThan},
			},
			v:    Version{major: 6},
			want: false,
		},
		{
			i: interval{
				&constraint{Version{major: 4}, greaterThanOrEqualTo},
				&constraint{Version{major: 6}, lessThan},
			},
			v:    Version{major: 5},
			want: true,
		},
		{
			i: interval{
				&constraint{Version{major: 4}, greaterThanOrEqualTo},
				&constraint{Version{major: 6}, lessThan},
			},
			v:    Version{major: 4},
			want: true,
		},
		{
			i: interval{
				&constraint{Version{major: 4}, greaterThanOrEqualTo},
				&constraint{Version{major: 6}, lessThan},
			},
			v:    Version{major: 3},
			want: false,
		},
		{
			i: interval{
				&constraint{Version{major: 4}, greaterThan},
				&constraint{Version{major: 6}, lessThanOrEqualTo},
			},
			v:    Version{major: 7},
			want: false,
		},
		{
			i: interval{
				&constraint{Version{major: 4}, greaterThan},
				&constraint{Version{major: 6}, lessThanOrEqualTo},
			},
			v:    Version{major: 6},
			want: true,
		},
		{
			i: interval{
				&constraint{Version{major: 4}, greaterThan},
				&constraint{Version{major: 6}, lessThanOrEqualTo},
			},
			v:    Version{major: 5},
			want: true,
		},
		{
			i: interval{
				&constraint{Version{major: 4}, greaterThan},
				&constraint{Version{major: 6}, lessThanOrEqualTo},
			},
			v:    Version{major: 4},
			want: false,
		},
		{
			i: interval{
				&constraint{Version{major: 4}, greaterThan},
				&constraint{Version{major: 6}, lessThanOrEqualTo},
			},
			v:    Version{major: 3},
			want: false,
		},
		{
			i: interval{
				&constraint{Version{major: 4}, greaterThan},
				&constraint{Version{major: 6}, lessThan},
			},
			v:    Version{major: 7},
			want: false,
		},
		{
			i: interval{
				&constraint{Version{major: 4}, greaterThan},
				&constraint{Version{major: 6}, lessThan},
			},
			v:    Version{major: 6},
			want: false,
		},
		{
			i: interval{
				&constraint{Version{major: 4}, greaterThan},
				&constraint{Version{major: 6}, lessThan},
			},
			v:    Version{major: 5},
			want: true,
		},
		{
			i: interval{
				&constraint{Version{major: 4}, greaterThan},
				&constraint{Version{major: 6}, lessThan},
			},
			v:    Version{major: 4},
			want: false,
		},
		{
			i: interval{
				&constraint{Version{major: 4}, greaterThan},
				&constraint{Version{major: 6}, lessThan},
			},
			v:    Version{major: 3},
			want: false,
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			if got := tt.i.Check(tt.v); got != tt.want {
				t.Errorf("%q.Check(%q) = %v, want %v", tt.i, tt.v.Short(), got, tt.want)
			}
		})
	}
}

func Test_interval_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		i    interval
		want string
	}{
		{
			i:    interval{},
			want: "*",
		},
		{
			i:    interval{&constraint{Version{major: 9}, lessThanOrEqualTo}},
			want: "<=9",
		},
		{
			i:    interval{&constraint{Version{major: 9}, lessThan}},
			want: "<9",
		},
		{
			i:    interval{&constraint{Version{major: 9}, greaterThan}},
			want: ">9",
		},
		{
			i:    interval{&constraint{Version{major: 9}, greaterThanOrEqualTo}},
			want: ">=9",
		},
		{
			i: interval{
				&constraint{Version{major: 9}, greaterThanOrEqualTo},
				&constraint{Version{major: 9}, lessThanOrEqualTo},
			},
			want: "9",
		},
		{
			i: interval{
				&constraint{Version{major: 7}, greaterThanOrEqualTo},
				&constraint{Version{major: 8}, lessThanOrEqualTo},
			},
			want: ">=7 <=8",
		},
		{
			i: interval{
				&constraint{Version{major: 7}, greaterThanOrEqualTo},
				&constraint{Version{major: 8}, lessThan},
			},
			want: ">=7 <8",
		},
		{
			i: interval{
				&constraint{Version{major: 7}, greaterThan},
				&constraint{Version{major: 8}, lessThanOrEqualTo},
			},
			want: ">7 <=8",
		},
		{
			i: interval{
				&constraint{Version{major: 7}, greaterThan},
				&constraint{Version{major: 8}, lessThan},
			},
			want: ">7 <8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()

			if got := tt.i.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
