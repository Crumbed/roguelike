package packet

import (
	"log"

	"google.golang.org/protobuf/proto"
)


func ReadPacket(bytes []byte) (*Packet, error) {
    p := &Packet{}
    return p, proto.Unmarshal(bytes, p)
}

func InitPacketBuffer(kind Type) proto.Message {
    var data proto.Message
    switch kind {
    case Type_CSProfile:
        data = &Profile{}
    default:
        log.Fatal("Invalid Packet Type:", kind)
    }

    return data
}

func NewPacket(kind Type, data proto.Message) (*Packet, error) {
    bytes, err := proto.Marshal(data)
    if err != nil { return nil, err }

    return &Packet {
        Type: kind,
        Data: bytes,
    }, nil
}

func CreatePacket(kind Type, data proto.Message) ([]byte, error) {
    packet, err := NewPacket(kind, data)
    if err != nil { return nil, err }

    bytes, err := proto.Marshal(packet)
    if err != nil { return nil, err }
    return bytes, nil
}




