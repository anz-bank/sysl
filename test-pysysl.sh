#!/bin/bash

set -e

SYSL_PYTHON_DIR="$1"

flake8
pytest test/e2e --syslexe=${SYSL_PYTHON_DIR}sysl --reljamexe=${SYSL_PYTHON_DIR}reljam
coverage run --source=src/sysl -m py.test
