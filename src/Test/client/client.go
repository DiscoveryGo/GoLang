package main

import (
	"./scenarioParser"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"runtime"
	"sync"
)

func sendSimulationData(airthreat scenarioParser.Airthreat, conn net.Conn) {
	enc := gob.NewEncoder(conn)
	if err := enc.Encode(airthreat); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("데이터 전송 완료")

	//defer conn.Close()
}

func connectServer() (conn net.Conn) {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
	 	log.Println("서버에 연결할 수 없습니다.")
	 }

	fmt.Println("서버에 연결되었습니다.")

	return conn
}

func main() {
	airthreats := scenarioParser.ReadScenarioFile()
	conn := connectServer()

	runtime.GOMAXPROCS(1)
	var wg sync.WaitGroup
	wg.Add(2)

	fmt.Println("goroutine go !")

	//	air 1 go
	go func() {
		defer wg.Done()

		for airthreats[0].PositionX < 1000 {
			airthreats[0].PositionX += 100

			sendSimulationData(airthreats[0], conn)
			fmt.Println(airthreats[0])
		}
	}()

	//	air 2 go
	go func() {
		defer wg.Done()

		for airthreats[1].PositionX < 1000 {
			airthreats[1].PositionX += 50

			sendSimulationData(airthreats[1], conn)
			fmt.Println(airthreats[1])
		}
	}()

	fmt.Println("waiting ...")
	wg.Wait()

	defer conn.Close()
}
