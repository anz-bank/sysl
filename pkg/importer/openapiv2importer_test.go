package importer

import (
	"testing"

	"github.com/alecthomas/assert"

	"github.com/sirupsen/logrus"

	"github.com/spf13/afero"
)

func TestNewOpenAPIV3Importer(t *testing.T) {
	fs := afero.NewMemMapFs()
	afero.WriteFile(fs, "/foo.yaml", []byte(`openapi: "3.0.0"
info:
  title: Fruit API
  version: 1.0.0
paths: {}
components:
  schemas:
    Basket:
      type: object
      properties:
        id:
          description: id of the basket
          type: string
        fruit:
          description: something about fruit
          $ref: '#/components/schemas/Fruit'
        veggies:
          description: something about veggies
          $ref: '#/components/schemas/Veggies'
    Fruit:
      type: object
      properties:
        id:
          description: id of the fruit
          type: string
        entity:
          type: string
          enum:
            - expense
            - income
        name:
          type: string
        details:
          type: array
          items:
            $ref: '#/components/schemas/Details'
          minItems: 1
    Veggies:
      type: object
      properties:
        id:
          description: id of the vegetable
          type: string
          pattern: hello
        entity:
          type: string
          enum:
            - expense
            - income
        name:
          type: string
        details:
          type: array
          items:
            $ref: '#/components/schemas/Details'
    Details:
      type: object
      properties:
        month:
          type: string
          enum:
            - Jan
            - Feb
            - Mar
            - Apr
            - May
            - Jun
            - Jul
            - Aug
            - Sep
            - Oct
            - Nov
            - Dec
          description: month of the planting
        monthlyTotal:
          type: int
        transDet:
          type: array
          items:
            type: string

`), 0755)
	//log, _ := test.NewNullLogger()
	i := NewOpenAPIV3Importer(logrus.New(), fs)
	res, err := i.Load("/foo.yaml")
	assert.NotNil(t, res)
	assert.NoError(t, err)
	assert.Equal(t, "", res)

}
