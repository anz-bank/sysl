from sysl.core.__main__ import main
from sysl.reljam.reljam import main as reljam

from sysl.util.file import filesAreIdentical

import pytest

from os import path, remove, listdir
from subprocess import call
import shutil

REPO_ROOT = path.normpath(path.join(path.dirname(__file__), '..', '..'))


def remove_file(fname):
    try:
        remove(fname)
    except OSError:
        pass


@pytest.mark.unit
@pytest.mark.parametrize("fname", [
    ('001_annotations'),
    ('002_annotations'),
    ('003_annotations'),
    ('004_annotations'),
    ('005_annotations'),
    ('001_annotations'),
    ('002_annotations'),
    ('003_annotations'),
    ('004_annotations'),
    ('005_annotations')
])
def test_e2e(fname, syslexe):
    e2e_dir = path.normpath(path.dirname(__file__))
    e2e_rel_dir = path.relpath(e2e_dir, start=REPO_ROOT)

    root = path.join(e2e_dir, 'input')
    model = '/' + fname
    out_dir = path.join(REPO_ROOT, 'tmp', e2e_rel_dir)
    out_fname = path.join(out_dir, 'actual_output', fname + '.txt')
    remove_file(out_fname)

    args = ['--root', root, 'textpb', '-o', out_fname, model]

    if syslexe:
        print 'Sysl exe call'
        call([syslexe] + args)
    else:
        print 'Sysl python function call'
        main(args)

    expected_fname = path.join(e2e_dir, 'expected_output', fname + '.txt')
    assert filesAreIdentical(expected_fname, out_fname)


@pytest.mark.unit
@pytest.mark.parametrize("mode, module, app, java_pkg, expected", [
    ('model', '/test/java/tuplecomplex', 'UserFormComplex', 'io/sysl/reljam/gen/tuple/complex/', 'test_reljam_1'),
    ('model', '/test/java/relationalmodel', 'UserModel', 'io/sysl/model/', 'test_reljam_2'),
    ('facade', '/test/java/relationalmodel', 'UserFacade', 'io/sysl/facade/', 'test_reljam_3'),
])
def test_reljam(mode, module, app, java_pkg, expected, reljamexe):
    java = 'tmp/src/main/java'
    expected_file = path.join(REPO_ROOT, 'test/e2e/expected_output', expected + '.txt')
    out_dir = path.join(REPO_ROOT, java)
    shutil.rmtree(out_dir, ignore_errors=True)

    args = ["--out", out_dir, mode, module, app]
    if reljamexe:
        print 'Reljam exe call'
        call([reljamexe] + args)
    else:
        print 'Reljam python function call'
        reljam(args)

    with open(expected_file) as f:
        expected = f.read().splitlines().sort()
    actual = listdir(path.join(out_dir, java_pkg)).sort()
    assert expected == actual
