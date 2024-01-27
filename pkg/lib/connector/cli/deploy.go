package cli

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/lib/connector/cli/builder"
	config2 "github.com/ischenkx/kantoku/pkg/lib/connector/cli/config"
	"github.com/spf13/cobra"
)

type deployFlags struct {
	config             string
	noScheduler        bool
	noProcessor        bool
	noStatus           bool
	noApi              bool
	noServiceDiscovery bool
	scheduler          bool
	processor          bool
	status             bool
	api                bool
	serviceDiscovery   bool
}

func NewDeploy() *cobra.Command {
	flags := &deployFlags{}
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy the application",
		Run: func(cmd *cobra.Command, args []string) {
			if !flags.scheduler && !flags.processor && !flags.status && !flags.api && !flags.serviceDiscovery {
				flags.scheduler = true
				flags.processor = true
				flags.status = true
				flags.api = true
				flags.serviceDiscovery = true
			}
			if flags.noScheduler {
				flags.scheduler = false
			}
			if flags.noProcessor {
				flags.processor = false
			}
			if flags.noStatus {
				flags.status = false
			}
			if flags.noApi {
				flags.api = false
			}
			if flags.noServiceDiscovery {
				flags.serviceDiscovery = false
			}

			var cfg config2.Config

			if flags.config != "" {
				var err error
				cfg, err = config2.FromFile(flags.config)
				if err != nil {
					cmd.PrintErrln("failed to parse config from file:", err)
					return
				}
			} else {
				cmd.PrintErrln("no config provided (please use environment variables / config file)")
				return
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			b := &builder.Builder{}

			cmd.Println("building system")
			sys, err := b.BuildSystem(ctx, cfg.System)
			if err != nil {
				cmd.PrintErrln("failed to create a system instance from config:", err)
				return
			}

			var deployer service.Deployer

			if flags.scheduler {
				cmd.Println("building scheduler")
				srvc, mws, err := b.BuildScheduler(ctx, sys, cfg.Services.Scheduler)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				deployer.Add(srvc, mws...)
			}

			if flags.processor {
				cmd.Println("building processor")
				srvc, mws, err := b.BuildProcessor(ctx, sys, cfg.Services.Processor)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				deployer.Add(srvc, mws...)
			}

			if flags.status {
				cmd.Println("building status")

				srvc, mws, err := b.BuildStatus(ctx, sys, cfg.Services.Status)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				deployer.Add(srvc, mws...)
			}

			if flags.api {
				cmd.Println("building api")

				srvc, mws, err := b.BuildHttpApi(ctx, sys, cfg.Services.HttpServer)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				deployer.Add(srvc, mws...)
			}

			if flags.serviceDiscovery {
				cmd.Println("building service discovery")

				srvc, mws, err := b.BuildDiscovery(ctx, sys, cfg.Services.Discovery)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				deployer.Add(srvc, mws...)
			}

			cmd.Println("deploying...")
			if err := deployer.Deploy(context.Background()); err != nil {
				cmd.PrintErrln(err)
				return
			}
			cmd.Println("DONE")
		},
	}

	cmd.Flags().StringVar(&flags.config, "config", "", "Specify config file path")
	cmd.Flags().BoolVar(&flags.noScheduler, "no-scheduler", false, "Disable scheduler")
	cmd.Flags().BoolVar(&flags.noProcessor, "no-processor", false, "Disable processor")
	cmd.Flags().BoolVar(&flags.noStatus, "no-status", false, "Disable status")
	cmd.Flags().BoolVar(&flags.noApi, "no-api", false, "Enable API")
	cmd.Flags().BoolVar(&flags.noServiceDiscovery, "no-service-discovery", false, "Enable API")
	cmd.Flags().BoolVar(&flags.scheduler, "scheduler", false, "Enable scheduler")
	cmd.Flags().BoolVar(&flags.processor, "processor", false, "Enable processor")
	cmd.Flags().BoolVar(&flags.status, "status", false, "Enable status")
	cmd.Flags().BoolVar(&flags.api, "api", false, "Enable API")
	cmd.Flags().BoolVar(&flags.serviceDiscovery, "service-discovery", false, "Enable API")

	return cmd
}
