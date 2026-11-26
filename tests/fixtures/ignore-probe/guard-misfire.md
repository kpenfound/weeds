# The ripgrep guard on Dagger Cloud

## What the failed checks show

On the pull request's Dagger Cloud run, 20 of 44 checks failed with
"weeds needs ripgrep". Every failure was on a path through `FileAccess.listed`:
`scanFiles`, `listedFiles`, `uploadedFiles`, `comments`, `go-lint` and
`markdown`, including the directory-backed fixture workspaces of weeds-tests,
where ripgrep runs inside the engine and is present. The checks that passed
include `weeds/suppressions`, `weeds-tests/suppressions`, `walk` and
`ignore-probes`. All of them go through the same `hostRipgrep` probe. So ripgrep
was present, and the failing step was a judgement of the listing's own
results.

Before the fix the listing searched for the zero-width pattern `\A` and judged
each result: a result with no submatch meant "grep fallback". A zero-width match
reports one all-zero submatch (empty text, offsets 0 and 0). The same
`ignore-probes` check passed on that run and asserts one submatch per result for
`\A`, so a fresh ripgrep search does return it. The failing shape is therefore a
listing result whose submatch was lost on the way, for example when a cached
result is stored without an all-zero submatch. That fits the one run that
flipped (`scale-suppressions` passed once and failed once at `scanFiles`) and
the checks that search with a consuming pattern passing. The failed traces carry
only the error text, so the shape itself was not observed.

## The change

- The ripgrep decision is made on the probe only (`.`, limit 1), a pattern
  whose every match consumes text. No search result is judged by its own
  submatches, so no result shape from a ripgrep backend can raise.
- The listing's pattern is `(?s-u:\A.)`, the first byte of the text, so its
  match is never zero-width. It still returns one result per non-empty file,
  including a file that starts with a newline.
- A grep fallback still raises before any policy-applying search returns,
  because `searchable` runs the probe first and the probe's grep results have no
  submatches.

`weeds-tests/file-layer` tests the decision (results with submatches pass,
results without fail, no results pass) and that a listing result without
submatches is listed.
