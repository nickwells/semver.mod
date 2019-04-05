package semver

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/nickwells/check.mod/check"
)

const (
	Name           = "semantic version ID"
	Names          = "semantic version IDs"
	GoodIDDesc     = "a non-empty string of letters, digits or hyphens"
	GoodVsnNumDesc = "greater than or equal to zero"
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
		return errors.New("the Build ID: '" + id + "' must be " + GoodIDDesc)
	}
	return nil
}

// CheckPreRelID returns nil if the id is a well-formed semver ID - suitable
// for a pre-release ID, an error otherwise
func CheckPreRelID(id string) error {
	if numericOnlyRE.MatchString(id) {
		if !goodNumericRE.MatchString(id) {
			return errors.New("the Pre-Rel ID: '" + id +
				"' must have no leading zero if it's all numeric")
		}
		return nil
	}
	if !idRE.MatchString(id) {
		return errors.New("the Pre-Rel ID: '" + id + "' must be " + GoodIDDesc)
	}
	return nil
}

// CheckAllPreRelIDs will return an error if any of the ids is not valid
// according to the rules for pre-release IDs
func CheckAllPreRelIDs(ids []string) error {
	for _, id := range ids {
		if err := CheckPreRelID(id); err != nil {
			return err
		}
	}
	return nil
}

// CheckAllBuildIDs will return an error if any of the ids is not valid
// according to the rules for build IDs
func CheckAllBuildIDs(ids []string) error {
	for _, id := range ids {
		if err := CheckBuildID(id); err != nil {
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
		return fmt.Errorf("bad major version: %d - it must be %s",
			sv.Major, GoodVsnNumDesc)
	}

	if sv.Minor < 0 {
		return fmt.Errorf("bad minor version: %d - it must be %s",
			sv.Minor, GoodVsnNumDesc)
	}

	if sv.Patch < 0 {
		return fmt.Errorf("bad patch version: %d - it must be %s",
			sv.Patch, GoodVsnNumDesc)
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

// CheckRules will confirm that all the checks are satisfied by the slice of
// IDs and return an error if not.
func CheckRules(ids []string, checks []check.StringSlice) error {
	for _, chk := range checks {
		err := chk(ids)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewSVWithIDRules returns a pointer to a properly constructed SV or an
// error if the IDs are not well-formed. As well as the standard checks it
// will also run the additional checks on the IDs and return an error if any
// of them fail. This allows you to enforce rules for the pre-release IDs and
// the build IDs. For instance you could prevent any semvers from having
// build IDs by passing a check the that slice of strings is empty.
func NewSVWithIDRules(major, minor, patch int,
	prIDs, buildIDs []string,
	prIDRules, bIDRules []check.StringSlice) (*SV, error) {
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
	if err := CheckRules(prIDs, prIDRules); err != nil {
		return nil, err
	}
	if err := CheckRules(buildIDs, bIDRules); err != nil {
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
			fmt.Errorf("bad %s - it does not start with a 'v'", Name)
	}
	var err error

	parts := strings.SplitN(s, "+", 2)
	if len(parts) == 2 {
		s = parts[0]
		sv.BuildIDs = strings.Split(parts[1], ".")
		if err = CheckAllBuildIDs(sv.BuildIDs); err != nil {
			return nil, fmt.Errorf("bad %s - %s", Name, err)
		}
	}

	parts = strings.SplitN(s, "-", 2)
	if len(parts) == 2 {
		s = parts[0]
		sv.PreRelIDs = strings.Split(parts[1], ".")
		if err = CheckAllPreRelIDs(sv.PreRelIDs); err != nil {
			return nil, fmt.Errorf("bad %s - %s", Name, err)
		}
	}

	parts = strings.SplitN(s, ".", 3)
	if len(parts) != 3 {
		return nil,
			fmt.Errorf("bad %s"+
				" - it cannot be split into major/minor/patch parts",
				Name)
	}

	sv.Major, err = strToVNum(parts[0], "major")
	if err != nil {
		return nil, fmt.Errorf("bad %s - %s", Name, err)
	}
	sv.Minor, err = strToVNum(parts[1], "minor")
	if err != nil {
		return nil, fmt.Errorf("bad %s - %s", Name, err)
	}
	sv.Patch, err = strToVNum(parts[2], "patch")
	if err != nil {
		return nil, fmt.Errorf("bad %s - %s", Name, err)
	}

	return &sv, nil
}

// strToVNum converts a string into a version number and reports any errors
// it finds
func strToVNum(s, name string) (int, error) {
	if len(s) > 1 && s[0] == '0' {
		return 0, fmt.Errorf("the %s version: '%s' has a leading 0", name, s)
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("the %s version: '%s' is not an integer", name, s)
	}
	if i < 0 {
		return 0, fmt.Errorf("the %s version: '%s' must be %s",
			name, s, GoodVsnNumDesc)
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

// IncrMajor calls the IncrMajor method on the passed SV
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

// IncrMinor calls the IncrMinor method on the passed SV
func IncrMinor(sv *SV) {
	sv.IncrMinor()
}

// IncrPatch increments the patch number. It also clears the pre-release IDs
// (if any) but not the build IDs
func (sv *SV) IncrPatch() {
	sv.Patch++
	sv.ClearPreRelIDs()
}

// IncrPatch calls the IncrPatch method on the passed SV
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
