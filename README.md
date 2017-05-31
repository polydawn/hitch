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
