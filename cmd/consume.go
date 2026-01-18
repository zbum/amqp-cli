package cmd

import (
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"amqp-cli/internal/rabbitmq"

	"github.com/spf13/cobra"
)

var (
	consumeQueue string
	autoAck      bool
	count        int
	verbose      bool
	hexDump      bool
)

var consumeCmd = &cobra.Command{
	Use:   "consume",
	Short: "Consume messages from a RabbitMQ queue",
	Long: `Consume messages from a RabbitMQ queue.

Examples:
  # Consume messages from a queue (continuous)
  amqp-cli consume -q myqueue

  # Consume with auto-acknowledge
  amqp-cli consume -q myqueue --auto-ack

  # Consume only N messages
  amqp-cli consume -q myqueue -n 10

  # Verbose mode (show Method/Header/Body frames)
  amqp-cli consume -q myqueue -V

  # Hex dump mode (show body as hex)
  amqp-cli consume -q myqueue --hex`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if consumeQueue == "" {
			return fmt.Errorf("--queue is required")
		}

		cfg := &rabbitmq.Config{
			Host:     host,
			Port:     port,
			Username: username,
			Password: password,
			Vhost:    vhost,
		}

		client, err := rabbitmq.NewClient(cfg)
		if err != nil {
			return err
		}
		defer client.Close()

		fmt.Printf("Consuming from queue '%s'... (Press Ctrl+C to stop)\n", consumeQueue)

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		msgCount := 0
		done := make(chan bool)

		go func() {
			err = client.Consume(consumeQueue, autoAck, func(msg rabbitmq.Message) error {
				msgCount++
				fmt.Printf("\n%s Message #%d %s\n", "===", msgCount, "===")

				if verbose {
					printVerbose(msg)
				} else {
					printSimple(msg)
				}

				if hexDump {
					printHexDump(msg.RawBody)
				}

				if count > 0 && msgCount >= count {
					done <- true
				}
				return nil
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error consuming: %v\n", err)
			}
			done <- true
		}()

		select {
		case <-sigChan:
			fmt.Printf("\nReceived %d message(s)\n", msgCount)
		case <-done:
			fmt.Printf("\nReceived %d message(s)\n", msgCount)
		}

		return nil
	},
}

func printSimple(msg rabbitmq.Message) {
	if msg.Exchange != "" {
		fmt.Printf("Exchange: %s\n", msg.Exchange)
	}
	if msg.RoutingKey != "" {
		fmt.Printf("Routing Key: %s\n", msg.RoutingKey)
	}
	if !msg.Timestamp.IsZero() {
		fmt.Printf("Timestamp: %s\n", msg.Timestamp.Format("2006-01-02 15:04:05"))
	}
	fmt.Printf("Body:\n%s\n", msg.Body)
}

func printVerbose(msg rabbitmq.Message) {
	// Method Frame (Basic.Deliver)
	fmt.Println("\n[Method Frame] Basic.Deliver")
	fmt.Printf("  ConsumerTag:  %s\n", msg.ConsumerTag)
	fmt.Printf("  DeliveryTag:  %d\n", msg.DeliveryTag)
	fmt.Printf("  Redelivered:  %t\n", msg.Redelivered)
	fmt.Printf("  Exchange:     %s\n", msg.Exchange)
	fmt.Printf("  RoutingKey:   %s\n", msg.RoutingKey)

	// Header Frame (Content Header)
	fmt.Println("\n[Header Frame] Content Header")
	fmt.Printf("  ContentType:     %s\n", defaultStr(msg.ContentType, "(not set)"))
	fmt.Printf("  ContentEncoding: %s\n", defaultStr(msg.ContentEncoding, "(not set)"))
	fmt.Printf("  DeliveryMode:    %s\n", deliveryModeStr(msg.DeliveryMode))
	fmt.Printf("  Priority:        %d\n", msg.Priority)
	if msg.CorrelationId != "" {
		fmt.Printf("  CorrelationId:   %s\n", msg.CorrelationId)
	}
	if msg.ReplyTo != "" {
		fmt.Printf("  ReplyTo:         %s\n", msg.ReplyTo)
	}
	if msg.Expiration != "" {
		fmt.Printf("  Expiration:      %s\n", msg.Expiration)
	}
	if msg.MessageId != "" {
		fmt.Printf("  MessageId:       %s\n", msg.MessageId)
	}
	if !msg.Timestamp.IsZero() {
		fmt.Printf("  Timestamp:       %s\n", msg.Timestamp.Format("2006-01-02 15:04:05"))
	}
	if msg.Type != "" {
		fmt.Printf("  Type:            %s\n", msg.Type)
	}
	if msg.UserId != "" {
		fmt.Printf("  UserId:          %s\n", msg.UserId)
	}
	if msg.AppId != "" {
		fmt.Printf("  AppId:           %s\n", msg.AppId)
	}
	if len(msg.Headers) > 0 {
		fmt.Println("  Headers:")
		for k, v := range msg.Headers {
			fmt.Printf("    %s: %v\n", k, v)
		}
	}

	// Body Frame
	fmt.Println("\n[Body Frame] Content Body")
	fmt.Printf("  Size: %d bytes\n", len(msg.RawBody))
	fmt.Printf("  Data:\n%s\n", msg.Body)
}

func printHexDump(data []byte) {
	fmt.Println("\n[Hex Dump]")
	fmt.Print(hex.Dump(data))
}

func defaultStr(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

func deliveryModeStr(mode uint8) string {
	switch mode {
	case 1:
		return "1 (Non-persistent)"
	case 2:
		return "2 (Persistent)"
	default:
		return fmt.Sprintf("%d (Unknown)", mode)
	}
}

func init() {
	rootCmd.AddCommand(consumeCmd)

	consumeCmd.Flags().StringVarP(&consumeQueue, "queue", "q", "", "Queue name to consume from")
	consumeCmd.Flags().BoolVar(&autoAck, "auto-ack", false, "Auto-acknowledge messages")
	consumeCmd.Flags().IntVarP(&count, "count", "n", 0, "Number of messages to consume (0 = unlimited)")
	consumeCmd.Flags().BoolVarP(&verbose, "verbose", "V", false, "Show detailed Method/Header/Body frame info")
	consumeCmd.Flags().BoolVar(&hexDump, "hex", false, "Show body as hex dump")

	consumeCmd.MarkFlagRequired("queue")
}
