package server

import (
	"fmt"
	"log"
	"net"
)



func StartServer() {
    fmt.Println("Running server")
    server := NewServer(":3000")
    go func() {
        for msg := range server.msgCh {
            fmt.Println("Message recieved:", string(msg))
        }
    }()

    log.Fatal(server.Start())
}



type GameServer struct {
    liAddr  string
    li      net.Listener
    quitCh  chan struct{}
    msgCh   chan []byte
}

func NewServer(listenerAddr string) *GameServer {
    return &GameServer {
        liAddr: listenerAddr,
        quitCh: make(chan struct{}),
        msgCh: make(chan []byte, 10),
    }
}

func (s *GameServer) Start() error {
    li, err := net.Listen("tcp", s.liAddr)
    if err != nil { return err }
    defer li.Close()
    s.li = li

    go s.listen()
    <-s.quitCh
    close(s.msgCh)

    return nil
}

func (s *GameServer) listen() {
    for {
        conn, err := s.li.Accept()
        if err != nil {
            fmt.Println("Accept error:", err)
            continue
        }

        fmt.Println("New connection from:", conn.RemoteAddr())
        go s.read(conn)
    }
}

func (s *GameServer) read(c net.Conn) {
    defer c.Close()
    buf := make([]byte, 2048)

    for {
        n, err := c.Read(buf)
        if err != nil {
            fmt.Println("Read err:", err)
            continue
        }

        s.msgCh <- buf[:n]
    }
}


