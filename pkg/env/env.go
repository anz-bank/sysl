package env

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

//  SYSL_CACHE
//      Cache location in current directory, defaults to "sysl_modules" if
//      SYSL_MODULES is enabled
//  SYSL_MODULES
//      Whether the sysl modules is enabled.
//      Enable by default, set to "off" to disable sysl modules.
//  SYSL_PLANTUML
//      URL of PlantUML server. Sysl depends upon
//      [PlantUML](http://plantuml.com/) for diagram generation.
//  SYSL_PROXY
//      Proxy service to use, won't use SYSL_PROXY if not set
//  SYSL_SSH_PASSPHRASE
// 	    Passphrase for SSH private key file if there's any
//  SYSL_SSH_PRIVATE_KEY
//      SSH private key file path for git/github credentials + domains to use
//      them on
//  SYSL_TOKENS
//      Tokens to use for git/github credentials + domains to use them on
//      eg: SYSL_TOKENS=github.com:1234,gitlab.com:567
//
//  The following development-only vars will only be reported by sysl env if
//  their values differ from their defaults...
//
//  SYSL_DEV_RENEST_FLATTENED_TYPES
//      off:    Don't renest flattened types.
//      retain: Renest flattened types, but retain the original.
//      move:   Renest flattened types, removing the original.

type Var string

//nolint:revive,stylecheck
const (
	SYSL_MODULES         Var = "SYSL_MODULES"
	SYSL_PLANTUML        Var = "SYSL_PLANTUML"
	SYSL_TOKENS          Var = "SYSL_TOKENS" //nolint:gosec
	SYSL_CACHE           Var = "SYSL_CACHE"
	SYSL_PROXY           Var = "SYSL_PROXY"
	SYSL_SSH_PRIVATE_KEY Var = "SYSL_SSH_PRIVATE_KEY"
	SYSL_SSH_PASSPHRASE  Var = "SYSL_SSH_PASSPHRASE" //nolint:gosec

	// Development-only vars
	SYSL_DEV_RENEST_FLATTENED_TYPES Var = "SYSL_DEV_RENEST_FLATTENED_TYPES"
	SYSL_DEV_UPDATE_GOLDEN_TESTS    Var = "SYSL_DEV_UPDATE_GOLDEN_TESTS"
)

var entries = map[Var]*entry{
	SYSL_CACHE:           newEntry("sysl_modules"),
	SYSL_MODULES:         newEntry("on", "off"),
	SYSL_PLANTUML:        newEntry("https://plantuml.com/plantuml"),
	SYSL_PROXY:           newEntry(""),
	SYSL_SSH_PASSPHRASE:  newEntry(""),
	SYSL_SSH_PRIVATE_KEY: newEntry(""),
	SYSL_TOKENS:          newEntry(""),

	SYSL_DEV_RENEST_FLATTENED_TYPES: newEntry("off", "retain", "move"),
	SYSL_DEV_UPDATE_GOLDEN_TESTS:    newEntry("off", "on"),
}

type entry struct {
	defaultValue string
	validValues  []string
	validationRE *regexp.Regexp
	once         sync.Once
	value        *string
}

func newEntry(defaultValue string, otherValidValues ...string) *entry {
	var validValues []string
	var validValuesRE *regexp.Regexp
	if len(otherValidValues) > 0 {
		validValues = append(make([]string, 0, 1+len(otherValidValues)), defaultValue)
		validValues = append(validValues, otherValidValues...)

		validValuesRE = regexp.MustCompile(
			fmt.Sprintf("^(%s)$", regexp.MustCompile(strings.Join(validValues, "|"))))
	}

	return &entry{
		defaultValue: defaultValue,
		validValues:  validValues,
		validationRE: validValuesRE,
	}
}

var Vars = func() VarSlice {
	ret := make(VarSlice, 0, len(entries))
	for e := range entries {
		ret = append(ret, e)
	}
	sort.Sort(ret)
	return ret
}()

func (e Var) Default() string {
	return entries[e].defaultValue
}

func (e Var) Name() string {
	return string(e)
}

func (e Var) On() bool {
	return e.Value() == "on"
}

func (e Var) Value() string {
	entry := entries[e]
	entry.once.Do(func() {
		var value string
		if value = os.Getenv(string(e)); value == "" {
			value = entry.defaultValue
		}
		if entry.validationRE != nil && !entry.validationRE.MatchString(value) {
			expectation := ""
			if len(entry.validValues) > 0 {
				expectation = fmt.Sprintf(" (expecting one of: %s)", strings.Join(entry.validValues, ", "))
			}
			logrus.Errorf("invalid %s=%q%s; assuming the default: %s",
				e, value, expectation, entry.defaultValue)
			value = entry.defaultValue
		}
		entry.value = &value
	})
	return *entry.value
}

type VarSlice []Var

func (v VarSlice) Len() int {
	return len(v)
}

func (v VarSlice) Less(i, j int) bool {
	return v[i] < v[j]
}

func (v VarSlice) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}
