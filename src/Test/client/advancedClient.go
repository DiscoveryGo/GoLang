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
type AirthreatSimulator struct {
	threats []Model

	//	thinking issue
	conn net.Conn 	
	wg sync.WaitGroup
}

//	todo	simulation using interface definition

/*
	AirthreatModel interface implementaion ==> custom implementation
*/

//	todo	if scn reloaded, need to reset []Model
func (ats *AirthreatSimulator) readScenario() {
	scenario := scenarioParser.ReadScenarioFile()
	temp := make([]Model, len(scenario))
	ats.threats = append(temp)

	for i, objectData := range scenario {
		ats.threats[i].ID = objectData.ID
		ats.threats[i].PosX = objectData.PositionX
		ats.threats[i].PosY = objectData.PositionY
		ats.threats[i].PosZ = objectData.PositionZ
	}
}


func (ats *AirthreatSimulator) simulationStart() {
	runtime.GOMAXPROCS(1)
	ats.wg.Add(1)

	fmt.Println("goroutine go !")
	go ats.objectUpdate()

	fmt.Println("waiting ...")
	ats.wg.Wait()
}


func (ats *AirthreatSimulator) update() bool {
	isSuccess := true

	for i := 0; i < len(ats.threats); i++ {
		ats.threats[i].PosX += 50
		ats.threats[i].PosY += 75
		ats.threats[i].PosZ += 100

		if ats.threats[i].PosX > 1000 {
			isSuccess = false
			fmt.Printf("ats %d's x position is over 1000\n", i)
		}
	}

	return isSuccess
}


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
	ats.conn.Write([]byte("bye bye"))
}


// //	I want to do with interface..
func (ats *AirthreatSimulator) sendData() {
	for _, threat := range ats.threats {
		data := scenarioParser.Airthreat{threat.ID, threat.PosX, threat.PosY, threat.PosZ}
		enc := gob.NewEncoder(ats.conn)

		if err := enc.Encode(data); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("데이터 전송 완료")
	}
}


func (ats *AirthreatSimulator) objectUpdate() {
	defer ats.wg.Done()

	ticker := time.NewTicker(time.Millisecond * 100)		//	not apply really. need to know ticker spec
	defer ticker.Stop()

	for tick := range ticker.C {
		_ = tick 	//	not use

		for i := 0; i < len(ats.threats); i++ {
			if ats.update() == true {
				ats.sendData()	
			} else {
				return
			}
		}
	}
}



func main() {
	ats := &AirthreatSimulator{}
	var userCommand string

	for userCommand != "4" {
		fmt.Println("\n[simulation client menu]")
		fmt.Println("1. connect server")
		fmt.Println("2. simulation start")
		fmt.Println("3. disconnect server")
		fmt.Println("4. quit")
		fmt.Printf("input command : ")
		fmt.Scan(&userCommand)

		switch userCommand {
			case "1":
				ats.connectTCPServer()
			case "2":
				ats.readScenario()		//	initialize data by scenario
				ats.simulationStart()	//	data update start

			//	disconnect by ICD message "bye bye"
			case "3":
				ats.disconnectTCPServer()
		}
	}

	if ats.conn != nil {
		ats.conn.Close()	
	}
}