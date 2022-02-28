#! /usr/bin/env bash

currentFolderPath=$(cd "$(dirname $0)" && pwd);
source "${currentFolderPath}/commons.sh"

expectedGoVersion=$(awk '/golang/ { print $2 }' .tool-versions);
echo "expected version: ${expectedGoVersion}";
errcode=0;

goSystemVersion=$(go version);
if echo "${goSystemVersion}" | grep " go${expectedGoVersion} " &> /dev/null; then
  success "system version (${goSystemVersion})"
else
  >&2 fail "system version (${goSystemVersion})"
  errcode=1
fi

githubActionVersions=$(awk '/go-version:/ { print NR "," $2 }' .github/workflows/test.yml);
while IFS=, read -r line version; do
  if echo "${version}" | grep "\"${expectedGoVersion}\"" &> /dev/null; then
    success "github action workflow version line ${line} (${version})"
  else
    >&2 fail "github action workflow version line ${line} (${version})"
    errcode=1
  fi
done <<< "${githubActionVersions}"

if [ ${errcode} -gt 0 ]; then
  exit ${errcode};
fi
