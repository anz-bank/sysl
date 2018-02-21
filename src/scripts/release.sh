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
UPSTREAM="anz-bank"
UPSTREAM_URL=$(echo "$ORIGIN_URL" | sed -E 's|(.*//github.com/)(.+)(/.+\.git)|\1'"$UPSTREAM"'\3|')

echo "------- Checkout master ---------"
git checkout master || fatal "Cannot checkout master"

echo "------- Pull upstream ---------"
git pull "$UPSTREAM_URL" master || fatal "Cannot pull  upstream master"

if [[ $COMMAND = "prepare" ]]; then
    RELEASE_BRANCH="release-v$VERSION"
    ORIGIN=$(echo "$ORIGIN_URL" | sed -E 's|.*//github.com/(.+)/.+\.git|\1|')
    REPO=$(echo "$ORIGIN_URL" | sed -E 's|.*//github.com/.+/(.+)\.git|\1|')
    echo "------- Create release branch $RELEASE_BRANCH ---------"
    git checkout -b "$RELEASE_BRANCH" || fatal "Cannot create release branch $RELEASE_BRANCH"

    echo "------- Update version and commit ---------"
    echo "__version__ = '$VERSION'" > 'src/sysl/__version__.py' || fatal "Could not override src/sysl/__version__.py"
    echo "syslVersion=$VERSION" > 'src/libs/java/gradle.properties' || fatal "Could not override src/sysl/__version__.py"
    git commit -am "Bump version to $VERSION" || fatal "Cannot commit version update"

    echo "------- Push ---------"
    git push -u origin "$RELEASE_BRANCH" || fatal "Cannot push release branch to origin"

    echo "------- Create PR ---------"
    CREDENTIALS=$(echo "host=github.com" | git credential fill)
    USERNAME=$(echo "$CREDENTIALS" | sed -n "s/.*username=\\([^ ]*\\).*/\\1/p")
    PASSWORD=$(echo "$CREDENTIALS" | sed -n "s/.*password=\\([^ ]*\\).*/\\1/p")
    if [ "$USERNAME" = "" ] || [ "$PASSWORD" = "" ]; then
        fatal "Could not get GitHub credentials, please create pull request manually."
    fi

    JSON="{\"title\":\"Bump version to $VERSION\",\"head\":\"$ORIGIN:$RELEASE_BRANCH\",\"base\":\"master\"}"
    RESPONSE=$(curl -X POST -s -S \
            -u "$USERNAME:$PASSWORD" \
            --header "Content-Type: application/json" \
            --header "Accept: application/vnd.github.v3+json" \
            -d "$JSON" "https://api.github.com/repos/$UPSTREAM/$REPO/pulls")

    CURL_STATUS=$?
    if [ $CURL_STATUS -eq 0 ]; then
        GITHUB_PR_URL=$(echo "$RESPONSE" | jq -r '.html_url')
        echo "Pull request opened:"
        echo "$GITHUB_PR_URL"
        open "$GITHUB_PR_URL"
    elif [ $CURL_STATUS -eq 6 ]; then
        fatal "Wrong username or password/token"
    else
        echo "$CURL_STATUS"
        echo "$RESPONSE"
        fatal "Unknown error"
    fi
elif [[ $COMMAND = "deploy" ]]; then
    PY_VERSION=$(sed -n "s/__version__ = '\\([^ ]*\\)'/\\1/p" src/sysl/__version__.py)
    RELEASE_TAG="v$VERSION"
    [[  "$VERSION" =  "$PY_VERSION" ]] || fatal "Version ($VERSION) different from __version__.py ($PY_VERSION)"
    echo "------- Tag release ---------"
    git tag  "$RELEASE_TAG" || fatal "Cannot create tag $RELEASE_TAG"
    echo "------- Push tag ---------"
    git push "$UPSTREAM_URL" "$RELEASE_TAG"
    echo "------- Deployment triggered on Travis and Appveyor ---------"
fi
