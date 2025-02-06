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
