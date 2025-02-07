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

func getProgramType(args []string) (ProgramType, string) {
    if len(args) == 1 {
        fmt.Println("No argument found, defaulting to client at localhost:3000")
        return Client, "localhost:3000"
    } else if args[1] == "server" || args[1] == "s" {
        if len(args) == 3 { return Server, args[2] }
        return Server, ""
    }

    return Client, args[1]
}

const Type = "tcp"
func main() {
    args := os.Args
    p_type, ip := getProgramType(args)

    if p_type == Server {
        server.StartServer(ip)
        return
    }

    c := client.NewClient()   
    c.AddPacketListener(packet.SCJoinResponse, client.SCJoinResponseListener)
    c.AddPacketListener(packet.BWGameStart, client.SCGameStartListener)
    c.AddPacketListener(packet.BWPaddleMove, client.CCPaddleMoveListener)
    c.AddPacketListener(packet.SCBallMove, client.SCBallMoveListener)
    c.AddPacketListener(packet.SCScore, client.SCScoreListener)
    tcpIp, err := net.ResolveTCPAddr(Type, ip)
    if err != nil {
        log.Fatal("Failed to resolve tcp addr: ", err)
    }
    c.Connect(tcpIp)
    c.Start()
}













