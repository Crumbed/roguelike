package main

import (
	"fmt"
	"log"
	"main/client"
	"main/packet"
	"main/server"
	"net"
	"os"
)

type ProgramType int
const (
    Server ProgramType = iota
    Client
)

func getProgramType() ProgramType {
    args := os.Args
    if len(args) == 1 {
        fmt.Println("No argument found, defaulting to client...")
        return Client
    } else if args[1] == "server" || args[1] == "s" {
        return Server
    } else {
        fmt.Println("Invalid argument, defaulting to client...")
        return Client
    }
}

const (
    Host = "localhost"
    Port = "3000"
    Type = "tcp"
)

func main() {
    fmt.Println("Hello, world")
    p_type := getProgramType()

    if p_type == Server {
        server.StartServer()
        return
    }

    c := client.NewClient()   
    c.AddPacketListener(packet.SCJoinResponse, client.SCJoinResponseListener)
    c.AddPacketListener(packet.SCGameStart, client.SCGameStartListener)
    // default server
    tcpIp, err := net.ResolveTCPAddr(Type, Host + ":" + Port)
    if err != nil {
        log.Fatal("Failed to resolve tcp addr: ", err)
    }
    c.Connect(tcpIp)
    c.Start()
}













