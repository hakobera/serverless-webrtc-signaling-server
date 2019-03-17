package common

type ApiGatewayManagementAPI interface {
	PostToConnection(connectionID, body string) error
}

type DB interface {
	Table(name string) Table
	TxPut(items ...TableItem) error

	RoomsTable() Table
	ConnectionsTable() Table
}

type Table interface {
	FindOne(column string, key, out interface{}) error
	Put(row interface{}) error
	Delete(column string, key interface{}) error
}

type TableItem struct {
	Table Table
	Item  interface{}
}
