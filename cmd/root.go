package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ryanrizky/kb-cli/client"
)

var (
	cfgFile string
	jsonOut bool
	kbClient *client.KanboardClient
)

var rootCmd = &cobra.Command{
	Use:   "kb-cli",
	Short: "A CLI tool for Kanboard",
	Long:  `kb-cli is a fast and flexible CLI for managing Kanboard projects and tasks.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		url := viper.GetString("url")
		user := viper.GetString("user")
		token := viper.GetString("token")

		if url == "" || token == "" {
			return fmt.Errorf("URL and Token must be provided via config file or flags")
		}
		if user == "" {
			user = "jsonrpc"
		}

		kbClient = client.NewKanboardClient(url, user, token)
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kb-cli.yaml)")
	rootCmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "Output results in JSON format")
	
	rootCmd.PersistentFlags().String("url", "", "Kanboard JSON-RPC URL (e.g. https://kanboard.example.com/jsonrpc.php)")
	rootCmd.PersistentFlags().String("user", "jsonrpc", "Kanboard API User")
	rootCmd.PersistentFlags().String("token", "", "Kanboard API Token")

	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".kb-cli")

		// Create default config file if it doesn't exist
		configPath := filepath.Join(home, ".kb-cli.yaml")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			os.WriteFile(configPath, []byte("url: \"\"\nuser: \"jsonrpc\"\ntoken: \"\"\n"), 0600)
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("KB")

	if err := viper.ReadInConfig(); err == nil {
		// config read successfully
	}
}
