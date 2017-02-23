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

type Position struct {
	x int
	y int
	z int
}

type ThreatModeler interface {
	Update()
}

type Threat struct {
	id int
	Position
}

func (airthreat *Threat) Update() {
	airthreat.x += 50
	airthreat.y += 75
	airthreat.z += 100
}

func sendSimulationData(airthreat scenarioParser.Airthreat, conn net.Conn) {
	enc := gob.NewEncoder(conn)
	if err := enc.Encode(airthreat); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("데이터 전송 완료")
}

func connectServer() (conn net.Conn) {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
	 	log.Println("서버에 연결할 수 없습니다.")
	 }

	fmt.Println("서버에 연결되었습니다.")

	return conn
}

//	scenario read and convert threat model data
//	connect tcp server
//	model data update using go-routine
func main() {
	scnAirs := scenarioParser.ReadScenarioFile()
	threats := make([]Threat, len(scnAirs))

	for i, air := range scnAirs {
		threats[i].id = air.ID
		threats[i].x = air.PositionX
		threats[i].y = air.PositionY
		threats[i].z = air.PositionZ
	}

	conn := connectServer()

	runtime.GOMAXPROCS(1)
	var wg sync.WaitGroup
	wg.Add(1)

	fmt.Println("goroutine go !")

	go func() {
		defer wg.Done()

		ticker := time.NewTicker(time.Millisecond * 100)
		overRange := false

		for tick := range ticker.C {
			for i := 0; i < len(threats); i++ {
				UpdatePosition(&threats[i])
				if threats[i].x > 5000 {
					fmt.Println("theat id '", i, "' x range is over 5000 ")
					ticker.Stop()
					overRange = true
					break

				} else {
					data := scenarioParser.Airthreat{i, 
						threats[i].x, threats[i].y, threats[i].z}
					sendSimulationData(data, conn)
					fmt.Println(data, tick)	
				}
			}

			if overRange == true {
				break;
			}
		}
	}()

	fmt.Println("waiting ...")
	wg.Wait()

	defer conn.Close()
}


func UpdatePosition(model ThreatModeler) {
	model.Update()
}
