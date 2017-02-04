// CSVProcessor
package csvlib

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
	ID        int64
	PositionX float64
	PositionY float64
	PositionZ float64
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

		_id, _ := strconv.ParseInt(threat[0], 0, 64)
		_posX, _ := strconv.ParseFloat(threat[1], 64)
		_posY, _ := strconv.ParseFloat(threat[2], 64)
		_posZ, _ := strconv.ParseFloat(threat[3], 64)

		airthreats = append(airthreats, Airthreat{ID: _id, PositionX: _posX, PositionY: _posY, PositionZ: _posZ})
	}

	fmt.Println("시나리오 파일을 읽습니다", airthreats)

	return airthreats
}
