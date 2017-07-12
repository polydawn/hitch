hitch concepts overview
=======================

`hitch` is a release tracking and publishing tool for the [Timeless Stack](https://github.com/polydawn/timeless).


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
