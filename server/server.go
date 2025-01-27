package server

import (
	"fmt"
	"io"
	"log"
	"net"
)



func StartServer() {
    fmt.Println("Running server")
    server := NewServer(":3000")
    go func() {
        for msg := range server.msgCh {
            fmt.Printf(
                "Message recieved from (%s):\n%s\n", 
                msg.from,
                string(msg.payload),
            )
        }
    }()

    log.Fatal(server.Start())
}


type Message struct {
    from    string
    payload []byte
}

type GameServer struct {
    addr        string              // listener address
    listener    net.Listener        
    quitCh      chan struct{}       // 0 byte channel (idk why)
    msgCh       chan Message        // 
    conns       map[net.Addr]string // map of ip to profile
}

func NewServer(listenerAddr string) *GameServer {
    return &GameServer {
        addr: listenerAddr,
        quitCh: make(chan struct{}),
        msgCh: make(chan Message, 10),
    }
}

func (s *GameServer) Start() error {
    li, err := net.Listen("tcp", s.addr)
    if err != nil { return err }
    defer li.Close()
    s.listener = li

    go s.listen()
    <-s.quitCh
    close(s.msgCh)

    return nil
}

func (s *GameServer) listen() {
    for {
        conn, err := s.listener.Accept()
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
            if err == io.EOF { return }
            fmt.Println("Read err:", err)
            continue
        }

        s.msgCh <- Message {
            from: c.RemoteAddr().String(),
            payload: buf[:n],
        }
    }
}


