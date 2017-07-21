package main

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"go.polydawn.net/hitch/core"
)

func main() {
	exitCode := Main(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)
	os.Exit(int(exitCode))
}

var (
	app = kingpin.New("hitch", "Repeatr release control")

	initCmd = app.Command("init", "Initialize a new hitch DB.")

	createCmd                 = app.Command("create", "Grouping command for creating new things, like catalogs.")
	createCatalogCmd          = createCmd.Command("catalog", "Create a new catalog for publishing releases in.")
	releaseCmd                = app.Command("release", "Create new releases.")
	releaseStartCmd           = releaseCmd.Command("start", "Start staging a new release.  (Use more release commands to add data, then commit when done.)")
	releaseStart_CatalogArg   = releaseStartCmd.Arg("catalogName", "The name to assign this new step.").Required().String()
	releaseStart_ReleaseArg   = releaseStartCmd.Arg("releaseName", "The path to the formula file that runs this step.").Required().String()
	releaseAddStepCmd         = releaseCmd.Command("add-step", "Add a step to the replay instructions in the release currently being staged.")
	releaseAddStep_NameArg    = releaseAddStepCmd.Arg("name", "The name to assign this new step.").Required().String()
	releaseAddStep_FormulaArg = releaseAddStepCmd.Arg("formula", "The path to the formula file that runs this step.").Required().String()
	releaseAddStep_ImportsArg = releaseAddStepCmd.Arg("imports", "The path to an imports file which explains the catalogs and release names for wares in the step formula.").String()
	releaseAddItemCmd         = releaseCmd.Command("add-item", "Add a ware to the set of data in the release, and assigns it a name.")
	releaseAddItem_ItemArg    = releaseAddItemCmd.Arg("itemName", "The label to map the ware.").Required().String()
	releaseAddItem_WareArg    = releaseAddItemCmd.Arg("ware", "The WareID this label will be mapped to.").Required().String()

	showCmd = app.Command("show", "Show release info objects, or specific content hashes.")
)

func Main(args []string, stdin io.Reader, stdout, stderr io.Writer) (exitCode ExitCode) {
	app.HelpFlag.Short('h')

	// Rigging kingpin to use our in/out/err/code.
	var kingpinTerminate bool
	app.Terminate(func(status int) {
		kingpinTerminate = true
		exitCode = ExitCode(status)
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

	// Bundle UI handles.
	ui := core.UI{stdin, stdout, stderr}

	// Switch for command to invoke.
	err = func() error {
		switch cmd {
		case initCmd.FullCommand():
			return core.Init(ui)
		case releaseStartCmd.FullCommand():
			return core.ReleaseStart(ui, *releaseStart_CatalogArg, *releaseStart_ReleaseArg)
		case releaseAddItemCmd.FullCommand():
			return core.ReleaseAddItem(ui, *releaseAddItem_ItemArg, *releaseAddItem_WareArg)
		default:
			panic(fmt.Errorf("hitch: unhandled command %q", cmd))
		}
	}()
	if err == nil {
		return EXIT_SUCCESS
	}
	fmt.Fprintf(stderr, "%s\n", err)
	return mapToExitCode(err)
}
