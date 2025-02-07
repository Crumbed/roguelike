package packet

import "bytes"

// 1 | PacketType 1
type GameStart struct {} // doesnt need any data
func (p *GameStart) GetType() PacketType { return BWGameStart }
func (p *GameStart) Serialize() ([]byte, error) {
    bytes := make([]byte, 1, 1)
    bytes[0] = byte(BWGameStart)
    return bytes, nil
}
func (p *GameStart) Deserialize(data []byte) error { return nil }


// 6 | PacketType 1, PlayerN 1, Pos 4
type PaddleMove struct {
    PlayerN uint8
    Pos     int32
}
func (p *PaddleMove) GetType() PacketType { return BWPaddleMove }
func (p *PaddleMove) Serialize() ([]byte, error) {
    buf := bytes.NewBuffer(make([]byte, 0, 6))
    err := buf.WriteByte(byte(BWPaddleMove))
    err = buf.WriteByte(p.PlayerN)
    err = SerializeInt(buf, p.Pos)

    return buf.Bytes(), err
}
func (p *PaddleMove) Deserialize(data []byte) error {
    buf := bytes.NewBuffer(data)
    n, err := buf.ReadByte()
    pos, err := DeserializeInt(buf)

    p.PlayerN = n
    p.Pos = pos
    return err
}


// 9 | PacketType 1, X 4, Y 4
type BallMove struct {
    X   float32
    Y   float32
}
func (p *BallMove) GetType() PacketType { return SCBallMove }
func (p *BallMove) Serialize() ([]byte, error) {
    buf := bytes.NewBuffer(make([]byte, 0, 9))
    err := buf.WriteByte(byte(SCBallMove))
    err = SerializeFloat(buf, p.X)
    err = SerializeFloat(buf, p.Y)

    return buf.Bytes(), err
}
func (p *BallMove) Deserialize(data []byte) error {
    buf := bytes.NewBuffer(data)
    x, err := DeserializeFloat(buf)
    y, err := DeserializeFloat(buf)

    p.X = x
    p.Y = y
    return err
}


// 3 | PacketType 1, P1 1, P2 1
type Score struct {
    P1  uint8
    P2  uint8
}
func (p *Score) GetType() PacketType { return SCScore }
func (p *Score) Serialize() ([]byte, error) {
    buf := bytes.NewBuffer(make([]byte, 0, 3))
    err := buf.WriteByte(byte(SCScore))
    err = buf.WriteByte(p.P1)
    err = buf.WriteByte(p.P2)

    return buf.Bytes(), err
}
func (p *Score) Deserialize(data []byte) error {
    buf := bytes.NewBuffer(data)
    p1, err := buf.ReadByte()
    p2, err := buf.ReadByte()

    p.P1 = p1
    p.P2 = p2
    return err
}







