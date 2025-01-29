package server

import (
	"main/server/packet"
	"net"

	"github.com/google/uuid"
)




type Profile struct {
    Ip      net.Addr
    Uuid    uuid.UUID
    Name    string
}

func NewProfile(ip net.Addr, packet *packet.Profile) *Profile {
    return &Profile {
        Ip: ip,
        Uuid: uuid.New(),
        Name: packet.Name,
    }
}







