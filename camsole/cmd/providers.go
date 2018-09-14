package cmd

import (
	"fmt"
	"os"

	humanize "github.com/dustin/go-humanize"
	"github.com/grrtrr/clccam"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	cmdProviders = &cobra.Command{ // Top-level command
		Use:     "providers",
		Aliases: []string{"prov"},
		Short:   "Manage providers",
	}

	// List one or more providers
	providerGet = &cobra.Command{
		Use:     "ls  [providerID1, ...]",
		Aliases: []string{"list", "show"},
		Short:   "List providers",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				if providers, err := client.GetProviders(); err != nil {
					die("failed to query provider list: %s", err)
				} else if cmd.Flags().Lookup("json").Value.String() != "true" {
					printProviders(providers)
				}
			} else {
				for _, providerId := range args {
					if provider, err := client.GetProvider(providerId); err != nil {
						die("failed to query provider %s: %s", providerId, err)
					} else if cmd.Flags().Lookup("json").Value.String() != "true" {
						printProviders([]clccam.Provider{provider})
						if len(provider.Services) > 0 {
							fmt.Printf("\n%s available services:\n", provider.Name)
							for _, s := range provider.Services {
								fmt.Println("  -", s.Name)
							}
						}
					}
				}
			}
		},
	}

	// providerDelete removes a provider
	providerDelete = &cobra.Command{
		Use:     "rm  providerId",
		Aliases: []string{"delete", "purge"},
		Short:   "Remove provider",
		PreRunE: checkArgs(1, "Need a provider ID"),
		Run: func(cmd *cobra.Command, args []string) {
			if err := client.DeleteProvider(args[0]); err != nil {
				die("failed to delete provider %s: %s", args[0], err)
			}
			fmt.Printf("Deleted provider %s.\n", args[0])
		},
	}
)

func init() {
	cmdProviders.AddCommand(providerGet, providerDelete)
	Root.AddCommand(cmdProviders)
}

// printProviders prints @providers in tabulated form.
func printProviders(providers []clccam.Provider) {
	if len(providers) == 0 {
		fmt.Println("No providers.")
	} else {
		var table = tablewriter.NewWriter(os.Stdout)

		table.SetAutoFormatHeaders(false)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(true)

		table.SetHeader([]string{
			"Name", "Type", "ID", "Owner", "Created", "Updated", "State",
		})
		for _, p := range providers {
			table.Append([]string{
				p.Name,
				p.Type,
				p.ID.String(),
				p.Owner,
				humanize.Time(p.Created.Local()), humanize.Time(p.Updated.Local()),
				p.State,
			})
		}
		table.Render()
	}
}
