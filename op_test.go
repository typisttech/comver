package comver

import "testing"

func Test_op_compare(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		o     op
		other op
		want  int
	}{
		{
			name:  "greaterThanOrEqualTo_compare_greaterThanOrEqualTo",
			o:     greaterThanOrEqualTo,
			other: greaterThanOrEqualTo,
			want:  0,
		},
		{
			name:  "greaterThanOrEqualTo_compare_greaterThan",
			o:     greaterThanOrEqualTo,
			other: greaterThan,
			want:  -1,
		},
		{
			name:  "greaterThanOrEqualTo_compare_lessThan",
			o:     greaterThanOrEqualTo,
			other: lessThan,
			want:  -1,
		},
		{
			name:  "greaterThanOrEqualTo_compare_lessThanOrEqualTo",
			o:     greaterThanOrEqualTo,
			other: lessThanOrEqualTo,
			want:  -1,
		},
		{
			name:  "greaterThan_compare_greaterThanOrEqualTo",
			o:     greaterThan,
			other: greaterThanOrEqualTo,
			want:  1,
		},
		{
			name:  "greaterThan_compare_greaterThan",
			o:     greaterThan,
			other: greaterThan,
			want:  0,
		},
		{
			name:  "greaterThan_compare_lessThan",
			o:     greaterThan,
			other: lessThan,
			want:  -1,
		},
		{
			name:  "greaterThan_compare_lessThanOrEqualTo",
			o:     greaterThan,
			other: lessThanOrEqualTo,
			want:  -1,
		},
		{
			name:  "lessThan_compare_greaterThanOrEqualTo",
			o:     lessThan,
			other: greaterThanOrEqualTo,
			want:  1,
		},
		{
			name:  "lessThan_compare_greaterThan",
			o:     lessThan,
			other: greaterThan,
			want:  1,
		},
		{
			name:  "lessThan_compare_lessThan",
			o:     lessThan,
			other: lessThan,
			want:  0,
		},
		{
			name:  "lessThan_compare_lessThanOrEqualTo",
			o:     lessThan,
			other: lessThanOrEqualTo,
			want:  -1,
		},
		{
			name:  "lessThanOrEqualTo_compare_greaterThanOrEqualTo",
			o:     lessThanOrEqualTo,
			other: greaterThanOrEqualTo,
			want:  1,
		},
		{
			name:  "lessThanOrEqualTo_compare_greaterThan",
			o:     lessThanOrEqualTo,
			other: greaterThan,
			want:  1,
		},
		{
			name:  "lessThanOrEqualTo_compare_lessThan",
			o:     lessThanOrEqualTo,
			other: lessThan,
			want:  1,
		},
		{
			name:  "lessThanOrEqualTo_compare_lessThanOrEqualTo",
			o:     lessThanOrEqualTo,
			other: lessThanOrEqualTo,
			want:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.o.compare(tt.other); got != tt.want {
				t.Errorf("compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
