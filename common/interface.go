package common

type ApiGatewayManagementAPI interface {
	PostToConnection(connectionID, body string) error
}
