package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/rwese/archivar/archivar"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	var (
		configFile         string
		debugging          bool
		quiet              bool
		profiler           bool
		defaultJobInterval int
		serviceConfig      archivar.Config
	)
	profilerPort := 6060
	logger := logrus.New()

	var cmdWatch = &cobra.Command{
		Use:   "watch",
		Short: "Start gathering and archiving until stopped",
		Long: `Starts the monitoring process processing all gathers 
and running the archivers until receiving an interrupt signal.`,
		Run: func(cmd *cobra.Command, args []string) {
			s := archivar.New(serviceConfig, logger)
			logger.Debugf("running watch with default interval: %d", defaultJobInterval)
			s.RunJobs(defaultJobInterval)
		},
	}

	cmdWatch.Flags().IntVarP(&defaultJobInterval, "interval", "i", defaultJobInterval, "default wait time between processing of all configured archivers, can be overriden by specifying it per job")

	var rootCmd = &cobra.Command{
		Use: "app",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			viper.SetConfigName(configFile)
			viper.SetConfigType("yaml")
			viper.AddConfigPath(".")
			viper.AddConfigPath("/etc/go-archivar/")
			err := viper.ReadInConfig()
			if err != nil {
				panic(fmt.Errorf("fatal error config file: %s", err))
			}

			err = viper.Unmarshal(&serviceConfig)
			if err != nil {
				panic(fmt.Errorf("fatal error config file: %s", err))
			}

			if debugging || serviceConfig.Settings.Log.Debugging {
				logger.SetLevel(logrus.DebugLevel)
			}

			if quiet {
				logger.SetLevel(logrus.ErrorLevel)
			}

			if defaultJobInterval == 0 {
				defaultJobInterval = serviceConfig.Settings.DefaultInterval
			}

			if profiler {
				go func() {
					listenHostPort := "0.0.0.0:" + fmt.Sprintf("%d", profilerPort)
					logger.Warnln("Run profiler", listenHostPort)
					logger.Warnln(http.ListenAndServe(listenHostPort, nil))
				}()
			}
		},
	}

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "archivar.yaml", "Configfile")
	rootCmd.PersistentFlags().BoolVarP(&debugging, "debug", "d", false, "enable verbose logging output")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress all non-error output")
	rootCmd.PersistentFlags().BoolVar(&profiler, "profiler", false, "run go profiler server")
	rootCmd.PersistentFlags().IntVar(&profilerPort, "profilerPort", profilerPort, "run go profiler server")
	rootCmd.AddCommand(cmdWatch)
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
