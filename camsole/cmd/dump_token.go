package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/grrtrr/clccam"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func init() {
	// Print details of a CAM service/user token to stdout
	Root.AddCommand(&cobra.Command{
		Use:     "token  <token | /path/to/token.file>",
		Aliases: []string{"t", "dump-token"},
		Short:   "Print details of a CAM token",
		PreRunE: checkArgs(1, "Need token contents or file path"),
		Run: func(cmd *cobra.Command, args []string) {
			if camToken, err := tokenFromStringOrFile(args[0]); err != nil {
				fmt.Printf("Token does not decode: %s\n", err)
			} else if cl, err := clccam.Token(camToken).Claims(); err != nil {
				fmt.Printf("Invalid CAM token %q: %s\n", args[0], err)
			} else {
				const timeFmt = `Mon Jan _2 15:04:05 MST 2006`
				var (
					table = tablewriter.NewWriter(os.Stdout)
					exp   = time.Unix(cl.Exp, 0).Format(timeFmt)
				)

				if cl.IsPermanent() {
					exp = "never (permanent token)"
				}

				fmt.Printf("%s:\n", cl)

				table.SetAutoFormatHeaders(false)
				table.SetAutoWrapText(false)
				table.SetHeader([]string{"Field", "Token Value"})

				table.AppendBulk([][]string{
					[]string{"exp", exp},
					[]string{"iat", time.Unix(cl.Iat, 0).Format(timeFmt)},
					[]string{"jti", cl.Jti.String()},
				})

				switch cl.Type {
				case "user":
					table.AppendBulk([][]string{
						[]string{"sub", cl.Subject},
						[]string{"name", cl.Name},
						[]string{"organization", cl.Organization},
					})
				case "service":
					table.AppendBulk([][]string{
						[]string{"instance", cl.InstanceId},
						[]string{"machine", cl.MachineId},
						[]string{"service", cl.ServiceId},
					})
				default:
					die("Unexpected token type %q.", cl.Type)
				}
				table.Render()
			}
		},
	})
}
