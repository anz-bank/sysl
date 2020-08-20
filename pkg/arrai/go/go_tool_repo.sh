#!/bin/bash

set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
ORG=$1
NAME=$2
DESCRIPTION=$3

mkdir -p $NAME
cd $DIR && \
    arrai run --out=dir:$NAME go_tool_repo.arrai $@ && \
    cd - \
    mv $DIR/$NAME ./
