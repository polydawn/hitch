
```
pipeline:

  formulas:
    "<setup-hash>": {...}
    "<setup-hash>": {...}
    "<setup-hash>": {...}

  modules:
    "<name>": {
      basis: "<formulas.setup-hash>"
      customize:
        wares:
          "<key>": "pipe:<name>"
          "<key>": "tar:<hash>"
        saves:
          "<key>":

  wheres:
    "<type>:<hash>": "<fetchable>"
    "<type>:<hash>": "<fetchable>"
    "<type>:<hash>": "<fetchable>"

results:
  memos:
    "<module-name>": "<run-record-id>"
    "<module-name>": "<run-record-id>"
    "<module-name>": "<run-record-id>"
  pipes:
    "<pipe-name>": "<ware-type>:<hash>"
    "<pipe-name>": "<ware-type>:<hash>"
    "<pipe-name>": "<ware-type>:<hash>"
  records:
    "<run-record-id>": {...}
    "<run-record-id>": {...}
    "<run-record-id>": {...}
  saves:
    "<key>": "pipe:<name>"
```

more exampley:

```
pipeline:
  formulas:
    "1234beefcake":
      action:  # hash should include this
        wares:
          "/":         "tar:589tu49yu984u5y"
          "/task/src": "git:afafafafafaffaf"
        script:
          - build git
        save:
          "/task/git/bin": "tar"
      expects: # Expected results included for result hashing
        "/build/output": "tar:123beefcake"
    "foobar":
      action:
        wares:
          "/":         "tar:589tu49yu984u5y"
          "/test/relocatable/git": "tar:faffyfloop"
          script:
            - run tests
  modules:
    "build-git":
      basis: "1234beefcake"
      customize:
        "/task/git/bin": "pipe:git-bin"
    "test-git":
      basis: "foobar"
      customize:
        "/test/relocatable/git": "pipe:git-bin"
```

... later ...


```
pipeline:
  formulas:
    "1234beefcake":
      action:
        wares:
          "/":            "tar:589tu49yu984u5y"
          "/task/go-src": "tar:aaaaaaaa"
          "/task/go":     "tar:fffffff"
        script:
          - build go
        save:
          "/new/go/bin": "tar"
      expects:
        "/build/output": "tar:123beefcake"
  modules:
    "go-1.4-to-1.8":
      basis: "1234beefcake"
      customize:
        wares:
          "/task/go":     "tar:a-brave-new-hash"
          "/task/go-src": "tar:boop"
        expects:
          "/new/go/bin": "tar:expected-hash"
  wheres:
    "tar:a-brave-new-hash": "https://storage.googleapis.com/golang/go1.4.src.tar.gz"
    "tar:boop": "https://storage.googleapis.com/golang/go1.8.1.src.tar.gz"
```