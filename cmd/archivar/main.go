package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/rwese/archivar/archivar"
	_ "github.com/rwese/archivar/internal/imap"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	archivarSvc archivar.Archivar
	logger      = logrus.New()
)

const DEFAULT_PROFILER_PORT = 6060

func main() {
	var cmdWatch = &cobra.Command{
		Use:   "watch",
		Short: "Start gathering and archiving until stopped",
		Long: `Starts the monitoring process processing all gathers 
and running the archivers until receiving an interrupt signal.`,
		Run: func(cmd *cobra.Command, args []string) {
			defaultJobInterval, _ := cmd.Flags().GetInt("interval")
			logger.Debugf("running watch with default interval: %d", defaultJobInterval)
			archivarSvc.RunJobs(defaultJobInterval)
		},
	}

	cmdWatch.Flags().IntP("interval", "i", 60, "default wait time between processing of all configured archivers, can be overriden by specifying it per job")

	var rootCmd = &cobra.Command{
		Use: "app",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			configFile, _ := cmd.Flags().GetString("config")
			viper.SetConfigName(configFile)
			viper.SetConfigType("yaml")
			viper.AddConfigPath(".")
			viper.AddConfigPath("/etc/go-archivar/")
			err := viper.ReadInConfig()
			if err != nil {
				panic(fmt.Errorf("fatal error config file: %s", err))
			}

			var serviceConfig archivar.Config
			err = viper.Unmarshal(&serviceConfig)
			if err != nil {
				panic(fmt.Errorf("fatal error config file: %s", err))
			}

			debugging, _ := cmd.Flags().GetBool("debug")
			if debugging || serviceConfig.Settings.Log.Debugging {
				logger.SetLevel(logrus.DebugLevel)
			}

			quiet, _ := cmd.Flags().GetBool("quiet")
			if quiet {
				logger.SetLevel(logrus.ErrorLevel)
			}

			// if defaultJobInterval == 0 {
			// 	defaultJobInterval = serviceConfig.Settings.DefaultInterval
			// }

			profiler, _ := cmd.Flags().GetBool("profiler")
			profilerPort, _ := cmd.Flags().GetInt("profilerPort")
			if profiler {
				go func() {
					listenHostPort := "0.0.0.0:" + fmt.Sprintf("%d", profilerPort)
					logger.Warnln("Run profiler", listenHostPort)
					logger.Warnln(http.ListenAndServe(listenHostPort, nil))
				}()
			}

			archivarSvc = archivar.New(serviceConfig, logger)
		},
	}

	rootCmd.PersistentFlags().StringP("config", "c", "archivar.yaml", "Configfile")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "enable verbose logging output")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "suppress all non-error output")
	rootCmd.PersistentFlags().Bool("profiler", false, "run go profiler server")
	rootCmd.PersistentFlags().Int("profilerPort", DEFAULT_PROFILER_PORT, "run go profiler server")
	rootCmd.AddCommand(cmdWatch)
	// rootCmd.AddCommand(imap.CmdImap)
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
