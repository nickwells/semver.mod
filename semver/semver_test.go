package semver_test

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/nickwells/semver.mod/semver"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestNewSV(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		major       int
		minor       int
		patch       int
		prIDs       []string
		bIDs        []string
		expSVString string
	}{
		{
			ID:          testhelper.MkID("good - nil version"),
			expSVString: "v0.0.0",
		},
		{
			ID:          testhelper.MkID("good - v1.2.3"),
			major:       1,
			minor:       2,
			patch:       3,
			expSVString: "v1.2.3",
		},
		{
			ID:          testhelper.MkID("good - v1.2.3-xxx.XXX"),
			major:       1,
			minor:       2,
			patch:       3,
			prIDs:       []string{"xxx", "XXX"},
			expSVString: "v1.2.3-xxx.XXX",
		},
		{
			ID:          testhelper.MkID("good - v1.2.3+yyy.YYY"),
			major:       1,
			minor:       2,
			patch:       3,
			bIDs:        []string{"yyy", "YYY"},
			expSVString: "v1.2.3+yyy.YYY",
		},
		{
			ID:          testhelper.MkID("good - v1.2.3-xxx.X-XX+yyy.YYY"),
			major:       1,
			minor:       2,
			patch:       3,
			prIDs:       []string{"xxx", "X-XX"},
			bIDs:        []string{"yyy", "YYY"},
			expSVString: "v1.2.3-xxx.X-XX+yyy.YYY",
		},
		{
			ID:    testhelper.MkID("bad - major version < 0"),
			major: -1,
			ExpErr: testhelper.MkExpErr(
				"bad major version: -1 - it must be " + semver.GoodVsnNumDesc),
		},
		{
			ID:    testhelper.MkID("bad - minor version < 0"),
			minor: -1,
			ExpErr: testhelper.MkExpErr(
				"bad minor version: -1 - it must be " + semver.GoodVsnNumDesc),
		},
		{
			ID:    testhelper.MkID("bad - patch version < 0"),
			patch: -1,
			ExpErr: testhelper.MkExpErr(
				"bad patch version: -1 - it must be " + semver.GoodVsnNumDesc),
		},
		{
			ID: testhelper.MkID(
				"bad - invalid Pre-Rel ID - non alphanumeric or '-'"),
			prIDs: []string{"aaa", "a$a", "bbb"},
			ExpErr: testhelper.MkExpErr(
				"the Pre-Rel ID: 'a$a' must be " + semver.GoodIDDesc),
		},
		{
			ID: testhelper.MkID(
				"bad - invalid Pre-Rel ID - numeric with leading 0"),
			prIDs: []string{"0", "012", "bbb"},
			ExpErr: testhelper.MkExpErr(
				"the Pre-Rel ID: '012' " +
					"must have no leading zero if it's all numeric"),
		},
		{
			ID: testhelper.MkID(
				"bad - invalid build ID - non alphanumeric or '-'"),
			bIDs: []string{"aaa", "a$a", "bbb"},
			ExpErr: testhelper.MkExpErr(
				"the Build ID: 'a$a' must be " + semver.GoodIDDesc),
		},
	}

	for _, tc := range testCases {
		sv, err := semver.NewSV(tc.major, tc.minor, tc.patch, tc.prIDs, tc.bIDs)
		if testhelper.CheckExpErr(t, err, tc) && err == nil {
			testhelper.CmpValString(t, tc.IDStr(), "semver string",
				sv.String(), tc.expSVString)
		}
	}
}

func TestParse(t *testing.T) {
	const badSemVer = "bad " + semver.Name
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		svStr      string
		svExpected semver.SV
	}{
		{
			ID:    testhelper.MkID("good - with build and pre-rel ID"),
			svStr: "v1.2.3-aaa.bbb.c-d.0.123.123a.0123a+xxx.yyy.z-z.01.0-1",
			svExpected: semver.SV{
				Major: 1,
				Minor: 2,
				Patch: 3,
				PreRelIDs: []string{"aaa", "bbb", "c-d",
					"0", "123", "123a", "0123a"},
				BuildIDs: []string{"xxx", "yyy", "z-z", "01", "0-1"},
			},
		},
		{
			ID:    testhelper.MkID("good - no build ID"),
			svStr: "v1.2.3-xxx",
			svExpected: semver.SV{
				Major:     1,
				Minor:     2,
				Patch:     3,
				PreRelIDs: []string{"xxx"},
			},
		},
		{
			ID:    testhelper.MkID("good - no pre-rel ID"),
			svStr: "v1.2.3+yyy",
			svExpected: semver.SV{
				Major:    1,
				Minor:    2,
				Patch:    3,
				BuildIDs: []string{"yyy"},
			},
		},
		{
			ID:    testhelper.MkID("good - no build or pre-rel ID"),
			svStr: "v1.2.3",
			svExpected: semver.SV{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
		},
		{
			ID:    testhelper.MkID("bad - no leading 'v'"),
			svStr: "1.2.3-xxx+yyy",
			ExpErr: testhelper.MkExpErr(
				badSemVer + " - it does not start with a 'v'"),
		},
		{
			ID:    testhelper.MkID("bad - invalid Pre-Rel ID - bad char"),
			svStr: "v1.2.3-x$xx+yyy",
			ExpErr: testhelper.MkExpErr(
				badSemVer + " - the Pre-Rel ID: 'x$xx' must be " +
					semver.GoodIDDesc),
		},
		{
			ID: testhelper.MkID(
				"bad - invalid Pre-Rel ID - number with leading 0"),
			svStr: "v1.2.3-012+yyy",
			ExpErr: testhelper.MkExpErr(
				badSemVer+" - the Pre-Rel ID: '012'",
				"must have no leading zero if it's all numeric"),
		},
		{
			ID:    testhelper.MkID("bad - empty Pre-Rel ID"),
			svStr: "v1.2.3-a..b+yyy",
			ExpErr: testhelper.MkExpErr(
				badSemVer + " - the Pre-Rel ID: '' must be " +
					semver.GoodIDDesc),
		},
		{
			ID:    testhelper.MkID("bad - invalid build ID - bad char"),
			svStr: "v1.2.3-xxx+y$yy",
			ExpErr: testhelper.MkExpErr(
				badSemVer + " - the Build ID: 'y$yy' must be " +
					semver.GoodIDDesc),
		},
		{
			ID:    testhelper.MkID("bad - empty build ID"),
			svStr: "v1.2.3-xxx+",
			ExpErr: testhelper.MkExpErr(
				badSemVer + " - the Build ID: '' must be " + semver.GoodIDDesc),
		},
		{
			ID:    testhelper.MkID("bad - too few version parts"),
			svStr: "v1.2",
			ExpErr: testhelper.MkExpErr(badSemVer +
				" - it cannot be split into major/minor/patch parts"),
		},
		{
			ID:    testhelper.MkID("bad - major part is not an integer"),
			svStr: "vX.2.3",
			ExpErr: testhelper.MkExpErr(
				badSemVer + " - the major version: 'X' is not an integer"),
		},
		{
			ID:    testhelper.MkID("bad - major part has a leading zero"),
			svStr: "v01.2.3",
			ExpErr: testhelper.MkExpErr(
				badSemVer + " - the major version: '01' has a leading 0"),
		},
		{
			ID:    testhelper.MkID("bad - minor part is not an integer"),
			svStr: "v1.X.3",
			ExpErr: testhelper.MkExpErr(
				badSemVer + " - the minor version: 'X' is not an integer"),
		},
		{
			ID:    testhelper.MkID("bad - minor part has a leading zero"),
			svStr: "v1.02.3",
			ExpErr: testhelper.MkExpErr(
				badSemVer + " - the minor version: '02' has a leading 0"),
		},
		{
			ID:    testhelper.MkID("bad - patch part is not an integer"),
			svStr: "v1.2.X",
			ExpErr: testhelper.MkExpErr(
				badSemVer + " - the patch version: 'X' is not an integer"),
		},
		{
			ID:    testhelper.MkID("bad - patch part has a leading zero"),
			svStr: "v1.2.03",
			ExpErr: testhelper.MkExpErr(
				badSemVer + " - the patch version: '03' has a leading 0"),
		},
	}

	for _, tc := range testCases {
		sv, err := semver.ParseSV(tc.svStr)
		if testhelper.CheckExpErr(t, err, tc) && err == nil {
			if !semver.Equals(sv, &tc.svExpected) {
				t.Log(tc.IDStr())
				t.Logf("\t: expected: %s", tc.svExpected)
				t.Logf("\t:      got: %s", sv)
				t.Errorf("\t: bad parsing\n")
			}
		}
	}
}

func TestIncr(t *testing.T) {
	major := 1
	minor := 2
	patch := 3
	prID := "Pre-Rel-ID"
	bID := "Build-ID"
	sv, err := semver.NewSV(major, minor, patch, []string{prID}, []string{bID})
	if err != nil {
		t.Fatal("Couldn't create the new semver: ", err)
	}

	testCases := []struct {
		testhelper.ID
		incrFunc   func(*semver.SV)
		svExpected semver.SV
	}{
		{
			ID:       testhelper.MkID("IncrMajor"),
			incrFunc: semver.IncrMajor,
			svExpected: semver.SV{
				Major:    major + 1,
				BuildIDs: []string{bID},
			},
		},
		{
			ID:       testhelper.MkID("IncrMinor"),
			incrFunc: semver.IncrMinor,
			svExpected: semver.SV{
				Major:    major,
				Minor:    minor + 1,
				BuildIDs: []string{bID},
			},
		},
		{
			ID:       testhelper.MkID("IncrPatch"),
			incrFunc: semver.IncrPatch,
			svExpected: semver.SV{
				Major:    major,
				Minor:    minor,
				Patch:    patch + 1,
				BuildIDs: []string{bID},
			},
		},
	}

	for _, tc := range testCases {
		localSV := new(semver.SV)
		sv.CopyInto(localSV)

		tc.incrFunc(localSV)
		if !semver.Equals(localSV, &tc.svExpected) {
			t.Log(tc.IDStr())
			t.Logf("\t: expected: %s", tc.svExpected)
			t.Logf("\t:      got: %s", localSV)
			t.Errorf("\t: bad increment\n")
		}
	}
}

func TestAllBadStrings(t *testing.T) {
	const fname = "testdata/badSemVers"
	file, err := os.Open(fname)
	if err != nil {
		t.Fatal("Cannot open the test file: ", fname, " - ", err)
	}

	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), "\t", 2)
		if len(parts) != 2 {
			continue
		}
		svStr, expectedErr := parts[0], strings.TrimSpace(parts[1])

		_, err := semver.ParseSV(svStr)
		if err == nil {
			t.Logf("parsing: %s:%d : %s", fname, lineNum, svStr)
			t.Errorf("\t: no error was reported, expected: %s", expectedErr)
		}
	}
}

func TestAllGoodStrings(t *testing.T) {
	const fname = "testdata/goodSemVers"
	file, err := os.Open(fname)
	if err != nil {
		t.Fatal("Cannot open the test file: ", fname, " - ", err)
	}

	scanner := bufio.NewScanner(file)
	lineNum := 0
	var prevSV *semver.SV
	for scanner.Scan() {
		svStr := scanner.Text()
		if svStr == "" {
			continue
		}

		sv, err := semver.ParseSV(svStr)
		if err != nil {
			t.Logf("parsing: %s:%d : %s", fname, lineNum, svStr)
			t.Errorf("\t: unexpected error: %s", err)
			continue
		}
		if prevSV != nil {
			if semver.Less(sv, prevSV) {
				t.Logf("checking order: %s:%d", fname, lineNum)
				t.Logf("\t:     this: %s", svStr)
				t.Logf("\t: previous: %s", prevSV)
				t.Errorf("\t: this should not be less than the previous value")
				continue
			}
		}
		prevSV = sv
	}
}
