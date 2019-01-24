package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// Upload the contents of a file as a new blob
	blobAdd = &cobra.Command{
		Use:     "blob  /path/to/file> [/path/to/another/file]",
		Aliases: []string{"upload", "up"},
		Short:   "Upload a new file as a blob",
		PreRunE: checkAtLeastArgs(1, "Need at least 1 file name"),
		Run: func(cmd *cobra.Command, args []string) {
			for i, fileName := range args {

				b, err := ioutil.ReadFile(fileName)
				if err != nil {
					die("failed to read %s: %s", fileName, err)
				}
				res, err := client.UploadFile(fileName, b)
				if err != nil {
					die("failed to upload %s: %s", fileName, err)
				}
				if i > 0 {
					fmt.Println("")
				}
				fmt.Printf("UUID:  %s\n", res.Url)
				fmt.Printf("URL:   https://%s%s\n", strings.TrimRight(rootFlags.url, "/"), res.Url)
			}
		},
	}
)

func init() {
	Root.AddCommand(blobAdd)
}
