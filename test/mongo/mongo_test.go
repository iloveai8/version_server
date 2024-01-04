package move_vsn

import (
	"fmt"
	"game_slots_vsn/pkg/mongo"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName("dev")
	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	mongo.SetUp()
}

func TestMongo(t *testing.T) {
	filter := bson.M{
		"archive_id": "10001",
	}

	srs := bson.M{}
	update := bson.D{{"$set", bson.D{
		{"a", 1},
		{"b", 2},
	}}}

	err := mongo.Mdb.FindOne("doc.account", filter, srs)
	if err != nil {
		return
	}

	//fmt.Println("rs:", srs)
	err = mongo.Mdb.FindOneAndUpdate("doc.account", filter, update, srs)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("rs:", srs)
}

func importUser() {

}
