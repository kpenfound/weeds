# Markdown discovery and global excludes

Probed with dagger v1.0.0-beta.27 and ripgrep 14.1.1. The question: can one
`Workspace.search`, limited to markdown files, honour every ignore source,
global git excludes included, for a non-empty markdown file? It cannot, so
`markdown` uses the shared file listing (the owner's ruling).

The in-module probe is `dagger call weeds-tests markdown-glob-probes`. The
host run is `markdown-global-excludes.sh`, which keeps HOME, and so the global
excludes, apart from the developer's settings.

## What limits a search to markdown

| Limit | Honours repo-local ignore sources | Honours global excludes | Keeps markdown in dotted or hidden directories |
| --- | --- | --- | --- |
| Positive glob (`**/*.[mM][dD]`) | No: it re-includes an ignored file | Directories yes, file patterns no | Yes |
| Negated globs for other extensions (`!**/*.[!mM]*`) | Yes | Yes | No |
| No globs (the shared listing) | Yes | Yes | Yes, but it lists every file type |

A negated glob matches directories as well as files, and ripgrep skips the
subtree of a directory it matches. An extension glob cannot tell the file
`a.go` from the directory `docs.v2` or `.github`, so every negated extension
glob that leaves out all non-markdown files also skips every directory whose
name has a dot, with the markdown inside it. The same globs do not exclude a
file with no extension, such as `Makefile`. Once a search has a positive
glob, a file that matches none is dropped and a file that matches one skips
every ignore rule. A negated glob after it can exclude only files it names, and
the files a global exclude hides are not known to weeds. A whitelisted
directory (`-g 'docs.v2/'`) whitelists none of its files: that search returns
nothing.

## What markdown does

`markdown` takes the markdown files out of the shared file listing
(`FileAccess.listed`, the search `comments` and `go-lint` also use), so every
check skips the same files on one host: repository ignore sources, global
excludes (file patterns and directories) and the standing skips. Ripgrep never
lists an empty file, so one more step adds the empty markdown files:

- one `Workspace.directory` listing limited to the markdown extensions, with
  `.gitignore` applied and the standing skips excluded, names the candidates;
- the candidates the shared listing did not return are either empty, hidden by
  an ignore source, or non-empty and hidden by a global exclude. One in-engine
  listing separates them. It runs on a tree of two copies of the repository's
  ignore files, each with its own `.git` marker so that each applies its rules
  by itself: `shown/` holds a placeholder per candidate, `held/` holds the real
  candidates. A candidate listed under `shown/` is not hidden by a repository
  ignore source; one also listed under `held/` holds text, so only a global
  exclude can have hidden it and it is dropped. A candidate in `shown/` only is
  empty and is kept.

The shared listing plus this step (the candidate names and the one in-engine
listing) are markdown's one pass. The in-engine listing runs on a tree of the
candidates and the ignore files, not on the workspace, and builds no second
host search.

A markdown file that ripgrep skips for its content, such as one holding a NUL
byte, is never listed with text, so it is treated like an empty file: it is
kept and reported unless a repository ignore source hides it. A non-ignored
markdown file is never dropped for its content.

## Recorded run

`markdown-global-excludes.sh` copies the module into a scratch workspace with
a minimal `.git` directory and runs `dagger call weeds markdown-report` with
`HOME` set to a scratch home, once with no global excludes and once with
`$HOME/.config/git/ignore` naming `global-excluded.*`,
`empty-global-excluded.*` and `globexcl-dir/`. The workspace holds `kept.md`,
`global-excluded.md`, `globexcl-dir/inner.md`, `.github/hidden.md`,
`docs.v2/dotted.md` and an empty `empty-global-excluded.md`.

```
== global excludes: none
.github/hidden.md:1: internal-markdown: ...
docs.v2/dotted.md:1: internal-markdown: ...
empty-global-excluded.md:1: internal-markdown: ...
global-excluded.md:1: internal-markdown: ...
globexcl-dir/inner.md:1: internal-markdown: ...
kept.md:1: internal-markdown: ...
== global excludes: xdg-default
.github/hidden.md:1: internal-markdown: ...
docs.v2/dotted.md:1: internal-markdown: ...
empty-global-excluded.md:1: internal-markdown: ...
kept.md:1: internal-markdown: ...
```

- A non-empty `.md` matched only by a global file pattern
  (`global-excluded.md`) gives no finding.
- A `.md` inside a globally excluded directory (`globexcl-dir/inner.md`) gives
  no finding.
- `.github/hidden.md` and `docs.v2/dotted.md` are still reported, so hidden
  and dotted directories are scanned.
- The one residual: the empty `empty-global-excluded.md` is still reported.
  Only the host-side search reads global excludes, and an empty file is never
  in it. CI has no such personal pattern, and CI decides.

## What a markdown-limited search does

`markdown-glob-probes` and the table above record why `markdown` does not run
its own limited `Workspace.search`:

```
== global excludes: xdg-default
shared listing: .github/hidden.md docs.v2/dotted.md kept.md
positive glob: .github/hidden.md docs.v2/dotted.md global-excluded.md kept.md
negated globs: kept.md
```

- The positive glob returns the non-empty `global-excluded.md` that the global
  file pattern hides.
- The negated globs honour every exclude, and lose `.github/hidden.md` and
  `docs.v2/dotted.md`.
