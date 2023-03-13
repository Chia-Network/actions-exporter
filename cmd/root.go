package cmd

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	gh "github.com/chia-network/actions-exporter/internal/github"
	"github.com/chia-network/actions-exporter/internal/http"
	"github.com/chia-network/actions-exporter/internal/prometheus"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "actions-exporter",
	Short: "Prometheus exporter for GitHub Actions",

	Run: func(cmd *cobra.Command, args []string) {
		gh.Init(viper.GetString("github-token"))
		prometheus.Init(viper.GetString("github-org"))

		m, err := http.NewServer(uint16(viper.GetUint("server-port")))
		if err != nil {
			log.Fatalln(err.Error())
		}

		log.Fatalln(m.StartServer())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: $HOME/.ci-exporter.yaml)")
	rootCmd.PersistentFlags().Uint16("server-port", 8080, "The port the metrics server binds to. (default: 8080)")
	rootCmd.PersistentFlags().String("github-token", "", "A GitHub API token")
	rootCmd.PersistentFlags().String("github-org", "", "A GitHub organization name")
	rootCmd.PersistentFlags().String("log-level", "info", "How verbose the logs should be. panic, fatal, error, warn, info, debug, trace (default: info)")

	err := viper.BindPFlag("server-port", rootCmd.PersistentFlags().Lookup("server-port"))
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = viper.BindPFlag("github-token", rootCmd.PersistentFlags().Lookup("github-token"))
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = viper.BindPFlag("github-org", rootCmd.PersistentFlags().Lookup("github-org"))
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	if err != nil {
		log.Fatalln(err.Error())
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".ci-exporter" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".actions-exporter")
	}

	viper.SetEnvPrefix("ACTIONS_EXPORTER")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}

	level, err := log.ParseLevel(viper.GetString("log-level"))
	if err != nil {
		log.Fatalf("Error parsing log level: %s\n", err.Error())
	}
	log.SetLevel(level)
}
