package client

import (
	"fmt"
	"io"
	"main/packet"
	"net"

	. "github.com/gen2brain/raylib-go/raylib"
	"google.golang.org/protobuf/proto"
)




func NewClient() *Client {
    return &Client {
        Conn: nil,
        p_listeners: make(map[packet.Type][]packet.PacketListener),
    }
}

type Client struct {
    Conn        net.Conn
    p_listeners map[packet.Type][]packet.PacketListener
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
        
        p, err := packet.ReadPacket(buf[:n])
        if err != nil {
            fmt.Println("Failed to read packet:", err)
            continue
        }

        self.handlePacket(context, p)
    }
}

func (c *Client) handlePacket(context *packet.PacketContext, p *packet.Packet) {
    //if s.stop { return }
    buf := packet.InitPacketBuffer(p.Type)
    err := proto.Unmarshal(p.Data, buf)
    if err != nil {
        fmt.Println("Unmarshal error:", err)
        return
    }

    listeners := c.p_listeners[p.Type]
    if listeners == nil { return }
    for _, listener := range listeners {
        listener(context, buf)
    }
}

func (c *Client) SendPacket(packet *packet.Packet) error {
    data, err := proto.Marshal(packet)
    if err != nil { return err }

    _, err = c.Conn.Write(data)
    return err
}

func (c *Client) render() {
    InitWindow(600, 400, "Game window")
    SetTargetFPS(60)
    black := NewColor(0, 0, 0, 0)

    for !WindowShouldClose() {
        BeginDrawing()
        ClearBackground(black)

        EndDrawing()
    }

    CloseWindow()
}
