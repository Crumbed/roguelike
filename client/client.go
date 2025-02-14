package client

import (
	"fmt"
	"io"
	"log"
	"main/packet"
	"math"
	"net"
	"os"
	"time"

	"github.com/gen2brain/raylib-go/raylib"
)

const (
    Width   int32           = 600
    Height  int32           = 400
    FSize   int32           = 50
    PW      int32           = 10
    PH      int32           = 80
    P1X     int32           = 5
    P2X     int32           = 600 - PW - 5
    CenterX int32           = 300
    CenterY int32           = 200
    BallS   float32         = 10 // hehe balls
    Update  time.Duration   = time.Millisecond * 10
)

type PlayerN uint8
const (
    Player1 PlayerN = 0
    Player2 PlayerN = 1
)


type Player struct {
    Target  float32
    Pos     float32
    Score   uint8
}
func (p *Player) interp() {
    if int32(p.Pos) == int32(p.Target) { return }
    p.Pos = rl.Lerp(p.Pos, p.Target, 0.5)
}

func (p *Player) render(n PlayerN) {
    var x int32
    if n == Player1 { 
        x = P1X 
    } else if n == Player2 {
        x = P2X 
    }

    rl.DrawRectangle(
        x, int32(p.Pos),
        PW, PH,
        rl.White)
}

var BallRect rl.Vector2 = rl.NewVector2(BallS, BallS)
type Ball struct {
    NewPos  rl.Vector2
    Pos     rl.Vector2
}
func (b *Ball) interp() {
    if rl.Vector2Equals(b.Pos, b.NewPos) { return }
    b.Pos = rl.Vector2Lerp(b.Pos, b.NewPos, 0.5)
}

func (b *Ball) render() {
    rl.DrawRectangleV(b.Pos, BallRect, rl.White)
}


func NewClient() *Client {
    return &Client {
        Conn: nil,
        listeners: make(map[packet.PacketType][]packet.PacketListener),
        screen: InitScreen(),
        serverIp: readLastIp(),
        Players: [2]Player{},
        Started: false,
        Ball: Ball {
            NewPos: rl.NewVector2(float32(CenterX) - 5, float32(CenterY) - 5),
            Pos: rl.NewVector2(float32(CenterX) - 5, float32(CenterY) - 5),
        },
    }
}

func (c *Client) Reset() {
    c.Conn = nil
    c.screen.disp = StartMenu
    c.Players = [2]Player{}
    c.Started = false
    c.Ball = Ball {
        NewPos: rl.NewVector2(float32(CenterX) - 5, float32(CenterY) - 5),
        Pos: rl.NewVector2(float32(CenterX) - 5, float32(CenterY) - 5),
    }
}

type Client struct {
    Conn        net.Conn
    listeners   map[packet.PacketType][]packet.PacketListener
    screen      Screen
    serverIp    []byte
    Iam         PlayerN
    Players     [2]Player
    Started     bool
    Ball        Ball
}
func (c *Client) GetOtherPlayer() *Player {
    switch c.Iam {
    case Player1: return &c.Players[Player2]
    case Player2: return &c.Players[Player1]
    default:
        fmt.Println("HOW DID THIS HAPPEN YOU FUCKING IDIOT")
        return nil
    }
}

func (c *Client) Start() {
    go c.UpdateServer()

    for !rl.WindowShouldClose() { // main loop
        switch c.screen.disp {
        case StartMenu: ipInput(c)
        case Game: moveInput(&c.Players[c.Iam]) 
        }

        if c.Started {
            c.Ball.interp()
            c.GetOtherPlayer().interp()
        }

        c.render()   
    }

    rl.UnloadFont(c.screen.font)
    rl.CloseWindow()
}

func (c *Client) UpdateServer() {
    lastPos := float32(0)
    var p *Player

    for {
        if c.Conn == nil { continue }
        p = &c.Players[c.Iam]
        if lastPos == p.Pos { continue }
        lastPos = p.Pos
        packet := &packet.PaddleMove {
            PlayerN: uint8(c.Iam),
            Pos: int32(p.Pos),
        }
        
        err := c.SendPacket(packet)
        if err != nil {
            fmt.Println("Failed to send movement update packet:", err)
        }
        time.Sleep(Update)
    }
}

func (c *Client) Connect() error {
    ip, err := net.ResolveTCPAddr("tcp", string(c.serverIp))
    if err != nil {
        fmt.Println("Failed to resolve tcp addr: ", err)
        return err
    }
    conn, err := net.DialTCP("tcp", nil, ip)
    fmt.Println("Establishing connection to:", ip.IP)
    if err != nil {
        fmt.Println("Failed to dial server: ", err)
        return err
    }
    writeServer(string(c.serverIp))
    c.screen.disp = Game

    if c.Conn != nil {
        c.Conn.Close()
    }
    c.Conn = conn
    go c.listen()       
    fmt.Println("Connected to:", c.Conn.RemoteAddr())

    connect := &packet.Connect { Name: "Player" }
    err = c.SendPacket(connect)
    if err != nil {
        fmt.Println(err) // jensen made past this point, then died
        // Everything is working now, even though i didnt change anything...
    }

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
                self.Reset()
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


func moveInput(p *Player) {
    if rl.IsKeyDown(rl.KeyJ) || rl.IsKeyDown(rl.KeyDown) || rl.IsKeyDown(rl.KeyS) {
        p.Pos += float32(math.Trunc(500 * float64(rl.GetFrameTime())))
        if int32(p.Pos) + PH >= Height {
            p.Pos = float32(Height - PH)
            return
        }

        if rl.IsCursorHidden() == false {
            rl.HideCursor()
        }
    }
    if rl.IsKeyDown(rl.KeyK) || rl.IsKeyDown(rl.KeyUp) || rl.IsKeyDown(rl.KeyW) {
        p.Pos += float32(math.Trunc(-500 * float64(rl.GetFrameTime())))
        if p.Pos <= 0 {
            p.Pos = 0
            return
        }

        if rl.IsCursorHidden() == false {
            rl.HideCursor()
        }
    }
}

func writeServer(ip string) {
    err := os.WriteFile("last_server", []byte(ip), 0644)
    if err != nil {
        log.Fatal("Failed to write last server...")
    }
}


func readLastIp() []byte {
    buf, err := os.ReadFile("last_server")
    if err != nil {
        return make([]byte, 0, 21)
    }

    return buf
}
