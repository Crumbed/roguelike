package server

import (
	"fmt"
	"io"
	"log"
	"net"
)



func Start() {
    fmt.Println("Running server")
    // start listening
    listener, err := net.Listen("tcp4", ":12345")
    if err != nil {
        log.Fatal(err)
    }
    defer listener.Close()

    // listen for connections
    for {
        c, err := listener.Accept()
        if err != nil {
            fmt.Println(err)
            return
        }

        go handleCon(c)
    }
}

func handleCon(c net.Conn) {
    fmt.Printf("Serving %s\n", c.RemoteAddr().String())
    packet := make([]byte, 4096)
    tmp := make([]byte, 4096)
    defer c.Close()

    for {
        _, err := c.Read(tmp)
        if err != nil {
            if err != io.EOF {
                fmt.Println("Read error:", err)
            }

            break
        }

        packet = append(packet, tmp...)
    }

    c.Write(packet)
}


type GameServer struct {
}



