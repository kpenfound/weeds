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

Run the checks with `dagger check`. Install the skills with `dagger generate`.
Both fetch Weeds and the skills repo without credentials.

Weeds pins its own tools the same way: ast-grep and golangci-lint by exact
version (the golangci-lint image by tag and digest), and one kpenfound/skills
commit by full SHA. A change to any of these pins is a breaking change and is
called out in the commit message and release notes.

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
linter = "<complexity linter>"
path = "internal/store/store.go"
symbol = "Store.Load"
```

The baseline may hold complexity findings only, and it may only shrink.

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
- Dagger finds modules only inside a git repository, so run `dagger check`
  from a git checkout.

## License

MIT. See `LICENSE`.
