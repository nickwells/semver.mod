package semver

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var idRE *regexp.Regexp
var numericOnlyRE *regexp.Regexp
var goodNumericRE *regexp.Regexp

func init() {
	idRE = regexp.MustCompile("^[-0-9A-Za-z]+$")
	numericOnlyRE = regexp.MustCompile("^[0-9]+$")
	goodNumericRE = regexp.MustCompile("^(0|[1-9][0-9]*)$")
}

// SV holds the parts of a semantic version number
type SV struct {
	Major     int
	Minor     int
	Patch     int
	PreRelIDs []string
	BuildIDs  []string
}

// CheckBuildID returns nil if the id is a well-formed semver ID - suitable
// for a build ID, an error otherwise
func CheckBuildID(id string) error {
	if !idRE.MatchString(id) {
		return errors.New(
			"Bad Build ID: '" + id +
				"' - must be a non-empty string of letters, digits or hyphens")
	}
	return nil
}

// CheckPreRelID returns nil if the id is a well-formed semver ID - suitable
// for a pre-release ID, an error otherwise
func CheckPreRelID(id string) error {
	if numericOnlyRE.MatchString(id) {
		if !goodNumericRE.MatchString(id) {
			return errors.New(
				"Bad Pre-Rel ID: '" + id +
					"' - if it's all numeric there must be no leading 0")
		}
		return nil
	}
	if !idRE.MatchString(id) {
		return errors.New(
			"Bad Pre-Rel ID: '" + id +
				"' - must be a non-empty string of letters, digits or hyphens")
	}
	return nil
}

// CheckAllPreRelIDs will return an error if any of the PreRelIDs is not valid
func CheckAllPreRelIDs(buildIDs []string) error {
	for _, buildID := range buildIDs {
		if err := CheckPreRelID(buildID); err != nil {
			return err
		}
	}
	return nil
}

// CheckAllBuildIDs will return an error if any of the BuildIDs is not valid
func CheckAllBuildIDs(buildIDs []string) error {
	for _, buildID := range buildIDs {
		if err := CheckBuildID(buildID); err != nil {
			return err
		}
	}
	return nil
}

// Check will return an error if any part of the SV is not valid
func (sv SV) Check() error {
	if err := CheckAllPreRelIDs(sv.PreRelIDs); err != nil {
		return err
	}
	if err := CheckAllBuildIDs(sv.BuildIDs); err != nil {
		return err
	}

	if sv.Major < 0 {
		return fmt.Errorf("bad major version: %d - it must be greater than 0",
			sv.Major)
	}

	if sv.Minor < 0 {
		return fmt.Errorf("bad minor version: %d - it must be greater than 0",
			sv.Minor)
	}

	if sv.Patch < 0 {
		return fmt.Errorf("bad patch version: %d - it must be greater than 0",
			sv.Patch)
	}

	return nil
}

// NewSV returns a pointer to a properly constructed SV or an error if the
// IDs are not well-formed
func NewSV(major, minor, patch int, prIDs, buildIDs []string) (*SV, error) {
	sv := &SV{
		Major:     major,
		Minor:     minor,
		Patch:     patch,
		PreRelIDs: prIDs,
		BuildIDs:  buildIDs,
	}

	if err := sv.Check(); err != nil {
		return nil, err
	}

	return sv, nil
}

// ParseSV will parse the semver string into an SV object. It will return a
// pointer to a properly constructed SV and a nil error if the semver is
// well-formed or a nil pointer and an error otherwise
func ParseSV(semver string) (*SV, error) {
	sv := SV{}

	s := strings.TrimPrefix(semver, "v")
	if s == semver {
		return nil,
			fmt.Errorf("Bad SemVer string: '%s' - it does not start with a 'v'",
				semver)
	}
	var err error

	parts := strings.SplitN(s, "+", 2)
	if len(parts) == 2 {
		s = parts[0]
		sv.BuildIDs = strings.Split(parts[1], ".")
		if err = CheckAllBuildIDs(sv.BuildIDs); err != nil {
			return nil,
				fmt.Errorf("Bad SemVer string: '%s' - %s", semver, err)
		}
	}

	parts = strings.SplitN(s, "-", 2)
	if len(parts) == 2 {
		s = parts[0]
		sv.PreRelIDs = strings.Split(parts[1], ".")
		if err = CheckAllPreRelIDs(sv.PreRelIDs); err != nil {
			return nil,
				fmt.Errorf("Bad SemVer string: '%s' - %s", semver, err)
		}
	}

	parts = strings.SplitN(s, ".", 3)
	if len(parts) != 3 {
		return nil,
			fmt.Errorf("Bad SemVer string: '%s' -"+
				" cannot split into major/minor/patch parts",
				semver)
	}

	sv.Major, err = strToVNum(parts[0], "major")
	if err != nil {
		return nil, fmt.Errorf("Bad SemVer string: '%s' - %s", semver, err)
	}
	sv.Minor, err = strToVNum(parts[1], "minor")
	if err != nil {
		return nil, fmt.Errorf("Bad SemVer string: '%s' - %s", semver, err)
	}
	sv.Patch, err = strToVNum(parts[2], "patch")
	if err != nil {
		return nil, fmt.Errorf("Bad SemVer string: '%s' - %s", semver, err)
	}

	return &sv, nil
}

// strToVNum converts a string into a version number and reports any errors
// it finds
func strToVNum(s, name string) (int, error) {
	if len(s) > 1 && s[0] == '0' {
		return 0, fmt.Errorf("bad %s version: %s - it has a leading 0", name, s)
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("bad %s version: %s - it is not a number", name, s)
	}
	if i < 0 {
		return 0, fmt.Errorf("bad %s version: %s - it must be >= 0", name, s)
	}
	return i, nil
}

// CopyInto copies from sv into target - it creates new slices and fills them
// with the pre-release and build IDs
func (sv SV) CopyInto(target *SV) {
	target.Major = sv.Major
	target.Minor = sv.Minor
	target.Patch = sv.Patch
	target.PreRelIDs = make([]string, len(sv.PreRelIDs))
	copy(target.PreRelIDs, sv.PreRelIDs)
	target.BuildIDs = make([]string, len(sv.BuildIDs))
	copy(target.BuildIDs, sv.BuildIDs)
}

// String returns a string representation of the semver
func (sv SV) String() string {
	b := ""
	pr := ""
	if len(sv.BuildIDs) > 0 {
		b = "+" + strings.Join(sv.BuildIDs, ".")
	}
	if len(sv.PreRelIDs) > 0 {
		pr = "-" + strings.Join(sv.PreRelIDs, ".")
	}
	return fmt.Sprintf("v%d.%d.%d%s%s", sv.Major, sv.Minor, sv.Patch, pr, b)
}

// IncrMajor increments the major version number and sets the minor and patch
// numbers to 0. It also clears the pre-release IDs (if any) but not the build
// IDs
func (sv *SV) IncrMajor() {
	sv.Major++
	sv.Minor = 0
	sv.Patch = 0
	sv.ClearPreRelIDs()
}
func IncrMajor(sv *SV) {
	sv.IncrMajor()
}

// IncrMinor increments the minor version number and sets the patch number to
// 0. It also clears the pre-release IDs (if any) but not the build IDs
func (sv *SV) IncrMinor() {
	sv.Minor++
	sv.Patch = 0
	sv.ClearPreRelIDs()
}
func IncrMinor(sv *SV) {
	sv.IncrMinor()
}

// IncrPatch increments the patch number. It also clears the pre-release IDs
// (if any) but not the build IDs
func (sv *SV) IncrPatch() {
	sv.Patch++
	sv.ClearPreRelIDs()
}
func IncrPatch(sv *SV) {
	sv.IncrPatch()
}

// ClearPreRelIDs clears the PreRelIDs
func (sv *SV) ClearPreRelIDs() {
	sv.PreRelIDs = []string{}
}

// ClearBuildIDs clears the BuildIDs
func (sv *SV) ClearBuildIDs() {
	sv.BuildIDs = []string{}
}
