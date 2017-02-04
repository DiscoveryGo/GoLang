package main

import (
	"csvlib"
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

func sendSimulationData(airthreats []csvlib.Airthreat) {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Println("서버에 연결할 수 없습니다.")
	}

	defer conn.Close()

	fmt.Println("서버에 연결되었습니다.")

	enc := gob.NewEncoder(conn)
	if err := enc.Encode(airthreats); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("데이터 전송 완료")
}

func main() {

	airthreats := csvlib.ReadScenarioFile()
	sendSimulationData(airthreats)
}
