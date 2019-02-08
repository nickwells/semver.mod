package semver_test

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/nickwells/semver.mod/semver"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestNewSV(t *testing.T) {
	testCases := []struct {
		name           string
		major          int
		minor          int
		patch          int
		prIDs          []string
		bIDs           []string
		errExpected    bool
		errMustContain []string
		expSVString    string
	}{
		{
			name:        "good - nil version",
			expSVString: "v0.0.0",
		},
		{
			name:        "good - v1.2.3",
			major:       1,
			minor:       2,
			patch:       3,
			expSVString: "v1.2.3",
		},
		{
			name:        "good - v1.2.3-xxx.XXX",
			major:       1,
			minor:       2,
			patch:       3,
			prIDs:       []string{"xxx", "XXX"},
			expSVString: "v1.2.3-xxx.XXX",
		},
		{
			name:        "good - v1.2.3+yyy.YYY",
			major:       1,
			minor:       2,
			patch:       3,
			bIDs:        []string{"yyy", "YYY"},
			expSVString: "v1.2.3+yyy.YYY",
		},
		{
			name:        "good - v1.2.3-xxx.XXX+yyy.YYY",
			major:       1,
			minor:       2,
			patch:       3,
			prIDs:       []string{"xxx", "XXX"},
			bIDs:        []string{"yyy", "YYY"},
			expSVString: "v1.2.3-xxx.XXX+yyy.YYY",
		},
		{
			name:        "bad - major version < 0",
			major:       -1,
			errExpected: true,
			errMustContain: []string{
				"bad major version: -1 - it must be greater than 0",
			},
		},
		{
			name:        "bad - minor version < 0",
			minor:       -1,
			errExpected: true,
			errMustContain: []string{
				"bad minor version: -1 - it must be greater than 0",
			},
		},
		{
			name:        "bad - patch version < 0",
			patch:       -1,
			errExpected: true,
			errMustContain: []string{
				"bad patch version: -1 - it must be greater than 0",
			},
		},
		{
			name:        "bad - invalid Pre-Rel ID - non alphanumeric or '-'",
			prIDs:       []string{"aaa", "a$a", "bbb"},
			errExpected: true,
			errMustContain: []string{
				"Bad Pre-Rel ID: 'a$a' - ",
				"must be a non-empty string of letters, digits or hyphens",
			},
		},
		{
			name:        "bad - invalid Pre-Rel ID - numeric with leading 0",
			prIDs:       []string{"0", "012", "bbb"},
			errExpected: true,
			errMustContain: []string{
				"Bad Pre-Rel ID: '012' - ",
				"if it's all numeric there must be no leading 0",
			},
		},
		{
			name:        "bad - invalid build ID",
			bIDs:        []string{"aaa", "a$a", "bbb"},
			errExpected: true,
			errMustContain: []string{
				"Bad Build ID: 'a$a' - ",
				"must be a non-empty string of letters, digits or hyphens",
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)
		sv, err := semver.NewSV(tc.major, tc.minor, tc.patch, tc.prIDs, tc.bIDs)
		testhelper.CheckError(t, tcID, err, tc.errExpected, tc.errMustContain)
		if err == nil && !tc.errExpected {
			if sv.String() != tc.expSVString {
				t.Log(tcID)
				t.Logf("\t: expected: %s", tc.expSVString)
				t.Logf("\t:      got: %s", sv.String())
				t.Errorf("\t: bad string representation\n")
			}
		}
	}
}

func TestParse(t *testing.T) {
	testCases := []struct {
		name           string
		svStr          string
		errExpected    bool
		errMustContain []string
		svExpected     semver.SV
	}{
		{
			name:  "good - with build and pre-rel ID",
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
			name:  "good - no build ID",
			svStr: "v1.2.3-xxx",
			svExpected: semver.SV{
				Major:     1,
				Minor:     2,
				Patch:     3,
				PreRelIDs: []string{"xxx"},
			},
		},
		{
			name:  "good - no pre-rel ID",
			svStr: "v1.2.3+yyy",
			svExpected: semver.SV{
				Major:    1,
				Minor:    2,
				Patch:    3,
				BuildIDs: []string{"yyy"},
			},
		},
		{
			name:  "good - no build or pre-rel ID",
			svStr: "v1.2.3",
			svExpected: semver.SV{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
		},
		{
			name:        "bad - no leading 'v'",
			svStr:       "1.2.3-xxx+yyy",
			errExpected: true,
			errMustContain: []string{
				"Bad SemVer string: '",
				"' - it does not start with a 'v'",
			},
		},
		{
			name:        "bad - invalid Pre-Rel ID - bad char",
			svStr:       "v1.2.3-x$xx+yyy",
			errExpected: true,
			errMustContain: []string{
				"Bad SemVer string: '",
				"' - Bad Pre-Rel ID: ",
				"must be a non-empty string of letters, digits or hyphens",
			},
		},
		{
			name:        "bad - invalid Pre-Rel ID - number with leading 0",
			svStr:       "v1.2.3-012+yyy",
			errExpected: true,
			errMustContain: []string{
				"Bad SemVer string: '",
				"' - Bad Pre-Rel ID: ",
				"if it's all numeric there must be no leading 0",
			},
		},
		{
			name:        "bad - empty Pre-Rel ID",
			svStr:       "v1.2.3-a..b+yyy",
			errExpected: true,
			errMustContain: []string{
				"Bad SemVer string: '",
				"' - Bad Pre-Rel ID: ",
				"must be a non-empty string of letters, digits or hyphens",
			},
		},
		{
			name:        "bad - invalid build ID - bad char",
			svStr:       "v1.2.3-xxx+y$yy",
			errExpected: true,
			errMustContain: []string{
				"Bad SemVer string: '",
				"' - Bad Build ID: ",
				"must be a non-empty string of letters, digits or hyphens",
			},
		},
		{
			name:        "bad - empty build ID",
			svStr:       "v1.2.3-xxx+",
			errExpected: true,
			errMustContain: []string{
				"Bad SemVer string: '",
				"' - Bad Build ID: ",
				"must be a non-empty string of letters, digits or hyphens",
			},
		},
		{
			name:        "bad - too few version parts",
			svStr:       "v1.2",
			errExpected: true,
			errMustContain: []string{
				"Bad SemVer string: '",
				"' - it cannot be split into major/minor/patch parts",
			},
		},
		{
			name:        "bad - major part is not a number",
			svStr:       "vX.2.3",
			errExpected: true,
			errMustContain: []string{
				"Bad SemVer string: '",
				"' - bad major version: X - it is not a number",
			},
		},
		{
			name:        "bad - major part has a leading zero",
			svStr:       "v01.2.3",
			errExpected: true,
			errMustContain: []string{
				"Bad SemVer string: '",
				"' - bad major version: 01 - it has a leading 0",
			},
		},
		{
			name:        "bad - minor part is not a number",
			svStr:       "v1.X.3",
			errExpected: true,
			errMustContain: []string{
				"Bad SemVer string: '",
				"' - bad minor version: X - it is not a number",
			},
		},
		{
			name:        "bad - minor part has a leading zero",
			svStr:       "v1.02.3",
			errExpected: true,
			errMustContain: []string{
				"Bad SemVer string: '",
				"' - bad minor version: 02 - it has a leading 0",
			},
		},
		{
			name:        "bad - patch part is not a number",
			svStr:       "v1.2.X",
			errExpected: true,
			errMustContain: []string{
				"Bad SemVer string: '",
				"' - bad patch version: X - it is not a number",
			},
		},
		{
			name:        "bad - patch part has a leading zero",
			svStr:       "v1.2.03",
			errExpected: true,
			errMustContain: []string{
				"Bad SemVer string: '",
				"' - bad patch version: 03 - it has a leading 0",
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)
		sv, err := semver.ParseSV(tc.svStr)
		testhelper.CheckError(t, tcID, err, tc.errExpected, tc.errMustContain)
		if err == nil && !tc.errExpected {
			if !semver.Equals(sv, &tc.svExpected) {
				t.Log(tcID)
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
		name       string
		incrFunc   func(*semver.SV)
		svExpected semver.SV
	}{
		{
			name:     "IncrMajor",
			incrFunc: semver.IncrMajor,
			svExpected: semver.SV{
				Major:    major + 1,
				BuildIDs: []string{bID},
			},
		},
		{
			name:     "IncrMinor",
			incrFunc: semver.IncrMinor,
			svExpected: semver.SV{
				Major:    major,
				Minor:    minor + 1,
				BuildIDs: []string{bID},
			},
		},
		{
			name:     "IncrPatch",
			incrFunc: semver.IncrPatch,
			svExpected: semver.SV{
				Major:    major,
				Minor:    minor,
				Patch:    patch + 1,
				BuildIDs: []string{bID},
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)

		localSV := new(semver.SV)
		sv.CopyInto(localSV)

		tc.incrFunc(localSV)
		if !semver.Equals(localSV, &tc.svExpected) {
			t.Log(tcID)
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
