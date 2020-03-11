package nodisk

import (
	"github.com/anz-bank/sysl/internal/lsp/foo"
)

func _() {
	foo.Foo() //@complete("F", Foo, IntFoo, StructFoo)
}
