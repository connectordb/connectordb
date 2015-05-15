package streamdb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthStreamIO(t *testing.T) {
	require.NoError(t, ResetTimeBatch())

	db, err := Open("postgres://127.0.0.1:52592/connectordb?sslmode=disable", "localhost:6379", "localhost:4222")
	require.NoError(t, err)
	defer db.Close()

	//Let's create a stream
	require.NoError(t, db.CreateUser("tst", "root@localhost", "mypass"))
	require.NoError(t, db.CreateDevice("tst/tst"))

	o, err := db.Operator("tst/tst")
	require.NoError(t, err)

	require.NoError(t, o.CreateStream("tst/tst/tst", `{"type": "integer"}`))

	//Now make sure that length is 0
	l, err := o.LengthStream("tst/tst/tst")
	require.NoError(t, err)
	require.Equal(t, int64(0), l)

	strm, err := o.ReadStream("tst/tst/tst")
	require.NoError(t, err)
	l, err = o.LengthStreamByID(strm.StreamId)

	data := []Datapoint{Datapoint{
		Timestamp: 1.0,
		Data:      1336,
	}}
	require.NoError(t, o.InsertStream("tst/tst/tst", data))

	l, err = o.LengthStream("tst/tst/tst")
	require.NoError(t, err)
	require.Equal(t, int64(1), l)

	dr, err := o.GetStreamTimeRange("tst/tst/tst", 0.0, 2.5, 0)
	require.NoError(t, err)

	dp, err := dr.Next()
	require.NoError(t, err)
	require.NotNil(t, dp)
	require.Equal(t, float64(1336), dp.Data)
	require.Equal(t, 1.0, dp.Timestamp)
	require.Equal(t, "", dp.Sender)

	dp, err = dr.Next()
	require.NoError(t, err)
	require.Nil(t, dp)

	dr.Close()

	dr, err = o.GetStreamIndexRange("tst/tst/tst", 0, 1)
	require.NoError(t, err)

	dp, err = dr.Next()
	require.NoError(t, err)
	require.NotNil(t, dp)
	require.Equal(t, float64(1336), dp.Data)
	require.Equal(t, 1.0, dp.Timestamp)
	require.Equal(t, "", dp.Sender)

	dp, err = dr.Next()
	require.NoError(t, err)
	require.Nil(t, dp)

	dr.Close()

	i, err := db.TimeToIndexStream("tst/tst/tst", 0.3)
	require.NoError(t, err)
	require.Equal(t, int64(0), i)

	//Now let's make sure that stuff is deleted correctly
	require.NoError(t, o.DeleteStream("tst/tst/tst"))
	require.NoError(t, db.CreateStream("tst/tst/tst", `{"type": "string"}`))
	l, err = db.LengthStream("tst/tst/tst")
	require.NoError(t, err)
	require.Equal(t, int64(0), l, "Timebatch has residual data from deleted stream")
}
