import platform
import pytest
import yaml

from sysl.importers.import_swagger import SwaggerTranslator, make_default_logger
from sysl.util import writer


class FakeLogger:
    def __init__(self):
        self.warnings = []

    def warn(self, msg):
        self.warnings.append(msg)


SIMPLE_SWAGGER_EXAMPLE = r"""
swagger: "2.0"
info:
  title: Simple
paths:
  /test:
    get:
      responses:
        200:
          description: 200 OK
          schema:
            $ref: '#/definitions/SimpleObj'
definitions:
  SimpleObj:

    type: object
    properties:
      name:
        type: string
"""

SWAGGER_TOP_LEVEL_ARRAY_EXAMPLE = r"""
swagger: "2.0"
info:
  title: Simple
paths: {}

definitions:
  TopLevelArray:
    type: array
    items:
      properties:
        name:
          type: string
"""

SWAGGER_TOP_LEVEL_ARRAY_EXAMPLE_EXPECTED_SYSL = r"""
 "Simple" [package=""]:
    @description =:
        | No description.

    #---------------------------------------------------------------------------
    # definitions

    !alias TopLevelArray:
        sequence of EXTERNAL_TopLevelArray_obj

    !alias EXTERNAL_TopLevelArray_obj:
        string
"""

SWAGGER_REQUIRED_AND_OPTIONAL_FIELDS_EXAMPLE = r"""
swagger: "2.0"
info:
  title: Simple
paths: {}

definitions:
  Goat:
    type: object
    required:
      - id
      - name
    properties:
      id:
        type: string
      name:
        type: string
      weight:
        type: number
"""

SWAGGER_REQUIRED_AND_OPTIONAL_FIELDS_EXAMPLE_EXPECTED_SYSL = r"""
    !type Goat:
        id <: string:
            @json_tag = "id"
        name <: string:
            @json_tag = "name"
        weight <: float?:
            @json_tag = "weight"
"""

SWAGGER_QUERY_PARAM_EXAMPLE = r"""
swagger: "2.0"
info:
  title: Simple
paths:
  /test:
    get:
      responses:
        200:
          description: 200 OK
      parameters:
        - name: key
          type: string
          required: false
          in: query
        - name: min_date
          type: string
          required: true
          in: query
          format: date
"""

SWAGGER_HEADER_AND_BODY_PARAM_EXAMPLE = r"""
swagger: "2.0"
info:
  title: Simple
paths:
  /test:
    post:
      responses:
        200:
          description: 200 OK
      parameters:
        - name: key
          type: integer
          required: false
          in: header
        - name: min_date
          type: string
          required: true
          in: header
          format: date
        - name: createrequest
          required: true
          in: body
          schema:
            $ref: '#/definitions/SimpleObj'
definitions:
  SimpleObj:
    type: object
    properties:
      name:
        type: string
"""

SWAGGER_HEADER_AND_BODY_PARAM_EXAMPLE_EXPECTED_SYSL = r"""
 "Simple" [package=""]:
    @description =:
        | No description.

    /test:
        POST (createrequest <: SimpleObj [~body], key <: int [~header, ~optional, name="key"], min_date <: string [~header, ~required, name="min_date"]):
            | No description.
            return 200

    #---------------------------------------------------------------------------
    # definitions

    !type SimpleObj:
        name <: string?:
            @json_tag = "name"
"""

SWAGGER_WITH_SYSL_KEYWORDS_EXAMPLE = r"""swagger: "2.0"
basePath: /fruit-basket
info:
  title: Fruit API
  version: 1.0.0
definitions:
  SimpleObj:
    type: object
    properties:
      date:
        type: string
      string:
        type: string
paths: {}
"""

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

SWAGGER_OBJECT_WITH_NO_PROPERTIES = r"""
swagger: "2.0"
basePath: /fruit-basket
info:
    title: Fruit API
    version: 1.0.0
definitions:
    MysteriousObject:
        type: object
paths: {}
"""

EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH = r"""swagger: "2.0"
basePath: /api/v1

host: goat.example.com

info:
  title: Goat CRUD API
  version: 1.2.3

definitions:
  Acknowledgement:
    additionalProperties: false
    description: Indicates if a request has succeeded or not.
    properties:
      message:
        type: string
    type: object

paths:
  /goat/delete-goat:
    post:
      consumes:
        - application/json
      description: Delete a goat.
      parameters:
        - name: goat_id
          in: query
          type: string
          required: true
      produces:
        - application/json
      responses:
        '200':
          description: ''
          schema:
            $ref: '#/definitions/Acknowledgement'
      summary: Delete a goat
"""

EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_EXPECTED_SYSL = r"""
 "Goat CRUD API" [package=""]:
    @version = "1.2.3"
    @host = "goat.example.com"
    @description =:
        | No description.

    /api/v1:

        /goat/delete-goat:
            POST ?goat_id=string:
                | Delete a goat.
                return Acknowledgement

    #---------------------------------------------------------------------------
    # definitions

    !type Acknowledgement:
        message <: string?:
            @json_tag = "message"
"""

SWAGGER_OBJECT_WITH_A_REQUIRED_PROPERTY = r"""swagger: "2.0"
basePath: /fruit-basket
info:
  title: Fruit API
  version: 1.0.0
definitions:
  Apple:
    properties:
      colour:
        type: string
    required:
      - colour
    type: object
paths: {}
"""

EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_RETURNING_ARRAY_OF_DEFINED_OBJECT_TYPE = r"""swagger: "2.0"
basePath: /api/v1

host: goat.example.com

info:
  title: Goat CRUD API
  version: 1.2.3

definitions:
  Goat:
    additionalProperties: false
    properties:
      name:
        type: string
      birthday:
        type: string
        format: date
    type: object

paths:
  /goat/get-goats:
    get:
      consumes:
        - application/json
      description: Gotta get goats.
      produces:
        - application/json
      responses:
        '200':
          description: ''
          schema:
            type: array
            items:
              $ref: '#/definitions/Goat'
      summary: Gotta get goats
"""

EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_RETURNING_ARRAY_OF_DEFINED_OBJECT_TYPE_EXPECTED_SYSL = r"""
 "Goat CRUD API" [package=""]:
    @version = "1.2.3"
    @host = "goat.example.com"
    @description =:
        | No description.

    /api/v1:

        /goat/get-goats:
            GET:
                | Gotta get goats.
                return sequence of Goat

    #---------------------------------------------------------------------------
    # definitions

    !type Goat:
        birthday <: date?:
            @json_tag = "birthday"
        name <: string?:
            @json_tag = "name"
"""

EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_201_LOCATION_HEADER_RESPONSE = r"""swagger: "2.0"
basePath: /api/v1

host: goat.example.com

info:
  title: Goat CRUD API
  version: 1.2.3

definitions:
  Goat:
    additionalProperties: false
    properties:
      name:
        type: string
      birthday:
        type: string
        format: date
    type: object

paths:
  /goat/create-goat:
    post:
      consumes:
        - application/json
      description: Creates a goat.
      produces:
        - application/json
      parameters:
        - name: name
          in: query
          type: string
          required: true
        - name: birthday
          in: query
          type: string
          required: true
      responses:
        '201':
          description: ''
          headers:
            Location:
              description: Location of the newly allocated goat.
      summary: Creates a goat.
"""

EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_201_LOCATION_HEADER_RESPONSE_EXPECTED_SYSL = r"""
 "Goat CRUD API" [package=""]:
    @version = "1.2.3"
    @host = "goat.example.com"
    @description =:
        | No description.

    /api/v1:

        /goat/create-goat:
            POST ?name=string&birthday=string:
                | Creates a goat.
                return 201

    #---------------------------------------------------------------------------
    # definitions

    !type Goat:
        birthday <: date?:
            @json_tag = "birthday"
        name <: string?:
            @json_tag = "name"
"""


EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_BODY_PARAMETER = r"""swagger: "2.0"
basePath: /api/v1

host: goat.example.com

info:
  title: Goat CRUD API
  version: 1.2.3

definitions:
  Goat:
    additionalProperties: false
    properties:
      name:
        type: string
      birthday:
        type: string
        format: date
    type: object

paths:
  /goat/create-goat:
    post:
      consumes:
        - application/json
      description: Creates a goat.
      produces:
        - application/json
      parameters:
        - name: body
          in: body
          schema:
            $ref: '#/definitions/Goat'
      responses:
        '201':
          description: ''
          headers:
            Location:
              description: Location of the newly allocated goat.
      summary: Creates a goat.
"""

EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_BODY_PARAMETER_EXPECTED_SYSL = r"""
 "Goat CRUD API" [package=""]:
    @version = "1.2.3"
    @host = "goat.example.com"
    @description =:
        | No description.

    /api/v1:

        /goat/create-goat:
            POST (body <: Goat [~body]):
                | Creates a goat.
                return 201

    #---------------------------------------------------------------------------
    # definitions

    !type Goat:
        birthday <: date?:
            @json_tag = "birthday"
        name <: string?:
            @json_tag = "name"
"""


EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_ERROR_RESPONSE = r"""
swagger: "2.0"
basePath: /api/v1

host: goat.example.com

info:
  title: Goat CRUD API
  version: 1.2.3

paths:
  /goat/status:
    get:
      description: Check goat status
      produces:
        - application/json
      responses:
        '200':
          description: 'here be status'
        '500':
          description: 'alas, the server is broken'
      summary: Check goat status
"""


EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_ERROR_RESPONSE_EXPECTED_SYSL = r"""
 "Goat CRUD API" [package=""]:
    @version = "1.2.3"
    @host = "goat.example.com"
    @description =:
        | No description.

    /api/v1:

        /goat/status:
            GET:
                | Check goat status
                return 200, 500

    #---------------------------------------------------------------------------
    # definitions
"""


EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_DEFAULT_RESPONSE = r"""swagger: "2.0"
basePath: /api/v1

host: goat.example.com

info:
  title: Goat CRUD API
  version: 1.2.3

paths:
  /goat/status:
    get:
      description: Check goat status
      produces:
        - application/json
      responses:
        'default':
          description: 'here be default response'
      summary: Check goat status
"""


EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_X_DASH_WHATEVER_RESPONSE = r"""swagger: "2.0"
basePath: /api/v1

host: goat.example.com

info:
  title: Goat CRUD API
  version: 1.2.3

paths:
  /goat/status:
    get:
      description: Check goat status
      produces:
        - application/json
      responses:
        'x-banana':
          description: 'here be an x-banana response'
      summary: Check goat status
"""

EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_200_RESPONSE_DESCRIPTION_ONLY = r"""
basePath: /api/v1

host: goat.example.com

info:
  title: Goat CRUD API
  version: 1.2.3

paths:
  /goat/status:
    get:
      description: Get goat status
      produces:
        - application/json
      responses:
        '200':
          description: 'okay'
      summary: Get goat status
"""

EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_200_RESPONSE_DESCRIPTION_ONLY_EXPECTED_SYSL = r"""
 "Goat CRUD API" [package=""]:
    @version = "1.2.3"
    @host = "goat.example.com"
    @description =:
        | No description.

    /api/v1:

        /goat/status:
            GET:
                | Get goat status
                return 200

    #---------------------------------------------------------------------------
    # definitions
"""

EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_201_RESPONSE_DESCRIPTION_ONLY = r"""
basePath: /api/v1

host: goat.example.com

info:
  title: Goat CRUD API
  version: 1.2.3

paths:
  /goat/status:
    post:
      description: Update goat status
      produces:
        - application/json
      responses:
        '201':
          description: 'created'
      summary: Update goat status
"""

EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_201_RESPONSE_DESCRIPTION_ONLY_EXPECTED_SYSL = r"""
 "Goat CRUD API" [package=""]:
    @version = "1.2.3"
    @host = "goat.example.com"
    @description =:
        | No description.

    /api/v1:

        /goat/status:
            POST:
                | Update goat status
                return 201

    #---------------------------------------------------------------------------
    # definitions
"""


def getOutputString(input):
    swag = yaml.load(input)
    w = writer.Writer('sysl')
    logger = FakeLogger()
    t = SwaggerTranslator(logger)
    t.translate(swag, appname='', package='', w=w)
    return str(w), logger


def test_importing_simple_swagger_with_json_tags():
    output, _ = getOutputString(SIMPLE_SWAGGER_EXAMPLE)
    assert 'name <: string?:\n            @json_tag = "name"' in output


def test_importing_swagger_with_top_level_array():
    output, _ = getOutputString(SWAGGER_TOP_LEVEL_ARRAY_EXAMPLE)
    assert SWAGGER_TOP_LEVEL_ARRAY_EXAMPLE_EXPECTED_SYSL in output


def test_importing_swagger_object_with_required_and_optional_fields():
    output, _ = getOutputString(SWAGGER_REQUIRED_AND_OPTIONAL_FIELDS_EXAMPLE)
    assert SWAGGER_REQUIRED_AND_OPTIONAL_FIELDS_EXAMPLE_EXPECTED_SYSL in output


def test_importing_swagger_with_query_params():
    output, _ = getOutputString(SWAGGER_QUERY_PARAM_EXAMPLE)
    assert 'GET ?key=string?&min_date=string:' in output


def test_importing_swagger_with_header_body_params():
    output, _ = getOutputString(SWAGGER_HEADER_AND_BODY_PARAM_EXAMPLE)
    assert SWAGGER_HEADER_AND_BODY_PARAM_EXAMPLE_EXPECTED_SYSL in output


def test_importing_swagger_with_sysl_keywords():
    output, _ = getOutputString(SWAGGER_WITH_SYSL_KEYWORDS_EXAMPLE)
    assert 'date_ <: string?' in output
    assert 'string_ <: string?' in output


def test_importing_swagger_array_type_with_example_produces_sysl_type():
    output, _ = getOutputString(SWAGGER_WITH_ARRAY_TYPE_WITH_EXAMPLE)
    expected_fragment = '    !type FruitBasket:\n        fruit <: sequence of EXTERNAL_FruitBasket_fruit_obj'
    assert expected_fragment in output


def test_importing_swagger_typeless_thing_with_items_produces_warning():
    _, logger = getOutputString(SWAGGER_WITH_TYPELESS_ITEMS)
    expected_warnings = ['Ignoring unexpected "items". Schema has "items" but did not have defined "type". Note: {\'items\': {\'type\': \'object\'}}']
    assert logger.warnings == expected_warnings


def test_importing_swagger_propertyless_object_works_without_warnings():
    output, logger = getOutputString(SWAGGER_OBJECT_WITH_NO_PROPERTIES)
    expected_fragment = '    !alias EXTERNAL_MysteriousObject:\n'
    assert expected_fragment in output
    expected_warnings = []
    assert logger.warnings == expected_warnings


def test_importing_swagger_spec_with_a_path_works_without_warnings():
    output, logger = getOutputString(EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH)
    assert EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_EXPECTED_SYSL in output
    expected_warnings = []
    assert logger.warnings == expected_warnings


def test_importing_swagger_object_with_required_field_produces_sysl_type_with_required_field():
    output, _ = getOutputString(SWAGGER_OBJECT_WITH_A_REQUIRED_PROPERTY)
    expected_fragment = '!type Apple:\n        colour <: string:\n'
    assert expected_fragment in output


def test_import_of_swagger_path_that_returns_array_of_defined_object_type():
    output, logger = getOutputString(EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_RETURNING_ARRAY_OF_DEFINED_OBJECT_TYPE)
    assert EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_RETURNING_ARRAY_OF_DEFINED_OBJECT_TYPE_EXPECTED_SYSL in output
    expected_warnings = []
    assert logger.warnings == expected_warnings


def test_import_of_swagger_path_that_has_a_defined_201_response():
    output, logger = getOutputString(EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_201_LOCATION_HEADER_RESPONSE)
    assert EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_201_LOCATION_HEADER_RESPONSE_EXPECTED_SYSL in output
    expected_warnings = []
    assert logger.warnings == expected_warnings


def test_import_of_swagger_path_that_has_a_body_parameter():
    output, logger = getOutputString(EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_BODY_PARAMETER)
    assert EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_BODY_PARAMETER_EXPECTED_SYSL in output
    expected_warnings = []
    assert logger.warnings == expected_warnings


def test_import_of_swagger_path_with_error_response():
    # Characterisation test. Who knows if this is what we actually want it to do.
    output, logger = getOutputString(EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_ERROR_RESPONSE)
    assert EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_ERROR_RESPONSE_EXPECTED_SYSL in output
    expected_warnings = []
    assert logger.warnings == expected_warnings


def test_import_of_swagger_path_with_default_response_is_not_implemented():
    _, logger = getOutputString(EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_DEFAULT_RESPONSE)
    expected_warnings = ['default responses and x-* responses are not implemented']
    assert logger.warnings == expected_warnings


def test_import_of_swagger_path_with_x_dash_whatever_response_is_not_implemented():
    _, logger = getOutputString(EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_X_DASH_WHATEVER_RESPONSE)
    expected_warnings = ['default responses and x-* responses are not implemented']
    assert logger.warnings == expected_warnings


def test_import_of_swagger_path_with_description_only_200_response():
    output, logger = getOutputString(EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_200_RESPONSE_DESCRIPTION_ONLY)
    assert EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_200_RESPONSE_DESCRIPTION_ONLY_EXPECTED_SYSL in output
    expected_warnings = []
    assert logger.warnings == expected_warnings


def test_import_of_swagger_path_with_description_only_201_response():
    output, logger = getOutputString(EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_201_RESPONSE_DESCRIPTION_ONLY)
    assert EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_201_RESPONSE_DESCRIPTION_ONLY_EXPECTED_SYSL in output
    expected_warnings = []
    assert logger.warnings == expected_warnings


def test_parse_typespec_boolean():
    t = SwaggerTranslator(None)
    assert t.parse_typespec({'type': 'boolean', 'description': 'foo'}) == ('bool', 'foo')


def test_parse_typespec_datetime():
    t = SwaggerTranslator(None)
    assert t.parse_typespec({'type': 'string', 'format': 'date-time', 'description': 'foo'}) == ('datetime', 'foo')


def test_parse_typespec_integer():
    t = SwaggerTranslator(None)
    assert t.parse_typespec({'type': 'integer', 'description': 'foo'}) == ('int', 'foo')


def test_parse_typespec_int32():
    t = SwaggerTranslator(None)
    assert t.parse_typespec({'type': 'integer', 'format': 'int32', 'description': 'foo'}) == ('int32', 'foo')


def test_parse_typespec_int64():
    t = SwaggerTranslator(None)
    assert t.parse_typespec({'type': 'integer', 'format': 'int64', 'description': 'foo'}) == ('int64', 'foo')


def test_parse_typespec_number_is_translated_to_float():
    t = SwaggerTranslator(None)
    assert t.parse_typespec({'type': 'number', 'description': 'foo'}) == ('float', 'foo')


def test_parse_typespec_float_is_translated_to_float():
    t = SwaggerTranslator(None)
    assert t.parse_typespec({'type': 'number', 'format': 'float', 'description': 'foo'}) == ('float', 'foo')


def test_parse_typespec_double_is_translated_to_float():
    t = SwaggerTranslator(None)
    assert t.parse_typespec({'type': 'number', 'format': 'double', 'description': 'foo'}) == ('float', 'foo')


def test_parse_typespec_object():
    t = SwaggerTranslator(None)
    assert t.parse_typespec({'type': 'object', 'description': 'foo'}, '', 'Object') == ('EXTERNAL_Object_obj', 'foo')


def test_parse_typespec_ref():
    t = SwaggerTranslator(None)
    assert t.parse_typespec({'$ref': '#/definitions/Barr', 'description': 'foo'}) == ('Barr', 'foo')


def test_parse_typespec_warns_and_ignores_type_if_array_items_type_has_both_type_and_ref():
    l = FakeLogger()
    t = SwaggerTranslator(logger=l)

    array_type = {
        'type': 'array',
        'items': {
            '$ref': '#/definitions/Barr',
            'type': 'Foo',
        },
        'description': 'this is where we keep our ill-specified things'
    }
    assert t.parse_typespec(array_type) == ('sequence of Barr', 'this is where we keep our ill-specified things')
    expected_warnings = ['Ignoring unexpected "type". Schema has "$ref" but also has unexpected "type". Note: {\'items\': {\'type\': \'Foo\', \'$ref\': \'#/definitions/Barr\'}, \'type\': \'array\'}']
    assert l.warnings == expected_warnings


def test_translate_path_template_params_leaves_paths_without_templates_unchanged():
    t = SwaggerTranslator(logger=None, vocabulary_factory=(lambda: ['x']))
    assert t.translate_path_template_params('/foo/barr/') == '/foo/barr/'


def test_translate_path_template_params_rewrites_dashed_template_names_as_camelcase_string_typed_parameters():
    t = SwaggerTranslator(logger=None, vocabulary_factory=(lambda: ['x']))
    assert t.translate_path_template_params('/foo/{fizz-buzz}/') == '/foo/{fizzBuzz<:string}/'


def test_translate_path_template_params_rewrites_names_of_things_that_look_like_a_dictionary_word_ending_with_id_suffix_as_camelcase():
    t = SwaggerTranslator(logger=None, vocabulary_factory=(lambda: ['bread']))
    assert t.translate_path_template_params('/foo/{breadid}/') == '/foo/{breadId<:string}/'


def test_translate_path_template_params_wont_rewrite_names_of_things_ending_with_id_suffix_as_camelcase_if_no_vocabulary_present():
    l = FakeLogger()
    t = SwaggerTranslator(logger=l, vocabulary_factory=(lambda: []))
    # perhaps breadid is a valid word. we dont know, we have no vocab.
    assert t.translate_path_template_params('/foo/{breadid}/') == '/foo/{breadid<:string}/'
    assert l.warnings == ['could not load any vocabulary, janky environment-specific heuristics for renaming path template names may fail']


def test_translate_path_template_params_doesnt_rewrite_nonwords_ending_in_id_typed_parameters():
    t = SwaggerTranslator(logger=None, vocabulary_factory=(lambda: ['bread']))
    assert t.translate_path_template_params('/foo/{braedid}/') == '/foo/{braedid<:string}/'


@pytest.mark.skipif(platform.system() not in ('Linux', 'Darwin'), reason='no defined source of vocabulary for this platform')
def test_default_vocabulary_containing_common_business_nouns_is_defined_for_non_windows_platforms():
    t = SwaggerTranslator(None)
    assert 'customer' in t.words()


def test_make_default_logger_returns_something_thats_probably_a_logger():
    logger = make_default_logger()
    assert hasattr(logger, 'warn')
