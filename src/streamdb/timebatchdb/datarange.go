package timebatchdb

import (
	"container/list"
)

//The DataRange interface - this is the object that is returned from different caches/stores - it represents
//a range of data values stored in a certain way, and Next() gets the next datapoint in the range.
type DataRange interface {
	Init() error               //Does the necessary steps to get the datarange ready for returning datapoints
	Next() (*Datapoint, error) //Returns the next datapoint in sequence - or nil if the sequence is finished
	Close()                    //Closes the datarange - can be called before Init. But Init does not have to work after close.
}

//The EmptyRange is a range that always returns nil - as if there were no datapoints left.
//It is the DataRange equivalent of nil
type EmptyRange struct{}

//Close does absolutely nothing
func (r EmptyRange) Close() {}

//Init just returns nil
func (r EmptyRange) Init() error {
	return nil
}

//Next always just returns nil,nil
func (r EmptyRange) Next() (*Datapoint, error) {
	return nil, nil
}

//The RangeList - it is a list of DataRanges, and acts as one large DataRange.
type RangeList struct {
	rlist *list.List //A list of DataRange objects
}

//Init initializes the RangeList as well as the first element in the actual range list
func (r *RangeList) Init() error {
	if r.rlist.Len() != 0 {
		//Initialize the first in the list
		return r.rlist.Front().Value.(DataRange).Init()
	}
	return nil
}

//Close calls close on all elements of the RangeList, and closes itself.
func (r *RangeList) Close() {
	if r.rlist.Len() > 0 {
		//Closes all child DataRanges
		elem := r.rlist.Front()
		for elem.Next() != nil {
			elem.Value.(DataRange).Close()
			elem = elem.Next()
		}
		elem.Value.(DataRange).Close()
	}
}

//Next returns the next available datapoint value from the list, initializing and closing the necessary stuff
func (r *RangeList) Next() (*Datapoint, error) {
	if r.rlist.Len() == 0 {
		return nil, nil
	}
	e := r.rlist.Front().Value.(DataRange)
	d, err := e.Next()
	if d != nil || err != nil {
		return d, err
	}

	//Okay, this element of the list is empty, we close it, remove it from the list,
	//initialize the next element, and repeat
	e.Close()
	r.rlist.Remove(r.rlist.Front())
	if r.rlist.Len() == 0 {
		return nil, nil
	}
	//Initialize the next element
	err = r.rlist.Front().Value.(DataRange).Init()
	if err != nil {
		return nil, err
	}

	//repeat the procedure
	return r.Next()

}

//Append to the end of the rangelist an uninitialized datarange
func (r *RangeList) Append(d DataRange) {
	r.rlist.PushBack(d)
}

//NewRangeList creates empty RangeList
func NewRangeList() *RangeList {
	return &RangeList{list.New()}
}

//A TimeRange is a Datarange which is time-bounded from both sides. That is, the datapoints allowed are only
//within the given time range. So if given a DataRange with range [a,b], and the timerange is (c,d], the
//TimeRange will return all datapoints within the Datarange which are within (c,d].
type TimeRange struct {
	dr        DataRange //The DataRange to wrap
	starttime int64     //The time at which to start the time range
	endtime   int64     //The time at which to stop returning datapoints
}

//Close closes the internal DataRange
func (r *TimeRange) Close() {
	r.dr.Close()
}

//Init runs Init on the internal DataRange
func (r *TimeRange) Init() error {
	return r.dr.Init()
}

//Next returns the next datapoint in sequence from the underlying DataRange, so long as it is within the
//correct timestamp bounds
func (r *TimeRange) Next() (*Datapoint, error) {
	dp, err := r.dr.Next()
	//Skip datapoints before the starttime
	for dp != nil && dp.Timestamp() <= r.starttime {
		dp, err = r.dr.Next()
	}
	//Return nil if the timestamp is beyond our range
	if dp != nil && r.endtime > 0 && dp.Timestamp() > r.endtime {
		//The datapoint is beyond our range.
		return nil, nil
	}
	return dp, err
}

//NewTimeRange creates a time range given the time range of valid datapoints
func NewTimeRange(dr DataRange, starttime int64, endtime int64) *TimeRange {
	return &TimeRange{dr, starttime, endtime}
}

//NumRange returns only the first given number of datapoints (with an optional skip param) from a DataRange
type NumRange struct {
	dr      DataRange
	numleft uint64 //The number of dtapoints left to return
}

//Close closes the internal DataRange
func (r *NumRange) Close() {
	r.dr.Close()
}

//Init runs Init on the internal DataRange
func (r *NumRange) Init() error {
	return r.dr.Init()
}

//Next returns the next datapoint from the underlying DataRange so long as the datapoint is within the
//amonut of datapoints to return.
func (r *NumRange) Next() (*Datapoint, error) {
	if r.numleft == 0 {
		return nil, nil
	}
	r.numleft--
	return r.dr.Next()
}

//Skip the given number of datapoints without changing the number of datapoints left to return
func (r *NumRange) Skip(num int) error {
	for i := 0; i < num; i++ {
		_, err := r.dr.Next()
		if err != nil {
			return err
		}
	}
	return nil
}

//NewNumRange initializes a new NumRange which will return up to the given amount of datapoints.
func NewNumRange(dr DataRange, datapoints uint64) *NumRange {
	return &NumRange{dr, datapoints}
}
