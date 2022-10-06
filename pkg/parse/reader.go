package parse

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/anz-bank/golden-retriever/retriever"
	"github.com/anz-bank/sysl/pkg/env"
	"github.com/anz-bank/sysl/pkg/syslutil"

	"github.com/anz-bank/golden-retriever/reader"
	"github.com/anz-bank/golden-retriever/reader/filesystem"
	"github.com/anz-bank/golden-retriever/reader/remotefs"
	"github.com/anz-bank/golden-retriever/retriever/git"
	"github.com/spf13/afero"
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
	tokens := make(map[string]string)
	if tokensStr := env.SYSL_TOKENS.Value(); tokensStr != "" {
		hostTokens := strings.Split(tokensStr, ",")
		for _, t := range hostTokens {
			arr := strings.Split(t, ":")
			if len(arr) != 2 {
				return nil, fmt.Errorf(
					"envvar %s is invalid, should be in format `hosta:<tokena>,hostb:<tokenb>`",
					env.SYSL_TOKENS,
				)
			}
			tokens[arr[0]] = arr[1]
		}
	}

	keys := make(map[string]git.SSHKey)
	if keysStr := env.SYSL_SSH_KEYS.Value(); keysStr != "" {
		hostKeys := strings.Split(keysStr, ",")
		for _, k := range hostKeys {
			arr := strings.Split(k, ":")
			if len(arr) != 3 {
				return nil, fmt.Errorf(
					"envvar %s is invalid, should be in format `hosta:<keya>:<passphrasea>,hostb:<keyb>:<passphraseb>`",
					env.SYSL_SSH_KEYS,
				)
			}
			keys[arr[0]] = git.SSHKey{
				PrivateKey:         arr[1],
				PrivateKeyPassword: arr[2],
			}
		}
	}

	auth := &git.AuthOptions{Tokens: tokens, SSHKeys: keys}
	// If a modules.yaml file already exists, use it. Otherwise use a caching Git retriever without
	// pinning remote imports.
	pinnerPath := filepath.Join(SyslRootDir(fs), "modules.yaml")
	if ok, err := afero.Exists(fs, pinnerPath); ok && err == nil {
		return remotefs.NewPinnerGitRetriever(pinnerPath, auth)
	}
	return git.NewWithCache(auth, git.NewPlainFscache(remotefs.CacheDir)), nil
}

func SyslRootDir(fs afero.Fs) string {
	root := "."
	if v, is := fs.(*syslutil.ChrootFs); is {
		root = v.Root()
	}
	return filepath.Join(root, SyslRootMarker)
}
