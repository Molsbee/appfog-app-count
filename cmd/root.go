package cmd

import (
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/molsbee/appfog-app-count/service"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	endpoints              = map[string]string{"useast": "https://api.useast.appfog.ctl.io", "uswest": "https://api.uswest.appfog.ctl.io"}
	username               string
	organizations          string
	organization           string
	deleteExternalServices bool
)

func init() {
	rootCommand.AddCommand(generateReportCommand, listApps, deleteApps)
	rootCommand.PersistentFlags().StringVarP(&username, "username", "u", "", "required - admin account")
	deleteApps.PersistentFlags().StringVarP(&organizations, "organizations", "o", "", "list of orphaned organizations where accounts were terminated")
	deleteApps.PersistentFlags().BoolVarP(&deleteExternalServices, "deleteExternalServices", "d", false, "deleted external services provisioned through marketplace like (RDBS)")
	listApps.PersistentFlags().StringVarP(&organization, "organization", "o", "", "appfog organization name")
}

var rootCommand = &cobra.Command{
	Use:   "appfog",
	Short: "automate appfog things",
	Long:  "Internal utility for generating reports and performing specific actions against AppFog.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var generateReportCommand = &cobra.Command{
	Use:     "generate-app-count-report",
	Short:   "Generates csv reports of all applications currently deployed in each region.",
	Long:    "Generates csv reports of all applications currently deployed in each region.",
	Example: "generate-app-count-report --username bah.t3n",
	Run: func(cmd *cobra.Command, args []string) {
		if username == "" {
			cmd.Usage()
			os.Exit(1)
		}

		fmt.Printf("Password: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			fmt.Println("unable to capture password properly")
			os.Exit(1)
		}

		for region, endpoint := range endpoints {
			if err := service.GenerateReport(username, string(pass), region, endpoint); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}
	},
}

var listApps = &cobra.Command{
	Use:     "list-apps",
	Short:   "Returns a list of apps and states.",
	Example: "list-apps --username bah.t3n --organization nbri",
	Run: func(cmd *cobra.Command, args []string) {
		if username == "" || organization == "" {
			cmd.Usage()
			os.Exit(1)
		}

		fmt.Printf("Password: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			fmt.Println("unable to capture password properly")
			os.Exit(1)
		}

		for region, endpoint := range endpoints {
			fmt.Printf("Collecting organization %s applications data for region %s\n", organization, region)
			client, err := service.NewClient(endpoint, username, string(pass))
			if err != nil {
				fmt.Println(fmt.Sprintf("unable to login with cloud foundry api at endpoint %s - error %s", endpoint, err.Error()))
				os.Exit(1)
			}

			org, err := client.GetOrgByName(organization)
			apps, _ := client.ListAppsByOrgGuid(org.Guid)
			for _, app := range apps {
				fmt.Printf("Name: %-40s State: %-10s Guid: %s\n", app.Name, app.State, app.Guid)
			}
		}
	},
}

var deleteApps = &cobra.Command{
	Use:     "delete-apps",
	Short:   "Stops/Deletes all apps associated with customer organizations.",
	Long:    "Takes a list of customer organizations and iterates through them to stop them.",
	Example: "delete-apps --username bah.t3n --organizations cw07,nbri",
	Run: func(cmd *cobra.Command, args []string) {
		if username == "" || len(organizations) == 0 {
			cmd.Usage()
			os.Exit(1)
		}

		fmt.Printf("Password: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			fmt.Println("unable to capture password properly")
			os.Exit(1)
		}

		orgs := strings.Split(organizations, ",")
		for _, endpoint := range endpoints {
			client, err := service.NewClient(endpoint, username, string(pass))
			if err != nil {
				fmt.Println(fmt.Sprintf("unable to login with cloud foundry api at endpoint %s - error %s", endpoint, err.Error()))
				os.Exit(1)
			}

			for _, orgName := range orgs {
				org, _ := client.GetOrgByName(orgName)
				if len(org.Name) != 0 {
					apps, _ := client.ListAppsByOrgGuid(org.Guid)
					for _, app := range apps {
						if app.State == "STARTED" {
							services, _ := client.ListAppServiceBindings(app.Guid)
							for _, service := range services {
								fmt.Printf("deleting service guid: %s\n", service.Guid)
								if err := client.DeleteServiceBinding(service.Guid); err != nil {
									fmt.Printf("error deleting service binding guid: %s error: %s\n", service.Guid, err.Error())
								}
							}

							routes, _ := client.GetAppRoutes(app.Guid)
							for _, route := range routes {
								fmt.Printf("deleting route - host name: %s\n", route.Host)
								if err := client.DeleteRoute(route.Guid); err != nil {
									fmt.Printf("error deleting route guid: %s error: %s\n", route.Guid, err.Error())
								}
							}

							fmt.Printf("deleting application guid: %s name: %s\n", app.Guid, app.Name)
							if err := client.DeleteApp(app.Guid); err != nil {
								fmt.Printf("error deleting app guid: %s error: %s", app.Guid, err.Error())
							}
						}
					}

					if deleteExternalServices {
						client.DeleteServiceInstancesByOrgGuid(org.Guid)
					}
				}
			}
		}
	},
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
