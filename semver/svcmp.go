package semver

import "strconv"

// lessPRIDs compares the preRelIDs of the two semver values
//
//nolint:cyclop
func lessPRIDs(a, b *SV) bool {
	if len(a.preRelIDs) > 0 && len(b.preRelIDs) == 0 {
		return true
	}

	if len(a.preRelIDs) == 0 && len(b.preRelIDs) > 0 {
		return false
	}

	for i, aID := range a.preRelIDs {
		if i >= len(b.preRelIDs) {
			break
		}

		bID := b.preRelIDs[i]

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

	return len(a.preRelIDs) < len(b.preRelIDs)
}

// Less returns true if a is less than b according to the ordering rules for
// semantic versions given in the Semantic Versioning Specification v2.0.0
// (spec item 11)
func Less(a, b *SV) bool {
	if a.major < b.major {
		return true
	}

	if a.major > b.major {
		return false
	}

	if a.minor < b.minor {
		return true
	}

	if a.minor > b.minor {
		return false
	}

	if a.patch < b.patch {
		return true
	}

	if a.patch > b.patch {
		return false
	}

	return lessPRIDs(a, b)
}

// Equals compares the two SemVers and returns true if they are identical,
// false otherwise
func Equals(a, b *SV) bool {
	if a.major != b.major {
		return false
	}

	if a.minor != b.minor {
		return false
	}

	if a.patch != b.patch {
		return false
	}

	if len(a.preRelIDs) != len(b.preRelIDs) {
		return false
	}

	for i, id := range a.preRelIDs {
		if id != b.preRelIDs[i] {
			return false
		}
	}

	if len(a.buildIDs) != len(b.buildIDs) {
		return false
	}

	for i, id := range a.buildIDs {
		if id != b.buildIDs[i] {
			return false
		}
	}

	return true
}
