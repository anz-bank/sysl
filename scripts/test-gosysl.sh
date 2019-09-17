#!/bin/bash

set -e
ROOT="sysl2/sysl/tests"
for f in $ROOT/*.sysl; do
 f=`basename $f`
 $GOPATH/bin/sysl -v pb --mode textpb --root $ROOT -o $ROOT/$f.out.txt /$f
done;

rm $ROOT/*.out.txt

$GOPATH/bin/sysl -v sd -a 'Project' $ROOT/sequence_diagram_project.sysl
rm _.png

$GOPATH/bin/sysl -v sd -s 'WebFrontend <- RequestProfile' -o sd.png $ROOT/sequence_diagram_project.sysl
rm sd.png

$GOPATH/bin/sysl -v ints -j 'Project' $ROOT/integration_test.sysl
rm _.png

version=`$GOPATH/bin/sysl --version 2>&1 >/dev/null`
if [[ ${version} = "unspecified" ]]; then
    echo "version is unspecified"
    exit 1
fi
echo "gosysl version is ${version}"

$GOPATH/bin/sysl info
