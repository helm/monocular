package version

// Version is the current version of Monocular.
// Update this whenever making a new release.
// The version is of the format Major.Minor.Patch[-Prerelease][+BuildMetadata]
//
// Increment major number for new feature additions and behavioral changes.
// Increment minor number for bug fixes and performance enhancements.
// Increment patch number for critical fixes to existing releases.
var Version = "0.6.0"

// GetUserAgent returns the User Agent string for Monocular
func GetUserAgent() string {
	return "monocular/" + Version
}
