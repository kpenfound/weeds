#!/bin/sh
# Shows which markdown files the real markdown check reports on a host
# workspace, with and without global git excludes: a non-empty file matched by
# a global file pattern, a file inside a globally excluded directory, an empty
# file matched by a global pattern, and markdown in hidden and dotted
# directories.
#
# Usage: markdown-global-excludes.sh PATH_TO_A_CHECKOUT_OF_THIS_REPOSITORY
#
# Kept apart from the repository and the developer's settings like
# global-excludes.sh: the module is copied into a scratch workspace with a
# minimal .git directory, and dagger runs with HOME pointing at a scratch home
# that holds the only git configuration ripgrep reads. Needs dagger and a
# ripgrep on PATH.
set -eu

checkout=$(cd "$1" && pwd)
scratch=$(mktemp -d)
trap 'rm -rf "$scratch"' EXIT

report() {
  label=$1
  home=$scratch/home-$label
  ws=$scratch/ws-$label
  mkdir -p "$home" "$ws/.git/info" "$ws/docs.v2" "$ws/.github" "$ws/globexcl-dir"
  cp "$checkout"/dagger.toml "$checkout"/dagger-module.toml "$checkout"/*.dang "$ws/"
  cp -R "$checkout/rules" "$checkout/tests" "$ws/"
  printf 'ref: refs/heads/main\n' >"$ws/.git/HEAD"
  printf '*.dang\nrules/\ntests/\ndagger.toml\ndagger-module.toml\n' >"$ws/.git/info/exclude"
  printf 'module example.com/host\n\ngo 1.22\n' >"$ws/go.mod"
  printf '# kept\n' >"$ws/kept.md"
  printf '# global\n' >"$ws/global-excluded.md"
  printf '# in an excluded directory\n' >"$ws/globexcl-dir/inner.md"
  printf '# dotted\n' >"$ws/docs.v2/dotted.md"
  printf '# hidden\n' >"$ws/.github/hidden.md"
  : >"$ws/empty-global-excluded.md"
  case $label in
    xdg-default)
      mkdir -p "$home/.config/git"
      printf 'global-excluded.*\nempty-global-excluded.*\nglobexcl-dir/\n' >"$home/.config/git/ignore"
      ;;
    none) ;;
  esac
  echo "== global excludes: $label"
  (
    cd "$ws"
    export HOME=$home
    dagger call weeds markdown-report || true
  )
}

report none
report xdg-default
