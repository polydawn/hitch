hitch concepts overview
=======================

`hitch` is a release tracking and publishing tool for the [Timeless Stack](https://github.com/polydawn/timeless).


Context: where Hitch fits in to the Timeless Stack
--------------------------------------------------

Go look at the ecosystem diagram in the [Timeless Stack docs repo](https://github.com/polydawn/timeless).

Got it?  Good.

Hitch is for tracking and publishing releases.
That is the *only* thing it does: maintains files with snapshots of structured release
information, and helps communicate them between different people and machines.
Hitch is basically a database.

To get real mileage out of Hitch, you probably want to use it together with other tools.
For example, you might want to use it with some other "planner" tool in the Timeless Stack:
the planner will be responsible for taking rough instructions of what do build and run,
ask Hitch to provide it with a bunch of version information so the planner can select *precise* versions of things,
and after then *running* those builds, the planner will come back around to talk to Hitch
again, to save and publish the results.

Situations in which you want to use Hitch directly as a user are mostly:

- to inspect published releases
- to sync your release database with someone
- to do customize manual details of making release (but again,
  mostly a planner tool should be doing this work for you)
- you're writing scripts to work as your own planner :)


What is a Release?
------------------

First, context:

In the Timeless Stack, everything is communicated as a FileSet -- simply a
slice of filesystem; files and dirs, complete with metadata -- which are packed into Wares,
and identified by [content-addressable](TODO:wikipedia link) WareID.
The trouble with WareIDs is they aren't human readable; content-addressable is
great for being precise (even to the point it becomes a security property),
but since a WareID can't be *chosen*, we need some

A `hitch` "release" is the process of picking a (human-readable!) name, and
declaring that it maps to a WareID.
This resolves the problem of opaqueness in WareIDs.
Making a `hitch` release is a way to publish semantic identities of Wares.
