package mod

import (
	"log"

	"github.com/anz-bank/sysl/pkg/env"

	"github.com/anz-bank/gop/pkg/cli"
	"github.com/spf13/afero"
)

func Retriever(fs afero.Fs) (cli.Retriever, error) {
	var retr cli.Retriever
	tokenmap, err := cli.NewTokenMap(env.SYSL_TOKENS.String(), "GIT_CREDENTIALS")
	if err != nil {
		return retr, err
	}

	var cache, proxy, privKey, passphrase string
	if moduleFlag := env.SYSL_MODULES.Value(); moduleFlag != "" && moduleFlag != "false" && moduleFlag != "off" {
		cache = env.SYSL_CACHE.Value()
		proxy = env.SYSL_PROXY.Value()
		privKey = env.SYSL_SSH_PRIVATE_KEY.Value()
		passphrase = env.SYSL_SSH_PASSPHRASE.Value()
	}

	return cli.Moduler(fs, "sysl_modules.yaml", cache, proxy, tokenmap, privKey, passphrase, log.Printf), nil
}
