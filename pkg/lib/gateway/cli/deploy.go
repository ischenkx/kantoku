package cli

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/common/logging/prefixed"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/lib/platform"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
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

			var cfg platform.Config

			logger := slog.New(
				prefixed.NewHandler(
					slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}),
					&prefixed.HandlerOptions{
						PrefixKeys:      nil,
						PrefixFormatter: nil,
					},
				),
			)

			if flags.config != "" {
				var err error
				cfg, err = platform.FromFile(flags.config)
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

			cmd.Println("building: system")
			sys, err := platform.BuildSystem(ctx, logger, cfg.Core.System)
			if err != nil {
				cmd.PrintErrln("failed to create a system instance from config:", err)
				return
			}

			cmd.Println("building: system")
			specifications, err := platform.BuildSpecifications(ctx, cfg.Core.Specifications)
			if err != nil {
				cmd.PrintErrln("failed to build a specifications manager:", err)
				return
			}

			var deployer service.Deployer

			if flags.scheduler {
				cmd.Println("building: scheduler")
				deployment, err := platform.BuildSchedulerDeployment(ctx, sys, logger, cfg.Services.Scheduler)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				deployer.Add(deployment.Service, deployment.Middlewares...)
			}

			if flags.processor {
				cmd.Println("building: processor")
				// TODO: add processor!
				deployment, err := platform.BuildProcessorDeployment(ctx, sys, nil, logger, cfg.Services.Processor)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				deployer.Add(deployment.Service, deployment.Middlewares...)
			}

			if flags.status {
				cmd.Println("building: status")

				deployment, err := platform.BuildStatusDeployment(ctx, sys, logger, cfg.Services.Status)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				deployer.Add(deployment.Service, deployment.Middlewares...)
			}

			if flags.api {
				cmd.Println("building: api")

				deployment, err := platform.BuildHttpApiDeployment(ctx, sys, specifications, logger, cfg.Services.HttpApi)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				deployer.Add(deployment.Service, deployment.Middlewares...)
			}

			if flags.serviceDiscovery {
				cmd.Println("building service discovery")

				deployemnt, err := platform.BuildDiscoveryDeployment(ctx, sys, logger, cfg.Services.Discovery)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				deployer.Add(deployemnt.Service, deployemnt.Middlewares...)
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
