package main

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/hakobera/serverless-webrtc-signaling-server/common"
	mock "github.com/hakobera/serverless-webrtc-signaling-server/mock_common"
)

func TestOndisconnectHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	myConnID := "Conn1"
	otherConnID := "Other1"
	closeMessage := `{"type":"close"}`
	roomID := "Room1"

	api := mock.NewMockApiGatewayManagementAPI(ctrl)
	db := mock.NewMockDB(ctrl)
	connectionsTable := mock.NewMockTable(ctrl)
	roomsTable := mock.NewMockTable(ctrl)

	api.EXPECT().PostToConnection(myConnID, closeMessage).Return(nil).Times(0)
	api.EXPECT().PostToConnection(otherConnID, closeMessage).Return(nil).Times(1)

	db.EXPECT().ConnectionsTable().Return(connectionsTable)
	db.EXPECT().RoomsTable().Return(roomsTable)

	connectionsTable.EXPECT().FindOne("connectionId", myConnID, gomock.Any()).DoAndReturn(
		func(column string, key, out interface{}) error {
			conn := out.(*common.Connection)
			conn.ConnectionID = myConnID
			conn.RoomID = roomID
			return nil
		})
	connectionsTable.EXPECT().Delete("connectionId", myConnID).Return(nil)

	roomsTable.EXPECT().FindOne("roomId", roomID, gomock.Any()).DoAndReturn(
		func(column string, key, out interface{}) error {
			room := out.(*common.Room)
			room.RoomID = roomID
			room.Clients = []common.Client{
				common.Client{
					ConnectionID: myConnID,
				},
				common.Client{
					ConnectionID: otherConnID,
				},
			}
			return nil
		})
	roomsTable.EXPECT().Put(gomock.Any()).DoAndReturn(
		func(item interface{}) error {
			room := item.(common.Room)
			if len(room.Clients) != 1 {
				t.Errorf("len(room.Clients) should be 1, but %d : %v", len(room.Clients), room.Clients)
			}

			for _, c := range room.Clients {
				if c.ConnectionID == myConnID {
					t.Errorf("room.Clients should not include myConnID")
				}
			}

			return nil
		})

	ondisconnectHandler(api, db, myConnID)
}

func TestOndisconnectHandlerConnectionNotFoundError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	myConnID := "Conn1"
	closeMessage := `{"type":"close"}`
	otherConnID := "Other1"
	expectedError := "connection not found"

	api := mock.NewMockApiGatewayManagementAPI(ctrl)
	db := mock.NewMockDB(ctrl)
	connectionsTable := mock.NewMockTable(ctrl)
	roomsTable := mock.NewMockTable(ctrl)

	api.EXPECT().PostToConnection(myConnID, closeMessage).Return(nil).Times(0)
	api.EXPECT().PostToConnection(otherConnID, closeMessage).Return(nil).Times(0)

	db.EXPECT().ConnectionsTable().Return(connectionsTable)
	db.EXPECT().RoomsTable().Return(roomsTable)

	connectionsTable.EXPECT().FindOne("connectionId", myConnID, gomock.Any()).DoAndReturn(
		func(column string, key, out interface{}) error {
			return fmt.Errorf(expectedError)
		})
	connectionsTable.EXPECT().Delete("connectionId", myConnID).Return(nil).Times(0)

	err := ondisconnectHandler(api, db, myConnID)
	if err.Error() != expectedError {
		t.Errorf("expect %s, but actual %s", expectedError, err.Error())
	}
}

func TestOndisconnectHandlerRoomNotFoundError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	myConnID := "Conn1"
	closeMessage := `{"type":"close"}`
	otherConnID := "Other1"
	roomID := "Room1"
	expectedError := "room not found"

	api := mock.NewMockApiGatewayManagementAPI(ctrl)
	db := mock.NewMockDB(ctrl)
	connectionsTable := mock.NewMockTable(ctrl)
	roomsTable := mock.NewMockTable(ctrl)

	api.EXPECT().PostToConnection(myConnID, closeMessage).Return(nil).Times(0)
	api.EXPECT().PostToConnection(otherConnID, closeMessage).Return(nil).Times(0)

	db.EXPECT().ConnectionsTable().Return(connectionsTable)
	db.EXPECT().RoomsTable().Return(roomsTable)

	connectionsTable.EXPECT().FindOne("connectionId", myConnID, gomock.Any()).DoAndReturn(
		func(column string, key, out interface{}) error {
			conn := out.(*common.Connection)
			conn.ConnectionID = myConnID
			conn.RoomID = roomID
			return nil
		})
	connectionsTable.EXPECT().Delete("connectionId", myConnID).Return(nil)

	roomsTable.EXPECT().FindOne("roomId", roomID, gomock.Any()).DoAndReturn(
		func(column string, key, out interface{}) error {
			return fmt.Errorf(expectedError)
		})

	err := ondisconnectHandler(api, db, myConnID)
	if err.Error() != expectedError {
		t.Errorf("expect %s, but actual %s", expectedError, err.Error())
	}
}

func TestRemove(t *testing.T) {
	cases := []struct {
		clients      []common.Client
		connectionID string
		expected     int
	}{
		{
			clients:      []common.Client{},
			connectionID: "conn1",
			expected:     0,
		},
		{
			clients: []common.Client{
				common.Client{
					ConnectionID: "conn1",
				},
			},
			connectionID: "conn1",
			expected:     0,
		},
		{
			clients: []common.Client{
				common.Client{
					ConnectionID: "conn0",
				},
				common.Client{
					ConnectionID: "conn1",
				},
			},
			connectionID: "conn1",
			expected:     1,
		},
		{
			clients: []common.Client{
				common.Client{
					ConnectionID: "conn0",
				},
				common.Client{
					ConnectionID: "conn1",
				},
			},
			connectionID: "conn3",
			expected:     2,
		},
	}

	for _, c := range cases {
		cs := remove(c.clients, c.connectionID)
		if len(cs) != c.expected {
			t.Errorf("remove() should remove matched client which have connectionID")
		}
	}
}
