package syslutil

import (
	"fmt"

	"github.com/pkg/errors"
)

func Assert(ok bool, format string, args ...interface{}) {
	if !ok {
		panic(fmt.Errorf(format, args...))
	}
}

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func PanicOnErrorf(err error, format string, args ...interface{}) {
	if err != nil {
		panic(errors.Wrapf(err, format, args...))
	}
}
