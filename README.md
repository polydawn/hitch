hitch [![GoDoc](https://godoc.org/github.com/polydawn/hitch?status.svg)](https://godoc.org/github.com/polydawn/hitch)
=====

`hitch` is two things:

- a metadata format specification, for tracking software releases (both source and artifacts), as well as their *build records*, *dependencies*, and *reproduce instructions*.
- a tool for manipulating, linting, and handling those release informations.

`hitch` is part of the [Repeatr](https://github.com/polydawn/repeatr) ecosystem.
While you might be able to integrate `hitch` with other projects as well, we usually expect it to Just Work in concert with operations like:

- `repeatr unpack` can take wares (e.g. filesystem snapshots) from `hitch` names and put them on your local disk;
- `repeatr scan` can look at remote URLs and emit a ware ID, which `hitch` can record, associate with a name, and sign for you;
- when `repeatr run` is used to build software or process data, `hitch` will record the Formula hash and RunRecords as well as associating the outputs with a name, and sign the whole bundle for you;
- and `hitch` can "replay" builds (of entire dependency trees, if desired!) by re-reading those recorded Formulas and invoking `repeatr run` again, so you can reproduce all your software from source!

---

What does hitch solve?
----------------------

Name stuff:

```
repeatr pack "./" | hitch name "projname.org/thing:1.8:amd64-linux"
```

Fetch stuff by name:

```
repeatr unpack "$(hitch show "projname.org/thing:1.8:amd64-linux")"
```

Share names:

```
hitch pull github.com/polydawn/hitch-records
```

Audit stuff:

```
[ ! -z "$(hitch explain "$WARE_HASH" --recursive --hazards)" ] \
  && {echo nope nope nope; exit 1; }
```

---

Ecosystem (or, What Hitch Is Not)
---------------------------------

`hitch` is not a package manager for end-user host maintenance -- use REDACTED for more automation of that process (it gets its information and updates from `hitch`).

`hitch` does not track upstream releases (source repos, etc) automatically -- but the REDACTED project does, and pushes that information into `hitch` (we use it to maintain our public builds tree).

`hitch` does not manage automatic builds -- but REDACTED does, by looking at a `hitch` database and formula templates and triggers, and calculating what new things need to be built!

`hitch` is not a CI service -- look to REDACTED, which combines `hitch` and REDACTED (and even REDACTED clusters!) to automate builds triggered directly from source pushes.

---

What are 'updates', really?
---------------------------


"How should I update my dependences" is a question that people can go to war over.

Therefore, we don't answer it.

Instead, `hitch` provides a framework for describing what *is* released, *how* it was built, and simple *metadata* for describing versions.
(And for a bonus, any *verification* and gating steps that required additional containers can be described in a release record as well.)

"Metadata" is interpreted broadly, and the string grabbag there can be used to bind with almost any system you can imagine.
(We have some tools for seeing things as "release tracks"; you can use "semver"-style concepts just as easily.)
However, we also have a specific eye toward one thing: "hazard"-flagging metadata.

"Hazard" flags are kind of metadata also allows free text,
but declared separately so that all tooling can see it clearly, and any "hazard" metadata at all can be taken to mean "this should not be used".
This allows all different kinds of version/update/resolver/templater tools to consistently agree to reject anything that is marked as having CVEs, for example --
without actually requiring them all to share an understanding of "CVE", or any other security notice format for that matter.
Making "hazard" metadata explicit means safety protocols are always followed, even as we have diverse tooling for how to update.

`hitch` provides information --
focused on keeping you safe and secure,
and ensuring recursive audit is available for posterity --
but stops short of forcing updates, and never tries to sneak in "small" changes, thus never getting in the way of how you want to develop.
