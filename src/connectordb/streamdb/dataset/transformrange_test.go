package dataset

import (
	"connectordb/streamdb/datastream"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransformRange(t *testing.T) {
	dpa := datastream.DatapointArray{
		datastream.Datapoint{Data: 1},
		datastream.Datapoint{Data: 10},
		datastream.Datapoint{Data: 7},
		datastream.Datapoint{Data: 1.0},
		datastream.Datapoint{Data: 3},
		datastream.Datapoint{Data: 2.0},
		datastream.Datapoint{Data: 3.14},
	}

	dpa2 := datastream.DatapointArray{
		datastream.Datapoint{Data: false},
		datastream.Datapoint{Data: false},
		datastream.Datapoint{Data: true},
		datastream.Datapoint{Data: false},
		datastream.Datapoint{Data: true},
	}

	dr := datastream.NewDatapointArrayRange(dpa, 0)

	tr, err := NewTransformRange(dr, "if get() < 5:get() >= 3")
	require.NoError(t, err)

	for i := 0; i < len(dpa2); i++ {
		dp, err := tr.Next()
		require.NoError(t, err)
		require.Equal(t, dpa2[i].String(), dp.String())
	}
	dp, err := tr.Next()
	require.NoError(t, err)
	require.Nil(t, dp)

	tr.Close()

	dr = datastream.NewDatapointArrayRange(dpa, 0)

	tr, err = NewTransformRange(dr, "if get() < 5:get() >= 3")
	require.NoError(t, err)

	da, err := tr.NextArray()
	require.NoError(t, err)
	require.NotNil(t, da)
	require.Equal(t, dpa2.String(), da.String())
}
