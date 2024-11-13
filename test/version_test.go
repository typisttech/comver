package main

import (
	"slices"
	"testing"

	"github.com/typisttech/comver"
)

var invalidWordFenceVersions = []string{ //nolint:gochecknoglobals
	"*",
	".47.1",
	".48.9",
	".51",
	".51.1",
	".52.5",
	".53.2",
	".53.4",
	"08-03-2018",
	"1.0 12319",
	"1.0f",
	"1.1(Beta)",
	"1.3 EN",
	"1.4.9.8.9",
	"1.6.49.6.2",
	"1.6.49.6.3",
	"1.6.5-6497609",
	"1.6.59.1.1",
	"1.6.59.1.2",
	"1.6.61.1.0",
	"1.6.61.1.1",
	"1.6.61.2.1",
	"1.7.f",
	"13-07-2019",
	"17-07-2019",
	"2.0.1.8.2",
	"2.0.5.4.1",
	"2.35.1.2.3",
	"2.5d",
	"2.9.9.2.8",
	"2.9.9.2.9",
	"2.9.9.3.4",
	"2.9.9.4.0",
	"3.0 (Beta r7)",
	"3.1.0.1.1",
	"3.1.1.4.2",
	"3.1.37.11.L",
	"3.1.37.12.L",
	"3.2.8.3.1",
	"3.5.5.5.1",
	"3.9.9.0.1",
	"4.1.7.3.2",
	"4.10.44.decaf",
	"5 alpha 2",
	"5.8 beta 1",
	"5.8 beta 2",
	"8..1",
	"p1.2.5",
}

//go:generate go run gen_wordfence.go gen.go
func TestNewVersion_Wordfence(t *testing.T) {
	t.Parallel()

	for _, v := range wordFenceVersions {
		t.Run(v, func(t *testing.T) {
			t.Parallel()

			_, err := comver.NewVersion(v)

			wantErr := slices.Contains(invalidWordFenceVersions, v)
			if (err != nil) != wantErr {
				t.Fatalf("NewVersion(%q) error = %v, wantErr %v", v, err, wantErr)
			}
		})
	}
}
