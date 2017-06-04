rebuilding repeat-thyself, with hitch
=====================================

This is a brainstorming document for what repeatr.git//repeat-thyself.sh should look like,
if it's upgraded to take into account hitch as part of its universe: both to keep
things updated, and also to release things.

For most projects, this will not be the path to follow: we expect to build more
"planners" as part of the ecosystem, which will be much more robust tools that
take on these reponsibilities.
We're exploring it here based on the theory that this *should* be possible to
do with a basic shell script instead (not necessarily graceful, but *possible*);
and the repeatr build process is an interesting example simply because bootstrapping
systems really should hew to the minimum complexity anyway.

Some features are speculative in the extreme (e.g. "--track", which likely
does *not* in fact belong in hitch's responsibilities).


overview
--------

There are several different phases:

- Stamping out the formula to run, and saving the corresponding upstream import labels doc.
- Running formulas!
- Choosing to release.

For repeat-thyself, we would likely group a couple of these phases together -- splitting some (!) and joining others:

1. Updating the import selection.
  - This is actually one half of stamping out the formula: it's *just* selecting the ReleaseItemID tuples.
  - This is likely to be run rarely: we don't need every single build to reconsider whether it wants to jump from go1.8 to 1.8.1 (and also, answering that question is *not* hitch's job!).
  - TODO: REVIEW this at LENGTH.  ISTM we still need some answers here: I *don't* actually want the ReleaseName for imports to auto-update, but I *do* want my track info to be stored in a way standard enough that I can run a command that checks for updates.
    - I guess you could run such a tool on the import-labels.json file.  Though I don't know why such a tool should know about the input paths.  (Does this suggest we're doing something mildly wrong, and imports should actually be at the scope of the Replay rather than each Step, and they're required to be "wire" objects in a step?  That would actually make some things much more consistent... just verbose: it would irritate me to write a local name for something I'm only going to use once if I have a single step.)

2. Doing run-the-first.
  - Stamp the full formula.  This includes shelling out to git to see the latest (we're not going to proxy that through hitch because it seems silly to bother it to do so).
  - Immediately run it, because that's what I want to see.
  - Save that runRecord in `./runrecord-1.json`.

3. Optional: re-run, and just say whether it reproduced in your exit code.
  - Literally, just run it again and output `./runrecord-2.json`, and check result equality.

4. Release.
  - Do run-the-second if it's not been done yet.  Do the result match check again.
  - If it passes, start invoking hitch commands to stage the release.

Steps 2 and 4 are the only ones you'd really expect to run with any frequency.


scripts
-------

### build one

```bash
RELEASEID_BASE="ports.repeatr.io/base:20170501:linux-amd64"
RELEASEID_GO="ports.repeatr.io/golang:1.8:linux-amd64"

WARE_BASE=`hitch select "$RELEASEID_BASE"`
WARE_GO=`hitch select "$RELEASEID_GO"`
WARE_SRC="git:"`git rev-parse HEAD`

for target in linux-amd64 darwin-amd64; do {
	mkdir "$target"

	cat <<EOF
		{
			"/":       "$RELEASEID_BASE",
			"/app/go": "$RELEASEID_GO",
		}
EOF > "$target/import-labels.json"

	cat <<EOF
		inputs:
			"/":        "$WARE_BASE"
			"/app/go/":	"$WARE_GO"
			"/task/":   "$WARE_SRC"
		action:
			command: ["/bin/bash", "./goad install"]
		outputs:
			"/task/bin/": "tar"
EOF > "$target/build.formula"

	# ...
}; done
```

### rebuild, check, & release

// todo
