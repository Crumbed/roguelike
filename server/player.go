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

    profile := NewProfile(context.Sender.(net.Conn), p_conn)
    resp := packet.JoinResponse(true)
    if server.players[0] != nil && server.players[1] != nil {
        resp = !resp
    } else if server.players[0] == nil {
        server.players[0] = profile
    } else {
        server.players[1] = profile
    }
    
    server.SendPacketTo(&resp, profile)
    if resp == false {
        return
    }

    server.ipconns[sender.RemoteAddr()] = profile
    server.Logf("Player connected:\n%s\n", *profile)
}



type Profile struct {
    Conn    net.Conn
    Uuid    uuid.UUID
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
        Name: packet.Name,
    }
}







