# weeds

Fight the weeds.

Bad code spreads like weeds. Comments and internal markdown let agents excuse
bad code instead of fixing it. Skills only advise, and they do not reach every
agent. A `dagger check` is the one enforcement every agent gets.

Weeds is a [Dagger](https://dagger.io) module, written in
[Dang](https://github.com/vito/dang), that enforces a "no excuses" code policy
with static analysis only: no model, no classifier, no secret, and no consumer
code leaving the machine. It also installs the matching agent skills from
[kpenfound/skills](https://github.com/kpenfound/skills) with `dagger generate`.

## Install

Weeds needs [ripgrep](https://github.com/BurntSushi/ripgrep) (`rg`) on the
`PATH` of the machine that runs `dagger`, a laptop or a CI runner. The scanning
checks find files with `Workspace.search`, which runs ripgrep on that machine
and applies `.gitignore`, `.ignore`, `.rgignore`, `.git/info/exclude` and the
machine's global git excludes. Without ripgrep Dagger falls back to `grep`,
which honours none of them, so Weeds stops with an error that names ripgrep
instead of scanning ignored files.

1. Run `dagger install github.com/kpenfound/weeds`
2. Give the repository somewhere to install skills: an `AGENTS.md` (skills go
   to `.agents/skills`), a `CLAUDE.md` (`.claude/skills`), or an existing
   skills directory.
3. Run `dagger generate -y` to install the skills and commit the result.
4. Run `dagger check` to run the five checks. `installed-skills` fails when
   the installed skill folders drift from what `dagger generate` writes.
5. Optionally add a `weeds.toml` at the repository root (see below) to name
   exceptions, a markdown allowlist, a go-lint baseline or a skills subset.

Both commands fetch Weeds and the skills repo without credentials. To move to
a newer Weeds, change the pinned SHA, run `dagger generate` and review the
diff.

## Checks

Every check reports findings as `path:line: rule: message` lines in one
deterministic order, so running a check twice on the same workspace gives
byte-identical output. Every check passes its findings through the same
exception evaluator driven by `weeds.toml`.

| Check | Enforces |
|---|---|
| `comments` | The comment policy for Go and TS/JS. |
| `suppressions` | No inline lint or type suppression directives in any file kind. |
| `markdown` | No internal markdown beyond the always-allowed set. |
| `go-lint` | A pinned golangci-lint preset, once per Go module. |
| `installed-skills` | The installed skill folders equal the installer's output. |

### suppressions

One `Workspace.search` for all eight directives over every file kind that
ripgrep reads (see Scanning). It reports each line that holds one of
`//nolint`, `lint:ignore`, `eslint-disable`, `@ts-expect-error`, `@ts-ignore`,
`@ts-nocheck`, `ast-grep-ignore` or `# noqa`, one finding per directive per
line, under a rule named for the directive (`nolint`, `lint-ignore`,
`eslint-disable`, `ts-expect-error`, `ts-ignore`, `ts-nocheck`,
`ast-grep-ignore`, `noqa`). Generated files are scanned like any other. Only a
named exception in `weeds.toml` clears a finding; no inline text does.

### markdown

Markdown files (`.md`, `.markdown`, `.mdx`, in any letter case) come from one
listing limited to those extensions, empty files included, with ignore files
applied. They are reported under the rule `internal-markdown`, except these,
which are always allowed:

- README, AGENTS, CLAUDE, CONTRIBUTING, LICENSE, changelog and release-notes
  files, matched on the lowercased file name without its extension, so
  `claude.md` passes and `AGENTS-design.md` does not;
- anything under a `docs/` directory;
- anything under `testdata/` or a test fixture directory (`fixtures`,
  `__fixtures__` and similar);
- Weeds-owned skill folders (see below).

Any other markdown file is a finding unless the `weeds.toml` markdown
allowlist names it. An allowlist entry that matches no markdown file needing
it, including one that names an always-allowed file, is stale and fails.

### go-lint

Runs a pinned golangci-lint image (`golangci/golangci-lint:v2.15.0`, pinned by
tag and sha256 digest) once per Go module: every `go.mod` outside `vendor`,
`testdata` and `node_modules` whose module holds Go files. Weeds writes its own
preset and passes it with `--config`, so a golangci-lint configuration in your
repository is never read. The preset enables:

- correctness: `errcheck`, `govet`, `ineffassign`, `staticcheck`;
- dead code: `unused`;
- `godox` (TODO, BUG and FIXME comments);
- complexity: `funlen` (60 lines, 40 statements), `gocognit` (15), `nestif`
  (4) and `dupl` (80 tokens).

Every issue is reported, none are capped or merged, and a `//nolint` directive
does not hide a go-lint finding (the `suppressions` check reports the
directive itself). The rule of a finding is the linter's name, and paths are
relative to the workspace root. A module that does not build, a golangci-lint
failure or output Weeds cannot read raises an error instead of passing. The Go
module and build caches are kept between runs, and go.mod's `go` directive must
not need a newer toolchain than the image carries.

Modules come from the files the workspace search sees, so a copy of a module in
an ignored directory (an agent worktree, say) is neither a module nor linted.
golangci-lint gets the Go sources, `go.mod`, `go.sum` and `go.work` files from
that same search, so an ignored one never reaches it. Its other build inputs,
such as embedded assets and cgo sources, come from the modules' directories with
`.gitignore` applied. golangci-lint reports only on Go files, so a build input
hidden only by another ignore source cannot produce a finding. Its JSON output
is decoded with Dang's built-in decoder.

A complexity finding ends with `(baseline symbol "Name")`, where `Name` is the
nearest preceding top-level declaration (`Recv.Method` for a method), so the
baseline entry can be copied from the finding. One `grep` in the golangci-lint
container reads the declaration lines of the files with complexity findings;
no file is read one at a time.

### Scanning

The scanning checks (`suppressions`, `markdown`, `comments`, `go-lint`) reach
files through the workspace, not through a copy of it. Nothing evaluates a
whole-workspace `Directory`, nothing reads files one at a time, and no check
walks the workspace in Dang.

- `suppressions` runs one `Workspace.search` for the eight directive patterns
  and turns only the matched lines into findings.
- `markdown` takes its candidates from one listing limited to markdown files.
- `comments` and `go-lint` upload only the files their tool needs, built from
  the same search that lists the files ripgrep sees. Checks that need that
  file list share one search.
- `scanFiles` and `goLintModules` stay public. They list the scanned files and
  the Go modules among them, and leave ignored files out.

Files that are left out:

- Files matched by `.gitignore`, `.ignore` or `.rgignore`, nested ones
  included, and by `.git/info/exclude`.
- Files matched by the global git excludes of the machine that runs
  `dagger` (`core.excludesFile`, or `$XDG_CONFIG_HOME/git/ignore`, which
  defaults to `~/.config/git/ignore`).
- `.git`, `vendor`, `node_modules` and Weeds-owned skill folders.

Hidden files that no ignore file matches are scanned, and so are the ignore
files themselves. A directive, a comment or internal markdown in an ignored
file is not reported, so an ignore file decides what Weeds checks. An ignore
file is committed and shows up in review, and a committed `.ignore` or
`.rgignore` can hide a directive from Weeds. Weeds has no rule about ignore
files.

How the ignore files reach the workspace:

- `Workspace.search` runs ripgrep on the machine that runs `dagger`, in your
  checkout, so it applies that machine's git configuration. A checkout with a
  `.git` directory gets `.gitignore` as well as `.ignore` and `.rgignore`.
  `.git/info/exclude` is honoured with or without one.
- `Workspace.directory(gitignore: true)` reads `.gitignore` files only. It
  does not read `.ignore`, `.rgignore`, `.git/info/exclude` or global
  excludes. Weeds therefore builds each upload from the search's file list, so
  an upload holds the files a search sees and no others. `go-lint` adds the
  other files its modules need to build (embedded assets, empty or binary
  test data) with `.gitignore` applied. golangci-lint reports on Go files
  only, so a build input hidden only by another ignore source cannot produce a
  finding.
- The global git excludes belong to the machine, not to the repository. Your
  laptop and a CI runner can differ, and CI decides. An empty markdown file
  hidden only by your personal global excludes is still reported, because
  ripgrep never lists an empty file and the in-engine check that closes this
  gap for repository ignore sources cannot read them.
- Ripgrep must be on the `PATH` (see Install). Dagger's `grep` fallback
  honours no ignore file, so Weeds raises an error naming ripgrep instead of
  scanning ignored files.

The skips apply to the scanning checks only. The installer and
`installed-skills` read the skills directories, their agent directories and
`AGENTS.md` and `CLAUDE.md`, including owned skill folders, their markers and
any hand edit inside them, whether or not an ignore file mentions them.
`tests/fixtures/ignore-probe` records how each source behaves, with a host run
that shows the global excludes.

### Comments

`comments` uploads only the Go and TS/JS files that no ignore file leaves out
and runs a pinned ast-grep (0.39.5) in a container once per Weeds rule pack
(`rules/go`, `rules/typescript` and the TSX pack), reading its `--json` output.
The packs report violations only: a doc comment that complies produces no
match, and no file is read to judge one. The consumer's ast-grep configuration
is never read, and ast-grep's own inline suppression comment is neutralised in
the scanned copy, so that comment is reported like any other. `.go` uses the Go
grammar, `.ts` the TypeScript grammar, and `.tsx`, `.js`, `.jsx`, `.mjs` and
`.cjs` the TSX grammar (the TSX pack is derived from the TypeScript pack at run
time).

The rules cannot see line breaks in a syntax tree, so the scanned copy gets a
`//weeds-nl` comment line before every line but the first. Two adjacent
markers are a blank line, and a comment whose previous sibling is not a marker
has code before it on its line. Generated files and licence headers are also
rules (`go-generated-file`, `go-licence-header` and the `ts-` pair), so their
facts come from the same scan. A licence header is a leading comment run that
holds a licence or copyright line; a block comment counts whole.

Allowed: a doc comment directly above an exported identifier in an importable
package (Go: any comment, including the package comment and fields or methods
of exported types; TS/JS: a `/** */` block above an `export`, a CommonJS
`module.exports`/`exports.x` assignment, or a public member of an exported
class, interface or enum); tool directives (Go: `//go:`, `// +build`, `//line`;
TS/JS: triple-slash references, `@ts-check`, `@jsx*`, `@flow` and source-map
comments); generated files; and a licence header, the leading comment block
that has a licence or copyright line. Everything else is reported:

| Rule | Reports |
|---|---|
| `go-comment-todo`, `ts-comment-todo` | `TODO` or `FIXME` anywhere. |
| `go-comment-main` | Any comment in `package main`. |
| `go-comment-internal`, `ts-comment-internal` | Any comment under `internal/`. |
| `go-comment-code`, `ts-comment-code` | A `//` comment that looks like code, outside function bodies. |
| `go-comment-not-doc`, `ts-comment-not-doc` | Comments in function bodies, on unexported code, or otherwise not an adjacent doc comment. |
| `go-comment-detached`, `ts-comment-detached` | A trailing comment, or a doc comment separated from its declaration by a blank line. |

Where forms overlap, in order: a generated file gets no findings; `TODO` and
`FIXME` are always reported, even in a licence header or an exported doc
comment; a directive or licence header is allowed anywhere (function bodies,
unexported code, `internal/`, `main`); a directive line between a doc comment
and its declaration does not detach it; any other comment must be an adjacent
doc comment on an exported identifier.

### installed-skills and `dagger generate`

`dagger generate` copies the selected skills, from the kpenfound/skills commit
this module pins by full SHA, into each target. `installed-skills` runs the
same code and fails when the installed folders differ from what the installer
would write, with one finding per path: `skill-added` (the installer would add
it), `skill-modified` (a hand edit), `skill-removed` (it would delete it), and
`skill-conflict` (a folder the installer would write exists without the
ownership marker, so `dagger generate` would abort; the finding names the
folder).

Targets are the real directories (not symlinks) among `.agents/skills`,
`.claude/skills`, `.cursor/skills`, `.gemini/skills`, `.github/skills` and
`.opencode/skills` at the repo root. With none, `AGENTS.md` creates
`.agents/skills` and `CLAUDE.md` creates `.claude/skills`. With neither,
`dagger generate` fails, and `installed-skills` passes with no findings because
there is nothing to install. A symlink or file is never replaced by a target.

Each installed folder gets the ownership marker (see below). If a folder to be
written exists in any target without a valid marker, the whole run aborts with
no change. Owned folders are replaced whole, and owned folders of skills that
are no longer selected are removed. Every selected skill needs
`<name>/SKILL.md` whose frontmatter `name` equals the folder name. The skill
source is a module input (`skillSource`); tests pass fixture skill trees, so
they need no network.

## weeds.toml

`weeds.toml` sits at the repository root. A missing file means defaults: no
exceptions, an empty markdown allowlist, an empty baseline and all seven
skills. An unknown key, an unknown check or rule name or a malformed entry
fails with a message that starts `weeds.toml is invalid:` and names the
problem. Every problem in the file is reported together.

An exception, allowlist entry or baseline entry that matches nothing is stale,
and a stale entry fails the check. Exceptions exist only here, by name. No
inline directive is ever honoured as an exception.

### Exceptions

```toml
[[exceptions]]
check = "suppressions"
paths = ["tests/fixtures/**"]
reason = "fixtures hold real directives"
```

An exception names exactly one `check` or one `rule`, plus `paths`: a
non-empty list of globs. `*` matches within a path segment, `**` matches any
number of segments and `?` matches one character. A path that names a
directory covers everything under it. The optional `reason` is for readers.

### Markdown allowlist

```toml
[markdown]
allowlist = ["notes/design.md"]
```

### Go lint baseline

```toml
[[go-lint.baseline]]
linter = "gocognit"
path = "internal/store/store.go"
symbol = "Store.Load"
reason = "split Load before the next release"
```

The baseline may hold complexity findings only (`funlen`, `gocognit`, `nestif`
and `dupl`) and it may only shrink. An entry for any other linter is a config
error. An entry matches a finding by linter, path and enclosing symbol, never
by line number, so it survives edits above the function. Leave `symbol` out to
match every finding of that linter in the file. A matching entry removes its
finding, and an entry that matches nothing is stale and fails the check.
Exceptions can name any go-lint linter as a `rule`.

### Skills subset

```toml
[skills]
subset = ["no-deltas", "behavioral-testing"]
```

The default is all seven skills: `no-comments`, `lint-discipline`,
`simple-code`, `typescript-best-practices`, `unslop`, `no-deltas` and
`behavioral-testing`. A subset must be non-empty and drawn from those seven.

## Skill folders and ownership

An installed skill folder carries an ownership marker: a `.weeds-managed` file
whose content is exactly `github.com/kpenfound/weeds` and a newline. A folder
with a valid marker is Weeds-owned. The scans skip it, the `markdown` check
allows it, and the installer replaces it whole. A skills-looking folder without
the marker is not covered, and the installer never overwrites or deletes it.

## Language coverage

The comment rules cover Go and TS/JS only. Other languages (Dang, TOML, YAML,
shell and the rest) get the suppression and markdown scans, and agents apply
the comment policy there by hand.

## Planned

Semantic checks (jint) are being considered. Nothing about them is promised,
and `weeds.toml` reserves no keys for them.

## Development

`dagger check` at the root runs the Weeds checks on this repository and every
`weeds-tests` suite. Test cases live under `tests/fixtures/<suite>/<case>`:
a `workspace` directory plus `expected.txt` (compared exactly after trimming)
or `expected-error.txt` (fragments the failure must contain, one per line, a
leading `!` for a fragment it must not contain). A `findings.txt` feeds finding
lines to the shared exception evaluator.

Fixture rules:

- Name an ignore file `gitignore.fixture`, `ignore.fixture` or
  `rgignore.fixture`, at the workspace root or in any directory. The harness
  renames it to `.gitignore`, `.ignore` or `.rgignore` beside it when the case
  runs, and `git-info-exclude.fixture` at the workspace root to
  `.git/info/exclude`. It also adds a `.git/HEAD` marker, because ripgrep
  applies `.gitignore` only inside a git repository. A fixture never holds a
  live ignore file or a `.git` directory, which would change this repository's
  own scans.
- The `ignore-files` suite is the end-to-end case for ignore handling: one
  workspace with a directive, bad comments, internal markdown and Go lint
  issues in a directory ignored by each of `.gitignore`, `.ignore`,
  `.rgignore`, nested ignore files and `.git/info/exclude`, plus violations
  that nothing ignores. The harness adds an ignored full copy of the
  workspace, with its own `go.mod`, as `worktrees/copy`. The case runs all
  five checks and the installer and asserts that no ignored file is
  reported, that `goLintModules` omits the copy, and that the search, the
  shared upload, the comments upload, the markdown listing and the go-lint
  upload agree. It also asserts that the same workspace without its ignore
  files does report the ignored files.
- Global git excludes cannot be set from inside a module, so no suite covers
  them. `tests/fixtures/ignore-probe/global-excludes.sh` runs the real
  entry points on a host workspace with a scratch `HOME` and records the
  result.
- Commit no symlinks. List them in a `symlinks.txt` beside `workspace`, one
  `link -> target` per line.
- The `generate` suite prints the installer's changeset: `+ path body` for an
  added file, `~ path body` for a modified one, `- path` for a removed one,
  with newlines in bodies written `\n`. `installed-skills` cases print the
  check's findings. `source.txt` names a skill tree under
  `tests/fixtures/skill-sources` (default `full`).
- The `consumer` suite is the end-to-end case: one workspace that uses every
  `weeds.toml` section, run through all five checks and the installer. Its
  case adds `hand-edits/` on top of the workspace as a hand-edited owned
  skill folder, and its golden `expected.txt` holds one section per check plus
  the generator's changeset.
- The `scale*` checks generate large inputs in memory and assert the exact
  output, never a time: about two thousand golangci-lint issues across two
  modules, and a workspace of about five thousand files (directives, markdown,
  Go and TS comments, and a skills target installed from a fixture skill
  source). The generators live in `tests/scale*.dang`; each expected output is
  computed from how its input was built. No large fixture is committed.
- The `golintRaises` and `golintBoundaries` checks feed golangci-lint process
  results (exit status, stdout, stderr) to `goLintReportFromResult`, the test
  entry point that runs the same validation as the go-lint check, and pin each
  raise message and trigger boundary.
- Dagger finds modules only inside a git repository, so run `dagger check`
  from a git checkout.

## License

MIT. See `LICENSE`.
