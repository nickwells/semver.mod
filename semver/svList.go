package semver

type SVList []*SV

// Less reports whether the element with index i should sort before the
// element with index j
func (svl SVList) Less(i, j int) bool {
	return Less(svl[i], svl[j])
}

// Len reports the number of elements in the collection
func (svl SVList) Len() int {
	return len(svl)
}

// Swap swaps the elements with indexes i and j
func (svl SVList) Swap(i, j int) {
	svl[i], svl[j] = svl[j], svl[i]
}
