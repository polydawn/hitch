// meta note: this is revisiting the core concept of catalogs with the following context:
//  - recent prototype work trying to build a distro highlights a huge (huge) practical need for referencing *multiple* formulas as part of the justification for something
//    - also internal multi-steps to a single advertised result are non-zero frequency
//  - i want to start accumulating permanent information, but not commit to an update system
//  - one of our winter epiphanies is that if the resolvers and planners can't themselves operate in a container with the factbase as a pinned input and have reproducible outputs, it's probably a bug.


---

factbase
========

```
[ ! -z "$(explain --recursive --hazards)" ] && {echo nope nope nope}
```



big goals
---------

"How should I update my dependences" is a question that people can go to war over.

Therefore, we don't answer it.

Instead, PROJ_NAME_HERE provides a framework for describing what *is* released, *how* it was built, and simple *metadata* for describing versions.
(And for a bonus, any *verification* and gating steps that required additional containers can be described in a release record as well.)

"Metadata" is interpreted broadly, and the string grabbag there can be used to bind with almost any system you can imagine.
(We have some tools for seeing things as "release tracks"; you can use "semver"-style concepts just as easily.)
However, we also have a specific eye toward one thing: "hazard"-flagging metadata.

"Hazard" flags are kind of metadata also allows free text,
but declared separately so that all tooling can see it clearly, and any "hazard" metadata at all can be taken to mean "this should not be used".
This allows all different kinds of version/update/resolver/templater tools to consistently agree to reject anything that is marked as having CVEs, for example --
without actually requiring them all to share an understanding of "CVE", or any other security notice format for that matter.
Making "hazard" metadata explicit means safety protocols are always followed, even as we have diverse tooling for how to update.



hard questions
--------------

### relating to whether or not we really want to bring along "why" into the auditable portions

- What should happen if ware "tar:abcd" is released in catalog "foobar.cat" *and* "quaffle.cat", then is marked hazard in "foobar.cat"?
  - If we don't bring commissions and catalog names into the auditable portion:
    - Then this is tricky to answer clearly.
  - If we do bring along commissions and catalog names:
    - The usual objection: this binds us very strongly to having One single implementation of commissioning.
  - There is a middle ground option: We can bring along the catalog name for each ware, but *not* any of the commission content (e.g. version update rules).
    - This would let us do O(1) checking for a lot of things, including making this hazard handling question clear.
    - While also leaving the "when should I compute a new pinned formula" out of band of the factbase manager.

### relating to "what is the primary key of a release record"

- If `{catalogName,wareID}->{runRecords:map[setupHash]results, metadata:map[string]string}`:
  - The answer to the multi-arch question:
    - Clear: different arches tend to have different wareIDs, no?
	- Not handled: making a "release" of binaries all at once for $n$ arches.  Doesn't need to be clear; but nice if it is.  Especially if they share an asset prep pipeline that *is* arch-independent, or, we want to show that release is gated on things happening correctly for multiple architectures.
	- Implies: if two different arches *do* have the same release for some reason?  We'd need to support *sets* of values for a tag like "arch"; either that or a mess of asciifucking in the keyspace results, e.g. `tags:{"arch-amd64":"yes"}`.
- If `{catalogName,releaseName}->{wareID, runRecords:map[setupHash]results, metadata:map[string]string}`:
  - Bonus: well, you can prettyprint it now.
    - Drawback on other side of same coin: if people hold divergent face bases, this is problematic.  But then, this hits the "don't" horizon.
  - The answer to the multi-arch question:
    - Clear, but requires namecatting: basically we broke the sets problem by moving it to namecatting the arch onto the releaseName.
	- Not handled: making a "release" of binaries all at once for $n$ arches.
- If `{catalogName,releaseName}->{wares:[]{wareID, metadata:map[string]string}, runRecords:map[setupHash]results}`:
  - In other words: if a single release can offer a range of wareIDs.
    - This makes things alarmingly complicated IMO because now the tuple has gotten larger?!
  - Weird.  This either results in the set problem again for multiarch, or, posits the existence of another layer of "why" strings that invokes the namecatting problem again.
    - Apparently this is an innate problem that comes from "i want to name a thing, but then figure out what's appropriate for $context".
	  - ... duh, now that you mention it.
	- OR: if wares are a slice inside like this, then you're golden.
	  - unresolved: how fucking complex it is to pick a thing now, because really your resolver is ranging over the sets of metadata, so... yuff.

- The namecatting vs sets issue also appears if we want the same ware in multiple tracks, btw.  Which is *extremely* common in fact.

- Review this above stuff again thinking from the side of a resolver.
  - What is it supposed to do if two different release records contain all the same metadata?
