package cmd

import (
	"fmt"
	"os"

	humanize "github.com/dustin/go-humanize"
	"github.com/grrtrr/clccam"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func init() {
	var cmdInstances = &cobra.Command{ // Top-level command
		Use:     "vm",
		Aliases: []string{"instance"},
		Short:   "Manage instances",
	}

	cmdInstances.AddCommand(instanceGet)
	Root.AddCommand(cmdInstances)
}

var (
	// List one or more instances (VMs)
	instanceGet = &cobra.Command{
		Use:     "ls  [instanceID1, ...]",
		Aliases: []string{"list", "show"},
		Short:   "List CAM instances",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				if instances, err := client.GetInstances(); err != nil {
					die("failed to query instance list: %s", err)
				} else if cmd.Flags().Lookup("json").Value.String() != "true" {
					printInstances(instances)
				}
			} else {
				for _, instanceId := range args {
					if instance, err := client.GetInstance(instanceId); err != nil {
						fmt.Fprintf(os.Stderr, "Failed to query instance %s: %s\n", instanceId, err)
					} else if cmd.Flags().Lookup("json").Value.String() != "true" {
						printInstances([]clccam.Instance{instance})

						if len(instance.Service.Machines) > 0 {
							var table = tablewriter.NewWriter(os.Stdout)

							table.SetAutoFormatHeaders(false)
							table.SetAlignment(tablewriter.ALIGN_LEFT)
							table.SetAutoWrapText(true)

							fmt.Printf("%s machines:\n", instance.Name)
							table.SetHeader([]string{"Name", "State"})
							for _, m := range instance.Service.Machines {
								table.Append([]string{m.Name, m.State.String()})
							}
							table.Render()
						}
					}
				}
			}
		},
	}
)

func printInstances(instances []clccam.Instance) {
	if len(instances) == 0 {
		fmt.Println("No instances.")
	} else {
		var table = tablewriter.NewWriter(os.Stdout)

		table.SetAutoFormatHeaders(false)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(true)

		table.SetHeader([]string{
			"Name", "ID", "Service", "Box", "Owner", "Updated", "Operation", "State",
		})
		for _, i := range instances {
			table.Append([]string{
				i.Name,
				i.ID,
				i.Service.ID,
				i.Box.String(),
				i.Owner,
				humanize.Time(i.Updated.Local()),
				i.Operation.Event.String(),
				i.State.String(),
			})
		}
		table.Render()
	}
}
