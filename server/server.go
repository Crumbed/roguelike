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
    server.AddPacketListener(packet.Type_CSProfile, CSProfileListener)

    log.Fatal(server.Start())
}

type PacketSender interface {
    RemoteAddr()    net.Addr
}

type PacketListener func(PacketContext, proto.Message)

type PacketContext struct {
    Sender  PacketSender
    Server  *GameServer
}

type Message struct {
    From    net.Conn
    Packet  *packet.Packet
}

type GameServer struct {
    addr        string              // listener address
    listener    net.Listener        
    quitCh      chan struct{}       // 0 byte channel (idk why)
    msgCh       chan Message        
    ipconns     map[net.Addr]*Profile 
    idconns     map[uuid.UUID]*Profile 
    p_listeners map[packet.Type][]PacketListener
}

func NewServer(listenerAddr string) *GameServer {
    return &GameServer {
        addr: listenerAddr,
        quitCh: make(chan struct{}),
        msgCh: make(chan Message, 10),
        ipconns: make(map[net.Addr]*Profile),
        idconns: make(map[uuid.UUID]*Profile),
        p_listeners: make(map[packet.Type][]PacketListener),
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
        
        p, err := packet.ReadPacket(buf[:n])
        if err != nil {
            fmt.Println("Failed to read packet:", err)
            continue
        }

        s.msgCh <- Message {
            From: c,
            Packet: p,
        }
    }
}

func (s *GameServer) SendPacket(packet *packet.Packet, profiles ...*Profile) error {
    var err error
    for _, p := range profiles {
        err = p.SendPacket(packet)       
    }

    return err
}

func (s *GameServer) handleMsgs() {
    context := PacketContext { Server: s }
    for msg := range s.msgCh {
        p := msg.Packet
        sender := s.ipconns[msg.From.RemoteAddr()]
        if sender == nil {
            context.Sender = msg.From
        } else {
            context.Sender = sender
        }

        buf := packet.InitPacketBuffer(p.Type)
        err := proto.Unmarshal(p.Data, buf)
        if err != nil {
            fmt.Println("Unmarshal error:", err)
            continue
        }

        listeners := s.p_listeners[p.Type]
        if listeners == nil { continue }
        for _, listener := range listeners {
            listener(context, buf)
        }
    }
}


func (s *GameServer) AddPacketListener(
    packet_type packet.Type,
    listener func(PacketContext, proto.Message),
) {
    listeners := s.p_listeners[packet_type]
    if listeners == nil {
        listeners = make([]PacketListener, 0, 10)
    }

    listeners = append(listeners, listener)
}

func (s *GameServer) RemovePlayerId(uuid uuid.UUID) {
    delete(s.idconns, uuid)
}

func (s *GameServer) RemovePlayerIp(ip net.Addr) {
    delete(s.ipconns, ip)
}


