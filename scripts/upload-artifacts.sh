#!/bin/bash

set -e

# Ensure that the GITHUB_TOKEN secret is included
if [[ -z "$GITHUB_TOKEN" ]]; then
  echo "Set the GITHUB_TOKEN env variable."
  exit 1
fi

# Ensure that the file path is present
if [[ "$#" -lt 1 ]]; then
  echo "Missing file patterns to upload."
  exit 1
fi

AUTH_HEADER="Authorization: token ${GITHUB_TOKEN}"
GIT_TAG="$(echo ${GITHUB_REF} | cut -d'/' -f3)"

# Check whether a release exists with given tag
HTTP_STATUS=$(curl -sS \
     -X GET \
     -o /dev/null -w '%{http_code}' \
     -H "${AUTH_HEADER}" \
     "https://api.github.com/repos/${GITHUB_REPOSITORY}/releases/tags/${GIT_TAG}")

if [[ ${HTTP_STATUS} -eq 200 ]]; then
    echo "A release already exists with tag: ${GIT_TAG}"
    exit 1
fi

JSON=$(cat <<-END
    {
        "tag_name": "${GIT_TAG}",
        "name": "Sysl ${GIT_TAG}",
        "body": "Sysl ${GIT_TAG}",
        "draft": true,
        "prerelease": true
    }
END
)

# Draft the release
RESPONSE=$(curl -sS \
     -X POST \
     -H "${AUTH_HEADER}" \
     -H "Content-Type: application/json" \
     -H "Accept: application/vnd.github.v3+json" \
     -w "HTTP_STATUS:%{http_code}" \
     -d "$JSON" "https://api.github.com/repos/${GITHUB_REPOSITORY}/releases")

HTTP_STATUS=$(echo "$RESPONSE" | tr -d '\n' | sed -e 's/.*HTTP_STATUS://')
if [[ ${HTTP_STATUS} -ne 201 ]]; then
    echo "Failed to create release."
    exit 1
fi

RELEASE_ID=$(echo "$RESPONSE" | sed -e 's/HTTP_STATUS:.*//g' | jq -r '.id')

# For each matching file
for file in $@; do
    FILENAME=$(basename ${file})
    UPLOAD_URL="https://uploads.github.com/repos/${GITHUB_REPOSITORY}/releases/${RELEASE_ID}/assets?name=${FILENAME}"
    echo "$UPLOAD_URL"

    # Upload the file
    curl \
        -sSL \
        -XPOST \
        -H "${AUTH_HEADER}" \
        --upload-file "${file}" \
        --header "Content-Type:application/octet-stream" \
        "${UPLOAD_URL}"
done
