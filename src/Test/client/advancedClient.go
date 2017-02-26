package main 

import (
	"./scenarioParser"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"runtime"
	"sync"
	"time"
)

//	simulation using type definition

/*
	common type
*/

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
	Event []string
	FriendsOrEnemy []string
	ID int
}

type Model struct {
	PhysicalAttribute
	LogicalAttribute
}

/*
	user type
*/

type AirthreatModel struct {
	Model
	//	speicial air attribute
}


//	simulation using interface definition

/*
	common interface
*/

type Parser interface {
	Parse()		//	raw data in and out to model current data
}

type Modeler interface {
	Update() bool	//	current state data to next state data	
}


/*
	AirthreatModel implementaion
*/

//	todo	need to efficient(read with csv item nums). lazy imp for velocity
func (airthreat *AirthreatModel) Parse() {
	rawData := scenarioParser.ReadScenarioFile()

	for _, data := range rawData {
		if airthreat.ID == data.ID {
			airthreat.PosX = data.PositionX
			airthreat.PosY = data.PositionY
			airthreat.PosZ = data.PositionZ
		}
	}
}

func (airthreat *AirthreatModel) Update() bool {
	airthreat.PosX += 50
	airthreat.PosY += 75
	airthreat.PosZ += 100

	isSuccess := true

	if airthreat.PosX > 1000 {
		isSuccess = false
	}

	return isSuccess
}

func Connect() net.Conn {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
	 	log.Println("서버에 연결할 수 없습니다.")
	 }

	fmt.Println("서버에 연결되었습니다.")

	return conn
}


//	wrapper
func DataParse(p Parser) {
	p.Parse()
}

func UpdateData(m Modeler) bool {
	if m.Update() != true {
		return false
	}

	return true
}

//	I want to do with interface..
func SendData(airthreat AirthreatModel, conn net.Conn) {
	data := scenarioParser.Airthreat{airthreat.ID, airthreat.PosX, airthreat.PosY, airthreat.PosZ}	

	enc := gob.NewEncoder(conn)
	if err := enc.Encode(data); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("데이터 전송 완료")
}

//	todo refactoring
func objectUpdate(airthreats []AirthreatModel, conn net.Conn) {
	defer wg.Done()

	ticker := time.NewTicker(time.Nanosecond * 100)		//	not apply really. need to know ticker spec
	defer ticker.Stop()

	for tick := range ticker.C {
		_ = tick 	//	not use

		for i := 0; i < len(airthreats); i++ {
			if UpdateData(&airthreats[i]) == true {
				SendData(airthreats[i], conn)	
			} else {
				return
			}
		}
	}
}


/*	Simulation Definition

	Read scenario raw data <- using 'go' read IO interface
	Parse read data
	update model
	model data send to server
*/

//	global variable. any idea?
var wg sync.WaitGroup


func main() {
	runtime.GOMAXPROCS(1)
	wg.Add(1)

	conn := Connect()
	airthreats := make([]AirthreatModel, 2)

	for i := 0; i < len(airthreats); i++ {
		airthreats[i].ID = i + 1
		DataParse(&airthreats[i])
	}

	fmt.Println("goroutine go !")

	go objectUpdate(airthreats, conn)

	fmt.Println("waiting ...")
	wg.Wait()

	defer conn.Close()
}