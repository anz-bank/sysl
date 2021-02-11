#!/bin/sh

# this script is designed to work in a docker image using the sysl-pb-version image as its base.
# gen.sh is specifically made so that the generation script can be reused.
set -e

SYSL_PB_VERSION=/scripts/version.sh

APPNAME=$1
LEFT_TAG=$2
RIGHT_TAG=$3

if [ -z "$APPNAME" ] || [ -z "$LEFT_TAG" ] || [ -z "$RIGHT_TAG" ]; then
    echo "Usage: [APPNAME] [VERSION1] [VERSION2]"
    exit 1
fi

GEN_FOLDER=/gen
mkdir -p $GEN_FOLDER

# generate the two necessary files for diff
cd /work
$SYSL_PB_VERSION -m pb -o $GEN_FOLDER/lhs.pb $APPNAME $LEFT_TAG
$SYSL_PB_VERSION -m pb -o $GEN_FOLDER/rhs.pb $APPNAME $RIGHT_TAG

# create diff and show result in stdout
cd /work
arrai run /scripts/version_diff.arraiz $APPNAME $GEN_FOLDER/lhs.pb $GEN_FOLDER/rhs.pb

# clean up generated files
rm -rfv $GEN_FOLDER > /dev/null 2>&1
