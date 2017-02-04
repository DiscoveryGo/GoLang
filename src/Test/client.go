package main 

import (
	"encoding/csv"
	"net"
	"fmt"
	"log"
	"io"
	"io/ioutil"
	"bytes"
	"strconv"
	"encoding/gob"
)

type Airthreat struct {
	ID				int64
	PositionX		float64
	PositionY		float64
	PositionZ		float64
}

func readScenarioFile() []Airthreat {
	dat, err := ioutil.ReadFile("airScenario.csv")
	if err != nil {
		log.Println(err)
	}

	threats := csv.NewReader(bytes.NewReader(dat))
	var airthreats []Airthreat

	 for 
	 {
	 	threat, err := threats.Read()
	 	
	 	if err == io.EOF {
	 		break
	 	}
	 	if err != nil {
	 		log.Fatal(err)
	 	}

	 	_id, _ := strconv.ParseInt(threat[0], 0, 64)
	 	_posX, _ := strconv.ParseFloat(threat[1], 64)
	 	_posY, _ := strconv.ParseFloat(threat[2], 64)
	 	_posZ, _ := strconv.ParseFloat(threat[3], 64)

	 	airthreats = append(airthreats, Airthreat{ID: _id, PositionX: _posX, PositionY: _posY, PositionZ: _posZ})
	 }

	 fmt.Println("시나리오 파일을 읽습니다", airthreats)

	 return airthreats
}


func sendSimulationData(airthreats []Airthreat) {
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
	airthreats := readScenarioFile()
	sendSimulationData(airthreats)
}

