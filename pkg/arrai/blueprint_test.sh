#!/bin/sh

# Runs the blueprint Makefile generation script and invokes make on the result.

set -e

BLUEPRINT_DIR=.blueprint
MAKEFILE=${BLUEPRINT_DIR}/Makefile

mkdir -p ${BLUEPRINT_DIR}
arrai run blueprint.arrai \
    | arrai eval '//{github.com/anz-bank/sysl/pkg/arrai/blueprint_make.arrai}' \
    > ${MAKEFILE}

>&2 echo "Running make on ${MAKEFILE}..."
make -f ${MAKEFILE} --debug=basic
