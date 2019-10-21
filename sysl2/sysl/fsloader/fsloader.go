package fsloader

import (
	"github.com/anz-bank/sysl/sysl2/sysl/rootfinder"
	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
	"github.com/spf13/afero"
)

func LazyFS(newFs func() (afero.Fs, error)) func() (afero.Fs, error) {
	var fs afero.Fs
	return func() (afero.Fs, error) {
		if fs == nil {
			var err error
			fs, err = newFs()
			if err != nil {
				return nil, err
			}
		}
		return fs, nil
	}
}

func LazyChrootFS(parentFs afero.Fs, rf rootfinder.RootFinder) func() (afero.Fs, error) {
	return LazyFS(func() (afero.Fs, error) {
		root, err := rf.Root()
		if err != nil {
			return nil, err
		}
		return syslutil.NewChrootFs(parentFs, root), nil
	})
}

func LazyOsFS() func() (afero.Fs, error) {
	return LazyFS(func() (afero.Fs, error) {
		return afero.NewOsFs(), nil
	})
}

func EagerFS(fs afero.Fs) func() (afero.Fs, error) {
	return func() (afero.Fs, error) {
		return fs, nil
	}
}
