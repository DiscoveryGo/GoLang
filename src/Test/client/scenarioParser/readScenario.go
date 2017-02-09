// CSVProcessor
package scenarioParser

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
)

type Airthreat struct {
	ID        int
	PositionX int
	PositionY int
	PositionZ int
}

func ReadScenarioFile() []Airthreat {
	dat, err := ioutil.ReadFile("airScenario.csv")
	if err != nil {
		log.Println(err)
	}

	threats := csv.NewReader(bytes.NewReader(dat))
	var airthreats []Airthreat

	for {
		threat, err := threats.Read()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		_id, _ := strconv.Atoi(threat[0])
		_posX, _ := strconv.Atoi(threat[1])
		_posY, _ := strconv.Atoi(threat[2])
		_posZ, _ := strconv.Atoi(threat[3])

		airthreats = append(airthreats, Airthreat{ID: _id, PositionX: _posX, PositionY: _posY, PositionZ: _posZ})
	}

	fmt.Println("시나리오 파일을 읽습니다", airthreats)

	return airthreats
}
