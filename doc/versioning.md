Versioning in Hitch
===================

First, an admonishment: Hitch publishes and tracks releases.
As covered in the [Labelling doc](./labels.md), you can call these releases anything you want.
Hitch does not inherently demand *any* versioning relationship between releases at all.

That said, there are a couple systems commonly used in the Timeless Stack ecosystem.
"Semver" -- probably familiar -- and another system called "tracks".

### Semantic Versioning ("semver")

Semantic Versioning -- TODO link spec site -- is sometimes used by release authors
to pick the names of their releases.

Hitch does not natively understand or do anything with this, but the convention is common.
Some pipelining and updater tools that work with the Timeless Stack and Hitch standards
will recognize sevmer release names and attempt to handle them per the SemVer spec.

### Tracks

Tracks are a new concept introduced by Hitch and the Timeless ecosystem.
Track info is often attached to hitch releases as metadata, and expresses how different
releases can be considered to fulfill various update lifecycles.

TODO

### Inherent Ordering

Okay, the admonishment at the top of the page ("hitch doesn't
have any inherent relationship between release") was a bit of a fib.

Releases inside a Hitch Catalog are ordered -- they're stored as an array.
We call this the "inherent ordering".
Whenever the hitch tools add a new release, they push it onto the top of the array.

We don't recommend *using* this for any serious purpose.
Hitch tools typically present information to the user in this same sorted order,
but any semantic usage -- e.g., affecting actual behavior, outside of the UI stuff --
is considered Bad Behavior.
The hitch specifications do not promise that the topmost release in the inherent ordering
is "better", "newer", "more secure", "more patched", or even "latest".
