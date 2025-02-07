package server

import (
	"fmt"
	"main/packet"
	"net"

	"github.com/google/uuid"
)


func CSConnectListener(context *packet.PacketContext, data packet.Packet) {
    sender := context.Sender.(net.Conn)
    server := context.Handler.(*GameServer)
    p_conn := data.(*packet.Connect)

    start := false
    profile := NewProfile(context.Sender.(net.Conn), p_conn)
    resp := &packet.JoinResponse { Response: true }
    if server.Players[0] != nil && server.Players[1] != nil {
        resp.Response = false 
    } else if server.Players[0] == nil {
        server.Players[0] = profile
        resp.PlayerN = 0
        fmt.Println("Player 1 joined")
        if server.Players[1] != nil { start = true }
    } else {
        server.Players[1] = profile
        resp.PlayerN = 1
        fmt.Println("Player 2 joined")
        if server.Players[0] != nil { start = true }
    }
    
    server.SendPacketTo(resp, profile)
    if !resp.IsOk() {
        server.Logf("Kicking %s because game is full", sender)
        return
    }

    server.ipconns[sender.RemoteAddr()] = profile
    server.Logf("Player connected:\n%s\n", profile)
    if start {
        server.State.Started = true
        err := server.SendPacket(&packet.GameStart{})
        if err != nil {
            fmt.Println("Failed to start game:", err)
        }
    }
}

func SSPaddleMoveListener(context *packet.PacketContext, data packet.Packet) {
    server := context.Handler.(*GameServer)
    move := data.(*packet.PaddleMove)

    var otherPlayer uint8
    if move.PlayerN == 0 {
        server.State.P1.Move(move.Pos)
        otherPlayer = 1
    } else {
        server.State.P2.Move(move.Pos)
    }

    server.SendPacketTo(data, server.Players[otherPlayer])
}

func SSGameStartListener(context *packet.PacketContext, data packet.Packet) {
    sender := context.Sender.(*Profile)
    sender.Started = true
}



type Profile struct {
    Conn    net.Conn
    Uuid    uuid.UUID
    Started bool
    Name    string
}

func (self *Profile) String() string {
    return fmt.Sprintf(
        "(Ip: %s) %s, uuid=%s",
        self.Conn.RemoteAddr().String(), 
        self.Name,
        self.Uuid.String(),
    )
}

func (self *Profile) RemoteAddr() net.Addr {
    return self.Conn.RemoteAddr()
}

func (self *Profile) SendPacketBytes(data []byte) error {
    _, err := self.Conn.Write(data)
    return err
}

func (self *Profile) SendPacket(packet packet.Packet) error {
    data, err := packet.Serialize()
    if err != nil { return err }

    _, err = self.Conn.Write(data)
    return err
}

func NewProfile(conn net.Conn, packet *packet.Connect) *Profile {
    return &Profile {
        Conn: conn,
        Uuid: uuid.New(),
        Started: false,
        Name: packet.Name,
    }
}







