package semver_test

import (
	"testing"

	"github.com/nickwells/semver.mod/v3/semver"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestLess(t *testing.T) {
	var (
		v100    = semver.NewSVOrPanic(1, 0, 0, nil, nil)
		v100A   = semver.NewSVOrPanic(1, 0, 0, []string{"alpha"}, nil)
		v100A1  = semver.NewSVOrPanic(1, 0, 0, []string{"alpha", "1"}, nil)
		v100AB  = semver.NewSVOrPanic(1, 0, 0, []string{"alpha", "beta"}, nil)
		v100B   = semver.NewSVOrPanic(1, 0, 0, []string{"beta"}, nil)
		v100B2  = semver.NewSVOrPanic(1, 0, 0, []string{"beta", "2"}, nil)
		v100B11 = semver.NewSVOrPanic(1, 0, 0, []string{"beta", "11"}, nil)
		v100RC1 = semver.NewSVOrPanic(1, 0, 0, []string{"rc", "1"}, nil)
		v200    = semver.NewSVOrPanic(2, 0, 0, nil, nil)
		v210    = semver.NewSVOrPanic(2, 1, 0, nil, nil)
		v211    = semver.NewSVOrPanic(2, 1, 1, nil, nil)
	)

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
			a:            v100A,
			b:            v100A1,
			shouldBeLess: true,
		},
		{
			ID: testhelper.MkID("prIDs - shorter is less, a>b"),
			a:  v100A1,
			b:  v100A,
		},
		{
			ID: testhelper.MkID(
				"prIDs - numeric is less than alphanumeric, a<b"),
			a:            v100A1,
			b:            v100AB,
			shouldBeLess: true,
		},
		{
			ID: testhelper.MkID(
				"prIDs - numeric is less than alphanumeric, a>b"),
			a: v100AB,
			b: v100A1,
		},
		{
			ID: testhelper.MkID(
				"prIDs - alphanumeric less by lexi order, a<b"),
			a:            v100AB,
			b:            v100B,
			shouldBeLess: true,
		},
		{
			ID: testhelper.MkID("prIDs - alphanumeric less by lexi order, a>b"),
			a:  v100B,
			b:  v100AB,
		},
		{
			ID: testhelper.MkID(
				"prIDs - numeric less by numeric order, a<b"),
			a:            v100B2,
			b:            v100B11,
			shouldBeLess: true,
		},
		{
			ID: testhelper.MkID("prIDs - numeric less by numeric order, a>b"),
			a:  v100B11,
			b:  v100B2,
		},
		{
			ID: testhelper.MkID(
				"prIDs - any prID less than none, a<b"),
			a:            v100RC1,
			b:            v100,
			shouldBeLess: true,
		},
		{
			ID: testhelper.MkID("prIDs - any prID less than none, a>b"),
			a:  v100,
			b:  v100RC1,
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

// setPreRelIDs sets the pre-release IDs and reports a fatal error if it
// cannot
func setPreRelIDs(t *testing.T, sv *semver.SV, prIDs []string) {
	t.Helper()

	if err := sv.SetPreRelIDs(prIDs); err != nil {
		t.Fatal("Error constructing the copy: ", err)
	}
}

// setBuildIDs sets the build IDs and reports a fatal error if it
// cannot
func setBuildIDs(t *testing.T, sv *semver.SV, buildIDs []string) {
	t.Helper()

	if err := sv.SetBuildIDs(buildIDs); err != nil {
		t.Fatal("Error constructing the copy: ", err)
	}
}

func TestEquals(t *testing.T) {
	baseSV := semver.NewSVOrPanic(1, 2, 3,
		[]string{"a", "b"}, []string{"a", "b"})

	var svCopies [10]semver.SV

	for i := range svCopies {
		baseSV.CopyInto(&svCopies[i])
	}

	(&svCopies[1]).IncrMajor()
	(&svCopies[2]).IncrMinor()
	(&svCopies[3]).IncrPatch()
	setPreRelIDs(t, &svCopies[4], []string{"a"})
	setPreRelIDs(t, &svCopies[5], []string{"a", "b", "c"})
	setPreRelIDs(t, &svCopies[6], []string{"b", "a"})
	setBuildIDs(t, &svCopies[7], []string{"a"})
	setBuildIDs(t, &svCopies[8], []string{"a", "b", "c"})
	setBuildIDs(t, &svCopies[9], []string{"b", "a"})

	testCases := []struct {
		testhelper.ID
		sv       *semver.SV
		expEqual bool
	}{
		{
			ID:       testhelper.MkID("should be equal"),
			sv:       &svCopies[0],
			expEqual: true,
		},
		{ID: testhelper.MkID("Major version differs"), sv: &svCopies[1]},
		{ID: testhelper.MkID("Minor version differs"), sv: &svCopies[2]},
		{ID: testhelper.MkID("Patch version differs"), sv: &svCopies[3]},
		{ID: testhelper.MkID("too few PreRelIDs"), sv: &svCopies[4]},
		{ID: testhelper.MkID("too many PreRelIDs"), sv: &svCopies[5]},
		{ID: testhelper.MkID("PreRelIDs in wrong order"), sv: &svCopies[6]},
		{ID: testhelper.MkID("too few BuildIDs"), sv: &svCopies[7]},
		{ID: testhelper.MkID("too many BuildIDs"), sv: &svCopies[8]},
		{ID: testhelper.MkID("BuildIDs in wrong order"), sv: &svCopies[9]},
	}

	for _, tc := range testCases {
		if semver.Equals(baseSV, tc.sv) {
			if tc.expEqual {
				continue
			}

			t.Log(tc.IDStr())
			t.Logf("\t: %s", baseSV)
			t.Logf("\t: %s", tc.sv)
			t.Errorf("\t: were not expected to be equal\n")
		} else {
			if !tc.expEqual {
				continue
			}

			t.Log(tc.IDStr())
			t.Logf("\t: %s", baseSV)
			t.Logf("\t: %s", tc.sv)
			t.Errorf("\t: were expected to be equal\n")
		}
	}
}
