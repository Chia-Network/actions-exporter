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
		org := viper.GetString("github-org")
		user := viper.GetString("github-user")

		if org == "" && user == "" {
			log.Fatalln("at least one of --github-org or --github-user must be provided")
		}

		gh.Init(viper.GetString("github-token"))
		prometheus.Init(org, user)

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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: $HOME/.actions-exporter.yaml)")
	rootCmd.PersistentFlags().Uint16("server-port", 8080, "The port the metrics server binds to. (default: 8080)")
	rootCmd.PersistentFlags().String("github-token", "", "A GitHub API token")
	rootCmd.PersistentFlags().String("github-org", "", "A GitHub organization name")
	rootCmd.PersistentFlags().String("github-user", "", "A GitHub username")
	rootCmd.PersistentFlags().String("log-level", "info", "How verbose the logs should be. panic, fatal, error, warn, info, debug, trace (default: info)")

	for _, flag := range []string{"server-port", "github-token", "github-org", "github-user", "log-level"} {
		err := viper.BindPFlag(flag, rootCmd.PersistentFlags().Lookup(flag))
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".actions-exporter")
	}

	viper.SetEnvPrefix("ACTIONS_EXPORTER")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}

	level, err := log.ParseLevel(viper.GetString("log-level"))
	if err != nil {
		log.Fatalf("Error parsing log level: %s\n", err.Error())
	}
	log.SetLevel(level)
}
