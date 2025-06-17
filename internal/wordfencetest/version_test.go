package main

import (
	"slices"
	"testing"

	"github.com/typisttech/comver"
)

var invalidWordFenceVersions = []string{ //nolint:gochecknoglobals
	".3.1",
	".47.1",
	".48.9",
	".51.1",
	".51",
	".52.5",
	".53.2",
	".53.4",
	"*",
	"* - v1.01",
	"0.1.2 β",
	"08-03-2018",
	"1.0 12319",
	"1.0f",
	"1.1(Beta)",
	"1.3 EN",
	"1.4.9.8.9",
	"1.6.2d",
	"1.6.49.6.2",
	"1.6.49.6.3",
	"1.6.5-6497609",
	"1.6.59.1.1",
	"1.6.59.1.2",
	"1.6.61.1.0",
	"1.6.61.1.1",
	"1.6.61.2.1",
	"1.7.f",
	"1.9.9.4.1",
	"1.9.9.5.1",
	"1.9.9.5.2",
	"1.9.9.5.3",
	"1.9.9.7.7",
	"1.dec.2012",
	"13-07-2019",
	"17-07-2019",
	"2.0.1.8.2",
	"2.0.2.0.1",
	"2.0.5.4.1",
	"2.0e",
	"2.24080000-WP6.6.1",
	"2.35.1.2.3",
	"2.35.1.3.0",
	"2.5d",
	"2.9.9.2.8",
	"2.9.9.2.9",
	"2.9.9.3.4",
	"2.9.9.4.0",
	"2.9.9.4.7",
	"2.9.9.5.0",
	"2.9.9.5.0",
	"2.9.9.5.1",
	"2.9.9.5.2",
	"2.9.9.5.3",
	"2.9.9.5.4",
	"2.9.9.9.9.9.5",
	"2025r1",
	"2025r2",
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
	"4.10.46.decaf",
	"4.2.6.8.1",
	"4.2.6.8.2",
	"4.2.6.9.3",
	"4.23.1.1.23.1",
	"4.3-revision-3",
	"44.0 (17-08-2023)",
	"47.0(20-11-2023)",
	"5 alpha 2",
	"5.0.28.decaf",
	"5.8 beta 1",
	"5.8 beta 2",
	"6.2-revision-5",
	"6.2-revision-9",
	"6.3-revision-0",
	"8..1",
	"p1.2.5",
	"v.1.1",
}

//go:generate go run gen_wordfence.go gen.go
func TestParse_Wordfence(t *testing.T) {
	t.Parallel()

	for _, v := range wordFenceVersions {
		t.Run(v, func(t *testing.T) {
			t.Parallel()

			_, err := comver.Parse(v)

			wantErr := slices.Contains(invalidWordFenceVersions, v)
			if (err != nil) != wantErr {
				t.Fatalf("Parse(%q) error = %v, wantErr %v", v, err, wantErr)
			}
		})
	}
}
