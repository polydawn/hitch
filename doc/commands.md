hitch command line interface
============================

Release creation commands:

```
hitch release start
hitch release add-step <stepname> <formula.filename> [<importlist.filename>]
hitch release add-runrecord <stepname> <runrecord.filename>...
hitch release add-label <labelname> <wareID>
hitch release add-label <labelname> "wire:<stepname>:<outputSlot>"
hitch release metadata set <key>=<value>
hitch release commit <projectname> <releasename>
hitch release reset
```

Inspection commands:

```
hitch show <projectname>
hitch show <projectname>:<releasename>
hitch show <projectname>:<releasename>:<labelname>
hitch show-latest <projectname> [<labelname>]
```
