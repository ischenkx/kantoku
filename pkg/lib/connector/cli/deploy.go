package cli

import (
	"github.com/ischenkx/kantoku/pkg/common/data/uid"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/core/services/scheduler/dependencies/simple"
	"github.com/ischenkx/kantoku/pkg/core/services/scheduler/dummy"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/spf13/cobra"
	"log/slog"
)

type deployFlags struct {
	env                bool
	noScheduler        bool
	noProcessor        bool
	noStatus           bool
	noApi              bool
	noServiceDiscovery bool
	config             string
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

			var config Config

			if flags.env {
				var err error
				config, err = configFromEnv()
				if err != nil {
					cmd.PrintErrln("failed to parse config from env:", err)
					return
				}
			} else if flags.config != "" {
				var err error
				config, err = configFromFile(flags.config)
				if err != nil {
					cmd.PrintErrln("failed to parse config from file:", err)
					return
				}
			} else {
				cmd.PrintErrln("no config provided (please use environment variables / config file)")
				return
			}

			sys, err := systemFromConfig(config.System)
			if err != nil {
				cmd.PrintErrln("failed to create a system instance from config:", err)
				return
			}

			var deployer service.Deployer

			if flags.scheduler {
				core := service.NewCore(
					"dependencies",
					uid.Generate(),
					slog.Default(),
				)

				var srvc service.Service

				switch config.Services.Scheduler.Type {
				case "dependencies":
					srvc = &simple.Service{
						System:  sys,
						Manager: nil,
						Core:,
					}
				case "dummy":
				default:
					cmd.PrintErrf("unsupported scheduler type: %s\n", config.Services.Scheduler.Type)
					return
				}
			}

			if flags.processor {

			}

			if flags.status {

			}

			if flags.api {

			}

			if flags.serviceDiscovery {

			}

		},
	}

	cmd.Flags().BoolVar(&flags.env, "env", false, "Set environment")
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

type serviceBuilder struct {
	system system.AbstractSystem
}

func (builder serviceBuilder) buildScheduler(config SchedulerConfig) (srvc service.Service, mws []service.Middleware, err error) {
	core := builder.buildServiceCore("scheduler", config.Options)

	switch config.Type {
	case "dependencies":
	case "dummy":
		srvc = &dummy.Service{
			System: builder.system,
			Core:   core,
		}

	}


}

func (builder serviceBuilder) buildServiceCore(defaultName string, options map[string]any) service.Core {

}

func (builder serviceBuilder) buildMiddlewares(options map[string]any) []service.Middleware {

}