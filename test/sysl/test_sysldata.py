from sysl.core import syslloader, sysldata
import unittest
import re
import os
import sys
from os import path
import traceback
import tempfile
import argparse as ap


class TestSetOf(unittest.TestCase):
    def setUp(self):
        self.outpath = tempfile.gettempdir()

    def test_set_of(self):

        try:
            (module, _, _) = syslloader.load('/test/data/test_data', True, '.')

            d = {
                'project': 'TestData :: Data Views',
                'output': path.join(self.outpath, 'test_set_of-data.png'),
                'plantuml': '',
                'verbose': '',
                'filter': ''}
            args = ap.Namespace(**d)

            out = sysldata.dataviews(module, args)

            setof_re = re.compile(r'_\d+\s+\*-- "0\.\.\*"\s+_\d+')

            self.assertTrue(setof_re.search(out[0]))

        except (IOError, Exception) as e:
            self.fail(traceback.format_exc())

    def test_at_prefixed_attr(self):

        try:
            (module, _, _) = syslloader.load(
                '/test/data/test_at_prefixed_attr', True, '.')

            #import pdb; pdb.set_trace()
            val_set = set(
                elt.s for elt in module.apps['TestData :: Top Level App'].endpoints['Second Level App'].attrs['bracketed_array_attr'].a.elt)
            self.assertTrue({'bval1', 'bval2'} & val_set)
            self.assertTrue('sla_attribute string' ==
                            module.apps['TestData :: Top Level App'].endpoints['Second Level App'].attrs['sla_attribute'].s)

            self.assertTrue(
                'test id' == module.apps['TestData :: Top Level App'].types['TestType'].attrs['id'].s)

        except (IOError, Exception) as e:
            self.fail(traceback.format_exc())


if __name__ == '__main__':
    unittest.main()
