package client

import (
	"fmt"
	"io"
	"main/packet"
	"net"

	. "github.com/gen2brain/raylib-go/raylib"
)




func NewClient() *Client {
    return &Client {
        Conn: nil,
    }
}

type Client struct {
    Conn    net.Conn
}

func (c *Client) Start() {
    go c.render()

    for c.Conn != nil {

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

    return nil
}

func (self *Client) listen() {
    c := self.Conn
    defer c.Close()
    buf := make([]byte, 2048)

    for {
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

        pktBuf := packet.InitPacketBuffer(p.Type)
        err = proto.Unmarshal(p.Data, pktBuf)
        if err != nil {
            fmt.Println("Unmarshal error:", err)
            continue
        }

        if p.Type == packet.Type_SCBGColor {
            fmt.Println("Changing color")
            c.changeColor(pktBuf.(*packet.BackgroundColor))
        }
    }
}

func (c *Client) handlePackets() {
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
