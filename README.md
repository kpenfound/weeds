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

### Scanning

The scanning checks (`suppressions`, `markdown`, `comments`) share one walk of
the workspace. It skips `.git`, `vendor`, `node_modules` and Weeds-owned skill
folders, and it includes hidden and gitignored files, so a directive cannot
hide in an ignored file. The skips apply to the walk only. The installer and
`installed-skills` read the whole workspace, including owned skill folders,
their markers and any hand edit inside them.

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
