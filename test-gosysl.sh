#!/bin/bash

set -e
ROOT="sysl2/sysl/tests"
for f in $ROOT/*.sysl; do
 f=`basename $f`
 $GOPATH/bin/sysl pb --mode textpb --root $ROOT -o $ROOT/$f.out.txt /$f -v
done;

rm $ROOT/*.out.txt

$GOPATH/bin/sysl sd -a 'Project' $ROOT/sequence_diagram_project.sysl -v
rm _.png

$GOPATH/bin/sysl sd -s 'WebFrontend <- RequestProfile' -o sd.png $ROOT/sequence_diagram_project.sysl -v
rm sd.png

$GOPATH/bin/sysl ints -j 'Project' $ROOT/integration_test.sysl -v
rm _.png

version=`$GOPATH/bin/sysl --version 2>&1 >/dev/null`
if [[ ${version} = "unspecified" ]]; then
    echo "version is unspecified"
    exit 1
fi
echo "gosysl version is ${version}"
