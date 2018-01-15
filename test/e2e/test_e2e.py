from sysl.core import syslloader
from sysl.util.file import filesAreIdentical

import pytest
import os

from os import path, remove
from subprocess import call

REPO_ROOT = path.normpath(path.join(path.dirname(__file__), '..', '..'))


def remove_file(fname):
    try:
        remove(fname)
    except OSError:
        pass


@pytest.mark.parametrize("fname", [
    '001_annotations',
    '002_annotations',
    '003_annotations',
    '004_annotations',
    '005_annotations',
])
def test_e2e(fname):
    e2e_dir = path.normpath(path.dirname(__file__))
    e2e_rel_dir = os.path.relpath(e2e_dir, start=REPO_ROOT)

    root = path.join(e2e_dir, 'input')
    model = '/' + fname
    out_dir = path.join(REPO_ROOT, 'tmp', e2e_rel_dir)
    out_fname = path.join(out_dir, 'actual_output', fname + '.txt')
    remove_file(out_fname)

    cmd = ['sysl', '--root', root, 'textpb', '-o', out_fname, model]
    print 'calling', ' '.join(cmd)
    call(cmd)

    expected_fname = path.join(e2e_dir, 'expected_output', fname + '.txt')
    assert filesAreIdentical(expected_fname, out_fname)
