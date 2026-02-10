package semver

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/nickwells/check.mod/v2/check"
)

// These constants provide consistent text values to describe semantic
// version IDs and parts of them.
const (
	ShortName      = "semver-ID"
	ShortNames     = "semver-IDs"
	Name           = "semantic version ID"
	Names          = "semantic version IDs"
	GoodIDDesc     = "a non-empty string of letters, digits or hyphens"
	GoodVsnNumDesc = "greater than or equal to zero"
)

const (
	semverMajorVsnIdx = iota
	semverMinorVsnIdx
	semverPatchVsnIdx
	semverVsnPartCount
)

const (
	semverPrefix             = "v"
	semverBuildIDsSeparator  = "+"
	semverPreRelIDsSeparator = "-"
	semverPartSeparator      = "."
)

var (
	idRE          *regexp.Regexp
	numericOnlyRE *regexp.Regexp
	goodNumericRE *regexp.Regexp
)

func init() {
	idRE = regexp.MustCompile("^[-0-9A-Za-z]+$")
	numericOnlyRE = regexp.MustCompile("^[0-9]+$")
	goodNumericRE = regexp.MustCompile("^(0|[1-9][0-9]*)$")
}

// SV holds the parts of a semantic version number
type SV struct {
	major     int
	minor     int
	patch     int
	preRelIDs []string
	buildIDs  []string

	hasBeenSet bool
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
	if err := CheckAllPreRelIDs(sv.preRelIDs); err != nil {
		return err
	}

	if err := CheckAllBuildIDs(sv.buildIDs); err != nil {
		return err
	}

	if sv.major < 0 {
		return fmt.Errorf("bad major version: %d - it must be %s",
			sv.major, GoodVsnNumDesc)
	}

	if sv.minor < 0 {
		return fmt.Errorf("bad minor version: %d - it must be %s",
			sv.minor, GoodVsnNumDesc)
	}

	if sv.patch < 0 {
		return fmt.Errorf("bad patch version: %d - it must be %s",
			sv.patch, GoodVsnNumDesc)
	}

	return nil
}

// NewSV returns a pointer to a properly constructed SV or an error if the
// IDs are not well-formed
func NewSV(major, minor, patch int, prIDs, buildIDs []string) (*SV, error) {
	sv := &SV{
		major:      major,
		minor:      minor,
		patch:      patch,
		preRelIDs:  prIDs,
		buildIDs:   buildIDs,
		hasBeenSet: true,
	}

	if err := sv.Check(); err != nil {
		return nil, err
	}

	return sv, nil
}

// NewSVOrPanic returns a pointer to a properly constructed SV. If there were
// any errors it will panic.
func NewSVOrPanic(major, minor, patch int, prIDs, buildIDs []string) *SV {
	sv, err := NewSV(major, minor, patch, prIDs, buildIDs)
	if err != nil {
		panic(err)
	}

	return sv
}

// CheckRules will confirm that all the checks are satisfied by the slice of
// IDs and return an error if not.
func CheckRules(ids []string, checks []check.ValCk[[]string]) error {
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
// build IDs by passing a check that the slice of strings is empty.
func NewSVWithIDRules(major, minor, patch int,
	prIDs, buildIDs []string,
	prIDRules, bIDRules []check.ValCk[[]string],
) (*SV, error) {
	sv := &SV{
		major:      major,
		minor:      minor,
		patch:      patch,
		preRelIDs:  prIDs,
		buildIDs:   buildIDs,
		hasBeenSet: true,
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

// NewSVWithIDRulesOrPanic returns a pointer to a properly constructed SV. If
// there were any errors it will panic.
func NewSVWithIDRulesOrPanic(major, minor, patch int,
	prIDs, buildIDs []string,
	prIDRules, bIDRules []check.ValCk[[]string],
) *SV {
	sv, err := NewSVWithIDRules(major, minor, patch,
		prIDs, buildIDs,
		prIDRules, bIDRules)
	if err != nil {
		panic(err)
	}

	return sv
}

// ParseSV will parse the semver string into an SV object. It will strip off
// the leading 'v' which must be present. It will return a pointer to a
// properly constructed SV and a nil error if the semver is well-formed or a
// nil pointer and an error otherwise
func ParseSV(semver string) (*SV, error) {
	s, ok := strings.CutPrefix(semver, semverPrefix)
	if !ok {
		return nil,
			fmt.Errorf("bad %s - it does not start with a 'v'", Name)
	}

	return ParseStrictSV(s)
}

// ParseStrictSV will parse the semver string into an SV object. It is
// expected to have no prefix. It will return a pointer to a properly
// constructed SV and a nil error if the semver is well-formed or a nil
// pointer and an error otherwise
func ParseStrictSV(semver string) (*SV, error) {
	sv := SV{}

	var err error

	var buildIDs, preRelIDs string

	var ok bool

	semver, buildIDs, ok = strings.Cut(semver, semverBuildIDsSeparator)
	if ok {
		sv.buildIDs = strings.Split(buildIDs, semverPartSeparator)
		if err = CheckAllBuildIDs(sv.buildIDs); err != nil {
			return nil, fmt.Errorf("bad %s - %s", Name, err)
		}
	}

	semver, preRelIDs, ok = strings.Cut(semver, semverPreRelIDsSeparator)
	if ok {
		sv.preRelIDs = strings.Split(preRelIDs, semverPartSeparator)
		if err = CheckAllPreRelIDs(sv.preRelIDs); err != nil {
			return nil, fmt.Errorf("bad %s - %s", Name, err)
		}
	}

	parts := strings.SplitN(semver, semverPartSeparator, semverVsnPartCount)
	if len(parts) != semverVsnPartCount {
		return nil,
			fmt.Errorf("bad %s"+
				" - it cannot be split into major/minor/patch parts",
				Name)
	}

	sv.major, err = strToVNum(parts[semverMajorVsnIdx], "major")
	if err != nil {
		return nil, fmt.Errorf("bad %s - %s", Name, err)
	}

	sv.minor, err = strToVNum(parts[semverMinorVsnIdx], "minor")
	if err != nil {
		return nil, fmt.Errorf("bad %s - %s", Name, err)
	}

	sv.patch, err = strToVNum(parts[semverPatchVsnIdx], "patch")
	if err != nil {
		return nil, fmt.Errorf("bad %s - %s", Name, err)
	}

	sv.hasBeenSet = true

	return &sv, nil
}

// strToVNum converts a string into a version number and reports any errors
// it finds
func strToVNum(s, name string) (int, error) {
	if len(s) > 1 && s[0] == '0' {
		return 0, fmt.Errorf("the %s version: %q has a leading 0", name, s)
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("the %s version: %q is not an integer", name, s)
	}

	if i < 0 {
		return 0, fmt.Errorf("the %s version: %q must be %s",
			name, s, GoodVsnNumDesc)
	}

	return i, nil
}

// CopyInto copies from sv into target - it creates new slices and fills them
// with the pre-release and build IDs
func (sv SV) CopyInto(target *SV) {
	target.major = sv.major
	target.minor = sv.minor
	target.patch = sv.patch

	target.preRelIDs = make([]string, len(sv.preRelIDs))
	copy(target.preRelIDs, sv.preRelIDs)

	target.buildIDs = make([]string, len(sv.buildIDs))
	copy(target.buildIDs, sv.buildIDs)

	target.hasBeenSet = sv.hasBeenSet
}

// Major returns the major version number part of the SemVer
func (sv SV) Major() int { return sv.major }

// Minor returns the minor version number part of the SemVer
func (sv SV) Minor() int { return sv.minor }

// Patch returns the patch version number part of the SemVer
func (sv SV) Patch() int { return sv.patch }

// PreRelIDs returns the preRelIDs version number part of the SemVer
func (sv SV) PreRelIDs() []string { return sv.preRelIDs }

// HasPreRelIDs returns true if the preRelIDs version number part of the
// SemVer is non-empty
func (sv SV) HasPreRelIDs() bool { return len(sv.preRelIDs) > 0 }

// BuildIDs returns the buildIDs version number part of the SemVer
func (sv SV) BuildIDs() []string { return sv.buildIDs }

// HasBuildIDs returns true if the buildIDs version number part of the
// SemVer is non-empty
func (sv SV) HasBuildIDs() bool { return len(sv.buildIDs) > 0 }

// HasBeenSet returns the value of the internal flag which is set if the
// SemVer has been set
func (sv SV) HasBeenSet() bool { return sv.hasBeenSet }

// String returns a string representation of the semver
func (sv SV) String() string {
	if !sv.hasBeenSet {
		return ""
	}

	buildIDs := ""
	prIDs := ""

	if len(sv.buildIDs) > 0 {
		buildIDs = semverBuildIDsSeparator +
			strings.Join(sv.buildIDs, semverPartSeparator)
	}

	if len(sv.preRelIDs) > 0 {
		prIDs = semverPreRelIDsSeparator +
			strings.Join(sv.preRelIDs, semverPartSeparator)
	}

	return semverPrefix +
		fmt.Sprintf("%d", sv.major) + semverPartSeparator +
		fmt.Sprintf("%d", sv.minor) + semverPartSeparator +
		fmt.Sprintf("%d", sv.patch) +
		prIDs + buildIDs
}

// IncrMajor increments the major version number and sets the minor and patch
// numbers to 0. It also clears the pre-release IDs (if any) but not the build
// IDs
func (sv *SV) IncrMajor() {
	sv.major++
	sv.minor = 0
	sv.patch = 0
	sv.ClearPreRelIDs()
}

// IncrMinor increments the minor version number and sets the patch number to
// 0. It also clears the pre-release IDs (if any) but not the build IDs
func (sv *SV) IncrMinor() {
	sv.minor++
	sv.patch = 0
	sv.ClearPreRelIDs()
}

// IncrPatch increments the patch number. It also clears the pre-release IDs
// (if any) but not the build IDs
func (sv *SV) IncrPatch() {
	sv.patch++
	sv.ClearPreRelIDs()
}

// ClearPreRelIDs clears the PreRelIDs
func (sv *SV) ClearPreRelIDs() {
	sv.preRelIDs = []string{}
}

// SetPreRelIDs sets the PreRelIDs
func (sv *SV) SetPreRelIDs(ids []string) error {
	err := CheckAllPreRelIDs(ids)
	if err != nil {
		return err
	}

	sv.preRelIDs = ids

	return nil
}

// ClearBuildIDs clears the BuildIDs
func (sv *SV) ClearBuildIDs() {
	sv.buildIDs = []string{}
}

// SetBuildIDs sets the BuildIDs
func (sv *SV) SetBuildIDs(ids []string) error {
	err := CheckAllBuildIDs(ids)
	if err != nil {
		return err
	}

	sv.buildIDs = ids

	return nil
}
