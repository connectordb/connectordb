package rediscache

import (
	"testing"

	"connectordb/streamdb/datastream"

	"github.com/stretchr/testify/require"
)

var (
	dpa1 = datastream.DatapointArray{datastream.Datapoint{1.0, "helloWorld", "me"}, datastream.Datapoint{2.0, "helloWorld2", "me2"}}
	dpa2 = datastream.DatapointArray{datastream.Datapoint{1.0, "helloWorl", "me"}, datastream.Datapoint{2.0, "helloWorld2", "me2"}}
	dpa3 = datastream.DatapointArray{datastream.Datapoint{1.0, "helloWorl", "me"}}

	dpa4 = datastream.DatapointArray{datastream.Datapoint{3.0, 12.0, ""}}

	//Warning: the map types change depending on marshaller/unmarshaller is used
	dpa5 = datastream.DatapointArray{datastream.Datapoint{3.0, map[string]interface{}{"hello": 2.0, "y": "hi"}, ""}}

	dpa6 = datastream.DatapointArray{datastream.Datapoint{1.0, 1.0, ""}, datastream.Datapoint{2.0, 2.0, ""}, datastream.Datapoint{3.0, 3., ""}, datastream.Datapoint{4.0, 4., ""}, datastream.Datapoint{5.0, 5., ""}}
	dpa7 = datastream.DatapointArray{
		datastream.Datapoint{1., "test0", ""},
		datastream.Datapoint{2., "test1", ""},
		datastream.Datapoint{3., "test2", ""},
		datastream.Datapoint{4., "test3", ""},
		datastream.Datapoint{5., "test4", ""},
		datastream.Datapoint{6., "test5", ""},
		datastream.Datapoint{6., "test6", ""},
		datastream.Datapoint{7., "test7", ""},
		datastream.Datapoint{8., "test8", ""},
	}

	dpa8 = datastream.DatapointArray{datastream.Datapoint{2.0, "helloWorld", "me"}, datastream.Datapoint{1.0, "helloWorld2", "me2"}}
)

func TestRedisBasics(t *testing.T) {

	require.NoError(t, rc.Clear())

	require.NoError(t, rc.DeleteStream("", "mystream"))

	i, err := rc.StreamLength("", "mystream", "")
	require.NoError(t, err)
	require.Equal(t, int64(0), i)

	i, err = rc.Insert("mybatcher", "hi", "mystream", "", dpa6, false)
	require.NoError(t, err)
	require.Equal(t, int64(5), i)

	i, err = rc.StreamLength("hi", "mystream", "")
	require.NoError(t, err)
	require.Equal(t, int64(5), i)

	i, err = rc.StreamLength("", "mystream", "")
	require.NoError(t, err)
	require.Equal(t, int64(0), i)

	require.NoError(t, rc.DeleteStream("hi", "mystream"))

	i, err = rc.StreamLength("hi", "mystream", "")
	require.NoError(t, err)
	require.Equal(t, int64(0), i)
}

func TestRedisInsert(t *testing.T) {

	require.NoError(t, rc.Clear())

	rc.BatchSize = 2

	_, err := rc.Insert("mybatcher", "", "mystream", "", dpa6, false)
	require.NoError(t, err)

	dpatest, err := rc.Get("", "mystream", "")
	require.NoError(t, err)
	require.True(t, dpa6.IsEqual(dpatest))

	writestrings, err := rc.GetList("mybatcher")
	require.NoError(t, err)
	require.Equal(t, 2, len(writestrings))
	require.Equal(t, writestrings[0], "{}mystream::2:4")
	require.Equal(t, writestrings[1], "{}mystream::0:2")

	_, err = rc.Insert("mybatcher", "", "mystream", "", dpa1, false)
	require.EqualError(t, err, ErrTimestamp.Error())

	dpz := datastream.DatapointArray{datastream.Datapoint{5.0, "helloWorld", "me"}, datastream.Datapoint{6.0, "helloWorld2", "me2"}}
	i, err := rc.Insert("mybatcher", "", "mystream", "", dpz, false)
	require.NoError(t, err)
	require.Equal(t, int64(7), i)

	writestrings, err = rc.GetList("mybatcher")
	require.NoError(t, err)
	require.Equal(t, 3, len(writestrings))
	require.Equal(t, writestrings[0], "{}mystream::4:6")
	require.Equal(t, writestrings[1], "{}mystream::2:4")
	require.Equal(t, writestrings[2], "{}mystream::0:2")

	i, err = rc.StreamLength("", "mystream", "")
	require.NoError(t, err)
	require.Equal(t, int64(7), i)

	//Now we must test an internal quirk in the redis lua code: inserting more than
	// 5k chunks.
	dpz = make(datastream.DatapointArray, 1, 6000)
	dpz[0] = datastream.Datapoint{9.0, "ol", ""}
	for iter := 1; iter < 6000; iter++ {
		dpz = append(dpz, datastream.Datapoint{10.0 + float64(iter), true, ""})
	}
	i, err = rc.Insert("mybatcher", "", "mystream", "", dpz, false)
	require.NoError(t, err)
	require.Equal(t, int64(6007), i)

	i, err = rc.StreamLength("", "mystream", "")
	require.NoError(t, err)
	require.Equal(t, int64(6007), i)
}

func TestRedisRestamp(t *testing.T) {

	require.NoError(t, rc.Clear())

	_, err := rc.Insert("mybatcher", "", "mystream", "", dpa6, false)
	require.NoError(t, err)
	_, err = rc.Insert("mybatcher", "", "mystream", "", dpa1, true)
	require.NoError(t, err)

	restampedDpa1 := make(datastream.DatapointArray, 2)
	copy(restampedDpa1, dpa1)

	restampedDpa1[0].Timestamp = 5.0
	restampedDpa1[1].Timestamp = 5.0

	dpatest, err := rc.Get("", "mystream", "")
	require.NoError(t, err)
	require.Equal(t, restampedDpa1.String(), dpatest[5:].String())
}

func TestRedisBatchWait(t *testing.T) {
	require.NoError(t, rc.Clear())

	rc.BatchSize = 2

	_, err := rc.Insert("mybatcher", "", "mystream", "", dpa6, false)
	require.NoError(t, err)

	writestrings, err := rc.GetList("mybatcher")
	require.NoError(t, err)
	require.Equal(t, 2, len(writestrings))
	require.Equal(t, writestrings[0], "{}mystream::2:4")
	require.Equal(t, writestrings[1], "{}mystream::0:2")

	s, err := rc.NextBatch("mybatcher", "donebatch")
	require.NoError(t, err)
	require.Equal(t, "{}mystream::0:2", s)

	writestrings, err = rc.GetList("donebatch")
	require.NoError(t, err)
	require.Equal(t, 1, len(writestrings))
	require.Equal(t, writestrings[0], "{}mystream::0:2")

	writestrings, err = rc.GetList("mybatcher")
	require.NoError(t, err)
	require.Equal(t, 1, len(writestrings))
	require.Equal(t, writestrings[0], "{}mystream::2:4")

	require.NoError(t, rc.DeleteKey("donebatch"))
	writestrings, err = rc.GetList("donebatch")
	require.NoError(t, err)
	require.Equal(t, 0, len(writestrings))
}

func TestRedisSubstream(t *testing.T) {

	require.NoError(t, rc.Clear())

	i, err := rc.Insert("mybatcher", "", "mystream", "s1", dpa6, false)
	require.NoError(t, err)
	require.Equal(t, int64(5), i)

	_, err = rc.Insert("mybatcher", "", "mystream", "s1", dpa1, false)
	require.EqualError(t, err, ErrTimestamp.Error())

	i, err = rc.StreamLength("", "mystream", "")
	require.NoError(t, err)
	require.Equal(t, int64(0), i)
	i, err = rc.StreamLength("", "mystream", "s1")
	require.NoError(t, err)
	require.Equal(t, int64(5), i)

	i, err = rc.Insert("mybatcher", "", "mystream", "", dpa1, false)
	require.NoError(t, err)
	require.Equal(t, int64(2), i)

	s, err := rc.GetList("{}mystream:s1")
	require.NoError(t, err)
	require.EqualValues(t, 5, len(s))

	require.NoError(t, rc.DeleteSubstream("", "mystream", "s1"))
	i, err = rc.StreamLength("", "mystream", "")
	require.NoError(t, err)
	require.Equal(t, int64(2), i)
	i, err = rc.StreamLength("", "mystream", "s1")
	require.NoError(t, err)
	require.Equal(t, int64(0), i)

	s, err = rc.GetList("{}mystream:s1")
	require.NoError(t, err)
	require.EqualValues(t, 0, len(s))

}

func TestRedisHashDelete(t *testing.T) {
	require.NoError(t, rc.Clear())

	i, err := rc.Insert("mybatcher", "h1", "mystream", "s1", dpa6, false)
	require.NoError(t, err)
	require.Equal(t, int64(5), i)

	_, err = rc.Insert("mybatcher", "h1", "my2stream", "", dpa1, false)
	require.NoError(t, err)

	_, err = rc.Insert("mybatcher", "h2", "my2stream", "", dpa1, false)
	require.NoError(t, err)

	require.NoError(t, rc.DeleteHash("h1"))

	i, err = rc.StreamLength("h1", "mystream", "s1")
	require.NoError(t, err)
	require.Equal(t, int64(0), i)
	i, err = rc.StreamLength("h1", "mystream", "")
	require.NoError(t, err)
	require.Equal(t, int64(0), i)
	i, err = rc.StreamLength("h1", "mystream", "")
	require.NoError(t, err)
	require.Equal(t, int64(0), i)
	i, err = rc.StreamLength("h2", "my2stream", "")
	require.NoError(t, err)
	require.Equal(t, int64(2), i)

	s, err := rc.GetList("{h1}mystream:s1")
	require.NoError(t, err)
	require.EqualValues(t, 0, len(s))
}

func TestRedisTrim(t *testing.T) {
	require.NoError(t, rc.Clear())
	i, err := rc.Insert("mybatcher", "", "mystream", "", dpa7, false)
	require.NoError(t, err)
	require.EqualValues(t, 9, i)

	require.NoError(t, rc.TrimStream("", "mystream", "", 2))

	dpa, err := rc.Get("", "mystream", "")
	require.NoError(t, err)
	require.Equal(t, dpa7[2:].String(), dpa.String())

	i, err = rc.StreamLength("", "mystream", "")
	require.NoError(t, err)
	require.EqualValues(t, 9, i)

	require.NoError(t, rc.TrimStream("", "mystream", "", 1))

	dpa, err = rc.Get("", "mystream", "")
	require.NoError(t, err)
	require.Equal(t, dpa7[2:].String(), dpa.String())

	require.NoError(t, rc.TrimStream("", "mystream", "", 2))

	dpa, err = rc.Get("", "mystream", "")
	require.NoError(t, err)
	require.Equal(t, dpa7[2:].String(), dpa.String())

	require.NoError(t, rc.TrimStream("", "mystream", "", 3))

	dpa, err = rc.Get("", "mystream", "")
	require.NoError(t, err)
	require.Equal(t, dpa7[3:].String(), dpa.String())
}

func TestRedisRange(t *testing.T) {
	require.NoError(t, rc.Clear())
	i, err := rc.Insert("mybatcher", "", "mystream", "", dpa7, false)
	require.NoError(t, err)
	require.EqualValues(t, 9, i)

	dpa, i1, i2, err := rc.Range("", "systream", "s1", 1, 8)
	require.Error(t, err)

	dpa, i1, i2, err = rc.Range("", "systream", "s1", 0, 8)
	require.NoError(t, err)
	require.EqualValues(t, 0, i1)
	require.EqualValues(t, 0, i2)

	dpa, i1, i2, err = rc.Range("", "mystream", "", 2, 8)
	require.NoError(t, err)
	require.EqualValues(t, 2, i1)
	require.EqualValues(t, 8, i2)
	require.Equal(t, dpa7[2:8].String(), dpa.String())

	dpa, i1, i2, err = rc.Range("", "mystream", "", 0, 0)
	require.NoError(t, err)
	require.EqualValues(t, 0, i1)
	require.EqualValues(t, 9, i2)
	require.Equal(t, dpa7.String(), dpa.String())

	dpa, i1, i2, err = rc.Range("", "mystream", "", -2, -1)
	require.NoError(t, err)
	require.EqualValues(t, 7, i1)
	require.EqualValues(t, 8, i2)
	require.Equal(t, dpa7[7:8].String(), dpa.String())

	dpa, i1, i2, err = rc.Range("", "mystream", "", -2, 20)
	require.NoError(t, err)
	require.EqualValues(t, 7, i1)
	require.EqualValues(t, 9, i2)
	require.Equal(t, dpa7[7:].String(), dpa.String())

	dpa, i1, i2, err = rc.Range("", "mystream", "", -20, 0)
	require.Error(t, err)

	//Now trim the range, to make sure that correct values
	//are returned if not all data is in redis
	require.NoError(t, rc.TrimStream("", "mystream", "", 3))
	dpa, i1, i2, err = rc.Range("", "mystream", "", 3, 0)
	require.NoError(t, err)
	require.EqualValues(t, 3, i1)
	require.EqualValues(t, 9, i2)
	require.Equal(t, dpa7[3:].String(), dpa.String())

	dpa, i1, i2, err = rc.Range("", "mystream", "", 2, 0)
	require.NoError(t, err)
	require.EqualValues(t, 2, i1)
	require.EqualValues(t, 9, i2)
	require.Nil(t, dpa)
}

func TestRedisReadBatch(t *testing.T) {
	rc.Clear()

	rc.BatchSize = 2

	_, err := rc.Insert("mybatcher", "", "mystream", "", dpa6, false)
	require.NoError(t, err)

	s, err := rc.NextBatch("mybatcher", "donebatch")
	require.NoError(t, err)
	require.Equal(t, "{}mystream::0:2", s)

	b, err := rc.ReadBatch(s)
	require.NoError(t, err)
	require.EqualValues(t, 0, b.StartIndex)
	require.EqualValues(t, 2, b.EndIndex())
	require.Equal(t, b.Data.String(), dpa6[:2].String())

	rc.BatchSize = 250

}

func BenchmarkRedis1Insert(b *testing.B) {
	rc.Clear()
	rc.BatchSize = 250
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		rc.Insert("mybatcher", "", "mystream", "", datastream.DatapointArray{datastream.Datapoint{float64(n), true, ""}}, false)
	}
}

func BenchmarkRedis1InsertRestamp(b *testing.B) {
	rc.Clear()
	rc.BatchSize = 250
	rc.Insert("mybatcher", "", "mystream", "", datastream.DatapointArray{datastream.Datapoint{2.0, true, ""}}, false)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		rc.Insert("mybatcher", "", "mystream", "", datastream.DatapointArray{datastream.Datapoint{1.0, true, ""}}, true)
	}
}

func BenchmarkRedis1InsertParallel(b *testing.B) {

	rc.Clear()
	rc.BatchSize = 250
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rc.Insert("mybatcher", "", "mystream", "", datastream.DatapointArray{datastream.Datapoint{1.0, true, ""}}, false)
		}
	})
}

func BenchmarkRedis1000Insert(b *testing.B) {
	rc.Clear()
	rc.BatchSize = 250
	dpa := make(datastream.DatapointArray, 1000)
	for i := 0; i < 1000; i++ {
		dpa[i] = datastream.Datapoint{1.0, true, ""}
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		rc.Insert("mybatcher", "", "mystream", "", dpa, false)
	}
}

func BenchmarkRedis1000InsertParallel(b *testing.B) {
	rc.Clear()
	rc.BatchSize = 250
	dpa := make(datastream.DatapointArray, 1000)
	for i := 0; i < 1000; i++ {
		dpa[i] = datastream.Datapoint{1.0, true, ""}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rc.Insert("mybatcher", "", "mystream", "", dpa, false)
		}
	})
}

func BenchmarkRedis1000InsertRestamp(b *testing.B) {
	rc.Clear()
	rc.BatchSize = 250
	dpa := make(datastream.DatapointArray, 1000)
	for i := 0; i < 1000; i++ {
		dpa[i] = datastream.Datapoint{1.0, true, ""}
	}

	rc.Insert("mybatcher", "", "mystream", "", datastream.DatapointArray{datastream.Datapoint{9000000.0, true, ""}}, false)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		rc.Insert("mybatcher", "", "mystream", "", dpa, true)
	}
}

func BenchmarkRedisStreamLength(b *testing.B) {
	rc.Clear()

	dpa := make(datastream.DatapointArray, 1000)
	for i := 0; i < 1000; i++ {
		dpa[i] = datastream.Datapoint{float64(i), true, ""}
	}
	rc.Insert("mybatcher", "", "mystream", "", dpa, false)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		rc.StreamLength("", "mystream", "")
	}
}

func BenchmarkRedis1000Get(b *testing.B) {
	rc.Clear()

	dpa := make(datastream.DatapointArray, 1000)
	for i := 0; i < 1000; i++ {
		dpa[i] = datastream.Datapoint{1.0, true, ""}
	}
	rc.Insert("mybatcher", "", "mystream", "", dpa, false)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		rc.Get("", "mystream", "")
	}
}

func BenchmarkRedis250Get(b *testing.B) {
	rc.Clear()

	dpa := make(datastream.DatapointArray, 250)
	for i := 0; i < 250; i++ {
		dpa[i] = datastream.Datapoint{1.0, true, ""}
	}
	rc.Insert("mybatcher", "", "mystream", "", dpa, false)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		rc.Get("", "mystream", "")
	}
}

func BenchmarkRedis250Range(b *testing.B) {
	rc.Clear()

	dpa := make(datastream.DatapointArray, 250)
	for i := 0; i < 250; i++ {
		dpa[i] = datastream.Datapoint{1.0, true, ""}
	}
	rc.Insert("mybatcher", "", "mystream", "", dpa, false)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		rc.Range("", "mystream", "", 0, 0)
	}

}

func BenchmarkRedis250RangeMiss(b *testing.B) {
	rc.Clear()

	dpa := make(datastream.DatapointArray, 250)
	for i := 0; i < 250; i++ {
		dpa[i] = datastream.Datapoint{1.0, true, ""}
	}
	rc.Insert("mybatcher", "", "mystream", "", dpa, false)

	rc.TrimStream("", "mystream", "", 4)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		rc.Range("", "mystream", "", 0, 0)
	}

}

func BenchmarkRedis10Range(b *testing.B) {
	rc.Clear()

	dpa := make(datastream.DatapointArray, 250)
	for i := 0; i < 250; i++ {
		dpa[i] = datastream.Datapoint{1.0, true, ""}
	}
	rc.Insert("mybatcher", "", "mystream", "", dpa, false)
	rc.TrimStream("", "mystream", "", 4)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		rc.Range("", "mystream", "", -10, 0)
	}

}
