package main

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	os.Exit(Main(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

var (
	app = kingpin.New("hitch", "Repeatr release control")
)

func Main(args []string, stdin io.Reader, stdout, stderr io.Writer) (exitCode int) {
	app.HelpFlag.Short('h')
	cmd, err := app.Parse(args)
	if err != nil {
		fmt.Fprintf(stderr, "hitch: invalid command: %s\n", err)
		return 1
	}
	switch cmd {
	// todo more
	default:
		fmt.Fprintf(stderr, "hitch: missing subcommand!  try 'hitch -h' for help.\n")
		return 1
	}
}
