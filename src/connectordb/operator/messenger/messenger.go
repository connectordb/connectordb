package messenger

/*
Package Messenger is a package that implements the pub/sub messaging system used for streaming uplinks and downlinks
as well as the messaging system that allows real-time low-latency data analysis.
*/

import (
	"bytes"
	"util"
	"strings"

	"github.com/nats-io/nats"
	"gopkg.in/vmihailenco/msgpack.v2"
)

//MessageEncoding is the encoding used for messages
const MessageEncoding string = "msgpack"

//Package messenger provides a simple messaging service using gnatsd, which can be used to
//send fast messages to a given user/device/stream from a given user/device

//The MsgPackEncoder encodes the data using msgpack (more wire-efficient than json)
type MsgPackEncoder struct {
}

//Encode is to fit the interface
func (mpe MsgPackEncoder) Encode(subject string, v interface{}) ([]byte, error) {
	b := new(bytes.Buffer)
	enc := msgpack.NewEncoder(b)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

//Decode is to fit the interface
func (mpe MsgPackEncoder) Decode(subject string, data []byte, vPtr interface{}) (err error) {
	//We need to use a special unmarshaller to be able to have map[string]
	return util.MsgPackUnmarshal(data, vPtr)
}

//Register the msgpack encoder
func init() {
	nats.RegisterEncoder("msgpack", MsgPackEncoder{})
}

//Messenger holds an open connection to the gnatsd daemon
type Messenger struct {
	SendConn  *nats.Conn        //The NATS connection
	SendEconn *nats.EncodedConn //The Encoded conn, ie, a data message
	RecvConn  *nats.Conn
	RecvEconn *nats.EncodedConn
}

//Close shuts down a Messenger
func (m *Messenger) Close() {
	m.SendEconn.Close()
	m.RecvEconn.Close()
	m.RecvConn.Close()
	m.SendConn.Close()
}

//ConnectMessenger initializes a connection with the gnatsd messenger. Allows daisy-chaining errors
func ConnectMessenger(opt *nats.Options, err error) (*Messenger, error) {
	if err != nil {
		return nil, err
	}

	sconn, err := opt.Connect()
	if err != nil {
		return nil, err
	}
	seconn, err := nats.NewEncodedConn(sconn, MessageEncoding)
	if err != nil {
		sconn.Close()
		return nil, err
	}

	rconn, err := opt.Connect()
	if err != nil {
		seconn.Close()
		sconn.Close()
		return nil, err
	}
	reconn, err := nats.NewEncodedConn(rconn, MessageEncoding)
	if err != nil {
		seconn.Close()
		sconn.Close()
		rconn.Close()
		return nil, err
	}

	return &Messenger{sconn, seconn, rconn, reconn}, nil
}

//Publish sends the given message over the connection
func (m *Messenger) Publish(routing string, msg Message) error {
	routing = strings.Replace(routing, "/", ".", -1)
	if routing[len(routing)-1] == '.' {
		routing = routing[0 : len(routing)-1]
	}
	return m.SendEconn.Publish(routing, msg)
}

//Subscribe creates a subscription for the given routing string. The routing string is of the format:
//  [user]/[device]/[stream]/[substream//]
//In order to skip something, you can use wildcards, and to skip "the rest" you can use ">" (this is literally the gnatsd routing)
//An example of subscribing to all posts by sender user user1:
//  msgr.Subscribe("user1/>",chn)
//An example of subscribing to everything is:
//	msgr.Subscribe(">",chn)
//Subscribing to a stream is:
// msgr.Subscribe("user/device/stream")
func (m *Messenger) Subscribe(routing string, chn chan Message) (*nats.Subscription, error) {
	routing = strings.Replace(routing, "/", ".", -1)
	if routing[len(routing)-1] == '.' {
		routing = routing[0 : len(routing)-1]
	}
	return m.RecvEconn.BindRecvChan(routing, chn)
}

//Flush makes sure all commands are acknowledged by the server
func (m *Messenger) Flush() {
	m.SendEconn.Flush()
	m.RecvEconn.Flush()
}
