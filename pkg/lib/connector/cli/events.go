package cli

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/common/data/uid"
	"github.com/ischenkx/kantoku/pkg/common/transport/broker"
	"github.com/ischenkx/kantoku/pkg/common/transport/queue"
	"github.com/ischenkx/kantoku/pkg/core/event"
	"github.com/ischenkx/kantoku/pkg/lib/connector/cli/builder"
	config2 "github.com/ischenkx/kantoku/pkg/lib/connector/cli/config"
	"github.com/spf13/cobra"
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

			var cfg config2.Config

			if sendFlags.config != "" {
				var err error
				cfg, err = config2.FromFile(sendFlags.config)
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

			b := &builder.Builder{}
			events, err := b.BuildEvents(ctx, cfg.System.Events)
			if err != nil {
				cmd.PrintErrln("failed to create events:", err)
				return
			}

			ev := event.New(sendFlags.name, []byte(sendFlags.payload))

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

			var cfg config2.Config

			if consumeFlags.config != "" {
				var err error
				cfg, err = config2.FromFile(consumeFlags.config)
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

			b := &builder.Builder{}
			events, err := b.BuildEvents(ctx, cfg.System.Events)
			if err != nil {
				cmd.PrintErrln("failed to create events:", err)
				return
			}

			channel, err := events.Consume(ctx, broker.TopicsInfo{
				Group:  consumeFlags.group,
				Topics: consumeFlags.names,
			})
			if err != nil {
				cmd.PrintErrln("failed to consume:", err)
				return
			}

			cmd.Println("Group:", consumeFlags.group)
			cmd.Println("Events:", strings.Join(consumeFlags.names, ", "))

			queue.Processor[event.Event]{
				Handler: func(ctx context.Context, ev event.Event) error {
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
