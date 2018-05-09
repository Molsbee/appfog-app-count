package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"github.com/howeyc/gopass"
	"github.com/molsbee/appfog-app-count/service"
)

var (
	endpoints = map[string]string{"useast": "https://api.useast.appfog.ctl.io", "uswest": "https://api.uswest.appfog.ctl.io"}
	username string
)


func init() {
	rootCommand.AddCommand(generateReportCommand, cleanupApps)
	rootCommand.PersistentFlags().StringVarP(&username, "username", "u", "", "required - admin account")
}

var rootCommand = &cobra.Command{
	Use: "appfog",
	Short: "automate appfog things",
	Long: "Internal utility for generating reports and performing specific actions against AppFog.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var generateReportCommand = &cobra.Command{
	Use: "generate-app-count-report",
	Short: "Generates csv reports of all applications currently deployed in each region.",
	Long: "Generates csv reports of all applications currently deployed in each region.",
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

var cleanupApps = &cobra.Command{
	Use: "cleanup",
	Short: "",
	Long: "",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

