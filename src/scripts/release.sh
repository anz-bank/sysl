#!/bin/bash

# print to stderr and exit
function exit_error {
    echo "$@" >&2
    exit 1
}
USAGE="Usage: release.sh <prepare|deploy> MAJOR.MINOR.PATCH"

if [ "$#" -ne 2 ]; then
    exit_error $USAGE
fi

VERSION=$2
if ! [[ $VERSION =~ ^[[:digit:]+\\.[:digit:]+\\.[:digit:]]+$ ]]; then
	exit_error $USAGE
fi

# TODO replace release with master
CURRENT_BRANCH=`git rev-parse --abbrev-ref HEAD`
if [[ "$CURRENT_BRANCH" != "release" ]]; then
 	 exit_error "Must be on release branch"
fi
if [[ $(git status --porcelain 2> /dev/null | tail -n1) != "" ]]; then
	exit_error "Repo is not clean please commit or delete files."
fi

git pull

BRANCH="release-$VERSION"
git co -b $BRANCH

echo "__version__ = '$VERSION'" > 'src/sysl/__version__.py' || { echo "Could not override src/sysl/__version__.py"; exit 1; }
git commit -am "Bump version to $VERSION"

CREDENTIALS=`echo "host=github.com" | git credential fill`
USERNAME=`echo $CREDENTIALS | sed -n "s/.*username=\([^ ]*\).*/\1/p"`
PASSWORD=`echo $CREDENTIALS | sed -n "s/.*password=\([^ ]*\).*/\1/p"`
if [ -z "$USERNAME" ] || [ -z "$PASSWORD" ]; then
    exit_error "Could not get GitHub credentials, please create pull request manually."
fi

ORIGIN=`git remote get-url origin | sed -n "s/.*\/\/github.com\/\(..*\)\/..*.git/\1/p"`
UPSTREAM="anz-bank"
REPO="sysl"
JSON=`jq -nc "{title:\"Bump version to $VERSION\",head:\"$ORIGIN:$BRANCH\",base:\"master\"}"`

# create a pull request through the github api
response=`wget --quiet --output-document=- --content-on-error \
               --user="$USERNAME" --password="$PASSWORD" --auth-no-challenge \
               --header="Content-Type: application/json" \
               --header="Accept: application/vnd.github.v3+json" \
               --post-data="$JSON" "https://api.github.com/repos/$UPSTREAM/$REPO/pulls"`

#open...
# git tag v$VERSION
# git push -u origin "$BRANCH" --follow-tags

# curl https://api.github.com/repos/anz-bank/sysl/pulls
# -d {"title": "Amazing new feature", "head": "juliaogris:", "base": "master"}
