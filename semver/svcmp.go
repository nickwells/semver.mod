package semver

import "strconv"

// lessPRIDs compares the PreRelIDs of the two semver values
func lessPRIDs(a, b *SV) bool {
	if len(a.PreRelIDs) > 0 && len(b.PreRelIDs) == 0 {
		return true
	}
	if len(a.PreRelIDs) == 0 && len(b.PreRelIDs) > 0 {
		return false
	}

	for i, aID := range a.PreRelIDs {
		if i >= len(b.PreRelIDs) {
			break
		}
		bID := b.PreRelIDs[i]
		if goodNumericRE.MatchString(aID) {
			if goodNumericRE.MatchString(bID) {
				aAsNum, _ := strconv.Atoi(aID)
				bAsNum, _ := strconv.Atoi(bID)
				if aAsNum < bAsNum {
					return true
				}
				if aAsNum > bAsNum {
					return false
				}
			} else {
				return true
			}
		} else if goodNumericRE.MatchString(bID) {
			return false
		} else if aID < bID {
			return true
		} else if aID > bID {
			return false
		}
	}

	return len(a.PreRelIDs) < len(b.PreRelIDs)
}

// Less returns true if a is less than b according to the ordering rules for
// semantic versions given in the Semantic Versioning Specification v2.0.0
// (spec item 11)
func Less(a, b *SV) bool {
	if a.Major < b.Major {
		return true
	}
	if a.Major > b.Major {
		return false
	}

	if a.Minor < b.Minor {
		return true
	}
	if a.Minor > b.Minor {
		return false
	}

	if a.Patch < b.Patch {
		return true
	}
	if a.Patch > b.Patch {
		return false
	}

	return lessPRIDs(a, b)
}

// Equals compares the two SemVers and returns true if they are identical,
// false otherwise
func Equals(a, b *SV) bool {
	if a.Major != b.Major {
		return false
	}
	if a.Minor != b.Minor {
		return false
	}
	if a.Patch != b.Patch {
		return false
	}
	if len(a.PreRelIDs) != len(b.PreRelIDs) {
		return false
	}
	if len(a.BuildIDs) != len(b.BuildIDs) {
		return false
	}
	for i, id := range a.PreRelIDs {
		if id != b.PreRelIDs[i] {
			return false
		}
	}
	for i, id := range a.BuildIDs {
		if id != b.BuildIDs[i] {
			return false
		}
	}

	return true
}
