#!/bin/bash

USAGE="Usage: release.sh <prepare|deploy> MAJOR.MINOR.PATCH"
COMMAND=$1
VERSION=$2

function fatal {
    echo "$@" >&2
    exit 1
}
if [ "$#" -ne 2 ]; then
    fatal "$USAGE"
fi
if ! [[ $VERSION =~ ^[[:digit:]+\\.[:digit:]+\\.[:digit:]]+$ ]]; then
  fatal "$USAGE"
fi
if [[ $COMMAND != "prepare" && $COMMAND != "deploy" ]]; then
  fatal "$USAGE"
fi
if [[ $(git status --porcelain 2> /dev/null | tail -n1) != "" ]]; then
  fatal "Repo is not clean please commit or delete dirty files."
fi

ORIGIN_URL=$(git remote get-url origin)
ORIGIN=$(echo "$ORIGIN_URL" | sed -n "s/.*\\/\\/github.com\\/\\(..*\\)\\/..*.git/\\1/p")
UPSTREAM="anz-bank"
UPSTREAM_URL=$(echo "$ORIGIN_URL" | sed -n "s/\\(.*\\/\\/github.com\\/\\)\\(..*\\)\\(\\/..*.git\\)/\\1$UPSTREAM\\3/p")
RELEASE_BRANCH="release-v$VERSION"

echo "------- Checkout master ---------"
git checkout master || fatal "Cannot checkout master"

echo "------- Pull upstream ---------"
git pull "$UPSTREAM_URL" master || fatal "Cannot pull  upstream master"

echo "------- Create release branch $RELEASE_BRANCH ---------"
git co -b "$RELEASE_BRANCH" || fatal "Cannot create release branch $RELEASE_BRANCH"

echo "------- Update version ---------"
echo "__version__ = '$VERSION'" > 'src/sysl/__version__.py' || fatal "Could not override src/sysl/__version__.py"

echo "------- Commit ---------"
git commit -am "Bump version to $VERSION" || fatal "Cannot commit version update"

echo "------- Push ---------"
git push -u origin || fatal "Cannot push release branch to origin"

echo "------- Create PR ---------"
REPO="sysl"
CREDENTIALS=$(echo "host=github.com" | git credential fill)
USERNAME=$(echo "$CREDENTIALS" | sed -n "s/.*username=\\([^ ]*\\).*/\\1/p")
PASSWORD=$(echo "$CREDENTIALS" | sed -n "s/.*password=\\([^ ]*\\).*/\\1/p")
if [ "$USERNAME" = "" ] || [ "$PASSWORD" = "" ]; then
    fatal "Could not get GitHub credentials, please create pull request manually."
fi

JSON="{\"title\":\"Bump version to $VERSION\",\"head\":\"$ORIGIN:$RELEASE_BRANCH\",\"base\":\"master\"}"
RESPONSE=$(wget --quiet --output-document=- --content-on-error \
               --user="$USERNAME" --password="$PASSWORD" --auth-no-challenge \
               --header="Content-Type: application/json" \
               --header="Accept: application/vnd.github.v3+json" \
               --post-data="$JSON" "https://api.github.com/repos/$UPSTREAM/$REPO/pulls")

WGET_STATUS=$?
if [ $WGET_STATUS -eq 0 ]; then
    GITHUB_PR_URL=$(echo "$RESPONSE" | jq -r '.html_url')
    echo "Pull request opened:"
    echo "$GITHUB_PR_URL"
    open "$GITHUB_PR_URL"
elif [ $WGET_STATUS -eq 6 ]; then
    fatal "Wrong username or password/token"
else
    fatal "Unknown error"
fi
