package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	host     string
	port     int
	username string
	password string
	vhost    string
)

var rootCmd = &cobra.Command{
	Use:   "amqp-cli",
	Short: "A CLI tool for RabbitMQ messaging",
	Long: `amqp-cli is a command line tool for publishing and consuming
messages to/from RabbitMQ queues and exchanges.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "localhost", "RabbitMQ host")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "P", 5672, "RabbitMQ port")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "guest", "RabbitMQ username")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "guest", "RabbitMQ password")
	rootCmd.PersistentFlags().StringVarP(&vhost, "vhost", "v", "", "RabbitMQ virtual host")
}
