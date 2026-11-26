# Changelog

## Unreleased

### Breaking

- The scanning checks need ripgrep (`rg`) on the `PATH` of the machine that
  runs `dagger`. They find files with `Workspace.search` and leave out files
  matched by ignore files; without ripgrep they fail with an error that names
  ripgrep.
- The skills source is repinned to the kpenfound/skills commit
  `295e8f78fe22a041b23674e96014c5d6e6778e02`, which carries the seven skills `no-comments`,
  `lint-discipline`, `simple-code`, `typescript-best-practices`, `unslop`,
  `no-deltas` and `behavioral-testing`. Changing the skills pin changes the
  installed skill folders, so the `installed-skills` check fails in a consumer
  repository until it runs `dagger generate` and commits the result.
