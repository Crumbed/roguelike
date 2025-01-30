package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"main/server/packet"
	"net"
	"os"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)



func StartServer() {
    fmt.Println("Running server")
    server := NewServer(":3000")
    server.AddPacketListener(packet.Type_CSProfile, CSProfileListener)
    err := server.Start()
    if err != nil { log.Fatal(err) }
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
    msgCh       chan Message        
    stop        bool
    cmd         string
    ipconns     map[net.Addr]*Profile 
    idconns     map[uuid.UUID]*Profile 
    p_listeners map[packet.Type][]PacketListener
    logs        []string
}

func NewServer(listenerAddr string) *GameServer {
    return &GameServer {
        addr: listenerAddr,
        msgCh: make(chan Message, 10),
        stop: false,
        cmd: "",
        ipconns: make(map[net.Addr]*Profile),
        idconns: make(map[uuid.UUID]*Profile),
        p_listeners: make(map[packet.Type][]PacketListener),
        logs: make([]string, 0, 10),
    }
}

func (s *GameServer) readCommand() string {
    cmd := s.cmd
    s.cmd = ""
    return cmd
}

func (s *GameServer) Start() error {
    li, err := net.Listen("tcp", s.addr)
    if err != nil { return err }
    defer li.Close()
    s.listener = li

    go s.startReading()
    go s.listen()
    go s.handleMsgs()

    InputLoop: for {
        if s.cmd == "" { continue }
        cmd := s.readCommand()

        switch cmd {
        case "stop":
            fmt.Println("Stopping server...")
            s.stop = true
            time.Sleep(5 * time.Second) // wait for all running threads to stop
            break InputLoop
        case "bg red":
            packet, err := packet.NewPacket(
                packet.Type_SCBGColor,
                packet.NewBgColor(255, 0, 0, 0))
            if err != nil { return err }
            s.SendAllPacket(packet)
        case "bg green":
            packet, err := packet.NewPacket(
                packet.Type_SCBGColor,
                packet.NewBgColor(0, 255, 0, 0))
            if err != nil { return err }
            s.SendAllPacket(packet)
        case "bg blue":
            packet, err := packet.NewPacket(
                packet.Type_SCBGColor,
                packet.NewBgColor(0, 0, 255, 0))
            if err != nil { return err }
            s.SendAllPacket(packet)
        case "debug":
            fmt.Println(*s)
        default: 
            fmt.Println("Invalid command:", cmd)
        }
    }

    close(s.msgCh)
    return nil
}

func (s *GameServer) startReading() {
    stdin := bufio.NewReader(os.Stdin)
    for {
        //fmt.Print("> ")
        input, err := stdin.ReadString('\n')
        if err != nil { log.Fatal("read err:", err) }
        input = input[:len(input)-1]
        s.cmd = input
        //fmt.Println("input:", input)
    }
}

func (s *GameServer) Log(a ...any) {
    str := fmt.Sprintln(a...)
    fmt.Println(str)
    s.logs = append(s.logs, str)
}

func (s *GameServer) Logf(format string, a ...any) {
    str := fmt.Sprintf(format, a...)
    fmt.Print(str)
    s.logs = append(s.logs, str)
}


func (s *GameServer) listen() {
    for {
        if s.stop { return }
        conn, err := s.listener.Accept()
        if err != nil {
            s.Log("Accept error:", err)
            continue
        }

        s.Log("New connection from:", conn.RemoteAddr())
        go s.read(conn)
    }
}

func (s *GameServer) read(c net.Conn) {
    defer c.Close()
    buf := make([]byte, 2048)

    for {
        if s.stop { return }
        n, err := c.Read(buf)
        if err != nil {
            if err == io.EOF {
                ip := c.RemoteAddr()
                s.RemovePlayerIp(ip)
                s.Logf("Player %s has disconnected\n", ip)
                return 
            }
            s.Log("Read err:", err)
            continue
        }
        
        p, err := packet.ReadPacket(buf[:n])
        if err != nil {
            s.Log("Failed to read packet:", err)
            continue
        }

        s.msgCh <- Message {
            From: c,
            Packet: p,
        }
    }
}

func (s *GameServer) SendAllPacket(packet *packet.Packet) error {
    s.Log("Sending packet to all connections")
    var err error
    for _, p := range s.ipconns {
        s.Log("Sending to", p.Conn.RemoteAddr())
        err = p.SendPacket(packet)
        if err != nil { break }
    }

    return err
}

func (s *GameServer) SendPacket(packet *packet.Packet, profiles ...*Profile) error {
    var err error
    for _, p := range profiles {
        err = p.SendPacket(packet)       
        if err != nil { break }
    }

    return err
}

func (s *GameServer) handleMsgs() {
    context := PacketContext { Server: s }
    for msg := range s.msgCh {
        if s.stop { return }
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
            s.Log("Unmarshal error:", err)
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
    s.p_listeners[packet_type] = listeners
}

func (s *GameServer) RemovePlayerId(uuid uuid.UUID) {
    delete(s.idconns, uuid)
}

func (s *GameServer) RemovePlayerIp(ip net.Addr) {
    delete(s.ipconns, ip)
}


