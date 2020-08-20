#!/bin/bash

# Runs the blueprint Makefile generation script and invokes make on the result.

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

BLUEPRINT_DIR=.blueprint
# TODO: Resolve remote import.
# BLUEPRINT_IMPORT='//{github.com/anz-bank/sysl/master/pkg/arrai/blueprint_make.arrai}'
BLUEPRINT_IMPORT="//{$(realpath --relative-to="$PWD" "$SCRIPT_DIR")/blueprint_make.arrai}"
MAKEFILE=${BLUEPRINT_DIR}/Makefile

>&2 echo "Working in ${BLUEPRINT_DIR}..."
mkdir -p ${BLUEPRINT_DIR}
>&2 echo "Fetching ${BLUEPRINT_IMPORT}..."
arrai run blueprint.arrai | arrai eval "${BLUEPRINT_IMPORT}" > ${MAKEFILE}

>&2 echo "Running make on ${MAKEFILE}..."
make -f ${MAKEFILE} --debug=basic
