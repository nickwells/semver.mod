[![GoDoc](https://godoc.org/github.com/nickwells/semver.mod?status.png)](https://godoc.org/github.com/nickwells/semver.mod)

# semver
funcs for parsing and manipulating semantic version numbers (semvers)

* You can parse a semver into a struct (`SV`) holding the major, minor and
  patch version numbers and the pre-release and build IDs.
* There are also funcs to correctly increment the various version numbers and
  safely perform other manipulations of the semver.
* There is a String function which will create a Semantic Version string from
  the `SV`.
* There is an `SVList` type which is a slice of pointers to `SV`s. This has
  the necessary methods for the slice to be sorted.
