hitch concepts overview
=======================

`hitch` is a release tracking and publishing tool for the [Timeless Stack](https://github.com/polydawn/timeless).


Context: where Hitch fits in to the Timeless Stack
--------------------------------------------------

### Core Concepts from the Timeless Stack

Remember that in the Timeless Stack, we already have some established ways of communicating data:
- everything is a Fileset;
- a Fileset is simply a slice of filesystem -- files and dirs, complete with metadata;
- and Filesets can be packed into Wares,
- and packing into Wares defines an immutable, content-addressable WareID.

WareIDs are very useful because they're immutable (so they're perfect for pinning precise versions of things),
and content-addressable (so they're perfect for identifying things, even in a decentralized system).

WareIDs also aren't human-readable.  (They're sort of like a git hash.)

So, just as git uses "branches" and "tags" to assign human-readable names to commits,
we also need something to help us map human-readable names onto WareIDs in the Timeless Stack.
**That's what Hitch does.**


### What Hitch Does

Go look at the ecosystem diagram in the [Timeless Stack docs repo](https://github.com/polydawn/timeless).

Got it?  Good.

Hitch is for tracking and publishing releases -- mapping human-readable names (and version info) onto WareIDs.
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

- to inspect published releases,
- to sync your release database with someone,
- to do customize manual details of making release (but again,
  mostly a planner tool should be doing this work for you),
- or, if you're writing scripts to work as your own planner. :)


What is a Release?
------------------

Long story short, a Hitch "release" maps a human-readable name to a WareID.
Once data is encapsulated into a Hitch release, Hitch tools can easily send,
recieve, and merge snapshots of release info with others.

Long story long, we do a few more things:

- Releases are gathered together in a "Catalog".
   We use this to represent authorship over time.
- Releases can refer to *several* WareIDs, each with their own "[Label](./labels.md)".
   We use this to describe release of e.g. a linux build and a mac build of the same project.
- Metadata to help with [Versioning](./versioning.md) can be attached to each release.
   Other parts of the Timeless Stack use this info to compute automatic updates.
- Instructions for building the named WareIDs can be included in the Release -- these are called [Replay](./replays.md) instructions,
   and they're in a standardized, container-based format which *you can evaluate automatically*.

So, when you make a release, you can include all this data.
Hitch will store it in a standard format, which makes it easy to query programmatically,
and makes the information consumable by the rest of the Timeless Stack ecosystem.
