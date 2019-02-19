import yaml

from sysl.importers.import_swagger import SwaggerTranslator
from sysl.util import writer


class FakeLogger:
    def __init__(self):
        self.warnings = []

    def warn(self, msg):
        self.warnings.append(msg)


SWAGGER_WITH_ARRAY_TYPE_WITH_EXAMPLE = r"""swagger: "2.0"
basePath: /fruit-basket
info:
    title: Fruit API
    version: 1.0.0
definitions:
    FruitBasket:
        additionalProperties: false
        properties:
            fruit:
                example: '[{"id":"banana"}, {"id":"mango"}]'
                items:
                    type: object
                type: array
paths: {}
"""

SWAGGER_WITH_TYPELESS_ITEMS = r"""swagger: "2.0"
basePath: /fruit-basket
info:
    title: Fruit API
    version: 1.0.0
definitions:
    FruitBasket:
        additionalProperties: false
        properties:
            fruit:
                items:
                    type: object
paths: {}
"""

SWAGGER_OBJECT_WITH_NO_PROPERTIES = r"""swagger: "2.0"
basePath: /fruit-basket
info:
    title: Fruit API
    version: 1.0.0
definitions:
    MysteriousObject:
        type: object
paths: {}
"""


def test_importing_swagger_array_type_with_example_produces_sysl_type():
    swag = yaml.load(SWAGGER_WITH_ARRAY_TYPE_WITH_EXAMPLE)
    w = writer.Writer('sysl')
    t = SwaggerTranslator(logger=FakeLogger())
    t.translate(swag, appname='', package='', w=w)
    output = str(w)
    expected_fragment = '    !type FruitBasket:\n        fruit <: set of {}'
    assert expected_fragment in output


def test_importing_swagger_typeless_thing_with_items_produces_warning():
    swag = yaml.load(SWAGGER_WITH_TYPELESS_ITEMS)
    w = writer.Writer('sysl')
    logger = FakeLogger()
    t = SwaggerTranslator(logger=logger)
    t.translate(swag, appname='', package='', w=w)
    expected_warnings = ['Ignoring unexpected "items". Schema has "items" but did not have defined "type". Note: {\'items\': {\'type\': \'object\'}}']
    assert logger.warnings == expected_warnings


def test_importing_swagger_propertyless_object_works_without_warnings():
    swag = yaml.load(SWAGGER_OBJECT_WITH_NO_PROPERTIES)
    w = writer.Writer('sysl')
    logger = FakeLogger()
    t = SwaggerTranslator(logger=logger)
    t.translate(swag, appname='', package='', w=w)
    output = str(w)

    expected_fragment = '    !type MysteriousObject:\n'
    assert expected_fragment in output

    expected_warnings = []
    assert logger.warnings == expected_warnings
