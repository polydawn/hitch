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

```bash
WARE_BASE=`hitch select ports.repeatr.io/base::linux-amd64 --track=latest`
WARE_GO=`hitch select ports.repeatr.io/golang:1.8:linux-amd64`
WARE_SRC="git:"`git rev-parse HEAD`

for target in linux-amd64 darwin-amd64; do {
	mkdir "$target"

	cat <<EOF
		{
			"/": "$WARE_BASE",
			"/app/go": "$WARE_GO",
			"/task/src": "$WARE_SRC",
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