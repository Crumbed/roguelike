package packet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"slices"
)


type PacketSender interface {
    RemoteAddr()    net.Addr
}

type PacketHandler interface {
    SendPacket(packet Packet)   error
}

type PacketListener func(*PacketContext, Packet)

type PacketContext struct {
    Sender  PacketSender
    Handler PacketHandler
}


type PacketType uint8
const (
    CSConnect PacketType = iota // request game join
    SCJoinResponse              // game join response
)

func (t *PacketType) InitPacket() Packet {
    var data Packet
    switch *t {
    case CSConnect: data = &Connect{}
    case SCJoinResponse: data = NewJoinResponse()
    default: 
        fmt.Println("Invalid PacketType:", *t)
        data = nil
    }

    return data
}

type RawPacket struct {
    Type    PacketType
    Data    []byte
}

func ReadPacket(bytes []byte) *RawPacket {
    return &RawPacket{
        Type: PacketType(bytes[0]),
        Data: bytes[1:],
    }
}

type Packet interface {
    GetType()           PacketType
    Serialize()         ([]byte, error) // including type
    Deserialize([]byte) error           // not including type
}



func SerializeString(buf *bytes.Buffer, str string) error {
    err := SerializeInt(buf, int32(len(str)))   
    data := []byte(str)
    slices.Reverse(data)
    _, err = buf.Write(data)
    //err = binary.Write(buf, binary.BigEndian, str)
    return err
}

func SerializeBool(buf *bytes.Buffer, b bool) error {
    err := binary.Write(buf, binary.BigEndian, b)
    return err
}

func SerializeInt(buf *bytes.Buffer, i int32) error {
    err := binary.Write(buf, binary.BigEndian, i)
    return err
}

func SerializeUInt(buf *bytes.Buffer, i uint32) error {
    err := binary.Write(buf, binary.BigEndian, i)
    return err
}

func SerializeFloat(buf *bytes.Buffer, f float32) error {
    err := binary.Write(buf, binary.BigEndian, f)
    return err
}

func DeserializeString(buf *bytes.Buffer) (string, error) {
    l, err := DeserializeInt(buf)
    if err != nil { return "", err }
    bytes := make([]byte, l, l)
    _, err = buf.Read(bytes)
    if err != nil { return "", err }

    slices.Reverse(bytes)
    return string(bytes), nil
}

func DeserializeBool(buf *bytes.Buffer) (bool, error) {
    b, err := buf.ReadByte()
    if err != nil { return false, err }
    return b != 0, nil
}

func DeserializeInt(buf *bytes.Buffer) (int32, error) {
    i, err := DeserializeUInt(buf)
    if err != nil { return 0, err }
    return int32(i), nil
}

func DeserializeFloat(buf *bytes.Buffer) (float32, error) {
    f, err := DeserializeUInt(buf)
    if err != nil { return 0, err }
    return math.Float32frombits(f), nil
}

func DeserializeUInt(buf *bytes.Buffer) (uint32, error) {
    bytes := [4]byte{}
    _, err := buf.Read(bytes[:])
    if err != nil { return 0, err }
    return binary.BigEndian.Uint32(bytes[:]), nil
}









