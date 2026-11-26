#!/bin/sh
# Shows how .git/info/exclude and the global git excludes of the machine that
# runs dagger reach Workspace.search and uploads, on a real host workspace.
#
# Usage: global-excludes.sh PATH_TO_A_CHECKOUT_OF_THIS_REPOSITORY
#
# The run is kept apart from the repository and from the developer's settings:
# it copies the module into a scratch directory that becomes the workspace, adds
# a minimal .git directory there, and runs dagger with HOME pointing at a
# scratch home that holds the only git configuration ripgrep can read. Needs
# dagger and a ripgrep on PATH.
set -eu

checkout=$(cd "$1" && pwd)
scratch=$(mktemp -d)
trap 'rm -rf "$scratch"' EXIT

report() {
  label=$1
  home=$scratch/home-$label
  ws=$scratch/ws-$label
  mkdir -p "$home" "$ws/.git/info"
  cp "$checkout"/dagger.toml "$checkout"/dagger-module.toml "$checkout"/*.dang "$ws/"
  cp -R "$checkout/rules" "$checkout/tests" "$ws/"
  printf 'ref: refs/heads/main\n' >"$ws/.git/HEAD"
  printf 'info-excluded.dat\n' >"$ws/.git/info/exclude"
  for name in kept info-excluded global-excluded; do
    printf 'hostmarker%s\n' 1 >"$ws/$name.dat"
  done
  case $label in
    xdg-default)
      mkdir -p "$home/.config/git"
      printf 'global-excluded.dat\n' >"$home/.config/git/ignore"
      ;;
    core-excludesfile)
      printf 'global-excluded.dat\n' >"$home/excludes"
      printf '[core]\n\texcludesFile = %s/excludes\n' "$home" >"$home/.gitconfig"
      ;;
    none) ;;
  esac
  echo "== global excludes: $label"
  (cd "$ws" && HOME=$home dagger call weeds-tests host-ignore-report --pattern 'hostmarker[0-9]')
}

report none
report xdg-default
report core-excludesfile
