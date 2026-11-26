# How ignore sources reach Workspace.search and uploads

Probed with dagger v1.0.0-beta.27 and ripgrep 14.1.1. The in-module probes
(`dagger call weeds-tests ignore-probes`, `scan-ignored`, `file-layer`,
`listed-with-empty`) cover every source a module can express. The global git
excludes cannot be set from inside a module, so `global-excludes.sh` shows them
on a real host workspace.

## Findings

| Question | Result |
| --- | --- |
| Search on a directory-backed workspace with no `.git` | Honours `.ignore` and `.rgignore` (nested too), not `.gitignore`. A `.git` entry at or above the root, such as a `.git/HEAD` file, makes it honour all three. The test harness adds one. |
| Search without ripgrep | Dagger falls back to `grep -rnH -E`: every ignored file is returned and no result carries a submatch. Weeds treats a result without a submatch as proof of the fallback and fails naming ripgrep. |
| `Workspace.directory(gitignore: true)` | Reads `.gitignore` only, nested files included, no `.git` needed. It does not read `.ignore`, `.rgignore`, `.git/info/exclude` or global excludes. |
| `.git/info/exclude` | Search honours it, with or without `.git/HEAD`. The upload option does not. |
| Global excludes (`core.excludesFile`, `$XDG_CONFIG_HOME/git/ignore`, `$HOME/.config/git/ignore`) | Search honours them when ripgrep runs on the host. The upload option does not. A directory-backed workspace runs ripgrep inside the engine, so a module test cannot configure them. |
| Positive search globs | They re-include ignored files, so include globs are applied in Dang. Negated globs keep ignore rules. |
| `filesOnly` and multiline `\A` | An empty file or a file holding a NUL byte is never listed. |

## Consequence for the shared layer

The upload helper builds its Directory from the search's file list, so uploads
hold exactly the files the search sees for every source. The listing helper
adds empty files that no repository ignore source hides, by checking them in a
small in-engine workspace that holds the repository's ignore files. An empty
file hidden only by a user's global excludes is still listed; CI decides.

For go-lint, `uploadWithSupport` adds the non-Go files a build needs (embedded
assets, empty or binary testdata) through a direct `Workspace.directory`
upload with `.gitignore` applied. `.go` files and `go.mod` discovery come only
from the search list. A support file hidden only by `.ignore`, `.rgignore`,
`.git/info/exclude` or a global exclude can reach the upload, but
golangci-lint reports on `.go` files only, so it cannot produce a finding.

## Mapping to a consumer's workspace

A consumer's workspace is a host directory. `Workspace.search` runs ripgrep on
the machine that runs dagger, with that machine's git configuration and global
ignore file, so a laptop and a CI runner can differ. `Workspace.directory`
filters with the engine's `.gitignore` handling, so weeds builds its uploads
from the search list and the two agree.

## Recorded run

`global-excludes.sh` copies the module into a scratch workspace with a minimal
`.git` directory whose `info/exclude` names `info-excluded.dat`, writes
`kept.dat`, `info-excluded.dat` and `global-excluded.dat` (each holding the
searched text), and runs `host-ignore-report` with `HOME` set to a scratch home.

```
== global excludes: none
search: global-excluded.dat:1 kept.dat:1
shared upload: global-excluded.dat kept.dat
plain upload with gitignore: global-excluded.dat info-excluded.dat kept.dat
== global excludes: xdg-default
search: kept.dat:1
shared upload: kept.dat
plain upload with gitignore: global-excluded.dat info-excluded.dat kept.dat
== global excludes: core-excludesfile
search: kept.dat:1
shared upload: kept.dat
plain upload with gitignore: global-excluded.dat info-excluded.dat kept.dat
```

`info-excluded.dat` is absent from the search and the shared upload in every
run. `global-excluded.dat` is absent from both once the scratch home names it
through `.config/git/ignore` or through `core.excludesFile`. The plain upload
keeps both files, which is why weeds does not use it.
