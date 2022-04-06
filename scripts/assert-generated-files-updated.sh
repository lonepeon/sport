#! /usr/bin/env bash

GO_BIN=go

currentFolderPath=$(cd "$(dirname $0)" && pwd);
source "${currentFolderPath}/commons.sh"

if [ 0 -ne $(git status --porcelain --untracked-files | wc -l) ]; then
  >&2 fail "aborting: git status is not clean";
  exit 1;
fi

$GO_BIN generate ./...;

if [ 0 -ne $(git status --porcelain --untracked-files | wc -l) ]; then
  >&2 fail "run go generate ./... to regenerate the following files:";
  git status --porcelain --untracked-files | awk '{print $2}' >&2
  exit 1;
else
  success "all generated files are up to date";
fi
