package server

import (
	"fmt"
	"io"
	"log"
	"main/server/packet"
	"net"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)



func StartServer() {
    fmt.Println("Running server")
    server := NewServer(":3000")

    log.Fatal(server.Start())
}


type Message struct {
    From    net.Addr
    Packet  *packet.Packet
}

type GameServer struct {
    addr        string              // listener address
    listener    net.Listener        
    quitCh      chan struct{}       // 0 byte channel (idk why)
    msgCh       chan Message        
    ipconns     map[net.Addr]*Profile 
    idconns     map[uuid.UUID]*Profile 
}

func NewServer(listenerAddr string) *GameServer {
    return &GameServer {
        addr: listenerAddr,
        quitCh: make(chan struct{}),
        msgCh: make(chan Message, 10),
        ipconns: make(map[net.Addr]*Profile),
        idconns: make(map[uuid.UUID]*Profile),
    }
}

func (s *GameServer) Start() error {
    li, err := net.Listen("tcp", s.addr)
    if err != nil { return err }
    defer li.Close()
    s.listener = li

    go s.listen()
    go s.handleMsgs()
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
            if err == io.EOF {
                ip := c.RemoteAddr()
                s.RemovePlayerIp(ip)
                fmt.Printf("Player %s has disconnected\n", ip)
                return 
            }
            fmt.Println("Read err:", err)
            continue
        }
        
        p := &packet.Packet{}
        err = proto.Unmarshal(buf[:n], p)
        if err != nil {
            fmt.Println("Failed to read packet:", err)
            continue
        }

        s.msgCh <- Message {
            From: c.RemoteAddr(),
            Packet: p,
        }
    }
}

func (s *GameServer) handleMsgs() {
    for msg := range s.msgCh {
        p := msg.Packet
        if p.Type == packet.Type_CSProfile {
            p_profile := &packet.Profile{}
            err := proto.Unmarshal(p.Data, p_profile)
            if err != nil {
                fmt.Println("Unmarshal error:", err)
                continue
            }
            
            profile := NewProfile(msg.From, p_profile)
            s.ipconns[msg.From] = profile
            s.idconns[profile.Uuid] = profile
            
            fmt.Printf("Player connected: %s\n", *profile)
        }

    }
}


func (s *GameServer) RemovePlayerId(uuid uuid.UUID) {
    profile := s.idconns[uuid]
    delete(s.idconns, uuid)
    delete(s.ipconns, profile.Ip)
}

func (s *GameServer) RemovePlayerIp(ip net.Addr) {
    profile := s.ipconns[ip]
    delete(s.ipconns, ip)
    delete(s.idconns, profile.Uuid)
}


