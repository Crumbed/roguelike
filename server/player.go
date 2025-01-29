package server

import (
	"fmt"
	"main/server/packet"
	"net"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)




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







