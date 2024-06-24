package internal

import (
	"os"

	"github.com/mattn/go-isatty"
)

var (
	EnvDev  = "dev"
	EnvProd = "prod"
)

// Env returns which environment we're running in. Currently this only checks if
// stdin is a tty, but in the future we might use environment variables.
func Env() string {
	if isatty.IsTerminal(os.Stdin.Fd()) {
		return EnvDev
	} else {
		return EnvProd
	}
}
