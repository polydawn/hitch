hitch command line interface
============================

### Release creation commands:

```
hitch create catalog <catalogname>
hitch release start <catalogname> <releasename>
hitch release add-step <stepname> <formula.filename> [<importlist.filename>]
hitch release add-runrecord <stepname> <runrecord.filename>...
hitch release add-item <itemname> <wareID>
hitch release add-item <itemname> "wire:<stepname>:<outputSlot>"
hitch release metadata set <key>=<value>
hitch release commit
hitch release reset
```

### Inspection commands:

```
hitch show <catalogname>
hitch show <catalogname>:<releasename>
hitch show <catalogname>:<releasename>:<itemname>
hitch show-latest <catalogname> [<itemname>]
hitch list catalogs
hitch list releases <catalogname>
hitch list items <catalogname>:<releasename>
hitch show-replay <catalogname>:<releasename>
```

### Exchanging catalogs with other hitch datasets:

```
hitch import raw <catalog.filename>
hitch import repo <otherhitch.path> <catalogname>
```



---

full descriptions
-----------------

### Release creation commands:

TODO

### Inspection commands:

`hitch show` commands dump json of all object info under your request.

- `hitch show someproject` dumps the top-level project name, release keys, and description, *all* the releases, all the metadata for each release, and all the release items: the wareIDs and their labels.
- `hitch show someproject:somerelease` dumps the information for the named release: all the metadata, and all the wares and their labels in that release.
- `hitch show someproject:somerelease:someitem` emits a single string: the WareID with for that specific label.

The `hitch show-latest` command will behave as per the second or third forms of `hitch show`, but selecting the latest release in the named project catalog.
(Be careful using this command: it is for convenience and debugging; if you can specify your desired versions in any better way than this, please do so!)

The `hitch show-replay` command dumps json describing the complete set of steps that built the the wares in a release:
each step contains a formula, complete with its pinned input WareIDs;
the import information of which other catalogs and releases provided those WareIDs;
and the run records and resulting WareIDs those formulas built.
(Replay instructions are not included in any of the `hitch show` command outputs, due to verbosity.)

The `hitch list` commands show a list of names: either
all the catalogs in the database,
all the releases in a catalog,
or all the item labels in a release in a catalog.
`hitch list` defaults to a simple line-separated strings format,
easy to loop over in shell scripting.
(The output of `hitch list` is a subset of `hitch show`.)

### Exchanging catalogs with other hitch datasets:

TODO



---

this is a draft
---------------

### Open questions

- We haven't discussed 'tracks'.  Should we?
  - The dream is to keep 'tracks', like sevmer, slightly out of core.  But...
  - The `hitch show-latest` command is *particularly* dubious.  If we're going to include it, we *have* to warn about it; but if we're going to warn about something, we should include a recommendation about what to use instead!
- Should `hitch show-replay` be broken into more subcommands?
  - It's probably useful if we can emit individual formulas, for example.
    - You probably shouldn't need `jq`-like glue if you want to offhandedly `repeatr run` a single step.
    - We've discussed whether the canonical hitch.db file layout should have formulas spread out this way for the same reasons.  The hitch CLI is even more important in this regard.

---
