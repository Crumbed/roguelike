package packet

import (
	"net"
)




type PacketSender interface {
    RemoteAddr()    net.Addr
}

type PacketHandler interface {
    SendPacket(packet *Packet)   error
}

type PacketListener func(*PacketContext, *Packet)

type PacketContext struct {
    Sender  PacketSender
    Handler PacketHandler
}

