package importer

import (
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMakeOpenAPI3Importer(t *testing.T) {
	logger := logrus.New()
	basePath := ""
	importer := MakeOpenAPI3Importer(logger, basePath, "")
	t.Log(importer)
}

func TestLoadOpenAPI3(t *testing.T) {
	spec := `openapi: "3.0"
info:
  title: Simple
paths:
  /test/:
    get:
      responses:
        200:
          description: "200 OK"
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/SimpleObj"
        500:
          $ref: "#/components/responses/500Response"
components:
  schemas:
    SimpleObj:
      type: object
      properties:
        name:
          type: string
    SimpleObj2:
      type: object
      properties:
        name:
          type: SimpleObj
  responses:
    500Response:
      description: Internal Server Error
      schema:
        $ref: "#/components/schemas/SimpleObj"
`
	logger := logrus.New()
	basePath := ""
	importer := MakeOpenAPI3Importer(logger, basePath, "")
	result, err := importer.Load(spec)
	assert.NoError(t, err)
	t.Log(result)
}

func TestLoadOpenAPI3RefSpecialChars(t *testing.T) {
	spec := `openapi: "3.0"
info:
  title: Simple
paths:
  /testSomeEndpoint:SpecialChars:
    get:
      responses:
        200:
          description: "200 OK"
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/SimpleObj:WithSpecialChar"
        "500":
          description: "Error"
          headers:
            content-type:
              schema:
                type: string
                example: "application/json"
            request-id:
              schema:
                type: string
                example: "12345678-1234-1234-1234-123456789012"
          content:
              application/json:
                schema:
                  $ref: "#/components/schemas/SimpleObj:WithSpecialChar"

components:
  schemas:
    SimpleObj:WithSpecialChar:
      type: object
      properties:
        name:
          type: string
  responses:
    500Response:
      description: Internal Server Error
      schema:
        $ref: "#/components/schemas/SimpleObj:WithSpecialChar"
`
	logger := logrus.New()
	basePath := ""
	importer := MakeOpenAPI3Importer(logger, basePath, "")
	result, err := importer.Load(spec)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(result, "!type _testSomeEndpoint%3ASpecialChars_500:"))
}

func TestLoadOpenAPI3MultipleErrorResponsesWithHeaders(t *testing.T) {
	spec := `openapi: "3.0"
info:
  title: Simple
paths:
  /testSomeEndpoint:
    get:
      responses:
        "200":
          description: "200 OK"
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Result"
        "404":
          description: "404 Not Found"
          headers:
            content-type:
              schema:
                type: string
                example: "application/json"
            request-id:
              schema:
                type: string
                example: "12345678-1234-1234-1234-123456789012"
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Error"
        "500":
          description: "Error"
          headers:
            content-type:
              schema:
                type: string
                example: "application/json"
            request-id:
              schema:
                type: string
                example: "12345678-1234-1234-1234-123456789012"
          content:
              application/json:
                schema:
                  $ref: "#/components/schemas/Error"

components:
  schemas:
    Error:
      type: object
      properties:
        code:
          type: string
        description:
          type: string
    Result:
      type: object
      properties:
        result:
          type: string
`
	logger := logrus.New()
	basePath := ""
	imp := MakeOpenAPI3Importer(logger, basePath, "")
	result, err := imp.Load(spec)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(result, "!type _testSomeEndpoint_404:"))
	assert.True(t, strings.Contains(result, "!type _testSomeEndpoint_500:"))
}
