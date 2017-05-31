hitch
=====

`hitch` is two things:

- a metadata format specification, for tracking software releases (both source and artifacts), *as well as their build records, dependencies, and reproduce instructions*.
- a tool for manipulating, linting, and handling those release informations.

`hitch` is part of the [Repeatr](https://github.com/polydawn/repeatr) ecosystem.
While you might be able to integrate `hitch` with other projects as well, we usually expect it to Just Work in concert with operations like:

- `repeatr unpack` can take wares (e.g. filesystem snapshots) from `hitch` names and put them on your local disk;
- `repeatr scan` can look at remote URLs and emit a ware ID, which `hitch` can record, associate with a name, and sign for you;
- when `repeatr run` is used to build software or process data, `hitch` will record the Formula hash and RunRecords as well as associating the outputs with a name, and sign the whole bundle for you;
- and `hitch` can "replay" builds (of entire dependency trees, if desired!) by re-reading those recorded Formulas and invoking `repeatr run` again, so you can reproduce all your software from source!


Ecosystem (or, What Hitch Is Not)
---------------------------------

`hitch` is not a package manager for end-user host maintinance -- use REDACTED for more automation of that process (it gets its information and updates from `hitch).

`hitch` does not track upstream releases (source repos, etc) automatically -- but the REDACTED project does, and pushes that information into `hitch` (we use it to maintain our public builds tree).

`hitch` does not manage automatic builds -- but REDACTED does, by looking at a `hitch` database and formula templates and triggers, and calculating what new things need to be built!

`hitch` not a CI service -- look to REDACTED, which combines `hitch` and REDACTED (and even REDACTED clusters!) to automate builds triggered directly from source pushes.
