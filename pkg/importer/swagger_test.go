package importer

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLoadOpenAPI2(t *testing.T) {
	spec := `swagger: "2.0"
version: 1.0.0
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
`
	logger := logrus.New()
	basePath := ""
	importer := MakeOpenAPI2Importer(logger, basePath, "")
	result, err := importer.Load(spec)
	assert.NoError(t, err)
	t.Log(result)
}
