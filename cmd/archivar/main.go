package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/rwese/archivar/archivar"
	_ "github.com/rwese/archivar/internal/imap"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	logger = log.New()
)

const DEFAULT_PROFILER_PORT = 6060

var rootCmd = &cobra.Command{
	Use: "app",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

	},
}

var cmdWatch = &cobra.Command{
	Use:   "watch",
	Short: "Start gathering and archiving until stopped",
	Long: `Starts the monitoring process processing all gathers 
and running the archivers until receiving an interrupt signal.`,
	Run: func(cmd *cobra.Command, args []string) {
		defaultJobInterval, _ := cmd.Flags().GetInt("interval")
		logger.Debugf("running watch with default interval: %d", defaultJobInterval)
		svc := setupArchivarSvc(cmd, args)
		svc.RunJobs()
	},
}

func main() {
	cobra.OnInitialize()
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func init() {
	cmdWatch.Flags().IntP("interval", "i", 60, "default wait time between processing of all configured archivers, can be overriden by specifying it per job")

	rootCmd.PersistentFlags().StringP("config", "c", "archivar.yaml", "Configfile")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "enable verbose logging output")
	rootCmd.PersistentFlags().Bool("trace", false, "enable tracing in log output")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "suppress all non-error output")
	rootCmd.PersistentFlags().Bool("profiler", false, "run go profiler server")
	rootCmd.PersistentFlags().Int("profilerPort", DEFAULT_PROFILER_PORT, "run go profiler server")
	rootCmd.AddCommand(cmdWatch)
}

func setupArchivarSvc(cmd *cobra.Command, args []string) archivar.Archivar {
	configFile, _ := cmd.Flags().GetString("config")
	viper.SetConfigName(configFile)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
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
	trace, _ := cmd.Flags().GetBool("trace")
	if debugging || serviceConfig.Settings.Log.Debugging {
		logger.SetLevel(log.DebugLevel)
		if trace {
			logger.SetReportCaller(true)
		}
		logger.SetFormatter(&log.JSONFormatter{
			// DisableColors: true,
			// FullTimestamp: true,
			FieldMap: log.FieldMap{
				log.FieldKeyTime:  "@timestamp",
				log.FieldKeyLevel: "@level",
				log.FieldKeyMsg:   "@message",
				log.FieldKeyFunc:  "@caller",
			},
		})
	}

	quiet, _ := cmd.Flags().GetBool("quiet")
	if quiet {
		logger.SetLevel(log.ErrorLevel)
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

	return archivar.New(serviceConfig, logger)
}
