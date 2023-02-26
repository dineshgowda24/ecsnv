package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/dineshgowda24/ecsnv/ecs"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ecsnv",
	Short: "Load AWS ECS envs locally",
	Long: `Application lets you download your AWS ECS envs locally.
The envs can be downloaded into a file or can be exported in the current terminal
session.`,
	Run: func(cmd *cobra.Command, args []string) {
		cluster, err := cmd.Flags().GetString("cluster")
		if err != nil {
			log.Fatal(err)
		}

		service, err := cmd.Flags().GetString("service")
		if err != nil {
			log.Fatal(err)
		}

		if len(service) > 0 && len(cluster) < 1 {
			fmt.Println(`Error: required flag(s) "cluster" not set. Service should be paired with cluster.`)
			cmd.Help()
			return
		}

		file, err := cmd.Flags().GetString("file")
		if err != nil {
			log.Fatal(err)
		}

		profile, err := cmd.Flags().GetString("profile")
		if err != nil {
			log.Fatal(err)
		}

		ecs.Run(cluster, service, file, profile)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("cluster", "c", "", "ECS cluster name")
	rootCmd.PersistentFlags().StringP("service", "s", "", "ECS service name, service can be paired with cluster")
	rootCmd.PersistentFlags().StringP("file", "f", "", "File to export")
	rootCmd.PersistentFlags().StringP("profile", "p", "", "AWS profile(overrides the default profile set in terminal)")
}
