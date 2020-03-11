package other

import (
	references "github.com/anz-bank/sysl/internal/lsp/references"
)

func _() {
	references.Q = "hello" //@mark(assignExpQ, "Q")
	bob := func(_ string) {}
	bob(references.Q) //@mark(bobExpQ, "Q")
}
