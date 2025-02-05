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
        p_listeners: make(map[packet.PacketType][]packet.PacketListener),
    }
}

type Client struct {
    Conn        net.Conn
    p_listeners map[packet.PacketType][]packet.PacketListener
    MyPos       Vector2
}

func (c *Client) Start() {
    go c.render()

    for { // main loop
        if c.Conn == nil { continue }
        
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

    listeners := c.p_listeners[p.Type]
    if listeners == nil { return }
    for _, listener := range listeners {
        listener(context, buf)
    }
}

func (c *Client) SendPacket(packet packet.Packet) error {
    data, err := packet.Serialize()
    if err != nil { return err }

    _, err = c.Conn.Write(data)
    return err
}

func (c *Client) render() {
    InitWindow(server.Width, server.Height, "Game window")
    SetTargetFPS(60)
    black := NewColor(0, 0, 0, 0)

    for !WindowShouldClose() {
        BeginDrawing()
        ClearBackground(black)

        EndDrawing()
    }

    CloseWindow()
}
