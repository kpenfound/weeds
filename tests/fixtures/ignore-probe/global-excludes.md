# How ignore sources reach Workspace.search and uploads

Probed with dagger v1.0.0-beta.27 and ripgrep 14.1.1. The in-module probes
(`dagger call weeds-tests ignore-probes`, `scan-ignored`, `file-layer`,
`listed-with-empty`, `ignore-files`) cover every source a module can express.
The global git excludes cannot be set from inside a module, so
`global-excludes.sh` shows them on a real host workspace, through the real
check entry points.

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
`.git` directory and runs dagger with `HOME` set to a scratch home, once with no
global excludes, once with `$HOME/.config/git/ignore` (the XDG default) and once
with `core.excludesFile`. The workspace's `.git/info/exclude` names
`info-excluded.*`, and also the module's own files so that the scans see only
the files below. For each of `kept`, `info-excluded` and `global-excluded` it
writes a `.dat` file holding the searched text, a `.ts` file with a
suppression directive and a bad comment, an internal `.md` file and a `.go`
file with an unchecked error, plus a root `go.mod`. It then runs the
`host-ignore-report` test entry point and the real entry points `weeds
go-lint-modules`, `suppressions-report`, `markdown-report`, `comments-report`
and `go-lint-report`. The three global-excludes runs differ only in
`global-excluded.*`.

```
== global excludes: none
search: global-excluded.dat:1 kept.dat:1
shared upload: global-excluded.dat kept.dat
plain upload with gitignore: global-excluded.dat info-excluded.dat kept.dat
-- go-lint-modules
.
-- suppressions-report
global-excluded.ts:1: ts-ignore: ...
kept.ts:1: ts-ignore: ...
-- markdown-report
global-excluded.md:1: internal-markdown: ...
kept.md:1: internal-markdown: ...
-- comments-report
global-excluded.ts:1: ts-comment-not-doc: ...
kept.ts:1: ts-comment-not-doc: ...
-- go-lint-report
global-excluded.go:7: errcheck: Error return value of `os.Remove` is not checked
kept.go:7: errcheck: Error return value of `os.Remove` is not checked
== global excludes: xdg-default      (core-excludesfile prints the same)
search: kept.dat:1
shared upload: kept.dat
plain upload with gitignore: global-excluded.dat info-excluded.dat kept.dat
-- go-lint-modules
.
-- suppressions-report
kept.ts:1: ts-ignore: ...
-- markdown-report
kept.md:1: internal-markdown: ...
-- comments-report
kept.ts:1: ts-comment-not-doc: ...
-- go-lint-report
kept.go:7: errcheck: Error return value of `os.Remove` is not checked
```

`info-excluded.*` gives no finding and no search or upload entry in any run.
`global-excluded.*` gives none once the scratch home names it through
`.config/git/ignore` or through `core.excludesFile`, and the real checks,
the search and the shared upload agree. The plain upload keeps both files,
which is why Weeds does not use it.
