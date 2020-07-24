package cmd

import (
	"github.com/xiayinchang/uk8sctl/config"
	"github.com/xiayinchang/uk8sctl/pkg/util"

	"github.com/spf13/cobra"
)

var globalConfig = &config.GlobalConfig{}

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uk8sctl",
		Short: "\ntool for uk8s",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			util.SetLogLevel(globalConfig.LogLevel)
		},
	}

	cmd.PersistentFlags().StringVar(&globalConfig.PublicKey, "publickey", "", "Required. ucloud publickey")
	cmd.PersistentFlags().StringVar(&globalConfig.PrivateKey, "privatekey", "", "Required. ucloud privatekey")
	cmd.PersistentFlags().StringVar(&globalConfig.Region, "region", "", "Required. The region choosed for creating resources")
	cmd.PersistentFlags().StringVar(&globalConfig.ProjectId, "project-id", "", "Required. Which project to use")
	cmd.MarkPersistentFlagRequired("publickey")
	cmd.MarkPersistentFlagRequired("privatekey")
	cmd.MarkPersistentFlagRequired("region")
	cmd.MarkPersistentFlagRequired("project-id")

	cmd.AddCommand(newCreateBaseUHostCmd())

	return cmd
}
