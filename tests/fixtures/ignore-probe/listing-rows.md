# Repeated listing rows on Dagger Cloud

## What the failed checks show

On the Dagger Cloud run of the assembled branch, the shared listing returned
several rows per file: the file-layer check got `kept.go` three times, the
ignore-files check got `app/app.go` nine times and `go.mod` three times from
`scanFiles`, and `scale-suppressions` listed 35843 files where 4912 exist.
go-lint reported each errcheck finding three times because every `go.mod`
candidate repeated. The same checks pass against a local engine.

The row counts match the number of lines in each file. That is the output of a
line-oriented ripgrep: without `--multiline`, `\A` matches at the start of every
line.

## Evidence

- Dagger `core/search.go` builds the ripgrep command from the search
  arguments: `multiline` becomes `--multiline`, `filesOnly` becomes
  `--files-with-matches` (and drops `--json`), and `limit` is applied while
  parsing results, so it caps the total and not the rows per file.
- Reproduced with ripgrep 14.1.1 on a three-line file and the listing pattern
  `(?s-u:\A.)`: `--json --multiline` gives 1 match record; `--json` without
  `--multiline` gives 3. `--files-with-matches` prints the file once with or
  without `--multiline`.
- The module pins engine `v1.0.0-beta.27`. Which engine and client build the
  Cloud run used was not observable from the traces, so the cause on Cloud is
  inferred from the row counts: the search ran without multiline mode.
- Without multiline mode `--files-with-matches` also lists a file that holds a
  NUL byte after its first line, which multiline mode does not. Such a file
  holds no text for any check, so listing it is harmless.

## The change

- The listing search sets `filesOnly: true` as well as `multiline: true`. A
  files-only search returns one row per matching file on every backend, so the
  result count follows the number of files and not of lines.
- `listingPaths` also makes each path unique, whatever the backend returns: it
  strips a leading `./`, sorts, and drops each path equal to its neighbour. It
  is one sort and one pass, with no lookup per element.
- Every consumer of `listed` (`scanFiles`, `goLintModules`, the shared upload,
  the empty-file step, markdown and the go-lint module list) therefore gets each
  file once. A go-lint module is a `go.mod` path of that list, so each module and
  each finding appears once.

`weeds-tests/listing-repeats` feeds `listingPaths` rows that repeat a path once
per line, with and without a leading `./`, and needs no Cloud backend.
