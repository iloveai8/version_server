package archive

import (
	"fmt"
	"game_slots_vsn/pkg/redis"
	"github.com/spf13/viper"
	"testing"
)

func init() {
	viper.AddConfigPath("../../conf")
	viper.SetConfigType("yaml")
	viper.SetConfigName("pre")
	viper.AutomaticEnv() // read in environment variables that match
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	redis.SetUp()
}
func TestGetArchive(t *testing.T) {

}
