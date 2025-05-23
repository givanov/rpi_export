package cmd

import (
	"fmt"
	"github.com/givanov/rpi_export/pkg/config"
	"github.com/givanov/rpi_export/pkg/mbox"
	"github.com/givanov/rpi_export/pkg/metrics"
	"github.com/givanov/rpi_export/pkg/server"
	"github.com/givanov/rpi_export/pkg/signals"
	"github.com/givanov/rpi_export/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

var c config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rpi_export",
	Short: "A tiny webserver that collects raspberry pi hardware metrics from /dev/vcio and exposes them as prometheus metrics",
	Run: func(cmd *cobra.Command, args []string) {
		err := config.ValidateConfig(&c)
		if err != nil {
			panic("Error validating config")
		}

		stopCtx := signals.SetupSignalHandlerCtx()

		mbox.Debug = c.MailboxDebug

		collector := metrics.NewRaspberryPiMboxCollector(&c)
		prometheus.MustRegister(collector)

		srv := server.New(&c)

		serverWg := new(sync.WaitGroup)
		serverWg.Add(1)

		go func() {
			defer func() {
				serverWg.Done()
			}()

			srv.Start()
		}()

		select {
		case <-stopCtx.Done():
			zap.L().Error("got stop signal, shutting down...", zap.Error(srv.Stop()))
		}
		serverWg.Wait()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	config.InitLogger(&c)

	rootCmd.Flags().StringVar(&c.LogLevel, "v", config.DefaultLogLevel, "Log level: all, debug, info, warn, error, panic, fatal, none")
	rootCmd.Flags().IntVarP(&c.Port, "port", "p", config.DefaultPort, "Webserver port")
	rootCmd.Flags().StringVar(&c.BindInterface, "bind-interface", config.DefaultBindInterface, "Webserver interface to bind to")
	rootCmd.Flags().BoolVar(&c.MailboxDebug, "mailbox-debug", false, "Whether to enable debug logging for all mailbox actions")

	hostname, err := os.Hostname()
	if err != nil {
		zap.L().Error("Error getting hostname. Using unknown as default instead", zap.Error(err))
		hostname = "unknown"
	}

	rootCmd.Flags().StringVar(&c.HostNameOverride, "hostname-override", hostname, "Override the hostname used in all metrics")

	if err := util.SetFlagsFromEnv(rootCmd.Flags(), ""); err != nil {
		panic(fmt.Sprintf("error setting flags from environment variables. %s ", err.Error()))
	}
}
