package main

import (
	"fmt"

	"github.com/rwese/archivar/archivar"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	var configFile string
	var debugging bool
	var watchInterval int
	var serviceConfig archivar.Config
	var logger = logrus.New()

	var cmdWatch = &cobra.Command{
		Use:   "watch",
		Short: "Start gathering and archiving until stopped",
		Long: `Starts the monitoring process processing all gathers 
and running the archivers until receiving an interrupt signal.`,
		Run: func(cmd *cobra.Command, args []string) {
			s := archivar.New(serviceConfig, logger)
			logger.Debugf("running watch with interval: %d", watchInterval)
			s.Watch(watchInterval)
		},
	}

	cmdWatch.Flags().IntVarP(&watchInterval, "interval", "i", 60, "wait time between processing of all configured archivers")

	var rootCmd = &cobra.Command{
		Use: "app",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if debugging {
				logger.SetLevel(logrus.DebugLevel)
			}

			viper.SetConfigName(configFile) // name of config file (without extension)
			viper.SetConfigType("yaml")
			viper.AddConfigPath("/etc/go-archivar/")
			err := viper.ReadInConfig()
			if err != nil {
				panic(fmt.Errorf("fatal error config file: %s", err))
			}

			err = viper.Unmarshal(&serviceConfig)
			if err != nil { // Handle errors reading the config file
				panic(fmt.Errorf("fatal error config file: %s", err))
			}
		},
	}

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "archivar.yaml", "Configfile")
	rootCmd.PersistentFlags().BoolVarP(&debugging, "debug", "d", false, "enable verbose logging output")
	rootCmd.AddCommand(cmdWatch)
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
