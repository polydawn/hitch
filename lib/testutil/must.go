package testutil

// Wrapper method you can use for fixture setup calls.
// (You could also `So(setupStep(), ShouldBeNil)`, but
// this way reduces green checkmark spam.)
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
