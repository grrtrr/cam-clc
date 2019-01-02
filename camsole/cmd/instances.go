package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/grrtrr/clccam"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	// Top-level instance command
	cmdInstances = &cobra.Command{
		Use:     "instance",
		Aliases: []string{"i"},
		Short:   "Manage instances",
	}

	// List one or more instances (VMs)
	instanceGet = &cobra.Command{
		Use:     "ls  [<instanceId> ...]",
		Aliases: []string{"list", "show"},
		Short:   "List CAM instance(s)",
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
						die("Failed to query instance %s: %s", instanceId, err)
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

	// Get the service of an instance
	instanceGetService = &cobra.Command{
		Use:     "service  <instanceId>",
		Aliases: []string{"s", "srv"},
		Short:   "List the instance service",
		PreRunE: checkArgs(1, "Need an instance ID"),
		Run: func(cmd *cobra.Command, args []string) {
			if srv, err := client.GetInstanceService(args[0]); err != nil {
				die("failed to query instance %s service: %s", args[0], err)
			} else if cmd.Flags().Lookup("json").Value.String() != "true" {
				fmt.Printf("%s/%s %s service %s, operation %s/%s:\n",
					srv.ClcAlias, srv.Organization, srv.Type, srv.ID, srv.Operation, srv.State)

				// 1. State history
				if len(srv.StateHistory) > 0 {
					var table = tablewriter.NewWriter(os.Stdout)

					table.SetAutoFormatHeaders(false)
					table.SetAlignment(tablewriter.ALIGN_LEFT)
					table.SetAutoWrapText(true)

					fmt.Printf("\n%s states:\n", srv.ID)

					sort.Slice(srv.StateHistory, func(i, j int) bool {
						return srv.StateHistory[i].Started.Time.Before(srv.StateHistory[j].Completed.Time)
					})

					table.SetHeader([]string{"State", "Operation", "Started", "Completed"})
					for _, s := range srv.StateHistory {
						var completed = "n/a"

						if !s.Completed.Time.IsZero() {
							completed = fmt.Sprintf("%s after start",
								s.Completed.Time.Sub(s.Started.Time).Round(time.Second))
						}
						table.Append([]string{
							s.State,
							srv.Operation.String(),
							humanize.Time(s.Started.Time.Local()),
							completed,
						})
					}
					table.Render()
				} else {
					fmt.Printf("\n%s: no state history.\n", srv.ID)
				}

				// 2. Machines
				if len(srv.Machines) > 0 {
					var table = tablewriter.NewWriter(os.Stdout)

					table.SetAutoFormatHeaders(false)
					table.SetAlignment(tablewriter.ALIGN_LEFT)
					table.SetAutoWrapText(true)

					fmt.Printf("\n%s (Virtual) Machines:\n", srv.ID)

					table.SetHeader([]string{"Host", "Provider ID", "IP", "State", "Agent Ping", "Agent Close"})
					for _, m := range srv.Machines {
						var (
							lastClose = "n/a"
							host      = m.Name
						)

						if t := m.LastAgentClose.Time.Local(); !t.IsZero() {
							lastClose = humanize.Time(t)
						}
						if h := m.Hostname; h != "" {
							host = h
						}

						table.Append([]string{
							host,
							m.ExternalID,
							m.Address.String(),
							m.State.String(),
							humanize.Time(m.LastAgentPing.Time.Local()),
							lastClose,
						})
					}
					table.Render()
				} else {
					fmt.Printf("\n%s: no VMs.\n", srv.ID)
				}
			}
		},
	}

	// Retrieve instance activities
	instanceGetActivity = &cobra.Command{
		Use:     "activity  <instanceId>",
		Aliases: []string{"act"},
		Short:   "Retrieve instance activity logs",
		PreRunE: checkArgs(1, "Need an instance ID"),
		Run: func(cmd *cobra.Command, args []string) {
			var filterByCmd string

			if op, err := cmd.Flags().GetString("op"); err == nil {
				filterByCmd = op
			}

			if activities, err := client.GetInstanceActivity(args[0], filterByCmd); err != nil {
				die("failed to query instance %s activities: %s", args[0], err)
			} else if cmd.Flags().Lookup("json").Value.String() == "true" {
			} else if len(activities) == 0 {
				fmt.Printf("No %s activities reported.\n", args[0])
			} else {
				fmt.Printf("%s activities:\n", args[0])
				printActivities(activities)
			}
		},
	}

	// Retrieve instance operations (includes activities).
	instanceGetOps = &cobra.Command{
		Use:     "ops  <instanceId>",
		Aliases: []string{"operations", "op"},
		Short:   "List instance operations",
		PreRunE: checkArgs(1, "Need an instance ID"),
		Run: func(cmd *cobra.Command, args []string) {
			if ops, err := client.GetInstanceOperations(args[0]); err != nil {
				die("failed to query instance %s operations: %s", args[0], err)
			} else if cmd.Flags().Lookup("json").Value.String() == "true" {
			} else if len(ops) == 0 {
				fmt.Println("No information on operations.")
			} else {
				for _, op := range ops {
					var state = op.State.String()

					if op.InstanceState != op.State {
						state = fmt.Sprintf("%s (instance state %s)", op.State, op.InstanceState)
					}
					fmt.Printf("%s %s => %s on workspace %s:\n",
						op.Created.Time.Local().Format("_2 Jan 15:04 MST"),
						op.Operation, state, op.Workspace)
					printActivities(op.Activity)
				}
			}
		},
	}

	// Retrieve machine logs
	instanceGetLogs = &cobra.Command{
		Use:     "logs  <instanceId> [<machineId>]",
		Aliases: []string{"get-machine-logs", "log"},
		Short:   "Retrieve VM log output",
		PreRunE: checkAtLeastArgs(1, "Need an instance ID and optionally a machine ID"),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				instanceId = args[0]
				machine    string
			)
			if len(args) == 2 {
				machine = args[1]
			} else {
				if instance, err := client.GetInstance(instanceId); err != nil {
					die("failed to query %s machines: %s", instanceId, err)
				} else if m := instance.Service.Machines; len(m) == 0 {
					fmt.Println("No machines available")
				} else if len(m) > 1 {
					die("unable to retrieve logs: %s has more than 1 machine", instanceId)
				} else {
					machine = m[0].Name
				}
			}
			if logs, err := client.GetInstanceMachineLogs(instanceId, machine); err != nil {
				die("failed to query instance %s activities: %s", instanceId, err)
			} else if cmd.Flags().Lookup("json").Value.String() == "true" {
			} else if len(logs) == 0 {
				fmt.Println("No log information.")
			} else {
				fmt.Println(logs)
			}
		},
	}

	// Get instance bindings
	instanceGetBindings = &cobra.Command{
		Use:     "bindings  <instanceId>",
		Aliases: []string{"bdgs", "bind"},
		Short:   "Query instance bindings",
		PreRunE: checkArgs(1, "Need an instance ID"),
		Run: func(cmd *cobra.Command, args []string) {
			if bindings, err := client.GetInstanceBindings(args[0]); err != nil {
				die("failed to query instance %s bindings: %s", args[0], err)
			} else if cmd.Flags().Lookup("json").Value.String() == "true" {
			} else if len(bindings) == 0 {
				fmt.Println("No binding information.")
			} else {
				// FIXME: not tested this in practice yet
				fmt.Println(bindings)
			}
		},
	}

	// Re-deploy an existing instance
	instanceDeploy = &cobra.Command{
		Use:     "deploy  <instanceId> [<instanceId1> ...]",
		Aliases: []string{"launch"},
		Short:   "Re-deploy instance(s)",
		PreRunE: checkAtLeastArgs(1, "Need at least 1 instance to deploy"),
		Run: func(cmd *cobra.Command, args []string) {
			for _, instanceId := range args {
				if err := client.DeployInstance(instanceId); err != nil {
					die("Failed to re-deploy %s: %s\n", instanceId, err)
				}
			}
		},
	}

	// Power up an existing instance
	instancePowerOn = &cobra.Command{
		Use:     "on  <instanceId> [<instanceId1> ...]",
		Aliases: []string{"power-on"},
		Short:   "Power-on instance(s)",
		PreRunE: checkAtLeastArgs(1, "Need at least 1 instance to power on"),
		Run: func(cmd *cobra.Command, args []string) {
			for _, instanceId := range args {
				if err := client.PowerOnInstance(instanceId); err != nil {
					die("Failed to power-on %s: %s\n", instanceId, err)
				}
			}
		},
	}

	// Shut down an existing instance
	instanceShutdown = &cobra.Command{
		Use:     "off  <instanceId> [<instanceId1> ...]",
		Aliases: []string{"power-off", "shutdown", "down"},
		Short:   "Shut down instance(s)",
		PreRunE: checkAtLeastArgs(1, "Need at least 1 instance to shut down"),
		Run: func(cmd *cobra.Command, args []string) {
			for _, instanceId := range args {
				if err := client.ShutdownInstance(instanceId); err != nil {
					die("Failed to shut down %s: %s\n", instanceId, err)
				}
			}
		},
	}

	// Re-install an existing instance
	instanceReinstall = &cobra.Command{
		Use:     "reinst  <instanceId> [<instanceId1> ...]",
		Aliases: []string{"reinstall", "install"},
		Short:   "Re-install instance(s)",
		PreRunE: checkAtLeastArgs(1, "Need at least 1 instance to re-install"),
		Run: func(cmd *cobra.Command, args []string) {
			for _, instanceId := range args {
				if err := client.ReinstallInstance(instanceId); err != nil {
					die("Failed to re-install %s: %s\n", instanceId, err)
				}
			}
		},
	}

	// Re-configure an existing instance
	instanceReconfigure = &cobra.Command{
		Use:     "reconf  <instanceId> [<instanceId1> ...]",
		Aliases: []string{"reconfigure", "config"},
		Short:   "Reconfigure instance(s)",
		PreRunE: checkAtLeastArgs(1, "Need at least 1 instance to re-configure"),
		Run: func(cmd *cobra.Command, args []string) {
			for _, instanceId := range args {
				if err := client.ReconfigureInstance(instanceId); err != nil {
					die("Failed to re-configure %s: %s\n", instanceId, err)
				}
			}
		},
	}

	// Try to (re-)import an unregistered instance
	instanceImport = &cobra.Command{
		Use:     "import  <instanceId> [<instanceId1> ...]",
		Aliases: []string{"re-import"},
		Short:   "(Re-)import instance(s)",
		PreRunE: checkAtLeastArgs(1, "Need at least 1 instance to import"),
		Run: func(cmd *cobra.Command, args []string) {
			for _, instanceId := range args {
				if err := client.ImportInstance(instanceId); err != nil {
					die("Failed to import %s: %s\n", instanceId, err)
				}
			}
		},
	}

	// Cancel a failed import of an unregistered instance.
	instanceCancelImport = &cobra.Command{
		Use:     "import-stop  <instanceId> [<instanceId1> ...]",
		Aliases: []string{"stop-import", "cancel-import"},
		Short:   "Cancel instance import",
		PreRunE: checkAtLeastArgs(1, "Need at least 1 instance ID"),
		Run: func(cmd *cobra.Command, args []string) {
			for _, instanceId := range args {
				if err := client.CancelImportInstance(instanceId); err != nil {
					die("Failed to cancel import %s: %s\n", instanceId, err)
				}
			}
		},
	}

	// Delegates managment of an existing instance to CenturyLink.
	// https://www.ctl.io/api-docs/cam/#instances-delegate-management-of-an-existing-instance
	instanceMakeManaged = &cobra.Command{
		Use:     "delegate  <instanceId>",
		Aliases: []string{"make-managed", "managed"},
		Short:   "Make Century(Link) managed instance",
		PreRunE: checkArgs(1, "Need an instance ID"),
		Run: func(cmd *cobra.Command, args []string) {
			if err := client.MakeManagedInstance(args[0]); err != nil {
				die("Failed to make managed instance: %s\n", err)
			}
		},
	}

	// Terminate instance(s)
	instanceTerminate = &cobra.Command{
		Use:     "term  <instanceId> [<instanceId1> ...]",
		Aliases: []string{"terminate"},
		Short:   "Terminate instance(s)",
		PreRunE: checkAtLeastArgs(1, "Need at least 1 instance to terminate"),
		Run: func(cmd *cobra.Command, args []string) {
			var op = "terminate"

			if yes, err := cmd.Flags().GetBool("force"); err == nil && yes {
				op = "force_terminate"
			}

			for _, instanceId := range args {
				if err := client.DeleteInstance(instanceId, op); err != nil {
					die("Failed to %s %s: %s\n", strings.Replace(op, "_", "-", 1), instanceId, err)
				}
			}
		},
	}

	// Delete instance(s)
	instanceDelete = &cobra.Command{
		Use:     "rm  <instanceId> [<instanceId1> ...]",
		Aliases: []string{"remove", "delete"},
		Short:   "Delete instance(s)",
		PreRunE: checkAtLeastArgs(1, "Need at least 1 instance to delete"),
		Run: func(cmd *cobra.Command, args []string) {
			for _, instanceId := range args {
				if err := client.DeleteInstance(instanceId, "delete"); err != nil {
					die("Failed to delete %s: %s\n", instanceId, err)
				}
			}
		},
	}
)

func init() {
	// Flags
	instanceGetActivity.Flags().String("op", "", "Filter by operation (optional)")
	instanceTerminate.Flags().BoolP("force", "f", false, "Whether to force-terminate the instance")

	cmdInstances.AddCommand(instanceGet,
		instanceGetService, instanceGetActivity, instanceGetOps, instanceGetLogs, instanceGetBindings,
		instanceDeploy, instancePowerOn, instanceShutdown, instanceReinstall, instanceReconfigure,
		instanceImport, instanceCancelImport, instanceMakeManaged,
		instanceTerminate, instanceDelete,
	)
	Root.AddCommand(cmdInstances)
}

func printInstances(instances []clccam.Instance) {
	if len(instances) == 0 {
		fmt.Println("No instances.")
	} else {
		var table = tablewriter.NewWriter(os.Stdout)

		table.SetAutoFormatHeaders(false)
		table.SetAlignment(tablewriter.ALIGN_RIGHT)
		table.SetAutoWrapText(true)

		table.SetHeader([]string{
			"Name", "ID", "Machines", "Service", "Box", "Updated", "Operation", "State",
		})
		for _, i := range instances {
			var machines []string

			for _, m := range i.Service.Machines {
				machines = append(machines, fmt.Sprintf("%25.25s", m.Name))
			}
			table.Append([]string{
				fmt.Sprintf("%25.25s", i.Name),
				i.ID,
				strings.Join(machines, ", "),
				i.Service.ID,
				i.Box.String(),
				humanize.Time(i.Updated.Local()),
				i.Operation.Event.String(),
				i.State.String(),
			})
		}
		table.Render()
	}
}

// printActivities prints instance @activities.
func printActivities(activities []clccam.InstanceActivity) {
	if len(activities) == 0 {
		fmt.Printf("No %s activities reported.\n")
	} else {
		var table = tablewriter.NewWriter(os.Stdout)

		table.SetAutoFormatHeaders(false)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})
		table.SetAutoWrapText(false)

		table.SetHeader([]string{"Time", "Text", "Level", "Event"})
		for _, a := range activities {
			table.Append([]string{
				a.Created.Time.Local().Format("_2 Jan  15:04:05.0"),
				fmt.Sprintf("%-.100s", a.Text), // Chop of at 100 characters
				a.Level,
				a.Event.String(),
			})
		}
		table.Render()
	}
}
