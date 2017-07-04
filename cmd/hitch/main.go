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

func init() {
	// Check for existing hitch.db root anywhere above this.  Reject if exists.  Nested repos would be silly.

	// Make hitch.db sigil file in cwd.
}

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
