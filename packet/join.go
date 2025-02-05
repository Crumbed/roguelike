package packet

import "bytes"



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


type JoinResponse bool
func NewJoinResponse() *JoinResponse {
    var p JoinResponse
    return &p
}

func (p *JoinResponse) GetType() PacketType { return SCJoinResponse }
func (p *JoinResponse) Serialize() ([]byte, error) { 
    buf := bytes.NewBuffer(make([]byte, 0, 2))
    _ = buf.WriteByte(byte(SCJoinResponse))
    err := SerializeBool(buf, bool(*p))
    return buf.Bytes(), err
}
func (p *JoinResponse) Deserialize(data []byte) error {
    buf := bytes.NewBuffer(data)
    b, err := DeserializeBool(buf)
    *p = JoinResponse(b)
    return err
}





