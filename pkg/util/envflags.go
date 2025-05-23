package util

import (
	"fmt"
	"github.com/spf13/pflag"
	"os"
	"strings"
)

// SetFlagsFromEnv parses all registered flags in the given flagset,
// and if they are not already set it attempts to set their values from
// environment variables. Environment variables take the name of the flag but
// are UPPERCASE, and any dashes are replaced by underscores. Environment
// variables additionally are prefixed by the given string followed by
// and underscore. For example, if prefix=PREFIX: some-flag => PREFIX_SOME_FLAG
func SetFlagsFromEnv(fs *pflag.FlagSet, prefix string) (err error) {
	alreadySet := make(map[string]bool)

	fs.Visit(func(f *pflag.Flag) {
		alreadySet[f.Name] = true
	})

	fs.VisitAll(func(f *pflag.Flag) {
		if !alreadySet[f.Name] {
			key := strings.ToUpper(strings.Replace(f.Name, "-", "_", -1))
			if prefix != "" {
				key = prefix + "_" + strings.ToUpper(strings.Replace(f.Name, "-", "_", -1))
			}
			val := os.Getenv(key)
			if val != "" {
				if serr := fs.Set(f.Name, val); serr != nil {
					err = fmt.Errorf("invalid value %q for %s: %v", val, key, serr)
				}
			}
		}
	})

	return err
}
