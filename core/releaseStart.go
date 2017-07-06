package core

func releaseStart(catalogName, releaseName string) {
	// Find hitch.db root.

	// Check for staging file.  Reject command if staging file already exists.

	// Check for catalog already existing.  Reject if not (this is fat-finger avoidance).
	// REVIEW : is this helpful?  skipping for now.

	// Check for catalog+release already existing.  Reject if released before.

	// If catalog has signing keys set up, check that we have those keys.

	// All checks passed.
	// Make staging file.  It's just a very skeletal catalog.
}
