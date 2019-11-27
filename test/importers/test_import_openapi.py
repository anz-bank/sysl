import platform
import pytest
import yaml
import re

from sysl.importers.import_openapi import OpenApiTranslator, make_default_logger
from sysl.util import writer


class FakeLogger:
    def __init__(self):
        self.warnings = []

    def warn(self, msg):
        self.warnings.append(msg)


def getOutputString(input):
    prefix = re.match('^ *', input).group(0)
    swag = yaml.load(re.sub(r'^' + prefix, '', input), Loader=yaml.FullLoader)
    w = writer.Writer('sysl')
    logger = FakeLogger()
    t = OpenApiTranslator(logger)
    t.translate(swag, appname='', package='', w=w)
    return str(w), logger


def test_importing_simple_openapi_with_error_type():
    output, _ = getOutputString(r"""
"openapi": "3.0"
info:
  title: Simple
paths:
  /test:
    get:
      responses:
        400:
          description: "client error"
          content:
              application/json:
                schema:
                  $ref: "#/components/schemas/SimpleObj"
components:
  schemas:
    SimpleObj:

      type: object
      properties:
        name:
          type: string
""")
    assert r"""
 "Simple" [package=""]:
    @description =:
        | No description.

    /test:
        GET:
            | No description.
            return error <: SimpleObj

    #---------------------------------------------------------------------------
    # definitions

    !type SimpleObj:
        name <: string?:
            @json_tag = "name"
""" in output


def test_importing_simple_openapi_with_json_tags():
    output, _ = getOutputString(r"""
"openapi": "3.0"
info:
  title: Simple
paths:
  /test:
    get:
      responses:
        200:
          description: "200 OK"
          content:
              application/json:
                schema:
                  $ref: "#/components/schemas/SimpleObj"
components:
  schemas:
    SimpleObj:

      type: object
      properties:
        name:
          type: string
""")
    assert r"""
 "Simple" [package=""]:
    @description =:
        | No description.

    /test:
        GET:
            | No description.
            return ok <: SimpleObj

    #---------------------------------------------------------------------------
    # definitions

    !type SimpleObj:
        name <: string?:
            @json_tag = "name"
""" in output


def test_importing_openapi_with_top_level_array():
    output, _ = getOutputString(r"""
openapi: "3.0"
info:
  title: Simple
paths: {}

components:
  schemas:
    TopLevelArray:
      type: array
      items:
        properties:
          name:
            type: string
""")
    assert r"""
 "Simple" [package=""]:
    @description =:
        | No description.

    #---------------------------------------------------------------------------
    # definitions

    !alias TopLevelArray:
        sequence of EXTERNAL_TopLevelArray_obj

    !alias EXTERNAL_TopLevelArray_obj:
        string
""" in output


def test_importing_openapi_object_with_required_and_optional_fields():
    output, _ = getOutputString(r"""
openapi: "3.0"
info:
  title: Simple
paths: {}

components:
  schemas:
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
""")
    assert r"""
    !type Goat:
        id <: string:
            @json_tag = "id"
        name <: string:
            @json_tag = "name"
        weight <: float?:
            @json_tag = "weight"
""" in output


def test_importing_openapi_with_query_params():
    output, _ = getOutputString(r"""
openapi: "3.0"
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
          schema:
            type: string
          required: false
          in: query
        - name: min_date
          schema:
            type: string
          required: true
          in: query
          format: date
""")
    assert 'GET ?key=string?&min_date=string:' in output


def test_importing_openapi_with_header_body_params():
    output, _ = getOutputString(r"""
openapi: "3.0"
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
          schema:
            type: integer
          required: false
          in: header
        - name: min_date
          schema:
            type: string
          required: true
          in: header
          format: date
      requestBody:
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/SimpleObj'
components:
  schemas:
    SimpleObj:
      type: object
      properties:
        name:
          type: string
""")
    assert r"""
 "Simple" [package=""]:
    @description =:
        | No description.

    /test:
        POST (SimpleObjRequest <: SimpleObj [~body], key <: int [~header, ~optional, name="key"], min_date <: string [~header, ~required, name="min_date"]):
            | No description.
            return

    #---------------------------------------------------------------------------
    # definitions

    !type SimpleObj:
        name <: string?:
            @json_tag = "name"
""" in output


def test_importing_openapi_with_sysl_keywords():
    output, _ = getOutputString(r"""openapi: "3.0"
basePath: /fruit-basket
info:
  title: Fruit API
  version: 1.0.0
components:
  schemas:
    SimpleObj:
      type: object
      properties:
        date:
          type: string
        string:
          type: string
paths: {}
""")
    assert 'date_ <: string?' in output
    assert 'string_ <: string?' in output


def test_importing_openapi_array_type_with_example_produces_sysl_type():
    output, _ = getOutputString(r"""openapi: "3.0"
basePath: /fruit-basket
info:
    title: Fruit API
    version: 1.0.0
components:
  schemas:
    FruitBasket:
        additionalProperties: false
        properties:
            fruit:
                example: '[{"id":"banana"}, {"id":"mango"}]'
                items:
                    type: object
                type: array
paths: {}
""")
    expected_fragment = '    !type FruitBasket:\n        fruit <: sequence of EXTERNAL_FruitBasket_fruit_obj'
    assert expected_fragment in output


def test_importing_openapi_typeless_thing_with_items_produces_warning():
    _, logger = getOutputString(r"""openapi: "3.0"
basePath: /fruit-basket
info:
    title: Fruit API
    version: 1.0.0
components:
  schemas:
    FruitBasket:
        additionalProperties: false
        properties:
            fruit:
                items:
                    type: object
paths: {}
""")
    expected_warnings = ['Ignoring unexpected "items". Schema has "items" but did not have defined "type". Note: {\'items\': {\'type\': \'object\'}}']
    assert logger.warnings == expected_warnings


def test_importing_openapi_propertyless_object_works_without_warnings():
    output, logger = getOutputString(r"""
openapi: "3.0"
basePath: /fruit-basket
info:
    title: Fruit API
    version: 1.0.0
components:
  schemas:
    MysteriousObject:
        type: object
paths: {}
""")
    expected_fragment = '    !alias EXTERNAL_MysteriousObject:\n'
    assert expected_fragment in output
    expected_warnings = []
    assert logger.warnings == expected_warnings


def test_importing_openapi_spec_with_a_path_works_without_warnings():
    output, logger = getOutputString(r"""openapi: "3.0"
basePath: /api/v1

host: goat.example.com

info:
  title: Goat CRUD API
  version: 1.2.3

components:
  schemas:
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
          schema:
            type: string
          required: true
      produces:
        - application/json
      responses:
        '200':
          description: ''
          content:
              application/json:
                schema:
                  $ref: '#/components/schemas/Acknowledgement'
      summary: Delete a goat
""")
    assert r"""
 "Goat CRUD API" [package=""]:
    @version = "1.2.3"
    @host = "goat.example.com"
    @description =:
        | No description.

    /api/v1:

        /goat/delete-goat:
            POST ?goat_id=string:
                | Delete a goat.
                return ok <: Acknowledgement

    #---------------------------------------------------------------------------
    # definitions

    !type Acknowledgement:
        message <: string?:
            @json_tag = "message"
""" in output
    expected_warnings = []
    assert logger.warnings == expected_warnings


def test_importing_openapi_object_with_required_field_produces_sysl_type_with_required_field():
    output, _ = getOutputString(r"""openapi: "3.0"
basePath: /fruit-basket
info:
  title: Fruit API
  version: 1.0.0
components:
  schemas:
    Apple:
      properties:
        colour:
          type: string
      required:
        - colour
      type: object
paths: {}
""")
    expected_fragment = '!type Apple:\n        colour <: string:\n'
    assert expected_fragment in output


def test_import_of_openapi_path_that_returns_array_of_defined_object_type():
    output, logger = getOutputString(r"""openapi: "3.0"
basePath: /api/v1

host: goat.example.com

info:
  title: Goat CRUD API
  version: 1.2.3

components:
  schemas:
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
          content:
              application/json:
                schema:
                  type: array
                  items:
                    $ref: '#/components/schemas/Goat'
      summary: Gotta get goats
""")
    assert r"""
 "Goat CRUD API" [package=""]:
    @version = "1.2.3"
    @host = "goat.example.com"
    @description =:
        | No description.

    /api/v1:

        /goat/get-goats:
            GET:
                | Gotta get goats.
                return ok <: sequence of Goat

    #---------------------------------------------------------------------------
    # definitions

    !type Goat:
        birthday <: date?:
            @json_tag = "birthday"
        name <: string?:
            @json_tag = "name"
""" in output
    expected_warnings = []
    assert logger.warnings == expected_warnings


def test_import_of_openapi_path_that_has_a_defined_201_response():
    output, logger = getOutputString(r"""openapi: "3.0"
basePath: /api/v1

host: goat.example.com

info:
  title: Goat CRUD API
  version: 1.2.3

components:
  schemas:
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
          schema:
            type: string
          required: true
        - name: birthday
          in: query
          schema:
            type: string
          required: true
      responses:
        '201':
          description: ''
          headers:
            Location:
              description: Location of the newly allocated goat.
      summary: Creates a goat.
""")
    assert r"""
 "Goat CRUD API" [package=""]:
    @version = "1.2.3"
    @host = "goat.example.com"
    @description =:
        | No description.

    /api/v1:

        /goat/create-goat:
            POST ?name=string&birthday=string:
                | Creates a goat.
                return

    #---------------------------------------------------------------------------
    # definitions

    !type Goat:
        birthday <: date?:
            @json_tag = "birthday"
        name <: string?:
            @json_tag = "name"
""" in output
    expected_warnings = []
    assert logger.warnings == expected_warnings


def test_import_of_openapi_path_that_has_a_body_parameter():
    output, logger = getOutputString(r"""openapi: "3.0"
basePath: /api/v1

host: goat.example.com

info:
  title: Goat CRUD API
  version: 1.2.3

components:
  schemas:
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
      requestBody:
        content:
            application/json:
              schema:
                $ref: '#/components/schemas/Goat'
      responses:
        '201':
          description: ''
          headers:
            Location:
              description: Location of the newly allocated goat.
      summary: Creates a goat.
""")
    assert r"""
 "Goat CRUD API" [package=""]:
    @version = "1.2.3"
    @host = "goat.example.com"
    @description =:
        | No description.

    /api/v1:

        /goat/create-goat:
            POST (GoatRequest <: Goat [~body]):
                | Creates a goat.
                return

    #---------------------------------------------------------------------------
    # definitions

    !type Goat:
        birthday <: date?:
            @json_tag = "birthday"
        name <: string?:
            @json_tag = "name"
""" in output
    expected_warnings = []
    assert logger.warnings == expected_warnings


def test_import_of_openapi_path_with_error_response():
    # Characterisation test. Who knows if this is what we actually want it to do.
    output, logger = getOutputString(r"""
openapi: "3.0"
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
""")
    assert r"""
 "Goat CRUD API" [package=""]:
    @version = "1.2.3"
    @host = "goat.example.com"
    @description =:
        | No description.

    /api/v1:

        /goat/status:
            GET:
                | Check goat status
                return

    #---------------------------------------------------------------------------
    # definitions
""" in output
    expected_warnings = []
    assert logger.warnings == expected_warnings


def test_import_of_openapi_path_with_default_response_is_not_implemented():
    _, logger = getOutputString(r"""openapi: "3.0"
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
""")
    expected_warnings = ['default and x-* responses are not implemented']
    assert logger.warnings == expected_warnings


def test_import_of_openapi_path_with_x_dash_whatever_response_is_not_implemented():
    _, logger = getOutputString(r"""openapi: "3.0"
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
""")
    expected_warnings = ['default and x-* responses are not implemented']
    assert logger.warnings == expected_warnings


def test_import_of_openapi_path_with_description_only_200_response():
    output, logger = getOutputString(r"""
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
""")
    assert r"""
 "Goat CRUD API" [package=""]:
    @version = "1.2.3"
    @host = "goat.example.com"
    @description =:
        | No description.

    /api/v1:

        /goat/status:
            GET:
                | Get goat status
                return

    #---------------------------------------------------------------------------
    # definitions
""" in output
    expected_warnings = []
    assert logger.warnings == expected_warnings


def test_import_of_openapi_path_with_description_only_201_response():
    output, logger = getOutputString(r"""
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
""")
    assert r"""
 "Goat CRUD API" [package=""]:
    @version = "1.2.3"
    @host = "goat.example.com"
    @description =:
        | No description.

    /api/v1:

        /goat/status:
            POST:
                | Update goat status
                return

    #---------------------------------------------------------------------------
    # definitions
""" in output
    expected_warnings = []
    assert logger.warnings == expected_warnings


def test_parse_typespec_boolean():
    t = OpenApiTranslator(FakeLogger())
    assert t.parse_typespec({'type': 'boolean', 'description': 'foo'}) == ('bool', 'foo')


def test_parse_typespec_datetime():
    t = OpenApiTranslator(FakeLogger())
    assert t.parse_typespec({'type': 'string', 'format': 'date-time', 'description': 'foo'}) == ('datetime', 'foo')


def test_parse_typespec_integer():
    t = OpenApiTranslator(FakeLogger())
    assert t.parse_typespec({'type': 'integer', 'description': 'foo'}) == ('int', 'foo')


def test_parse_typespec_int32():
    t = OpenApiTranslator(FakeLogger())
    assert t.parse_typespec({'type': 'integer', 'format': 'int32', 'description': 'foo'}) == ('int32', 'foo')


def test_parse_typespec_int64():
    t = OpenApiTranslator(FakeLogger())
    assert t.parse_typespec({'type': 'integer', 'format': 'int64', 'description': 'foo'}) == ('int64', 'foo')


def test_parse_typespec_number_is_translated_to_float():
    t = OpenApiTranslator(FakeLogger())
    assert t.parse_typespec({'type': 'number', 'description': 'foo'}) == ('float', 'foo')


def test_parse_typespec_float_is_translated_to_float():
    t = OpenApiTranslator(FakeLogger())
    assert t.parse_typespec({'type': 'number', 'format': 'float', 'description': 'foo'}) == ('float', 'foo')


def test_parse_typespec_double_is_translated_to_float():
    t = OpenApiTranslator(FakeLogger())
    assert t.parse_typespec({'type': 'number', 'format': 'double', 'description': 'foo'}) == ('float', 'foo')


def test_parse_typespec_object():
    t = OpenApiTranslator(FakeLogger())
    assert t.parse_typespec({'type': 'object', 'description': 'foo'}, '', 'Object') == ('EXTERNAL_Object_obj', 'foo')


def test_parse_typespec_ref():
    t = OpenApiTranslator(FakeLogger())
    assert t.parse_typespec({'$ref': '#/components/schemas/Barr', 'description': 'foo'}) == ('Barr', 'foo')


def test_parse_typespec_warns_and_ignores_type_if_array_items_type_has_both_type_and_ref():
    l = FakeLogger()
    t = OpenApiTranslator(logger=l)

    array_type = {
        'type': 'array',
        'items': {
            '$ref': '#/components/schemas/Barr',
            'type': 'Foo',
        },
        'description': 'this is where we keep our ill-specified things'
    }
    assert t.parse_typespec(array_type) == ('sequence of Barr', 'this is where we keep our ill-specified things')
    expected_warnings = ['Ignoring unexpected "type". Schema has "$ref" but also has unexpected "type". Note: {\'items\': {\'type\': \'Foo\', \'$ref\': \'#/components/schemas/Barr\'}, \'type\': \'array\'}']
    assert l.warnings == expected_warnings


def test_translate_path_template_params_leaves_paths_without_templates_unchanged():
    l = FakeLogger()
    t = OpenApiTranslator(logger=l)
    expected_warnings = []
    assert t.translate_path_template_params('/foo/barr/', []) == '/foo/barr/'
    assert l.warnings == expected_warnings


def test_translate_path_template_params_rewrites_dashed_template_names_as_camelcase_string_typed_parameters():
    l = FakeLogger()
    t = OpenApiTranslator(logger=l)
    assert t.translate_path_template_params('/foo/{fizz-buzz}/', []) == '/foo/{fizzBuzz<:string}/'
    expected_warnings = ['not enough path params path: /foo/{fizz-buzz}/', 'could not find type for path param: {fizz-buzz} in params[]']
    assert l.warnings == expected_warnings


def test_make_default_logger_returns_something_thats_probably_a_logger():
    logger = make_default_logger()
    assert hasattr(logger, 'warn')


def test_import_of_openapi_path_with_path_var_type_in_api():
    output, logger = getOutputString(r"""
openapi: "3.0"
info:
  title: Sample API
  description: API description in Markdown.
  version: 1.0.0
host: api.example.com
basePath: /v1
schemes:
  - https
paths:
  /users/{id}:
    parameters:
      - in: path
        name: id
        schema:
          type: integer
        required: true
        description: The user ID.
      - in: header
        name: request-id
        schema:
          type: string
        required: true
        description: the request ID.
    # GET/users/{id}?metadata=true
    get:
      summary: Gets a user by ID
      # Note we only define the query parameter, because the {id} is defined at the path level.
      parameters:
        - in: query
          name: metadata
          schema:
            type: boolean
          required: false
          description: If true, the endpoint returns only the user metadata.
      responses:
        '200':
          description: OK
""")
    assert r"""
 "Sample API" [package=""]:
    @version = "1.0.0"
    @host = "api.example.com"
    @description =:
        | API description in Markdown.

    /v1:

        /users/{id<:int}:
            GET (request_id <: string [~header, ~required, name="request-id"]) ?metadata=bool?:
                | No description.
                return

    #---------------------------------------------------------------------------
    # definitions
""" in output
    expected_warnings = []
    assert logger.warnings == expected_warnings


def test_import_of_openapi_path_with_path_var_type_in_global_parameters():
    output, logger = getOutputString(r"""
openapi: "3.0"
info:
  title: Sample API
  description: API description in Markdown.
  version: 1.0.0
host: api.example.com
basePath: /v1
schemes:
  - https
paths:
  /users/{id}:
    get:
      summary: Gets a user by ID.
      parameters:
        - in: query
          name: metadata
          schema:
            type: boolean
          required: false
          description: If true, the endpoint returns only the user metadata.
        - $ref: '#parameters/id'
        - $ref: '#parameters/request-id'
      responses:
        '200':
          description: OK
parameters:
  id:
    in: path
    name: id
    schema:
      type: integer
    required: true
    description: The user ID.
  request-id:
    in: header
    name: request-id
    schema:
      type: string
    required: true
    description: the request ID.
""")
    assert r"""
 "Sample API" [package=""]:
    @version = "1.0.0"
    @host = "api.example.com"
    @description =:
        | API description in Markdown.

    /v1:

        /users/{id<:int}:
            GET (request_id <: string [~header, ~required, name="request-id"]) ?metadata=bool?:
                | No description.
                return

    #---------------------------------------------------------------------------
    # definitions
""" in output
    expected_warnings = []
    assert logger.warnings == expected_warnings


def test_import_of_openapi_path_with_header_var_overridden_in_method():
    output, logger = getOutputString(r"""
openapi: "3.0"
info:
  title: Sample API
  description: API description in Markdown.
  version: 1.0.0
host: api.example.com
basePath: /v1
schemes:
  - https
paths:
  /users/{id}:
    parameters:
      - in: path
        name: id
        schema:
          type: integer
        required: true
        description: The user ID.
      - in: header
        name: metadata
        schema:
          type: boolean
        required: false
    get:
      summary: Gets a user by ID
      parameters:
       -  in: header
          name: metadata
          schema:
            type: string
          enum:
           - public
           - personal
           - all
          required: true
      responses:
        '200':
          description: OK
""")
    assert r"""
 "Sample API" [package=""]:
    @version = "1.0.0"
    @host = "api.example.com"
    @description =:
        | API description in Markdown.

    /v1:

        /users/{id<:int}:
            GET (metadata <: string [~header, ~required, name="metadata"]):
                | No description.
                return

    #---------------------------------------------------------------------------
    # definitions
""" in output
    expected_warnings = []
    assert logger.warnings == expected_warnings


def test_import_of_openapi_path_with_paths_var_referring_global_params_object():
    output, logger = getOutputString(r"""
openapi: "3.0"
info:
  title: Sample API
  description: API description in Markdown.
  version: 1.0.0
host: api.example.com
basePath: /v1
schemes:
  - https
paths:
  /users/{id}:
    parameters:
      - $ref: '#/parameters/id'
      - in: header
        name: metadata
        schema:
          type: boolean
        required: false
    get:
      summary: Gets a user by ID
      responses:
        '200':
          description: OK
parameters:
  id:
    in: path
    name: id
    schema:
      type: integer
    required: true
    description: The user ID.
""")
    assert r"""
 "Sample API" [package=""]:
    @version = "1.0.0"
    @host = "api.example.com"
    @description =:
        | API description in Markdown.

    /v1:

        /users/{id<:int}:
            GET (metadata <: bool [~header, ~optional, name="metadata"]):
                | No description.
                return

    #---------------------------------------------------------------------------
    # definitions
""" in output
    expected_warnings = []
    assert logger.warnings == expected_warnings


def test_import_of_openapi_path_with_paths_var_type_overridden_in_second_method():
    output, logger = getOutputString(r"""
openapi: "3.0"
info:
  title: Sample API
  description: API description in Markdown.
  version: 1.0.0
host: api.example.com
basePath: /v1
schemes:
  - https
paths:
  /users/{id}:
    parameters:
      - $ref: '#/parameters/id'
      - in: header
        name: metadata
        schema:
          type: boolean
        required: false
    get:
      summary: Gets a user by ID
      parameters:
       -  in: header
          name: metadata
          schema:
            type: string
          enum:
           - public
           - personal
           - all
          required: true
      responses:
        '200':
          description: OK
    delete:
      summary: Gets a user by ID
      parameters:
       -  in: path
          name: id
          schema:
            type: string
          required: true
      responses:
        '200':
          description: OK
parameters:
  id:
    in: path
    name: id
    schema:
      type: integer
    required: true
    description: The user ID.
""")
    assert r"""
 "Sample API" [package=""]:
    @version = "1.0.0"
    @host = "api.example.com"
    @description =:
        | API description in Markdown.

    /v1:

        /users/{id<:int}:
            GET (metadata <: string [~header, ~required, name="metadata"]):
                | No description.
                return

        /users/{id<:string}:
            DELETE (metadata <: bool [~header, ~optional, name="metadata"]):
                | No description.
                return

    #---------------------------------------------------------------------------
    # definitions
""" in output
    expected_warnings = []
    assert logger.warnings == expected_warnings
