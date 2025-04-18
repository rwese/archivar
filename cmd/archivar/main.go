package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/rwese/archivar/archivar"
	encrypter "github.com/rwese/archivar/internal/encrypter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DEFAULT_PROFILER_PORT = 6060
	DEFAULT_CONFIG_FILE   = "archivar.yaml"
)

var (
	logger = log.New()
)

var rootCmd = &cobra.Command{
	Use: "app",
}

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump existing modules and configuration details",
	Run: func(cmd *cobra.Command, args []string) {
		svc := setupArchivarSvc(cmd)
		fmt.Println("Dumping configured modules")
		svc.Dump()
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
		svc := setupArchivarSvc(cmd)
		svc.RunJobs()
	},
}

func main() {
	cobra.OnInitialize()
	if err := rootCmd.Execute(); err != nil {
		logger.Fatalf("Error executing command: %v", err)
	}
}

func init() {
	cmdWatch.Flags().IntP("interval", "i", 60, "Default wait time between processing of all configured archivers, can be overridden per job")

	rootCmd.PersistentFlags().StringP("config", "c", DEFAULT_CONFIG_FILE, "Path to the configuration file")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable verbose logging output")
	rootCmd.PersistentFlags().Bool("trace", false, "Enable tracing in log output")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Suppress all non-error output")
	rootCmd.PersistentFlags().Bool("profiler", false, "Run Go profiler server")
	rootCmd.PersistentFlags().Int("profilerPort", DEFAULT_PROFILER_PORT, "Port for the Go profiler server")
	rootCmd.AddCommand(cmdWatch)
	rootCmd.AddCommand(dumpCmd)
	rootCmd.AddCommand(encrypter.CmdEncrypter)
}

func setupArchivarSvc(cmd *cobra.Command) archivar.Archivar {
	configFile, _ := cmd.Flags().GetString("config")
	viper.SetConfigName(configFile)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		logger.Fatalf("Error reading config file: %v", err)
	}

	var serviceConfig archivar.Config
	if err := viper.Unmarshal(&serviceConfig); err != nil {
		logger.Fatalf("Error unmarshalling config file: %v", err)
	}

	configureLogging(cmd, serviceConfig)

	if profiler, _ := cmd.Flags().GetBool("profiler"); profiler {
		profilerPort, _ := cmd.Flags().GetInt("profilerPort")
		go startProfiler(profilerPort)
	}

	return archivar.New(serviceConfig, logger)
}

func configureLogging(cmd *cobra.Command, serviceConfig archivar.Config) {
	debugging, _ := cmd.Flags().GetBool("debug")
	trace, _ := cmd.Flags().GetBool("trace")
	quiet, _ := cmd.Flags().GetBool("quiet")

	if quiet {
		logger.SetLevel(log.ErrorLevel)
		return
	}

	if debugging || serviceConfig.Settings.Log.Debugging {
		logger.SetLevel(log.DebugLevel)
		if trace {
			logger.SetReportCaller(true)
		}
		logger.SetFormatter(&log.JSONFormatter{
			FieldMap: log.FieldMap{
				log.FieldKeyTime:  "@timestamp",
				log.FieldKeyLevel: "@level",
				log.FieldKeyMsg:   "@message",
				log.FieldKeyFunc:  "@caller",
			},
		})
	}
}

func startProfiler(port int) {
	listenHostPort := fmt.Sprintf("0.0.0.0:%d", port)
	logger.Warnf("Starting profiler on %s", listenHostPort)
	if err := http.ListenAndServe(listenHostPort, nil); err != nil {
		logger.Errorf("Profiler server error: %v", err)
	}
}
