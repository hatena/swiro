#!/bin/bash
set -e

USAGE="Usage: $0 <RELEASE_VERSION>"

if [ -z "$1" ]; then
  echo "$USAGE"
  exit 1
fi

RELEASE_VERSION=$(echo $1 | sed 's/^\([0-9]\)/v\1/')

echo "STEP1: tag the version"
if [[ $(git tag) =~ "${RELEASE_VERSION}" ]]; then
  echo "${RELEASE_VERSION} is already existed. You must delete this tag and try again."
  exit 1
fi
BODY=$(ghch --format=markdown --next-version=${RELEASE_VERSION})
git tag ${RELEASE_VERSION}

echo "STEP2: build"
make cross-build

echo "STEP3: push binaries"
ghr -u hatena -b "${BODY}" ${RELEASE_VERSION} dist

echo "Completed"
