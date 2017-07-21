#!/bin/bash
##
## Hello!
## This is a demo script that shows of use of hitch, interactively.
## The script follows a straight-and-narrow path, but pauses frequently;
## take your time to observe what's going on, and press enter to continue
## when you're ready.
##
## You can also run the script as `./demo.sh -t` to run straight through
## with no pauses.
##
## (This script is also used as a smoke test in our standard CI tests.
## Thus, the whole thing should exit without errors -- sometimes
## we demonstrate error paths, and those require some more complicated
## shell script features like subshells, etc.  It's *just* for our testing;
## normal users of hitch don't need to get fancy like this; it's always
## correct to simply stop when hitch errors.)
##
set -euo pipefail

## Normalize path.  Include helper libraries (mostly, prettyprinters).
cd "$( dirname "${BASH_SOURCE[0]}" )"
source ./meta/demo/util.shlib

## If there's a local fresh build of hitch, use that.
if [ -x bin/hitch ]; then PATH=$PWD/bin/:$PATH; fi
if [ -z "$(which hitch)" ]; then die "missing \`hitch\` command on \$PATH!"; fi

demodir="/tmp/hitch-demo";
rm -rf "$demodir"
mkdir -p "$demodir" && cd "$demodir" && demodir="$(pwd)"

(
	loudly hitch init
	loudly hitch release start "proj.net/team/thing" "v0.1"
	loudly hitch release add-item "thingy" "tar:ah4sh28vBa39DVuwI"
) 2>&1| section "Using hitch is easy!" "See?"

awaitack
