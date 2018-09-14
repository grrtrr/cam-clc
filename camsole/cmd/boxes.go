package cmd

import (
	"fmt"
	"os"
	"sort"

	humanize "github.com/dustin/go-humanize"
	"github.com/grrtrr/clccam"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	cmdBoxes = &cobra.Command{ // Top-level command
		Use:   "box",
		Short: "Manage boxes",
	}

	// boxList lists one or more boxes
	boxList = &cobra.Command{
		Use:     "ls  [boxId, ...]",
		Aliases: []string{"list", "show"},
		Short:   "List box(es)",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				if boxes, err := client.GetBoxes(); err != nil {
					die("failed to query %s box list: %s", args[0], err)
				} else if cmd.Flags().Lookup("json").Value.String() != "true" {
					listBoxes(boxes)
				}
			} else {
				for _, boxId := range args {
					if box, err := client.GetBox(boxId); err != nil {
						fmt.Fprintf(os.Stderr, "Failed to query box %s: %s\n", boxId, err)
					} else if cmd.Flags().Lookup("json").Value.String() != "true" {
						listBoxes([]clccam.Box{box})
					}
				}
			}
		},
	}

	// boxStack prints the box stack
	boxStack = &cobra.Command{
		Use:     "stack  boxId",
		Short:   "List the stack of @boxId (if any)",
		PreRunE: checkArgs(1, "Need a box ID"),
		Run: func(cmd *cobra.Command, args []string) {
			if boxes, err := client.GetBoxStack(args[0]); err != nil {
				die("failed to query %s box stack: %s", args[0], err)
			} else if cmd.Flags().Lookup("json").Value.String() != "true" {
				var filtered []clccam.Box

				// The output is unsorted. Move box in question to the top of the list.
				for i, box := range boxes {
					if box.ID.String() == args[0] && i > 0 {
						filtered = append([]clccam.Box{box}, filtered...)
					} else {
						filtered = append(filtered, box)
					}
				}
				listBoxes(filtered)
			}
		},
	}

	// boxVersions prints box versions
	boxVersions = &cobra.Command{
		Use:     "ver  boxId",
		Aliases: []string{"version"},
		Short:   "List version(s) of @boxId",
		PreRunE: checkArgs(1, "Need a box ID"),
		Run: func(cmd *cobra.Command, args []string) {
			if boxes, err := client.GetBoxVersions(args[0]); err != nil {
				die("failed to query %s box versions: %s", args[0], err)
			} else if cmd.Flags().Lookup("json").Value.String() == "true" {
			} else if len(boxes) == 0 {
				fmt.Printf("No versions available.\n")
			} else {
				var table = tablewriter.NewWriter(os.Stdout)

				table.SetAutoFormatHeaders(false)
				table.SetAlignment(tablewriter.ALIGN_LEFT)
				table.SetAutoWrapText(true)

				table.SetHeader([]string{
					"Name", "Version", "ID", "Owner", "Visibility", "Created", "Updated",
				})
				sort.Slice(boxes, func(i, j int) bool {
					return boxes[i].Version().LessThan(boxes[j].Version())
				})
				for _, b := range boxes {
					table.Append([]string{
						b.Name,
						b.Version().String(),
						b.ID.String(),
						b.Owner,
						b.Visibility.String(),
						humanize.Time(b.Created.Local()), humanize.Time(b.Updated.Local()),
					})
				}
				table.Render()
			}
		},
	}

	// boxDiff prints differences between boxes
	boxDiff = &cobra.Command{
		Use:     "diff  boxId",
		Short:   "Print the differences of @boxId",
		PreRunE: checkArgs(1, "Need a box ID"),
		Run: func(cmd *cobra.Command, args []string) {
			// NOTE/FIXME: seems to be privileged or not used, getting 405 response; no other documentation.
			if err := client.WithJsonResponse().GetBoxDiff(args[0]); err != nil {
				die("failed to query %s box differences: %s", args[0], err)
			}
		},
	}

	// boxBindings prints bindings
	boxBindings = &cobra.Command{
		Use:     "bindings  boxId",
		Aliases: []string{"bdg"},
		Short:   "List the bindings of @boxId",
		PreRunE: checkArgs(1, "Need a box ID"),
		Run: func(cmd *cobra.Command, args []string) {
			if bindings, err := client.GetBoxBindings(args[0]); err != nil {
				die("failed to query %s box bindings: %s", args[0], err)
			} else if cmd.Flags().Lookup("json").Value.String() != "true" {
				var table = tablewriter.NewWriter(os.Stdout)

				table.SetAutoFormatHeaders(false)
				table.SetAlignment(tablewriter.ALIGN_LEFT)
				table.SetAutoWrapText(true)
				table.SetHeader([]string{"Name", "ID", "Icon"})

				for _, b := range bindings {
					table.Append([]string{b.Name, b.ID.String(), b.Icon.String()})
				}
				table.Render()
			}
		},
	}

	// boxDelete removes a box
	boxDelete = &cobra.Command{
		Use:     "rm  boxId",
		Aliases: []string{"delete", "get-rid-of"},
		Short:   "Remove box",
		PreRunE: checkArgs(1, "Need a box ID"),
		Run: func(cmd *cobra.Command, args []string) {
			if err := client.DeleteBox(args[0]); err != nil {
				die("failed to delete box %s: %s", args[0], err)
			}
			fmt.Printf("Deleted box %s.\n", args[0])
		},
	}
)

// listBoxes prints a table with a subset of box information for each of the @boxes.
func listBoxes(boxes []clccam.Box) {
	if len(boxes) == 0 {
		fmt.Println("No boxes.")
	} else {
		var table = tablewriter.NewWriter(os.Stdout)
		const timefmt = `Mon Jan  _2 15:04 MST 2006`

		table.SetAutoFormatHeaders(false)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(true)

		table.SetHeader([]string{
			"Name", "ID", "Owner", "Visibility", "Created", "Updated",
		})

		for _, b := range boxes {
			table.Append([]string{
				b.Name,
				b.ID.String(),
				b.Owner,
				b.Visibility.String(),
				humanize.Time(b.Created.Local()), humanize.Time(b.Updated.Local()),
			})
		}
		table.Render()
	}
}

func init() {
	cmdBoxes.AddCommand(boxList, boxStack, boxVersions, boxDiff, boxBindings, boxDelete)
	Root.AddCommand(cmdBoxes)
}
