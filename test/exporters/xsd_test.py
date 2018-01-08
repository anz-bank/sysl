from sysl.reljam import reljam
from sysl.core import syslloader
from sysl.util.file import areFilesIdentical
import unittest
import re
import os
import sys
import tempfile

class TestXsd(unittest.TestCase):

    # def setUp(self):
    # self.outpath  = '.'
    # self.package_prefix  = 'io.sysl.test.data'
    # self.entities = {'TopLevelType'}
    # (self.module, _, _) = syslloader.load('/test/data/test_type_xsd', True, '.')

    # export has the side effect of removing entities from the list, so the member list
    # is not passed in

    def test_table_xsd(self):
        self.genAndCompare("/test/data/test_table_xsd",
                           "TestTableXsdModel",
                           os.path.join("test", "data", "test_table.xsd"))

    def test_simple_type(self):
        self.genAndCompare("/test/data/test_type_xsd",
                           "TestTypeXsdModel",
                           os.path.join("test", "data", "test_type.xsd"))

    def test_type_set(self):
        self.genAndCompare("/test/data/test_type_set_xsd",
                           "TestTypeSetXsdModel",
                           os.path.join("test", "data", "test_type_set.xsd"))

    def test_type_attribute(self):
        self.genAndCompare("/test/data/test_type_attr_xsd",
                           "TestTypeAttrXsdModel",
                           os.path.join("test", "data", "test_type_attr.xsd"))

    def test_table_attribute(self):
        self.genAndCompare("/test/data/test_table_attr_xsd",
                           "TestTableAttrXsdModel",
                           os.path.join("test", "data", "test_table_attr.xsd"))

    def genAndCompare(self, sysl_module, model, xsd_comparison_file):
        outpath = tempfile.gettempdir()
        package_prefix = 'io.sysl.test.data'
        (module, _, _) = syslloader.load(sysl_module, True, '.')

        reljam.export('xsd', module, model, outpath, package_prefix, {}, [])
        expected = os.path.join('.', xsd_comparison_file)
        real = os.path.join(outpath, model + ".xsd")
        self.assertTrue(areFilesIdentical(expected,real))

if __name__ == '__main__':
    unittest.main()
