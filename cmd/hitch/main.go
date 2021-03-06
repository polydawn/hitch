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

	// all the authoring commands:

	catalogCmd                = app.Command("catalog", "Manage and create catalogs for releases.")
	catalogCreateCmd          = catalogCmd.Command("create", "Create a new catalog, representing a project, and ready for publishing releases into.")
	catalogCreate_CatalogArg  = catalogCreateCmd.Arg("catalogName", "The name for the new catalog.  (Typically, patterned like \"yourorg.net/team/projname\".)").Required().String()
	releaseCmd                = app.Command("release", "Create new releases.")
	releaseStartCmd           = releaseCmd.Command("start", "Start staging a new release.  (Use more release commands to add data, then commit when done.)")
	releaseStart_CatalogArg   = releaseStartCmd.Arg("catalogName", "The catalog this release will be committed to (when finished).").Required().String()
	releaseStart_ReleaseArg   = releaseStartCmd.Arg("releaseName", "The path to the formula file that runs this step.").Required().String()
	releaseAddStepCmd         = releaseCmd.Command("add-step", "Add a step to the replay instructions in the release currently being staged.")
	releaseAddStep_NameArg    = releaseAddStepCmd.Arg("name", "The name to assign this new step.").Required().String()
	releaseAddStep_FormulaArg = releaseAddStepCmd.Arg("formula", "The path to the formula file that runs this step.").Required().String()
	releaseAddStep_ImportsArg = releaseAddStepCmd.Arg("imports", "The path to an imports file which explains the catalogs and release names for wares in the step formula.").String()
	releaseAddItemCmd         = releaseCmd.Command("add-item", "Add a ware to the set of data in the release, and assigns it a name.")
	releaseAddItem_ItemArg    = releaseAddItemCmd.Arg("itemName", "The label to map the ware.").Required().String()
	releaseAddItem_WareArg    = releaseAddItemCmd.Arg("ware", "The WareID this label will be mapped to.").Required().String()
	releaseCommitCmd          = releaseCmd.Command("commit", "Finish the release currently being staged: commit it to the hitch database.")

	// inspection commands:

	showCmd      = app.Command("show", "Show release info objects, or specific content hashes.")
	show_nameArg = showCmd.Arg("name", "What to show.  Amount of data shown depends on how specific a name is given: a catalog name alone yields a large amount of data; a catalog:release:item tuple shows just that one WareID.").Required().String()

	// import and sync commands:
	// FUTURE :)
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
		case catalogCreateCmd.FullCommand():
			return core.CatalogCreate(ui, *catalogCreate_CatalogArg)
		case releaseStartCmd.FullCommand():
			return core.ReleaseStart(ui, *releaseStart_CatalogArg, *releaseStart_ReleaseArg)
		case releaseAddItemCmd.FullCommand():
			return core.ReleaseAddItem(ui, *releaseAddItem_ItemArg, *releaseAddItem_WareArg)
		case releaseCommitCmd.FullCommand():
			return core.ReleaseCommit(ui)
		case showCmd.FullCommand():
			return core.Show(ui, *show_nameArg)
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
