package main_test

import (
	"bytes"
	"io"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"."
)

func callMain(args []string, stdin io.Reader) (string, string, main.ExitCode) {
	if stdin == nil {
		stdin = &bytes.Buffer{}
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := main.Main(args, stdin, stdout, stderr)
	return stdout.String(), stderr.String(), code
}

func TestCLI(t *testing.T) {
	Convey("Test CLI", t, func() {
		Convey("'hitch -h' should emit help", func() {
			stdout, stderr, code := callMain(
				[]string{"-h"}, nil,
			)
			So(stderr, ShouldStartWith, "usage: hitch")
			So(stderr, ShouldNotContainSubstring, "missing subcommand")
			So(stdout, ShouldEqual, "")
			So(code, ShouldEqual, 0)
		})
		Convey("'hitch -x' should error about unknown flag", func() {
			stdout, stderr, code := callMain(
				[]string{"-x"}, nil,
			)
			So(stderr, ShouldContainSubstring, "invalid command")
			So(stdout, ShouldEqual, "")
			So(code, ShouldEqual, 1)
		})
		Convey("'hitch asdf' should error about unknown command", func() {
			stdout, stderr, code := callMain(
				[]string{"asdf"}, nil,
			)
			So(stderr, ShouldContainSubstring, "invalid command")
			So(stdout, ShouldEqual, "")
			So(code, ShouldEqual, 1)
		})
	})
}
