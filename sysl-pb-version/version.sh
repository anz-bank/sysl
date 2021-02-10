#!/bin/sh
# the version.sh runs the sysl-pb command on a specific version based on the
# version configuration in a sysl repository. It is designed to run in a docker
# image with specific directory structure.

set -e

while getopts ":m:o:" opt; do
    case "$opt" in
        m)
            MODE=$OPTARG
            ;;
        o)
            OUTPUT=$OPTARG
            ;;
        \?)
            echo "Invalid option: -$OPTARG"
            exit 1
            ;;
    esac
done

if [ ! -z "$MODE" ]; then
    MODE_FLAG="--mode $MODE"
fi

if [ ! -z "$OUTPUT" ]; then
    OUTPUT_FLAG="--output $OUTPUT"
fi

shift `expr $OPTIND - 1` || :

APPNAME=$1
VERSION=$2
REPO_DIR="/opt/repository"

if [ ! -d $REPO_DIR ]; then
    mkdir $REPO_DIR
fi

if [ "$(ls -A $REPO_DIR)" ]; then
    rm -rfv $REPO_DIR/* > /dev/null 2>&1
fi

cp -R ./ $REPO_DIR/ && cd $REPO_DIR

git reset --hard HEAD > /dev/null 2>&1

TAG=$(arrai run /scripts/config.arraiz tag $APPNAME $VERSION)
git checkout $TAG > /dev/null 2>&1

FILENAME=$REPO_DIR/$(arrai run /scripts/config.arraiz file $APPNAME $VERSION)

cd /work
sysl pb $MODE_FLAG $OUTPUT_FLAG $FILENAME

rm -rfv $REPO_DIR/* > /dev/null 2>&1
