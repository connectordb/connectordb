package users

/**

Provides the ability to count the number of users/devices/streams in the database

Copyright 2015 - Joseph Lewis <joseph@josephlewis.net>

All Rights Reserved

**/

func (userdb *SqlUserDatabase) CountUsers() (uint64, error) {
	var output uint64
	err := userdb.Get(&output, "SELECT COUNT(UserId) FROM Users;")
	return output, err
}

func (userdb *SqlUserDatabase) CountStreams() (uint64, error) {
	var output uint64
	err := userdb.Get(&output, "SELECT COUNT(StreamId) FROM Streams;")
	return output, err
}

func (userdb *SqlUserDatabase) CountDevices() (uint64, error) {
	var output uint64
	err := userdb.Get(&output, "SELECT COUNT(DeviceId) FROM Devices;")
	return output, err
}
