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

type DynamoDB struct {
	db *dynamo.DB
}

type DynamoTable struct {
	table dynamo.Table
}

func (d DynamoDB) Table(name string) Table {
	return &DynamoTable{
		table: d.db.Table(name),
	}
}

func (d *DynamoDB) TxPut(items ...*TableItem) error {
	tx := d.db.WriteTx()
	for _, item := range items {
		t := item.Table.(DynamoTable)
		tx.Put(t.table.Put(item.Item))
	}
	return tx.Run()
}

func (d DynamoDB) RoomsTable() Table {
	return &DynamoTable{
		table: d.db.Table(os.Getenv("ROOM_TABLE_NAME")),
	}
}

func (d DynamoDB) ConnectionsTable() Table {
	return &DynamoTable{
		table: d.db.Table(os.Getenv("CONNECTION_TABLE_NAME")),
	}
}

func (t DynamoTable) FindOne(column string, key, out interface{}) error {
	return t.table.Get(column, key).One(out)
}

func (t DynamoTable) Put(item interface{}) error {
	return t.table.Put(item).Run()
}

func (t DynamoTable) Delete(column string, key interface{}) error {
	return t.table.Delete(column, key).Run()
}

func NewDB() DB {
	return &DynamoDB{
		db: dynamo.New(session.New(), aws.NewConfig()),
	}
}
