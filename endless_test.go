package comver

import (
	"reflect"
	"testing"
)

func TestEndless_Check(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		endless Endless
		version Version
		want    bool
	}{
		{
			name:    "lessThan_satisfied",
			endless: NewLessThan(MustParse("2")),
			version: MustParse("1"),
			want:    true,
		},
		{
			name:    "lessThan_just_satisfied",
			endless: NewLessThan(MustParse("2")),
			version: MustParse("2.rc"),
			want:    true,
		},
		{
			name:    "lessThan_just_not_satisfied",
			endless: NewLessThan(MustParse("2")),
			version: MustParse("2"),
			want:    false,
		},
		{
			name:    "lessThan_not_satisfied",
			endless: NewLessThan(MustParse("2")),
			version: MustParse("3"),
			want:    false,
		},

		{
			name:    "lessThanOrEqualTo_satisfied",
			endless: NewLessThanOrEqualTo(MustParse("2")),
			version: MustParse("1"),
			want:    true,
		},
		{
			name:    "lessThanOrEqualTo_just_satisfied",
			endless: NewLessThanOrEqualTo(MustParse("2")),
			version: MustParse("2"),
			want:    true,
		},
		{
			name:    "lessThanOrEqualTo_just_not_satisfied",
			endless: NewLessThanOrEqualTo(MustParse("2")),
			version: MustParse("2.patch"),
			want:    false,
		},
		{
			name:    "lessThanOrEqualTo_not_satisfied",
			endless: NewLessThanOrEqualTo(MustParse("2")),
			version: MustParse("3"),
			want:    false,
		},

		{
			name:    "greaterThan_satisfied",
			endless: NewGreaterThan(MustParse("2")),
			version: MustParse("3"),
			want:    true,
		},
		{
			name:    "greaterThan_just_satisfied",
			endless: NewGreaterThan(MustParse("2")),
			version: MustParse("2.patch"),
			want:    true,
		},
		{
			name:    "greaterThan_just_not_satisfied",
			endless: NewGreaterThan(MustParse("2")),
			version: MustParse("2"),
			want:    false,
		},
		{
			name:    "greaterThan_not_satisfied",
			endless: NewGreaterThan(MustParse("2")),
			version: MustParse("1"),
			want:    false,
		},

		{
			name:    "greaterThanOrEqualTo_satisfied",
			endless: NewGreaterThanOrEqualTo(MustParse("2")),
			version: MustParse("3"),
			want:    true,
		},
		{
			name:    "greaterThanOrEqualTo_just_satisfied",
			endless: NewGreaterThanOrEqualTo(MustParse("2")),
			version: MustParse("2"),
			want:    true,
		},
		{
			name:    "greaterThanOrEqualTo_just_not_satisfied",
			endless: NewGreaterThanOrEqualTo(MustParse("2")),
			version: MustParse("2.rc"),
			want:    false,
		},
		{
			name:    "greaterThanOrEqualTo_not_satisfied",
			endless: NewGreaterThanOrEqualTo(MustParse("2")),
			version: MustParse("1"),
			want:    false,
		},
		{
			name:    "wildcard",
			endless: NewWildcard(),
			version: MustParse("1"),
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.endless.Check(tt.version); got != tt.want {
				t.Errorf("Check() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndless_Ceiling(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		endless Endless
		want    Endless
	}{
		{
			name:    "lessThan",
			endless: NewLessThan(MustParse("1")),
			want:    NewLessThan(MustParse("1")),
		},
		{
			name:    "lessThanOrEqualTo",
			endless: NewLessThanOrEqualTo(MustParse("2")),
			want:    NewLessThanOrEqualTo(MustParse("2")),
		},
		{
			name:    "greaterThanOrEqualTo",
			endless: NewGreaterThanOrEqualTo(MustParse("3")),
			want:    NewWildcard(),
		},
		{
			name:    "greaterThan",
			endless: NewGreaterThan(MustParse("4")),
			want:    NewWildcard(),
		},
		{
			name:    "wildcard",
			endless: NewWildcard(),
			want:    NewWildcard(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.endless.ceiling(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ceiling() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndless_Floor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		endless Endless
		want    Endless
	}{
		{
			name:    "lessThan",
			endless: NewLessThan(MustParse("1")),
			want:    NewWildcard(),
		},
		{
			name:    "lessThanOrEqualTo",
			endless: NewLessThanOrEqualTo(MustParse("2")),
			want:    NewWildcard(),
		},
		{
			name:    "greaterThanOrEqualTo",
			endless: NewGreaterThanOrEqualTo(MustParse("3")),
			want:    NewGreaterThanOrEqualTo(MustParse("3")),
		},
		{
			name:    "greaterThan",
			endless: NewGreaterThan(MustParse("4")),
			want:    NewGreaterThan(MustParse("4")),
		},
		{
			name:    "wildcard",
			endless: NewWildcard(),
			want:    NewWildcard(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.endless.floor(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("floor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndless_wildcard(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		endless Endless
		want    bool
	}{
		{
			name:    "lessThan",
			endless: NewLessThan(MustParse("1")),
			want:    false,
		},
		{
			name:    "lessThanOrEqualTo",
			endless: NewLessThanOrEqualTo(MustParse("2")),
			want:    false,
		},
		{
			name:    "greaterThan",
			endless: NewGreaterThan(MustParse("3")),
			want:    false,
		},
		{
			name:    "greaterThanOrEqualTo",
			endless: NewGreaterThanOrEqualTo(MustParse("4")),
			want:    false,
		},
		{
			name:    "wildcard",
			endless: NewWildcard(),
			want:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.endless.wildcard(); got != tt.want {
				t.Errorf("wildcard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndless_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		endless Endless
		want    string
	}{
		{
			endless: NewLessThanOrEqualTo(MustParse("2")),
			want:    "<=2",
		},
		{
			endless: NewLessThan(MustParse("2")),
			want:    "<2",
		},
		{
			endless: NewGreaterThan(MustParse("2")),
			want:    ">2",
		},
		{
			endless: NewGreaterThanOrEqualTo(MustParse("2")),
			want:    ">=2",
		},
		{
			endless: NewWildcard(),
			want:    "*",
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()

			if got := tt.endless.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
