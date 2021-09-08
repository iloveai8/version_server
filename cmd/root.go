package cmd

import (
	`fmt`
	`game_slots_vsn/internal/config`
	"github.com/spf13/cobra"
	`github.com/spf13/viper`
	`log`
	`os`
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use: "game slots vsn",
		Run: func(cmd *cobra.Command, args []string) {
			//internal.Run()
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "conf/dev.yaml", "config file (default is $HOME/conf/dev.yaml)")
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Search config in home directory with name ".slots_robot" (without extension).
		viper.AddConfigPath(pwd)
		viper.AddConfigPath("conf")
		viper.SetConfigName("dev")
		viper.SetConfigType("yaml")
	}
	//viper.AutomaticEnv()
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("read config failed: %v", err)
		fmt.Println(err)
		os.Exit(1)
	}
	var c config.Config
	viper.Unmarshal(&c)

	fmt.Println(c)
	//fmt.Println(viper.GetString(consts.AppServerIp))
	//fmt.Println(viper.GetInt(consts.AppServerPort))
	//fmt.Println("Using config file:", viper.ConfigFileUsed())
	config.Init(cfgFile)
}
