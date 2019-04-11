#!/bin/bash

set -e
ROOT="sysl2/sysl/tests"
for f in $ROOT/*.sysl; do
 f=`basename $f`
 $GOPATH/bin/sysl -mode textpb -root $ROOT -o $ROOT/$f.out.txt /$f
done;

rm $ROOT/*.out.txt
