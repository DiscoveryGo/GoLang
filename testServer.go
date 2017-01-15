package main

import (
    "fmt"
    "net"
    "os"
    "encoding/gob"
)

type Airthreat struct {
    ID              int64
    PositionX       float64
    PositionY       float64
    PositionZ       float64
}

func main() {

    service := "0.0.0.0:8000"
    tcpAddr, err := net.ResolveTCPAddr("tcp", service)
    checkError(err)

    listener, err := net.ListenTCP("tcp", tcpAddr)
    checkError(err)

    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }

        //encoder := gob.NewEncoder(conn)
        decoder := gob.NewDecoder(conn)

        var airthreats []Airthreat
        decoder.Decode(&airthreats)
        fmt.Println(airthreats)
        conn.Close() // we're finished
    }
}

func checkError(err error) {
    if err != nil {
        fmt.Println("Fatal error ", err.Error())
        os.Exit(1)
    }
}