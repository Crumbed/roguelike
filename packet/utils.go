package packet

import (
	"net"

	"google.golang.org/protobuf/proto"
)




type PacketSender interface {
    RemoteAddr()    net.Addr
}

type PacketHandler interface {
    SendPacket(packet *Packet)   error
}

type PacketListener func(*PacketContext, proto.Message)

type PacketContext struct {
    Sender  PacketSender
    Handler PacketHandler
}

