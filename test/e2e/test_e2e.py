from sysl.core.__main__ import main
from sysl.reljam.reljam import main as reljam

from sysl.util.file import filesAreIdentical

import pytest

from os import path, remove, listdir
from subprocess import call

REPO_ROOT = path.normpath(path.join(path.dirname(__file__), '..', '..'))


def remove_file(fname):
    try:
        remove(fname)
    except OSError:
        pass


@pytest.mark.parametrize("fname, subprocess", [
    ('001_annotations', True),
    ('002_annotations', True),
    ('003_annotations', True),
    ('004_annotations', True),
    ('005_annotations', True),
    ('001_annotations', False),
    ('002_annotations', False),
    ('003_annotations', False),
    ('004_annotations', False),
    ('005_annotations', False)
])
def test_e2e(fname, subprocess):
    e2e_dir = path.normpath(path.dirname(__file__))
    e2e_rel_dir = path.relpath(e2e_dir, start=REPO_ROOT)

    root = path.join(e2e_dir, 'input')
    model = '/' + fname
    out_dir = path.join(REPO_ROOT, 'tmp', e2e_rel_dir)
    out_fname = path.join(out_dir, 'actual_output', fname + '.txt')
    remove_file(out_fname)

    args = ['--root', root, 'textpb', '-o', out_fname, model]

    if subprocess:
        cmd = ['sysl'] + args
        print 'subprocess call: ', ' '.join(cmd)
        call(cmd)
    else:
        print 'python function call: main([{}])'.format(', '.join(args))
        main(args)

    expected_fname = path.join(e2e_dir, 'expected_output', fname + '.txt')
    assert filesAreIdentical(expected_fname, out_fname)


@pytest.mark.parametrize("mode, module, app, output, expected", [
    ('model', '/test/java/tuplecomplex', 'UserFormComplex', 'io/sysl/reljam/gen/tuple/complex/', 'test_reljam_1'),
    ('model', '/test/java/relationalmodel', 'UserModel', 'io/sysl/model/', 'test_reljam_2'),
    ('facade', '/test/java/relationalmodel', 'UserFacade', 'io/sysl/facade/', 'test_reljam_3'),
])
def test_reljam(mode, module, app, output, expected):
    java = 'tmp/src/main/java'
    expected_file = path.join(REPO_ROOT, 'test/e2e/expected_output', expected + '.txt')
    out_dir = path.join(REPO_ROOT, java)
    output = path.join(REPO_ROOT, java, output)

    args = ["--out", out_dir, mode, module, app]
    reljam(args)
    with open(expected_file) as f:
        expected = f.read().splitlines().sort()
    out = listdir(output).sort()
    assert expected == out
