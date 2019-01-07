package semver_test

import (
	"fmt"
	"testing"

	"github.com/nickwells/semver.mod/semver"
)

func TestLess(t *testing.T) {
	v100 := &semver.SV{
		Major: 1,
		Minor: 0,
		Patch: 0,
	}
	v100_alpha := &semver.SV{
		Major:     1,
		Minor:     0,
		Patch:     0,
		PreRelIDs: []string{"alpha"},
	}
	v100_alpha_1 := &semver.SV{
		Major:     1,
		Minor:     0,
		Patch:     0,
		PreRelIDs: []string{"alpha", "1"},
	}
	v100_alpha_beta := &semver.SV{
		Major:     1,
		Minor:     0,
		Patch:     0,
		PreRelIDs: []string{"alpha", "beta"},
	}
	v100_beta := &semver.SV{
		Major:     1,
		Minor:     0,
		Patch:     0,
		PreRelIDs: []string{"beta"},
	}
	v100_beta_2 := &semver.SV{
		Major:     1,
		Minor:     0,
		Patch:     0,
		PreRelIDs: []string{"beta", "2"},
	}
	v100_beta_11 := &semver.SV{
		Major:     1,
		Minor:     0,
		Patch:     0,
		PreRelIDs: []string{"beta", "11"},
	}
	v100_rc_1 := &semver.SV{
		Major:     1,
		Minor:     0,
		Patch:     0,
		PreRelIDs: []string{"rc", "1"},
	}
	v200 := &semver.SV{
		Major: 2,
		Minor: 0,
		Patch: 0,
	}
	v210 := &semver.SV{
		Major: 2,
		Minor: 1,
		Patch: 0,
	}
	v211 := &semver.SV{
		Major: 2,
		Minor: 1,
		Patch: 1,
	}
	testCases := []struct {
		name         string
		a, b         *semver.SV
		shouldBeLess bool
	}{
		{
			name: "equal - no prIDs",
			a:    v100,
			b:    v100,
		},
		{
			name:         "major versions a<b",
			a:            v100,
			b:            v200,
			shouldBeLess: true,
		},
		{
			name: "major versions a>b",
			a:    v200,
			b:    v100,
		},
		{
			name:         "minor versions a<b",
			a:            v200,
			b:            v210,
			shouldBeLess: true,
		},
		{
			name: "minor versions a>b",
			a:    v210,
			b:    v200,
		},
		{
			name:         "patch versions a<b",
			a:            v210,
			b:            v211,
			shouldBeLess: true,
		},
		{
			name: "patch versions a>b",
			a:    v211,
			b:    v210,
		},
		{
			name:         "prIDs - shorter is less, a<b",
			a:            v100_alpha,
			b:            v100_alpha_1,
			shouldBeLess: true,
		},
		{
			name: "prIDs - shorter is less, a>b",
			a:    v100_alpha_1,
			b:    v100_alpha,
		},
		{
			name:         "prIDs - numeric is less than alphanumeric, a<b",
			a:            v100_alpha_1,
			b:            v100_alpha_beta,
			shouldBeLess: true,
		},
		{
			name: "prIDs - numeric is less than alphanumeric, a>b",
			a:    v100_alpha_beta,
			b:    v100_alpha_1,
		},
		{
			name:         "prIDs - alphanumeric less by lexi order, a<b",
			a:            v100_alpha_beta,
			b:            v100_beta,
			shouldBeLess: true,
		},
		{
			name: "prIDs - alphanumeric less by lexi order, a>b",
			a:    v100_beta,
			b:    v100_alpha_beta,
		},
		{
			name:         "prIDs - numeric less by numeric order, a<b",
			a:            v100_beta_2,
			b:            v100_beta_11,
			shouldBeLess: true,
		},
		{
			name: "prIDs - numeric less by numeric order, a>b",
			a:    v100_beta_11,
			b:    v100_beta_2,
		},
		{
			name:         "prIDs - any prID less than none, a<b",
			a:            v100_rc_1,
			b:            v100,
			shouldBeLess: true,
		},
		{
			name: "prIDs - any prID less than none, a>b",
			a:    v100,
			b:    v100_rc_1,
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s :", i, tc.name)
		isLess := semver.Less(tc.a, tc.b)
		if isLess != tc.shouldBeLess {
			t.Log(tcID)
			t.Logf("\t: %s", tc.a)
			t.Logf("\t: %s", tc.b)
			t.Errorf("\t: is less? %t should be less? %t\n",
				isLess, tc.shouldBeLess)
		}
	}
}
func TestEquals(t *testing.T) {
	baseSV := semver.SV{
		Major:     1,
		Minor:     2,
		Patch:     3,
		PreRelIDs: []string{"a", "b"},
		BuildIDs:  []string{"a", "b"},
	}
	var svCopies [10]semver.SV
	for i := range svCopies {
		baseSV.CopyInto(&(svCopies[i]))
	}
	svCopies[1].Major = 9
	svCopies[2].Minor = 9
	svCopies[3].Patch = 9
	svCopies[4].PreRelIDs = []string{"a"}
	svCopies[5].PreRelIDs = []string{"a", "b", "c"}
	svCopies[6].PreRelIDs = []string{"b", "a"}
	svCopies[7].BuildIDs = []string{"a"}
	svCopies[8].BuildIDs = []string{"a", "b", "c"}
	svCopies[9].BuildIDs = []string{"b", "a"}

	testCases := []struct {
		name     string
		sv1      semver.SV
		sv2      semver.SV
		expEqual bool
	}{
		{name: "should be equal",
			sv1: baseSV, sv2: svCopies[0], expEqual: true},
		{name: "Major version differs", sv1: baseSV, sv2: svCopies[1]},
		{name: "Minor version differs", sv1: baseSV, sv2: svCopies[2]},
		{name: "Patch version differs", sv1: baseSV, sv2: svCopies[3]},
		{name: "too few PreRelIDs", sv1: baseSV, sv2: svCopies[4]},
		{name: "too many PreRelIDs", sv1: baseSV, sv2: svCopies[5]},
		{name: "PreRelIDs in wrong order", sv1: baseSV, sv2: svCopies[6]},
		{name: "too few BuildIDs", sv1: baseSV, sv2: svCopies[7]},
		{name: "too many BuildIDs", sv1: baseSV, sv2: svCopies[8]},
		{name: "BuildIDs in wrong order", sv1: baseSV, sv2: svCopies[9]},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s :", i, tc.name)

		if semver.Equals(&tc.sv1, &tc.sv2) {
			if tc.expEqual {
				continue
			}
			t.Log(tcID)
			t.Logf("\t: %s", tc.sv1)
			t.Logf("\t: %s", tc.sv2)
			t.Errorf("\t: were not expected to be equal\n")
		} else {

			if !tc.expEqual {
				continue
			}
			t.Log(tcID)
			t.Logf("\t: %s", tc.sv1)
			t.Logf("\t: %s", tc.sv2)
			t.Errorf("\t: were expected to be equal\n")
		}
	}
}
