package cmd

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var config = &Config{}

func NewRootCommand() *cobra.Command {
	app := new(App)

	cmd := &cobra.Command{
		Use:   "exporter-merger",
		Short: "merges Prometheus metrics from multiple sources",
		Run:   app.run,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if app.viper.GetBool("verbose") {
				log.SetLevel(log.DebugLevel)
			} else {
				log.SetLevel(log.InfoLevel)
			}
		},
	}

	app.Bind(cmd)

	cmd.AddCommand(NewVersionCommand())

	return cmd
}

type App struct {
	viper *viper.Viper
}

func (app *App) Bind(cmd *cobra.Command) {
	app.viper = viper.New()
	app.viper.SetEnvPrefix("MERGER")
	app.viper.AutomaticEnv()

	configPath := cmd.PersistentFlags().StringP(
		"config-path", "c", "/etc/exporter-merger/config.yaml",
		"Path to the configuration file.")
	cobra.OnInitialize(func() {
		var err error
		if configPath != nil && *configPath != "" {
			config, err = ReadConfig(*configPath)
			if err != nil {
				log.WithField("error", err).Errorf("failed to load config file '%s'", *configPath)
				os.Exit(1)
				return
			}
		}
	})

	cmd.PersistentFlags().Int(
		"listen-port", 8080,
		"Listen port for the HTTP server. (ENV:MERGER_PORT)")
	app.viper.BindPFlag("port", cmd.PersistentFlags().Lookup("listen-port"))

	cmd.PersistentFlags().String(
		"listen-ip", "0.0.0.0",
		"Listen IP for the HTTP server.(ENV:MERGER_IP)")
	app.viper.BindPFlag("ip", cmd.PersistentFlags().Lookup("listen-ip"))

	cmd.PersistentFlags().Int(
		"exporters-timeout", 10,
		"HTTP client timeout for connecting to exporters. (ENV:MERGER_EXPORTERSTIMEOUT)")
	app.viper.BindPFlag("exporterstimeout", cmd.PersistentFlags().Lookup("exporters-timeout"))

	cmd.PersistentFlags().BoolP(
		"verbose", "v", false,
		"Include debug messages to output (ENV:MERGER_VERBOSE)")
	app.viper.BindPFlag("verbose", cmd.PersistentFlags().Lookup("verbose"))

}

func (app *App) run(cmd *cobra.Command, args []string) {
	http.Handle("/metrics", Handler{
		Exporters:            config.Exporters,
		ExportersHTTPTimeout: app.viper.GetInt("exporterstimeout"),
	})

	port := app.viper.GetInt("port")
	ip := app.viper.GetString("ip")
	log.Infof("starting HTTP server on %s:%d", ip, port)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
