/*
Package semver offers methods for parsing and manipulating semantic
version numbers (SVs).

An SV can be constructed programatically and then printed out or it can be
parsed from a string form. Additionally there are functions for incrementing
the major, minor or patch versions. The functions for creating a new SV
object will validate the parts - the major, minor and patch numbers and the
pre-release and build IDs and will return an error if any part is invalid.

There is an SVList type (a slice of pointers to SVs) which has member
functions (Len, Less and Swap) which make it able to be sorted.

*/
package semver
