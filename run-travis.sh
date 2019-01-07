#!/bin/bash

set -e -x

BUILDID=sysl-travis-1

docker rm -f $BUILDID || true
docker build -f travis.Dockerfile -t $BUILDID .
docker run --name $BUILDID -dit -v `pwd`/sysl2:/go/src/github.com/anz-bank/sysl/sysl2 $BUILDID /sbin/init
docker exec -it $BUILDID ./build.sh
docker rm -f $BUILDID
