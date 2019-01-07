from sysl.core.__main__ import main
from sysl.reljam.reljam import main as reljam

from sysl.util.file import filesAreIdentical

import pytest

from os import path, remove, listdir
from subprocess import check_call
import shutil

REPO_ROOT = path.normpath(path.join(path.dirname(__file__), '..', '..'))
IN_DIR = path.join(path.normpath(path.dirname(__file__)), 'input')
EXPECTED_DIR = path.join(path.normpath(path.dirname(__file__)), 'expected_output')
ACTUAL_DIR = path.join(REPO_ROOT, 'tmp', 'e2e_actual_output')


def run(program, f, exe, args):
    if exe:
        print program + ' exe call'
        check_call([exe] + args)
    else:
        print program + ' python call'
        f(args)


def remove_file(fname):
    try:
        remove(fname)
    except OSError:
        pass


@pytest.mark.unit
@pytest.mark.parametrize("name", [
    ('000_annotations'),
    ('001_annotations'),
    ('002_annotations'),
    ('003_annotations'),
    ('004_annotations'),
])
def test_e2e(name, syslexe):
    actual = path.join(ACTUAL_DIR, name + '.txt')
    remove_file(actual)

    args = ['--root', IN_DIR, 'textpb', '-o', actual, name + '.sysl']

    run('Sysl', main, syslexe, args)

    expected = path.join(EXPECTED_DIR, name + '.txt')
    assert filesAreIdentical(expected, actual)


@pytest.mark.unit
@pytest.mark.parametrize("mode, module, app, java_pkg, expected", [
    ('model', '/010_reljam_tuplecomplex', 'UserFormComplex', 'io/sysl/reljam/gen/tuple/complex/', '010_reljam'),
    ('model', '/011_reljam_relationalmodel', 'UserModel', 'io/sysl/model/', '011_reljam'),
    ('facade', '/011_reljam_relationalmodel', 'UserFacade', 'io/sysl/facade/', '012_reljam'),
])
def test_reljam(mode, module, app, java_pkg, expected, reljamexe):
    expected_file = path.join(REPO_ROOT, 'test/e2e/expected_output', expected + '.txt')
    out_dir = path.join(ACTUAL_DIR, 'reljam')
    shutil.rmtree(out_dir, ignore_errors=True)

    args = ["--root", IN_DIR, "--out", out_dir, mode, module, app]
    run('Reljam', reljam, reljamexe, args)

    with open(expected_file) as f:
        expected = frozenset(f.read().splitlines())
    actual = frozenset(listdir(path.join(out_dir, java_pkg)))
    assert expected == actual


@pytest.mark.unit
def test_sysl_seq_diagramm(syslexe):
    name = '020_diagram'
    actual_pattern = path.join(ACTUAL_DIR, name + '-%(epname).svg')
    fname = name + '-SEQ-ATM.svg'
    actual = path.join(ACTUAL_DIR, fname)
    remove_file(actual)
    args = ['--root', IN_DIR, 'sd', '-o', actual_pattern, '/' + name, '-a', 'Bank :: Sequences']

    run('Sysl', main, syslexe, args)

    with open(actual, 'r') as f:
        svg = f.read()

    assert svg.startswith('<?xml version="1.0" encoding="UTF-8" standalone="no"?>')
    assert 'SEQ-ATM: Submit Application</text>' in svg
    assert 'ATM</text>' in svg
    assert 'AccountTransactionApi</text>' in svg
    assert 'BankDatabase</text>' in svg
    assert 'GetBalance</text>' in svg
    assert 'GET /accounts/{account_number}</text>' in svg
    assert "/accounts/{account_number}/deposit</text>" in svg


@pytest.mark.unit
def test_sysl_tuple_data_diagramm(syslexe):
    name = '020_diagram'
    actual_pattern = path.join(ACTUAL_DIR, name + '-%(epname).svg')
    fname = name + '-AccountTransactionApi.svg'
    actual = path.join(ACTUAL_DIR, fname)
    remove_file(actual)
    args = ['--root', IN_DIR, 'data', '-o', actual_pattern, '/' + name, '-j', 'Bank :: Tuple Views']

    run('Sysl', main, syslexe, args)

    with open(actual, 'r') as f:
        svg = f.read()

    assert svg.startswith('<?xml version="1.0" encoding="UTF-8" standalone="no"?>')
    assert 'Account</text>' in svg
    assert 'account_number : int</text>' in svg
    assert 'account_type : string</text>' in svg
    assert 'account_status : string</text>' in svg
    assert 'account_balance : int</text>' in svg
    assert 'Transaction</text>' in svg
    assert 'transaction_id : int</text>' in svg
    assert 'from_account_number :</text>' in svg
    assert 'to_account_number :</text>' in svg


@pytest.mark.unit
def test_sysl_relational_data_diagramm(syslexe):
    name = '020_diagram'
    actual_pattern = path.join(ACTUAL_DIR, name + '-%(epname).svg')
    fname = name + '-BankModel.svg'
    actual = path.join(ACTUAL_DIR, fname)
    remove_file(actual)
    args = ['--root', IN_DIR, 'data', '-o', actual_pattern, '/' + name, '-j', 'Bank :: Relational Views']

    run('Sysl', main, syslexe, args)

    with open(actual, 'r') as f:
        svg = f.read()

    assert svg.startswith('<?xml version="1.0" encoding="UTF-8" standalone="no"?>')
    assert 'Account</text>' in svg
    assert 'account_number : int</text>' in svg
    assert 'account_type : string</text>' in svg
    assert 'account_status : string</text>' in svg
    assert 'account_balance : int</text>' in svg
    assert 'Transaction</text>' in svg
    assert 'transaction_id : int</text>' in svg
    assert 'from_account_number :</text>' in svg
    assert 'to_account_number :</text>' in svg
