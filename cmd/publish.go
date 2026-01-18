package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"amqp-cli/internal/rabbitmq"

	"github.com/spf13/cobra"
)

var (
	exchange   string
	routingKey string
	queue      string
	message    string
)

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish a message to RabbitMQ",
	Long: `Publish a message to a RabbitMQ queue or exchange.

Examples:
  # Publish to a queue
  amqp-cli publish -q myqueue -m "Hello World"

  # Publish to an exchange with routing key
  amqp-cli publish -e myexchange -r mykey -m "Hello World"

  # Read message from stdin
  echo "Hello World" | amqp-cli publish -q myqueue`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if queue == "" && exchange == "" {
			return fmt.Errorf("either --queue or --exchange must be specified")
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

		msg := message
		if msg == "" {
			msg, err = readFromStdin()
			if err != nil {
				return fmt.Errorf("failed to read from stdin: %w", err)
			}
		}

		if msg == "" {
			return fmt.Errorf("message cannot be empty")
		}

		if queue != "" {
			if err := client.PublishToQueue(queue, msg); err != nil {
				return fmt.Errorf("failed to publish to queue: %w", err)
			}
			fmt.Printf("Message published to queue '%s'\n", queue)
		} else {
			if err := client.Publish(exchange, routingKey, msg); err != nil {
				return fmt.Errorf("failed to publish to exchange: %w", err)
			}
			fmt.Printf("Message published to exchange '%s' with routing key '%s'\n", exchange, routingKey)
		}

		return nil
	},
}

func readFromStdin() (string, error) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", nil
	}

	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return strings.Join(lines, "\n"), nil
}

func init() {
	rootCmd.AddCommand(publishCmd)

	publishCmd.Flags().StringVarP(&exchange, "exchange", "e", "", "Exchange name")
	publishCmd.Flags().StringVarP(&routingKey, "routing-key", "r", "", "Routing key")
	publishCmd.Flags().StringVarP(&queue, "queue", "q", "", "Queue name (publishes directly to queue)")
	publishCmd.Flags().StringVarP(&message, "message", "m", "", "Message body")
}
