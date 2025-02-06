package packet


// 1 | PacketType 1
type GameStart struct {} // doesnt need any data
func (p *GameStart) GetType() PacketType { return SCGameStart }
func (p *GameStart) Serialize() ([]byte, error) {
    bytes := make([]byte, 1, 1)
    bytes[0] = byte(SCGameStart)
    return bytes, nil
}
func (p *GameStart) Deserialize(data []byte) error { return nil }


