package server

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"log"
	"main/packet"
	"net"
	"os"
	"slices"
)



func StartServer(port string) {
    fmt.Println("Running server")
    server := NewServer(port)
    server.AddPacketListener(packet.CSConnect, CSConnectListener)
    server.AddPacketListener(packet.BWPaddleMove, SSPaddleMoveListener)
    server.AddPacketListener(packet.BWGameStart, SSGameStartListener)

    server.AddUpdateFn(&UpdateBall)
    server.AddUpdateFn(&ConfirmReady)
    server.AddUpdateFn(&SendBallMove)
    err := server.Start()
    if err != nil { log.Fatal(err) }
}


type Message struct {
    From    net.Conn
    Packet  *packet.RawPacket
}

type GameServer struct {
    addr        string          // listener address
    listener    net.Listener        
    msgCh       chan Message        
    stop        bool
    cmd         string
    ipconns     map[net.Addr]*Profile 
    Players     [2]*Profile     // Players 1 & 2 profiles
    logs        []string
    State       *GameState
    p_listeners map[packet.PacketType][]packet.PacketListener
    updateFns   *list.List
    DeltaTime   float64         // DeltaTime seconds
}

func NewServer(listenerAddr string) *GameServer {
    return &GameServer {
        addr: listenerAddr,
        msgCh: make(chan Message, 10),
        stop: false,
        cmd: "",
        ipconns: make(map[net.Addr]*Profile),
        p_listeners: make(map[packet.PacketType][]packet.PacketListener),
        updateFns: list.New(),
        logs: make([]string, 0, 10),
        Players: [2]*Profile { nil, nil },
        State: NewGame(),
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

    firstStart := true
    for {
        if firstStart && s.State.Started {
            firstStart = false
            go s.UpdateClients()
        }
    }

    //fmt.Println("WHY TF IS THIS HAPPENING")
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
    s.Log("Listening for connections on:", s.listener.Addr().String())
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
                s.Log("Stopping server...")
                os.Exit(0)
                return 
            }
            s.Log("Read err:", err)
            continue
        }
        
        p := packet.ReadPacket(buf[:n])
        s.msgCh <- Message {
            From: c,
            Packet: p,
        }
    }
}

func (s *GameServer) SendPacket(packet packet.Packet) error {
    //s.Log("Sending packet to all connections")
    var err error
    for _, p := range s.ipconns {
        //s.Log("Sending to", p.Conn.RemoteAddr())
        err = p.SendPacket(packet)
        if err != nil { break }
    }

    return err
}

func (s *GameServer) SendPacketEx(packet packet.Packet, exclude ...*Profile) error {
    var err error
    for _, p := range s.ipconns {
        if slices.Contains(exclude, p) { continue }
        //s.Log("Sending to", p.Conn.RemoteAddr())
        err = p.SendPacket(packet)
        if err != nil { break }
    }

    return err
}

func (s *GameServer) SendPacketTo(packet packet.Packet, profiles ...*Profile) error {
    var err error
    for _, p := range profiles {
        err = p.SendPacket(packet)       
        if err != nil { break }
    }

    return err
}

func (s *GameServer) handleMsgs() {
    context := packet.PacketContext { Handler: s }
    for msg := range s.msgCh {
        if s.stop { return }
        p := msg.Packet
        sender := s.ipconns[msg.From.RemoteAddr()]
        if sender == nil {
            context.Sender = msg.From
        } else {
            context.Sender = sender
        }

        buf := p.Type.InitPacket()
        err := buf.Deserialize(p.Data)
        if err != nil {
            s.Log("Read Packet error:", err)
            continue
        }

        listeners := s.p_listeners[p.Type]
        if listeners == nil { continue }
        for _, listener := range listeners {
            listener(&context, buf)
        }
    }
}

func (s *GameServer) AddUpdateFn(up *UpdateFn) {
    s.updateFns.PushBack(up)
}

func (s *GameServer) AddPacketListener(
    packet_type packet.PacketType,
    listener packet.PacketListener,
) {
    listeners := s.p_listeners[packet_type]
    if listeners == nil {
        listeners = make([]packet.PacketListener, 0, 10)
    }

    listeners = append(listeners, listener)
    s.p_listeners[packet_type] = listeners
}

func (s *GameServer) RemovePlayerIp(ip net.Addr) {
    delete(s.ipconns, ip)
}


