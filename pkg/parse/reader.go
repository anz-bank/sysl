package parse

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/anz-bank/golden-retriever/reader"
	"github.com/anz-bank/golden-retriever/reader/filesystem"
	"github.com/anz-bank/golden-retriever/reader/remotefs"
	"github.com/anz-bank/golden-retriever/retriever"
	"github.com/anz-bank/golden-retriever/retriever/git"
	"github.com/spf13/afero"

	"github.com/anz-bank/sysl/pkg/env"
	"github.com/anz-bank/sysl/pkg/syslutil"
)

const SyslRootMarker = ".sysl"
const GitRootMarker = ".git"

func NewReader(fs afero.Fs) (reader.Reader, error) {
	pinner, err := NewPinner(fs)
	if err != nil {
		return nil, err
	}

	return remotefs.NewWithRetriever(
		filesystem.New(fs),
		pinner,
	), nil
}

func NewPinner(fs afero.Fs) (retriever.Retriever, error) {
	auth := &git.AuthOptions{}

	if tokensStr := env.SYSL_TOKENS.Value(); tokensStr != "" {
		var err error
		auth, err = auth.WithTokensFromString(tokensStr)
		if err != nil {
			return nil, fmt.Errorf("envvar %s is invalid: %w`", env.SYSL_TOKENS, err)
		}
	}

	if keysStr := env.SYSL_SSH_KEYS.Value(); keysStr != "" {
		auth.SSHKeys = make(map[string]git.SSHKey)
		hostKeys := strings.Split(keysStr, ",")
		for _, k := range hostKeys {
			arr := strings.Split(k, ":")
			if len(arr) != 3 {
				return nil, fmt.Errorf(
					"envvar %s is invalid, should be in format `hosta:<keya>:<passphrasea>,hostb:<keyb>:<passphraseb>`",
					env.SYSL_SSH_KEYS,
				)
			}
			auth.SSHKeys[arr[0]] = git.SSHKey{
				PrivateKey:         arr[1],
				PrivateKeyPassword: arr[2],
			}
		}
	}

	// If a modules.yaml file already exists, use it. Otherwise use a caching Git retriever without
	// pinning remote imports.
	pinnerPath := filepath.Join(SyslRootDir(fs), "modules.yaml")
	if ok, err := afero.Exists(fs, pinnerPath); ok && err == nil {
		return remotefs.NewPinnerGitRetriever(pinnerPath, auth)
	}
	return git.NewWithOptions(&git.NewGitOptions{
		AuthOptions:   auth,
		Cacher:        git.NewPlainFscache(remotefs.CacheDir),
		NoForcedFetch: remotefs.NoForcedFetch,
	}), nil
}

func SyslRootDir(fs afero.Fs) string {
	root := "."
	if v, is := fs.(*syslutil.ChrootFs); is {
		root = v.Root()
	}
	return filepath.Join(root, SyslRootMarker)
}
