package common

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type Room struct {
	RoomID  string    `dynamo:"roomId"`
	Clients []Client  `dynamo:"clients"`
	Created time.Time `dynamo:"created"`
}

type Client struct {
	ConnectionID string    `dynamo:"connectionId"`
	ClientID     string    `dynamo:"clientId"`
	Joined       time.Time `dynamo:"joined"`
}

type Connection struct {
	ConnectionID string `dynamo:"connectionId"`
	RoomID       string `dynamo:"roomId"`
}

func DB() *dynamo.DB {
	return dynamo.New(session.New(), aws.NewConfig())
}

func RoomsTable(db *dynamo.DB) dynamo.Table {
	return db.Table(os.Getenv("ROOM_TABLE_NAME"))
}

func ConnectionsTable(db *dynamo.DB) dynamo.Table {
	return db.Table(os.Getenv("CONNECTION_TABLE_NAME"))
}
