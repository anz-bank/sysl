#!/bin/bash

set -e
ROOT="sysl2/sysl/tests"
for f in $ROOT/*.sysl; do
 f=`basename $f`
 $GOPATH/bin/sysl pb --mode textpb --root $ROOT -o $ROOT/$f.out.txt /$f
done;

rm $ROOT/*.out.txt

ln -s $GOPATH/bin/sysl $GOPATH/bin/sd
$GOPATH/bin/sd -a 'Project' $ROOT/sequence_diagram_project.sysl
rm _.png
$GOPATH/bin/sd -s 'WebFrontend <- RequestProfile' -o sd.png $ROOT/sequence_diagram_project.sysl
rm sd.png
rm $GOPATH/bin/sd

ln -s $GOPATH/bin/sysl $GOPATH/bin/ints
$GOPATH/bin/ints -j 'Project' $ROOT/integration_test.sysl
rm _.png
rm $GOPATH/bin/ints
