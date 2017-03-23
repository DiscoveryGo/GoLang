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

/*
	user type
*/
type AirthreatSimulator struct {
	ThreatModels []SimulationModel

	//	thinking issue
	conn net.Conn 	
	wg sync.WaitGroup
}

type ShipthreatSimulator struct {
	ThreatModels []SimulationModel
}


/*
	AirthreatModel interface implementaion ==> custom implementation
*/

//	Scenario Control
//	todo: if scn reloaded, need to reset []SimulationModel
func (ats *AirthreatSimulator) readScenario() {
	scenario := scenarioParser.ReadScenarioFile()
	temp := make([]SimulationModel, len(scenario))
	ats.ThreatModels = append(temp)

	for i, objectData := range scenario {
		ats.ThreatModels[i].ID = objectData.ID
		ats.ThreatModels[i].PosX = objectData.PositionX
		ats.ThreatModels[i].PosY = objectData.PositionY
		ats.ThreatModels[i].PosZ = objectData.PositionZ
	}
}


//	Simulation Control

func (ats *AirthreatSimulator) simulationStart() {
	ats.sendEvent(SimulationEvent{10001, "simulation start ..."})

	runtime.GOMAXPROCS(1)
	ats.wg.Add(1)

	fmt.Println("goroutine go !")
	go ats.objectUpdate()

	fmt.Println("waiting ...")
	ats.wg.Wait()
}


func (ats *AirthreatSimulator) simulationStop() {
	ats.sendEvent(SimulationEvent{10003, "simulation stop ..."})
}


func (ats *AirthreatSimulator) update() bool {
	isSuccess := true

	for i := 0; i < len(ats.ThreatModels); i++ {
		ats.ThreatModels[i].PosX += 50
		ats.ThreatModels[i].PosY += 75
		ats.ThreatModels[i].PosZ += 100

		if ats.ThreatModels[i].PosX > 1000 {
			isSuccess = false
			fmt.Printf("ats %d's x position is over 1000\n", i)
		}
	}

	return isSuccess
}


func (ats *AirthreatSimulator) objectUpdate() {
	defer ats.wg.Done()

	ticker := time.NewTicker(time.Millisecond * 500)		//	not apply really. need to know ticker spec
	defer ticker.Stop()

	for tick := range ticker.C {
		_ = tick 	//	not use

		for i := 0; i < len(ats.ThreatModels); i++ {
			if ats.update() == true {
				ats.sendData()	
			} else {
				return
			}
		}
	}
}



//	simulator communication

func (ats *AirthreatSimulator) sendEvent(event SimulationEvent) {
	data := CommunicationMessage{100, []SimulationModel{}, SimulationEvent{event.EventID, event.EventMessage}}	//	100 is Opcode "threat event"
	enc := gob.NewEncoder(ats.conn)
	
	if err := enc.Encode(data); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(data, "이벤트 전송 완료")
}


func (ats *AirthreatSimulator) sendData() {
	data := CommunicationMessage{200,  ats.ThreatModels, SimulationEvent{}}	//	200 is Opcode "threat object data"
	enc := gob.NewEncoder(ats.conn)
	
	if err := enc.Encode(data); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(data, "데이터 전송 완료")	
}


//	tcp Client Server control

func (ats *AirthreatSimulator) connectTCPServer() {
	if ats.conn != nil {
		fmt.Println("서버와 이미 연결된 상태입니다.")
		return
	}

	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
	 	log.Println("서버에 연결할 수 없습니다.", err)
	 } else {
	 	fmt.Println("서버에 연결되었습니다.")
		ats.conn = conn
	 }
}


func (ats *AirthreatSimulator) disconnectTCPServer() {
	ats.sendEvent(SimulationEvent{20001, "TCP client request discoonection"})
	ats.conn.Close()
}


func main() {
	ats := &AirthreatSimulator{}
	var userCommand string

	for userCommand != "6" {
		fmt.Println("\n[simulation client menu]")
		fmt.Println("1. connect server")
		fmt.Println("2. simulation start")
		fmt.Println("3. simulation event send")
		fmt.Println("4. simulation stop")
		fmt.Println("5. disconnect server")
		fmt.Println("6. quit")
		fmt.Printf("input command : ")
		fmt.Scan(&userCommand)

		switch userCommand {
			case "1":
				ats.connectTCPServer()
			case "2":
				ats.readScenario()		//	initialize data by scenario
				ats.simulationStart()	//	data update start
			case "3":
				ats.sendEvent(SimulationEvent{10002, "event test"})
			case "4":
				ats.simulationStop()
			case "5":
				ats.disconnectTCPServer()
		}
	}

	if ats.conn != nil {
		ats.conn.Close()	
	}
}