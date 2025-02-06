package packet

import "bytes"


// >5 | PacketType 1, NameLen 4, Name ~
type Connect struct {
    Name    string
}

func (p *Connect) GetType() PacketType { return CSConnect }
func (p *Connect) Serialize() ([]byte, error) {
    buf := bytes.NewBuffer(make([]byte, 0, len(p.Name) + 5))
    _ = buf.WriteByte(byte(CSConnect))
    err := SerializeString(buf, p.Name)
    return buf.Bytes(), err
}
func (p *Connect) Deserialize(data []byte) error {
    buf := bytes.NewBuffer(data)
    name, err := DeserializeString(buf)
    p.Name = name
    return err
}


// 3 | PacketType 1, Response 1, PlayerN 1
type JoinResponse struct {
    Response    bool
    PlayerN     uint8
}
func (p *JoinResponse) IsOk() bool { return p.Response == true }

func (p *JoinResponse) GetType() PacketType { return SCJoinResponse }
func (p *JoinResponse) Serialize() ([]byte, error) { 
    buf := bytes.NewBuffer(make([]byte, 0, 3))
    _ = buf.WriteByte(byte(SCJoinResponse))
    err := SerializeBool(buf, p.Response)
    err = buf.WriteByte(p.PlayerN)
    return buf.Bytes(), err
}
func (p *JoinResponse) Deserialize(data []byte) error {
    buf := bytes.NewBuffer(data)
    r, err := DeserializeBool(buf)
    w, err := buf.ReadByte()

    p.Response = r
    p.PlayerN = w
    return err
}

// 2 | PacketType 1, PlayerN 1
type AddPlayer struct {
    PlayerN uint8
}

/*
func (p *AddPlayer) GetType() PacketType { return SCAddPlayer }
func (p *AddPlayer) Serialize() ([]byte, error) {
    buf := bytes.NewBuffer(make([]byte, 0, 2))
    _ = buf.WriteByte(byte(SCAddPlayer))
    err := buf.WriteByte(p.PlayerN)
    return buf.Bytes(), err
}
func (p *AddPlayer) Deserialize(data []byte) error {
    buf := bytes.NewBuffer(data)
    n, err := buf.ReadByte()

    p.PlayerN = n
    return err
}
*/




