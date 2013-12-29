package roomba_api

import (
	"roomba"
)

const DUMMY_PORT_NAME string = "DummyPort"

var DummyRoomba *roomba.Roomba

func MakeDummyRoomba() *roomba.Roomba {
	if DummyRoomba == nil {
		DummyRoomba = roomba.MakeTestRoomba()
	}
	//DummyRoomba.S.(*roomba.CloseableRWBuffer).WriteReadBuffer([]byte{214,11})
	return DummyRoomba
}

func ClearDummyRoomba() {
	if DummyRoomba != nil {
		DummyRoomba = nil
	}
}
