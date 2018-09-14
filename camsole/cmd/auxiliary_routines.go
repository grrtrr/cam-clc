package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/grrtrr/clccam"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

/*
 * Helper Functions
 */

// die is like die in Perl
func die(format string, a ...interface{}) {
	format = fmt.Sprintf("%s: %s\n", path.Base(os.Args[0]), strings.TrimSpace(format))
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

// checkArgs returns a cobra-compatible PreRunE argument-validation function
func checkArgs(nargs int, errMsg string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != nargs {
			return errors.Errorf(errMsg)
		}
		return nil
	}
}

// checkAtLeastArgs is analogous to checkArgs
func checkAtLeastArgs(nargs int, errMsg string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < nargs {
			return errors.Errorf(errMsg)
		}
		return nil
	}
}

// tokenFromStringOrFile loads a CAM token from file or string
// @s: contents of the base64-encoded CAM JWT, or path to a file containing the token
// On error returns an empty (zero-valued) token along with the error.
func tokenFromStringOrFile(s string) (token clccam.Token, err error) {
	if _, err := os.Stat(s); err == nil { // @s points to a file
		contents, err := ioutil.ReadFile(s)
		if err != nil {
			return token, errors.Wrapf(err, "failed to read %s", s)
		}
		return clccam.Token(string(bytes.TrimSpace(contents))), nil
	} else if s == "" { // @s probably contains a token
		return token, errors.Errorf("empty token string")
	} else if _, err := clccam.Token(s).Decode(); err != nil {
		return token, errors.Wrapf(err, "failed to decode %q", s)
	}
	return clccam.Token(s), nil
}

// truncate ensures that the length of @s does not exceed @maxlen
func truncate(s string, maxlen int) string {
	if len(s) >= maxlen {
		s = s[:maxlen]
	}
	return s
}
