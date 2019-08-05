#!/bin/bash

set -e
ROOT="sysl2/sysl/tests"
for f in $ROOT/*.sysl; do
 f=`basename $f`
 $GOPATH/bin/sysl pb --mode textpb --root $ROOT -o $ROOT/$f.out.txt /$f
done;

rm $ROOT/*.out.txt

$GOPATH/bin/sysl sd -a 'Project' $ROOT/sequence_diagram_project.sysl
rm _.png

$GOPATH/bin/sysl sd -s 'WebFrontend <- RequestProfile' -o sd.png $ROOT/sequence_diagram_project.sysl
rm sd.png

$GOPATH/bin/sysl ints -j 'Project' $ROOT/integration_test.sysl
rm _.png
