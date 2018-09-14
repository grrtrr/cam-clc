// Cobra commandline console driver
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/grrtrr/clccam/camsole/cmd"
	"github.com/spf13/cobra"
)

func main() {
	// Logging format - we don't need date/file
	log.SetFlags(log.Ltime)

	// Do sort the commands alphabetically
	cobra.EnableCommandSorting = true

	if err := cmd.Root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
