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

## Pinning

Pin Weeds by full commit SHA in your `dagger.toml`. Never follow a branch or a
tag.

```toml
[modules.weeds]
source = "github.com/kpenfound/weeds@<full commit sha>"
```

Weeds pins its own tools the same way: ast-grep and golangci-lint by exact
version (the golangci-lint image by tag and digest), and one kpenfound/skills
commit by full SHA. A change to any of these pins is a breaking change and is
called out in the commit message and release notes.

## Install

1. Add the module to your `dagger.toml` as shown above, pinned by full commit
   SHA.
2. Give the repository somewhere to install skills: an `AGENTS.md` (skills go
   to `.agents/skills`), a `CLAUDE.md` (`.claude/skills`), or an existing
   skills directory.
3. Run `dagger generate` to install the skills and commit the result.
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

A text scan over every file kind. It reports each line that holds one of
`//nolint`, `lint:ignore`, `eslint-disable`, `@ts-expect-error`, `@ts-ignore`,
`@ts-nocheck`, `ast-grep-ignore` or `# noqa`, one finding per directive per
line, under a rule named for the directive (`nolint`, `lint-ignore`,
`eslint-disable`, `ts-expect-error`, `ts-ignore`, `ts-nocheck`,
`ast-grep-ignore`, `noqa`). Generated files are scanned like any other. Only a
named exception in `weeds.toml` clears a finding; no inline text does.

### markdown

Markdown files (`.md`, `.markdown`, `.mdx`) are reported under the rule
`internal-markdown`, except these, which are always allowed:

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

A complexity finding ends with `(baseline symbol "Name")`, where `Name` is the
nearest preceding top-level declaration (`Recv.Method` for a method), so the
baseline entry can be copied from the finding.

### Scanning

The scanning checks (`suppressions`, `markdown`, `comments`) share one walk of
the workspace. It skips `.git`, `vendor`, `node_modules` and Weeds-owned skill
folders, and it includes hidden and gitignored files, so a directive cannot
hide in an ignored file. The skips apply to the walk only. The installer and
`installed-skills` read the whole workspace, including owned skill folders,
their markers and any hand edit inside them.

### Comments

`comments` runs a pinned ast-grep (0.39.5) in a container with the Weeds rule
packs in `rules/go` and `rules/typescript` and reads its `--json` output. The
consumer's ast-grep configuration is never read, and ast-grep's own inline
suppression comment is neutralised in the scanned copy, so that comment is
reported like any other. `.go` uses the Go grammar, `.ts` the TypeScript
grammar, and `.tsx`, `.js`, `.jsx`, `.mjs` and `.cjs` the TSX grammar (the TSX
pack is derived from the TypeScript pack at run time).

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

- Name an ignore file `gitignore.fixture` at the workspace root. The harness
  renames it to `.gitignore` when the case runs.
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
- Dagger finds modules only inside a git repository, so run `dagger check`
  from a git checkout.

## License

MIT. See `LICENSE`.
