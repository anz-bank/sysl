package roothandler

import (
	"github.com/spf13/afero"
)

type FsDecider interface {
	GetTopFsLoader() func() (afero.Fs, error)
	GetImportFsLoader() func() (afero.Fs, error)
	GenerateAppropriateFs() error
}


