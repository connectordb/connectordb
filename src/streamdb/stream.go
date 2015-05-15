package streamdb

import (
	"encoding/json"
	"errors"
	"streamdb/schema"
	"streamdb/timebatchdb"
	"streamdb/users"
)

var (
	//ErrSchema is thrown when schemas don't match
	ErrSchema = errors.New("The datapoints did not match the stream's schema")
)

//Stream is a wrapper for the users.Stream object which encodes the schema and other parts of a stream
type Stream struct {
	users.Stream
	Schema map[string]interface{} `json:"schema"` //This allows the JsonSchema to be directly unmarshalled

	//These are used internally for the stream to work out
	s *schema.Schema //The schema associated with the stream
}

//NewStream returns a new stream object
func NewStream(s *users.Stream, err error) (Stream, error) {
	if err != nil {
		return Stream{}, err
	}

	strmschema, err := schema.NewSchema(s.Type)
	if err != nil {
		return Stream{}, err
	}
	var schemamap map[string]interface{}

	err = json.Unmarshal([]byte(s.Type), &schemamap)

	return Stream{*s, schemamap, strmschema}, err
}

//Validate ensures the array of datapoints conforms to the schema and such
func (s *Stream) Validate(data []Datapoint) bool {
	for i := range data {
		if !s.s.IsValid(data[i].Data) {
			return false
		}
	}
	return true
}

//Converts a datapoint array to the timebatch equivalent, which is based on byte arrays
func (s *Stream) convertDatapointArray(data []Datapoint) (*timebatchdb.DatapointArray, error) {
	if !s.Validate(data) {
		return nil, ErrSchema
	}

	tbdpa := make([]timebatchdb.Datapoint, len(data))
	for i := range data {
		dpbytes, err := s.s.Marshal(data[i].Data)
		if err != nil {
			return nil, err
		}
		tbdpa[i] = timebatchdb.NewDatapoint(data[i].IntTimestamp(), dpbytes, data[i].Sender)
	}

	return timebatchdb.NewDatapointArray(tbdpa), nil
}
