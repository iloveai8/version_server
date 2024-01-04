package mongo

import (
	"context"
	"fmt"
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
	DataBase string `yaml:"database"`
}

type Mongodb struct {
	ctx context.Context
	db  *mongo.Database
}

var Mdb = &Mongodb{}

func SetUp() {
	rs := &setting{}
	err := viper.UnmarshalKey(consts.ConfigMongo, rs)
	if err != nil {
		_ = fmt.Errorf("error mongo uri: %v", err)
		return
	}
	Mdb.setup(rs)
}

func (mDB *Mongodb) NewMDB() {
	rs := &setting{}
	err := viper.UnmarshalKey(consts.ConfigMongo, rs)
	if err != nil {
		panic(err)
	}
	mDB.setup(rs)
}

func (mDB *Mongodb) setup(s *setting) {
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
	mDB.ctx = ctx
	mDB.db = client.Database(s.DataBase)
}

func (mDB *Mongodb) InsertOne(collectName string, doc interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return mDB.db.Collection(collectName).InsertOne(mDB.ctx, doc, opts...)
}

func (mDB *Mongodb) InsertMany(collectName string, docs []interface{},
	opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {

	return mDB.db.Collection(collectName).InsertMany(mDB.ctx, docs, opts...)
}

func (mDB *Mongodb) UpdateOne(collectName string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return mDB.db.Collection(collectName).UpdateOne(mDB.ctx, filter, update, opts...)
}

func (mDB *Mongodb) UpdateOrInsertOne(collectName string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	opts := options.Update().SetUpsert(true)
	return mDB.db.Collection(collectName).UpdateOne(mDB.ctx, filter, update, opts)
}

func (mDB *Mongodb) UpdateMany(collectName string, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {

	return mDB.db.Collection(collectName).UpdateMany(mDB.ctx, filter, update, opts...)
}

func (mDB *Mongodb) DeleteOne(collectName string, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {

	return mDB.db.Collection(collectName).DeleteOne(mDB.ctx, filter, opts...)
}

func (mDB *Mongodb) DeleteMany(collectName string, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {

	return mDB.db.Collection(collectName).DeleteMany(mDB.ctx, filter, opts...)
}

func (mDB *Mongodb) FindOneAndUpdate(collectName string, filter interface{}, update interface{}, result interface{},
	opts ...*options.FindOneAndUpdateOptions) (err error) {

	return mDB.db.Collection(collectName).FindOneAndUpdate(mDB.ctx, filter, update, opts...).Decode(result)
}

func (mDB *Mongodb) FindOne(collectName string, filter interface{}, result interface{},
	opts ...*options.FindOneOptions) error {

	return mDB.db.Collection(collectName).FindOne(mDB.ctx, filter, opts...).Decode(result)
}

func (mDB *Mongodb) Find(collectName string, filter interface{}, result interface{},
	opts ...*options.FindOptions) error {

	cursor, err := mDB.db.Collection(collectName).Find(mDB.ctx, filter, opts...)
	if err != nil {
		return err
	}
	return cursor.All(mDB.ctx, result)
}

func (mDB *Mongodb) Aggregate(collectName string, pipeline interface{},
	opts ...*options.AggregateOptions) ([]map[string]interface{}, error) {

	result := make([]map[string]interface{}, 0)
	var cursor *mongo.Cursor
	var err error
	cursor, err = mDB.db.Collection(collectName).Aggregate(mDB.ctx, pipeline, opts...)
	if err != nil {
		return nil, err
	}
	err = cursor.All(mDB.ctx, &result)
	return result, err
}
