package main

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"go.polydawn.net/hitch/core"
)

func main() {
	os.Exit(Main(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

var (
	app = kingpin.New("hitch", "Repeatr release control")

	initCmd = app.Command("init", "Initialize a new hitch DB.")

	releaseCmd                = app.Command("release", "Create new releases.")
	releaseStartCmd           = releaseCmd.Command("start", "Start staging a new release.  (Use more release commands to add data, then commit when done.)")
	releaseStart_CatalogArg   = releaseStartCmd.Arg("catalogName", "The name to assign this new step.").Required().String()
	releaseStart_ReleaseArg   = releaseStartCmd.Arg("releaseName", "The path to the formula file that runs this step.").Required().String()
	releaseAddStepCmd         = releaseCmd.Command("add-step", "Add a step to the replay instructions in the release currently being staged.")
	releaseAddStep_NameArg    = releaseAddStepCmd.Arg("name", "The name to assign this new step.").Required().String()
	releaseAddStep_FormulaArg = releaseAddStepCmd.Arg("formula", "The path to the formula file that runs this step.").Required().String()
	releaseAddStep_ImportsArg = releaseAddStepCmd.Arg("imports", "The path to an imports file which explains the catalogs and release names for wares in the step formula.").String()

	showCmd = app.Command("show", "Show release info objects, or specific content hashes.")
)

func Main(args []string, stdin io.Reader, stdout, stderr io.Writer) (exitCode int) {
	app.HelpFlag.Short('h')

	// Rigging kingpin to use our in/out/err/code.
	var kingpinTerminate bool
	app.Terminate(func(status int) {
		kingpinTerminate = true
		exitCode = status
	})
	app.UsageWriter(stderr)
	app.ErrorWriter(stderr)

	// Parse flags; raise any immediate parse errors.
	cmd, err := app.Parse(args)
	if kingpinTerminate {
		return
	}
	if err != nil {
		fmt.Fprintf(stderr, "hitch: invalid command: %s\n", err)
		return 1
	}

	// Switch for command to invoke.
	switch cmd {
	case initCmd.FullCommand():
		core.Init()
		return 0
	case releaseStartCmd.FullCommand():
		return 0
	case releaseAddStepCmd.FullCommand():
		fmt.Fprintf(stdout, "whee!\n%q\n%q\n%q\n", *releaseAddStep_NameArg, *releaseAddStep_FormulaArg, *releaseAddStep_ImportsArg)
		return 0
	default:
		panic(fmt.Errorf("hitch: unhandled command %q", cmd))
		return 1
	}
}
