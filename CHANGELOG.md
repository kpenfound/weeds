# Changelog

## Unreleased

### Breaking

- The scanning checks need ripgrep (`rg`) on the `PATH` of the machine that
  runs `dagger`. They find files with `Workspace.search` and leave out files
  matched by ignore files; without ripgrep they fail with an error that names
  ripgrep.
- `suppressions`, `markdown`, `comments` and `go-lint` skip files matched by
  `.gitignore`, `.ignore` or `.rgignore` (nested ones included),
  `.git/info/exclude` and the global git excludes of the machine that runs
  `dagger`, on top of `.git`, `vendor`, `node_modules` and Weeds-owned skill
  folders. Hidden files that no ignore file matches are still scanned. This
  reverses the earlier rule that ignored files are scanned: ignored agent
  worktrees and their `go.mod` files are no longer scanned or linted, and
  `scanFiles` and `goLintModules` leave ignored files out. Two tradeoffs are
  accepted: a committed `.ignore` or `.rgignore` can hide a
  directive from Weeds, and each person's global git excludes can make a
  laptop result differ from CI, where CI decides. Weeds has no rule about
  ignore files.
- The skills source is repinned to the kpenfound/skills commit
  `295e8f78fe22a041b23674e96014c5d6e6778e02`, which carries the seven skills `no-comments`,
  `lint-discipline`, `simple-code`, `typescript-best-practices`, `unslop`,
  `no-deltas` and `behavioral-testing`. Changing the skills pin changes the
  installed skill folders, so the `installed-skills` check fails in a consumer
  repository until it runs `dagger generate` and commits the result.
