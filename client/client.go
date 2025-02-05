package client

import (
	"fmt"
	"io"
	"main/packet"
	"main/server"
	"net"

	. "github.com/gen2brain/raylib-go/raylib"
)



func NewClient() *Client {
    return &Client {
        Conn: nil,
        listeners: make(map[packet.PacketType][]packet.PacketListener),
    }
}

type Client struct {
    Conn        net.Conn
    listeners   map[packet.PacketType][]packet.PacketListener
    MyPos       Vector2
}

func (c *Client) Start() {
    InitWindow(server.Width, server.Height, "Game window")
    SetTargetFPS(60)

    for !WindowShouldClose() { // main loop
        c.render()   
        if c.Conn == nil { continue }
    }

    CloseWindow()
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

func (c *Client) render() {
    BeginDrawing()
    ClearBackground(Black)

    EndDrawing()
}
