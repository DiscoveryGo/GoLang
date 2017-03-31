package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

type Position struct {
	PosX int
	PosY int
	PosZ int
}

type Velocity struct {
	VelX int
	VelY int
	VelZ int
}

type PhysicalAttribute struct {
	Position
	Velocity
}

type LogicalAttribute struct {
	ID int
	Category string
	FriendsOrEnemy int
}

type SimulationModel struct {
	PhysicalAttribute
	LogicalAttribute
}

type SimulationEvent struct {
	EventID int
	EventMessage string
}

type CommunicationMessage struct {
	OpCode int
	ThreatModels []SimulationModel
	Event SimulationEvent
}


func main() {
	service := "0.0.0.0:8000"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	var isConnected bool = false
	var conn net.Conn 	//	todo refactoring

	for !isConnected {
		_conn, err := listener.Accept()

		if err != nil {
			continue
		}

		fmt.Println("tcp connected")
		isConnected = true
		conn = _conn
	}
	
	dec := gob.NewDecoder(conn)
	for {
		// dec := gob.NewDecoder(conn)
		var msg CommunicationMessage

		if err := dec.Decode(&msg); err != nil {
			fmt.Println(err)
			return
		}

		switch msg.OpCode {
		case 100:
			switch msg.Event.EventID {
			case 10001, 10002, 10003:
				fmt.Println(msg.Event)
			case 20001:
				fmt.Println(msg.Event)
				conn.Close()
			}

		case 200:
			fmt.Println(msg.ThreatModels)

		default:
			fmt.Println("communication message parsing error!")
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
