package packet

import "google.golang.org/protobuf/proto"


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




