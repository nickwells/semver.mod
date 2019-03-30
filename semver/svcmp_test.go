package semver_test

import (
	"testing"

	"github.com/nickwells/semver.mod/semver"
	"github.com/nickwells/testhelper.mod/testhelper"
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
		testhelper.ID
		a, b         *semver.SV
		shouldBeLess bool
	}{
		{
			ID: testhelper.MkID("equal - no prIDs"),
			a:  v100,
			b:  v100,
		},
		{
			ID:           testhelper.MkID("major versions a<b"),
			a:            v100,
			b:            v200,
			shouldBeLess: true,
		},
		{
			ID: testhelper.MkID("major versions a>b"),
			a:  v200,
			b:  v100,
		},
		{
			ID:           testhelper.MkID("minor versions a<b"),
			a:            v200,
			b:            v210,
			shouldBeLess: true,
		},
		{
			ID: testhelper.MkID("minor versions a>b"),
			a:  v210,
			b:  v200,
		},
		{
			ID:           testhelper.MkID("patch versions a<b"),
			a:            v210,
			b:            v211,
			shouldBeLess: true,
		},
		{
			ID: testhelper.MkID("patch versions a>b"),
			a:  v211,
			b:  v210,
		},
		{
			ID:           testhelper.MkID("prIDs - shorter is less, a<b"),
			a:            v100_alpha,
			b:            v100_alpha_1,
			shouldBeLess: true,
		},
		{
			ID: testhelper.MkID("prIDs - shorter is less, a>b"),
			a:  v100_alpha_1,
			b:  v100_alpha,
		},
		{
			ID: testhelper.MkID(
				"prIDs - numeric is less than alphanumeric, a<b"),
			a:            v100_alpha_1,
			b:            v100_alpha_beta,
			shouldBeLess: true,
		},
		{
			ID: testhelper.MkID(
				"prIDs - numeric is less than alphanumeric, a>b"),
			a: v100_alpha_beta,
			b: v100_alpha_1,
		},
		{
			ID: testhelper.MkID(
				"prIDs - alphanumeric less by lexi order, a<b"),
			a:            v100_alpha_beta,
			b:            v100_beta,
			shouldBeLess: true,
		},
		{
			ID: testhelper.MkID("prIDs - alphanumeric less by lexi order, a>b"),
			a:  v100_beta,
			b:  v100_alpha_beta,
		},
		{
			ID: testhelper.MkID(
				"prIDs - numeric less by numeric order, a<b"),
			a:            v100_beta_2,
			b:            v100_beta_11,
			shouldBeLess: true,
		},
		{
			ID: testhelper.MkID("prIDs - numeric less by numeric order, a>b"),
			a:  v100_beta_11,
			b:  v100_beta_2,
		},
		{
			ID: testhelper.MkID(
				"prIDs - any prID less than none, a<b"),
			a:            v100_rc_1,
			b:            v100,
			shouldBeLess: true,
		},
		{
			ID: testhelper.MkID("prIDs - any prID less than none, a>b"),
			a:  v100,
			b:  v100_rc_1,
		},
	}

	for _, tc := range testCases {
		isLess := semver.Less(tc.a, tc.b)
		if isLess != tc.shouldBeLess {
			t.Log(tc.IDStr())
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
		testhelper.ID
		sv1      semver.SV
		sv2      semver.SV
		expEqual bool
	}{
		{ID: testhelper.MkID("should be equal"),
			sv1: baseSV, sv2: svCopies[0], expEqual: true},
		{ID: testhelper.MkID("Major version differs"),
			sv1: baseSV, sv2: svCopies[1]},
		{ID: testhelper.MkID("Minor version differs"),
			sv1: baseSV, sv2: svCopies[2]},
		{ID: testhelper.MkID("Patch version differs"),
			sv1: baseSV, sv2: svCopies[3]},
		{ID: testhelper.MkID("too few PreRelIDs"),
			sv1: baseSV, sv2: svCopies[4]},
		{ID: testhelper.MkID("too many PreRelIDs"),
			sv1: baseSV, sv2: svCopies[5]},
		{ID: testhelper.MkID("PreRelIDs in wrong order"),
			sv1: baseSV, sv2: svCopies[6]},
		{ID: testhelper.MkID("too few BuildIDs"),
			sv1: baseSV, sv2: svCopies[7]},
		{ID: testhelper.MkID("too many BuildIDs"),
			sv1: baseSV, sv2: svCopies[8]},
		{ID: testhelper.MkID("BuildIDs in wrong order"),
			sv1: baseSV, sv2: svCopies[9]},
	}

	for _, tc := range testCases {
		if semver.Equals(&tc.sv1, &tc.sv2) {
			if tc.expEqual {
				continue
			}
			t.Log(tc.IDStr())
			t.Logf("\t: %s", tc.sv1)
			t.Logf("\t: %s", tc.sv2)
			t.Errorf("\t: were not expected to be equal\n")
		} else {

			if !tc.expEqual {
				continue
			}
			t.Log(tc.IDStr())
			t.Logf("\t: %s", tc.sv1)
			t.Logf("\t: %s", tc.sv2)
			t.Errorf("\t: were expected to be equal\n")
		}
	}
}
