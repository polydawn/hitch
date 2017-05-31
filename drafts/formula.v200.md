drafts of a big break in formula, separating out hashable bits vs context better:

```
mounts:
  "/":         "tar:589tu49yu984u5y"
  "/app/go":   "tar:349345345345345"
  "/task/src": "git:afafafafafaffaf"
action:
  hostname: "..."
  cwd: "..."
  env: {"key": "value"}
  hostname: "..."
  policy: "..."
  cradle: "..."
  escapes: "..."
  script:
    - oi
    - oioioi
    - mkdir /build/output
    - echo wow > /build/output
save:
  "/build/output": "tar"
results:
  "/build/output": "tar:123beefcake"
config: ## Nothing in here is hashable
  mounts:
    default: "directory.repeatr.io"
    saveto: "./wares" ## this is actually a default anyway
    for:
      "/app/go": "github.com/wizbang/wow.git"
  save:
    default: "directory.repeatr.io"
    saveto: "./wares" ## this is actually a default anyway
    for:
      "/build/output/": "./output"
```

... later ...


```
config: ## Nothing in here is hashable
  wares:
    default: "directory.repeatr.io"
    for:
      "/app/go": "github.com/wizbang/wow.git"
  save:
    default: "directory.repeatr.io"
    saveto: "./wares" ## this is actually a default anyway
    for:
      "/build/output/": "./output"
action:  # hash should include this
  wares:
    "/":         "tar:589tu49yu984u5y"
    "/app/go":   "tar:349345345345345"
    "/task/src": "git:afafafafafaffaf"
  hostname: "..."
  cwd: "..."
  env: {"key": "value"}
  hostname: "..."
  policy: "..."
  cradle: "..."
  escapes: "..."
  script:
    - oi
    - oioioi
    - mkdir /build/output
    - echo wow > /build/output
  save:
    "/build/output": "tar"
expects: # Expected results included for result hashing
  "/build/output": "tar:123beefcake"
```
