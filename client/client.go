package client

import (
	"fmt"
	"io"
	"main/packet"
	"main/server"
	"net"
	"time"

	"github.com/gen2brain/raylib-go/raylib"
)

const (
    Width   int32           = 600
    Height  int32           = 400
    PW      int32           = 10
    PH      int32           = 60
    P1X     int32           = 5
    P2X     int32           = 600 - PW - 5
    CenterX int32           = 300
    CenterY int32           = 200
    BallS   float32         = 10 // hehe balls
    Update  time.Duration   = time.Millisecond * 50
)

type PlayerN uint8
const (
    Player1 PlayerN = 0
    Player2 PlayerN = 1
)


type Player struct {
    NewPos  int32
    Pos     int32
    Score   uint8
}

func NewClient() *Client {
    return &Client {
        Conn: nil,
        listeners: make(map[packet.PacketType][]packet.PacketListener),
        Players: [2]Player{},
        Started: false,
        BallPos: rl.NewVector2(float32(CenterX) - 5, float32(CenterY) - 5),
    }
}

type Client struct {
    Conn        net.Conn
    listeners   map[packet.PacketType][]packet.PacketListener
    Iam         PlayerN
    Players     [2]Player
    Started     bool
    BallPos     rl.Vector2      
}

func (c *Client) Start() {
    rl.InitWindow(server.Width, server.Height, "Game window")
    rl.SetTargetFPS(60)
    me := &c.Players[c.Iam]
    var other *Player
    if c.Iam == 0 { 
        other = &c.Players[1] 
    } else { other = &c.Players[0] }

    var p2t int32 = 0
    firstStart := true
    for !rl.WindowShouldClose() { // main loop
        if c.Conn == nil { continue }
        if c.Started && firstStart {
            firstStart = false
            go c.UpdateServer()
        }

        if c.Started { 
            keyInput(me) 
            if p2t == 0 && other.Pos != other.NewPos {
                p2t = other.NewPos - other.Pos    
            }
            if other.Pos == other.NewPos { p2t = 0 }
            if p2t != 0 {
                other.Pos += p2t / 2
            }
        }
        c.render()   
    }

    rl.CloseWindow()
}

func (c *Client) UpdateServer() {
    lastPos := int32(0)
    p := &c.Players[c.Iam]

    for {
        if lastPos == p.Pos { continue }
        lastPos = p.Pos
        packet := &packet.PaddleMove {
            PlayerN: uint8(c.Iam),
            Pos: p.Pos,
        }
        
        err := c.SendPacket(packet)
        if err != nil {
            fmt.Println("Failed to send movement update packet:", err)
        }
        time.Sleep(Update)
    }
}

func (c *Client) Connect(ip *net.TCPAddr) error {
    conn, err := net.DialTCP("tcp", nil, ip)
    fmt.Println("Establishing connection to:", ip.IP)
    if err != nil {
        fmt.Println("Failed to dial server: ", err)
        return err
    }

    if c.Conn != nil {
        c.Conn.Close()
    }
    c.Conn = conn
    go c.listen()       
    fmt.Println("Connected to:", c.Conn.RemoteAddr())

    connect := &packet.Connect { Name: "Player 1" }
    c.SendPacket(connect)

    return nil
}

func (self *Client) listen() {
    c := self.Conn
    defer c.Close()
    context := &packet.PacketContext {
        Sender: self.Conn,
        Handler: self,
    }
    buf := make([]byte, 2048)

    for {
        if c == nil { return }
        n, err := c.Read(buf)
        if err != nil {
            if err == io.EOF {
                fmt.Println("Connection lost")
                return 
            }
            fmt.Println("Read err:", err)
            continue
        }
        
        p := packet.ReadPacket(buf[:n])
        self.handlePacket(context, p)
    }
}

func (c *Client) handlePacket(context *packet.PacketContext, p *packet.RawPacket) {
    //if s.stop { return }
    buf := p.Type.InitPacket()
    err := buf.Deserialize(p.Data)
    if err != nil {
        fmt.Println("Packet Read error:", err)
        return
    }

    listeners := c.listeners[p.Type]
    if listeners == nil { return }
    for _, listener := range listeners {
        listener(context, buf)
    }
}

func (c *Client) SendPacket(packet packet.Packet) error {
    data, err := packet.Serialize()
    if err != nil { 
        fmt.Println("Serialize error:", err)
        return err 
    }

    _, err = c.Conn.Write(data)
    return err
}

func (c *Client) AddPacketListener(
    t           packet.PacketType,
    listener    packet.PacketListener,
) {
    listeners := c.listeners[t]
    if listeners == nil {
        listeners = make([]packet.PacketListener, 0, 10)
    }

    listeners = append(listeners, listener)
    c.listeners[t] = listeners
}


func keyInput(p *Player) {
    if rl.IsKeyDown(rl.KeyJ) || rl.IsKeyDown(rl.KeyDown) {
        p.Pos += int32(500 * rl.GetFrameTime())
        if p.Pos + PH >= Height {
            p.Pos = Height - PH
            return
        }
    }
    if rl.IsKeyDown(rl.KeyK) || rl.IsKeyDown(rl.KeyUp) {
        p.Pos += int32(-500 * rl.GetFrameTime())
        if p.Pos <= 0 {
            p.Pos = 0
            return
        }
    }
}

func (c *Client) render() {
    rl.BeginDrawing()
    rl.ClearBackground(rl.Black)
    rl.DrawRectangle(
        CenterX - 1, 0,
        2, Height,
        rl.Gray)

    // player 1 paddle
    rl.DrawRectangle(
        P1X, c.Players[0].Pos,
        PW, PH,
        rl.White)

    if !c.Started { 
        rl.EndDrawing()
        return
    }

    // player 2 paddle
    rl.DrawRectangle(
        P2X, c.Players[1].Pos,
        PW, PH,
        rl.White)

    // ball
    rl.DrawRectangleV(
        c.BallPos,
        rl.NewVector2(BallS, BallS),
        rl.White)

    rl.EndDrawing()
}
