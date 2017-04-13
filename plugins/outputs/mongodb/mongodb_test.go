package mongodb

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/influxdata/telegraf/testutil"
)

func TestMongoInsertWithCorrectSettings(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	testHost := testutil.GetLocalHost() + ":27017"
	e := &Mongodb{
		Hosts : []string{testHost},
		Db : "test",
		Collection : "telegrafMetric",
	}

	err := e.Connect()
	require.NoError(t, err)

	err = e.Write(testutil.MockMetrics())
	require.NoError(t, err)
}