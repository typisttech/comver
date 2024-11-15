//go:build ignore

package main

import (
	"encoding/json"
	"net/http"
	"slices"
)

const productionFeed = "https://www.wordfence.com/api/intelligence/v2/vulnerabilities/production"
const scannerFeed = "https://www.wordfence.com/api/intelligence/v2/vulnerabilities/scanner"

type Vulnerability struct {
	Software []struct {
		AffectedVersions map[string]struct {
			FromVersion string `json:"from_version"`
			ToVersion   string `json:"to_version"`
		} `json:"affected_versions"`
	} `json:"software"`
}

func main() {
	pVulns, err := getVulnerabilities(productionFeed)
	if err != nil {
		panic(err)
	}

	sVulns, err := getVulnerabilities(scannerFeed)
	if err != nil {
		panic(err)
	}

	vulns := slices.Concat(pVulns, sVulns)

	err = gen("wordfence_test.go", fileData{
		VariableName: "wordFenceVersions",
		GeneratedBy:  "gen_wordfence.go",
		Sources:      []string{productionFeed, scannerFeed},
		Versions:     getVersions(vulns...),
	})
	if err != nil {
		panic(err)
	}
}

func getVersions(vulns ...Vulnerability) []string {
	vers := make([]string, 0, len(vulns))
	for _, v := range vulns {
		for _, s := range v.Software {
			for _, a := range s.AffectedVersions {
				vers = append(vers, a.ToVersion, a.FromVersion)
			}
		}
	}

	slices.Sort(vers)
	return slices.Compact(vers)
}

func getVulnerabilities(url string) ([]Vulnerability, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var vulns map[string]Vulnerability
	if err := json.NewDecoder(resp.Body).Decode(&vulns); err != nil {
		return nil, err
	}

	vs := make([]Vulnerability, 0, len(vulns))
	for _, v := range vulns {
		vs = append(vs, v)
	}
	return vs, nil
}
