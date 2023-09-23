package move_vsn

import (
	"fmt"
	"game_slots_vsn/pkg/mongo"
	"github.com/spf13/viper"
	"testing"
)

var fMdb = &mongo.Mongodb{}
var tMdb = &mongo.Mongodb{}

func init() {
	init1()
	init2()
}

func init1() {
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName("dev")
	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	fMdb.NewMDB()
}

func init2() {
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName("dev")
	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	tMdb.NewMDB()
}

func TestMongo(t *testing.T) {
}

func TestClusterData(t *testing.T) {
}
