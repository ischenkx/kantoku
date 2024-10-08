package cli

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/common/data/uid"
	"github.com/ischenkx/kantoku/pkg/common/logging/prefixed"
	"github.com/ischenkx/kantoku/pkg/common/transport/broker"
	"github.com/ischenkx/kantoku/pkg/core"
	"github.com/ischenkx/kantoku/pkg/lib/builder"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"strings"
)

func NewEvents() *cobra.Command {
	var eventsCmd = &cobra.Command{
		Use:   "events",
		Short: "Manage events",
	}

	var sendFlags struct {
		name    string
		payload string
		config  string
	}
	var sendCmd = &cobra.Command{
		Use:   "send",
		Short: "Send an event",
		Run: func(cmd *cobra.Command, args []string) {
			if sendFlags.name == "" {
				cmd.PrintErrln("name is empty")
				return
			}

			logger := slog.New(
				prefixed.NewHandler(
					slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}),
					&prefixed.HandlerOptions{
						PrefixKeys:      nil,
						PrefixFormatter: nil,
					},
				),
			)

			var cfg builder.Config

			if sendFlags.config != "" {
				var err error
				cfg, err = builder.FromFile(sendFlags.config)
				if err != nil {
					cmd.PrintErrln("failed to parse config from file:", err)
					return
				}
			} else {
				cmd.PrintErrln("no config provided (please, use --config)")
				return
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			events, err := builder.BuildEvents(ctx, logger, cfg.Core.System.Events)
			if err != nil {
				cmd.PrintErrln("failed to create events:", err)
				return
			}

			ev := core.NewEvent(sendFlags.name, []byte(sendFlags.payload))

			cmd.Println("ID:", ev.ID)
			cmd.Println("Name:", ev.Topic)
			cmd.Println("Data:", string(ev.Data))

			if err := events.Send(ctx, ev); err != nil {
				cmd.PrintErrln("failed to send:", err)
				return
			}
			cmd.Println("Success")
		},
	}
	sendCmd.Flags().StringVar(&sendFlags.name, "name", "", "name")
	sendCmd.Flags().StringVar(&sendFlags.payload, "data", "", "payload")
	sendCmd.Flags().StringVar(&sendFlags.config, "config", "config.yaml", "config path")

	var consumeFlags struct {
		names  []string
		group  string
		config string
	}
	var consumeCmd = &cobra.Command{
		Use:   "consume",
		Short: "Consume events",
		Run: func(cmd *cobra.Command, args []string) {
			if len(consumeFlags.names) == 0 {
				cmd.PrintErrln("no events provided")
				return
			}

			logger := slog.New(
				prefixed.NewHandler(
					slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}),
					&prefixed.HandlerOptions{
						PrefixKeys:      nil,
						PrefixFormatter: nil,
					},
				),
			)

			var cfg builder.Config

			if consumeFlags.config != "" {
				var err error
				cfg, err = builder.FromFile(consumeFlags.config)
				if err != nil {
					cmd.PrintErrln("failed to parse config from file:", err)
					return
				}
			} else {
				cmd.PrintErrln("no config provided (please, use --config)")
				return
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			events, err := builder.BuildEvents(ctx, logger, cfg.Core.System.Events)
			if err != nil {
				cmd.PrintErrln("failed to create events:", err)
				return
			}

			channel, err := events.Consume(ctx,
				consumeFlags.names,
				broker.ConsumerSettings{Group: consumeFlags.group},
			)
			if err != nil {
				cmd.PrintErrln("failed to consume:", err)
				return
			}

			cmd.Println("Group:", consumeFlags.group)
			cmd.Println("Events:", strings.Join(consumeFlags.names, ", "))

			broker.Processor[core.Event]{
				Handler: func(ctx context.Context, ev core.Event) error {
					cmd.Printf("id='%s' event='%s' data='%s'\n", ev.ID, ev.Topic, string(ev.Data))
					return nil
				},
			}.Process(ctx, channel)
		},
	}
	consumeCmd.Flags().StringArrayVar(&consumeFlags.names, "event", nil, "name")
	consumeCmd.Flags().StringVar(&consumeFlags.group, "group", uid.Generate(), "payload")
	consumeCmd.Flags().StringVar(&consumeFlags.config, "config", "config.yaml", "config path")

	eventsCmd.AddCommand(sendCmd)
	eventsCmd.AddCommand(consumeCmd)

	return eventsCmd
}
