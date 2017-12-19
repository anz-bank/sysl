from sysl.reljam import reljam
from sysl.core import syslloader
import unittest
import re
import os
import sys
import tempfile


class TestReljam(unittest.TestCase):

    def setUp(self):
        self.outpath = tempfile.gettempdir()
        self.package_prefix = 'io.sysl.test.data.petshop.'
        self.entities = {'Employee', 'Breed', 'Pet', 'EmployeeTendsPet'}
        (self.module, _, _) = syslloader.load(
            '/test/data/petshop/petshop', True, '.')
        # export has the side effect of removing entities from the list, so the member list
        # is not passed in
        reljam.export('model', self.module, 'PetShopModel', self.outpath,
                      self.package_prefix +
                      'model', {'Employee', 'Breed',
                                'Pet', 'EmployeeTendsPet'},
                      ['json_out', 'xml_out'])
        reljam.export('facade', self.module, 'PetShopFacade', self.outpath,
                      self.package_prefix + 'facade', {'Employee', 'Breed', 'Pet', 'EmployeeTendsPet'}, [])

    def test_module_export(self):
        self.checkCounts(self.package_prefix + 'model',
                         {'Employee': 0, 'Breed': 0,
                             'Pet': 0, 'EmployeeTendsPet': 0},
                         {'Employee': 44, 'Breed': 44, 'Pet': 48, 'EmployeeTendsPet': 46})

    def test_facade_export(self):
        self.checkCounts(self.package_prefix + 'facade',
                         {'PetShopFacade': 0}, {'PetShopFacade': 9})

    def test_serializers(self):
        self.checkCounts(self.package_prefix + 'model',
                         {'PetShopModelJsonSerializer': 0,
                             'PetShopModelXmlSerializer': 0},
                         {'PetShopModelJsonSerializer': 5, 'PetShopModelXmlSerializer': 5})

    def checkCounts(self, package, expected_root_counts, expected_model_counts):
        root_pattern = re.compile('\W_?root\W')
        model_pattern = re.compile(r'[\( ]model[.;]|\.?model[; \)]')

        keys = set(expected_root_counts.iterkeys())
        assert keys == set(expected_model_counts.iterkeys()), (
            expected_root_counts, expected_model_counts)

        for entity in keys:
            root_count = 0
            model_count = 0
            file_name = self.outpath + '/' + \
                package.replace('.', '/') + '/' + entity + '.java'
            with open(file_name, 'r') as java_text:
                for line in java_text:
                    root_count += len(root_pattern.findall(line))
                    model_count += len(model_pattern.findall(line))
                self.assertEquals(expected_root_counts[entity], root_count, )
                self.assertEquals(expected_model_counts[entity], model_count)


if __name__ == '__main__':
    unittest.main()
