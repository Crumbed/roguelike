package main

import (
	"fmt"
	"log"
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

    // client   
    tcpServer, err := net.ResolveTCPAddr(Type, Host + ":" + Port)
    if err != nil {
        log.Fatal("Failed to resolve tcp addr: ", err)
    }

    conn, err := net.DialTCP(Type, nil, tcpServer)
    if err != nil {
        log.Fatal("Failed to dial server: ", err)
    }
    defer conn.Close()

    _, err = conn.Write([]byte("Test"))
    if err != nil {
        log.Fatal("Failed to write data to connection: ", err)
    }
    
    for {}
}

