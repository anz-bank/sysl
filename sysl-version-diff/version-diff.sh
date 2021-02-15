#!/bin/sh

# this script is designed to work in a docker. It takes a version diff argument,
# a filename and another optional filename for the rhs version if a different
# filename is used in rhs version.
set -e

############################# handle arguments #################################

# allow outputting diff report to a file
while getopts ":o:" opt; do
    case "$opt" in
        o)
            OUTPUT=$OPTARG
            ;;
        \?)
            echo "Invalid option: -$OPTARG"
            exit 1
            ;;
    esac
done

if [ ! -z "$OUTPUT" ]; then
    OUTPUT_FLAG="--out=file:$OUTPUT"
fi

shift `expr $OPTIND - 1` || :

VERSIONS=$1
LHS_FILENAME=$2
RHS_FILENAME=$3

if [ -z "$VERSIONS" ] || [ -z "$LHS_FILENAME" ]; then
    echo "Usage: VERSION_LHS..VERSION_RHS FILENAME [FILENAME_RHS]"
    exit 1
fi

# If rhs filename is not provided, it defaults to lhs filename. This means in
# both versions, the filenames are different.
if [ -z "$RHS_FILENAME" ]; then
    RHS_FILENAME=$LHS_FILENAME
fi

VERSION_SEP=".."

# removes version separator and everything after it
LHS_VERSION=${VERSIONS%$VERSION_SEP*}
# removes version separator and everything before it
RHS_VERSION=${VERSIONS#*$VERSION_SEP}

if [ -z "$LHS_VERSION" ] || [ "$LHS_VERSION" = "HEAD" ]; then
    LHS_VERSION="$(git rev-parse --short HEAD)"
fi

if [ -z "$RHS_VERSION" ] || [ "$RHS_VERSION" = "HEAD" ]; then
    RHS_VERSION="$(git rev-parse --short HEAD)"
fi

################################# create diff ##################################

GENDIR=/gen
WORKDIR=/workdir
REPODIR="/opt/repository"
DIFF_SCRIPT="/scripts/version_diff.arraiz"
LHS_PB=$GENDIR/lhs.pb
RHS_PB=$GENDIR/rhs.pb

clean_dir() {
    DIR=$1
    mkdir -p $DIR
    if [ ! -z "$(ls -A $DIR)" ]; then
        rm -rfv $DIR > /dev/null 2>&1
    fi
}

# clean up all the important directories
clean_dir $GENDIR
clean_dir $REPODIR

# copy the repository and clean it
cd $WORKDIR && cp -R ./ $REPODIR/
cd $REPODIR && git reset --hard HEAD > /dev/null 2>&1

# checkout to a version and generate a pb file
gen_pb_version() {
    OUTPUT=$1
    VERSION=$2
    FILENAME=$3

    git checkout $VERSION > /dev/null 2>&1 || {
        echo "unable to checkout to version $VERSION"
        exit 1
    }

    if [ ! -f "$FILENAME" ]; then
        echo "$FILENAME does not exist in version $VERSION"
        exit 1
    fi

    sysl pb --mode pb --output $OUTPUT $FILENAME
}

# generate the two necessary files for diff
cd $REPODIR
gen_pb_version $LHS_PB $LHS_VERSION $LHS_FILENAME
gen_pb_version $RHS_PB $RHS_VERSION $RHS_FILENAME

# create diff and show result in stdout
cd $WORKDIR
arrai run $OUTPUT_FLAG $DIFF_SCRIPT $LHS_PB $RHS_PB

# clean up generated files
clean_dir $REPODIR
clean_dir $GENDIR

################################################################################
