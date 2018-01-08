from sysl.core import syslloader
from sysl.util.file import areFilesIdentical

import pytest

from os import path, remove
from subprocess import call


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
	dname = path.dirname(__file__)
	root = path.join(dname, 'input')
	model = '/' + fname
	out_fname = path.join(dname, 'actual_output', fname + '.txt')
	remove_file(out_fname)

	cmd = ['sysl', '--root', root, 'textpb', '-o', out_fname, model ]
	print 'calling', ' '.join(cmd)
	call(cmd)

	expected_fname = path.join(dname, 'expected_output', fname + '.txt')
	assert areFilesIdentical(expected_fname, out_fname)
