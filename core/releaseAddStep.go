package core

func releaseAddStep(stepName, formulaFilename, importsFilename string) {
	// TODO programming style: do we need a bag for all the IO and code stuff?  seems yes.
	// Because you might still want to flip the "-" args to a single read of stdin here.

	// Find hitch.db root.

	// Check for staging file.  Must exist; reject command if no release has been started.

	// Check that this step name is not already in use.  Reject if is.

	// Check both files exist and parse.
	// If the args are '-', read stdin. (If both filenames are dash, only do so once; use the same buffer twice.)

	// Conditionally: if per-op validation is enabled:
	//  Look over all the imports; if any are "wire", check those are already available.

	// Append the replay object.  Write back out.
}
