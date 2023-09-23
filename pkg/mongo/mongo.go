package mongo

import (
	"context"
	"game_slots_vsn/pkg/consts"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type setting struct {
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Mongodb struct {
	Ctx    context.Context
	Client *mongo.Client
}

var Mdb = &Mongodb{}

func SetUp() {
	rs := &setting{}
	err := viper.UnmarshalKey(consts.ConfigMongo, rs)
	if err != nil {
		panic(err)
	}
	Mdb.setup(rs)
}

func (mdb *Mongodb) NewMDB() {
	rs := &setting{}
	err := viper.UnmarshalKey(consts.ConfigMongo, rs)
	if err != nil {
		panic(err)
	}
	mdb.setup(rs)
}

func (mdb *Mongodb) setup(s *setting) {
	ctx := context.Background()

	clientOptions := options.Client().ApplyURI(s.URL)
	if len(s.Username) > 0 && len(s.Password) > 0 {
		clientOptions.SetAuth(options.Credential{Username: s.Username, Password: s.Password})
	}
	clientOpts := options.Client().ApplyURI(s.URL)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		panic(err)
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	mdb.Ctx = ctx
	mdb.Client = client
}
