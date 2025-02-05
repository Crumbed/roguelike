package server

import (
	"fmt"
	"main/packet"
	"net"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)


func CSProfileListener(context packet.PacketContext, data proto.Message) {
    sender := context.Sender.(net.Conn)
    server := context.Handler.(*GameServer)

    profile := NewProfile(context.Sender.(net.Conn), data.(*packet.Profile))
    server.ipconns[sender.RemoteAddr()] = profile
    server.idconns[profile.Uuid] = profile

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
        self.Conn.RemoteAddr(), 
        self.Name,
        self.Uuid,
    )
}

func (self *Profile) RemoteAddr() net.Addr {
    return self.Conn.RemoteAddr()
}

func (self *Profile) SendPacketBytes(data []byte) error {
    _, err := self.Conn.Write(data)
    return err
}

func (self *Profile) SendPacket(packet *packet.Packet) error {
    data, err := proto.Marshal(packet)
    if err != nil { return err }

    _, err = self.Conn.Write(data)
    return err
}

func NewProfile(conn net.Conn, packet *packet.Profile) *Profile {
    return &Profile {
        Conn: conn,
        Uuid: uuid.New(),
        Name: packet.Name,
    }
}







