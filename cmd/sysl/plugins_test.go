package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegrationModelPlugins(t *testing.T) {
	t.Parallel()

	expectedOut := `{
  "edges": null,
  "nodes": [
    {
      "expanded": null,
      "group": "Namespace1",
      "key": "Namespace1 :: App1",
      "label": "App1",
      "visible": true
    },
    {
      "expanded": null,
      "group": "Namespace1",
      "key": "Namespace1 :: App2",
      "label": "App2",
      "visible": true
    },
    {
      "expanded": true,
      "group": null,
      "key": "Namespace1",
      "label": "Namespace1",
      "visible": true
    }
  ]
}
`
	output := runSyslWithOutput(t, ".sysl",
		"transform", "../../tests/simple.sysl", "--script", "../../pkg/plugins/integration_model_plugin.arraiz")
	assert.Equal(t, expectedOut, output)
}
