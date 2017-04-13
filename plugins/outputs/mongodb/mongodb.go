package mongodb

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/outputs"
	"gopkg.in/mgo.v2"
	"fmt"
	"time"
)

type Mongodb struct {
	Hosts      []string  `toml:"hosts"`
	Db         string   `toml:"db"`
	Collection string    `toml:"collection"`
	Username   string    `toml:"username"`
	Password   string    `toml:"password"`
	ServerName   string    `toml:"server_name"`
	session    *mgo.Session
}

var sampleConfig = `
  ## MongoDB host. If port is not specified, then default (27017) will be used for a connection.
  hosts = ["localhost:27017"]    #required (at least one)
  ## MongoDB database name   #required
  db = "mongodb_name"   #required
  ## MongoDB collection name, where documents are or will be stored.
  collection = "telegraf_metric"  #required
  ## User credentials. Not required.
  username = "username"
  password = "password"
  ## Just a unique field for environment in case you want to send similar by type data to one mongodb from different servers
  server_name = "server_name"
`

func (a *Mongodb) Connect() error {

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    a.Hosts,
		Timeout:  60 * time.Second,
		Database: a.Db,
		Username: a.Username,
		Password: a.Password,
	}

	if (len(mongoDBDialInfo.Addrs) == 0) {
		return fmt.Errorf("Host is not defined")
	}

	if a.Db == "" {
		return fmt.Errorf("Database name is not defined")
	}

	if a.Collection == "" {
		return fmt.Errorf("Collection name is not defined")
	}

	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)

	if err != nil {
		return fmt.Errorf("Couldn't connect to database")
	}

	a.session = mongoSession;
	return nil;
}

func (a *Mongodb) Write(metrics []telegraf.Metric) error {

	collection := a.session.DB(a.Db).C(a.Collection)

	for _, telegrafMetric := range metrics {

		var name = telegrafMetric.Name()
		m := make(map[string]interface{})

		m["timestamp"] = telegrafMetric.Time()
		m["measurement_name"] = name
		m["tag"] = telegrafMetric.Tags()
		m["server_name"] = a.ServerName
		m[name] = telegrafMetric.Fields()

		collection.Insert(m);
	}

	return nil
}

func (a *Mongodb) SampleConfig() string {
	return sampleConfig
}

func (a *Mongodb) Description() string {
	return "Configuration for Mongodb to send metrics to."
}

func (a *Mongodb) Close() error {
	a.session.Close();
	return nil;
}

func init() {
	outputs.Add("mongodb", func() telegraf.Output {
		return &Mongodb{}
	})
}