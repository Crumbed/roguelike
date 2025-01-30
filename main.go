package main

import (
	"fmt"
	"image/color"
	"io"

	//"io"
	"log"
	"main/server"
	"main/server/packet"
	"net"
	"os"

	. "github.com/gen2brain/raylib-go/raylib"
	"google.golang.org/protobuf/proto"
	//"google.golang.org/protobuf/proto"
)

type ProgramType int
const (
    Server ProgramType = iota
    Client
)

func getProgramType() ProgramType {
    args := os.Args
    if len(args) == 1 {
        fmt.Println("No argument found, defaulting to client...")
        return Client
    } else if args[1] == "server" || args[1] == "s" {
        return Server
    } else {
        fmt.Println("Invalid argument, defaulting to client...")
        return Client
    }
}

const (
    Host = "localhost"
    Port = "3000"
    Type = "tcp"
)

func main() {
    fmt.Println("Hello, world")
    p_type := getProgramType()

    if p_type == Server {
        server.StartServer()
        return
    }

    // client   
    tcpServer, err := net.ResolveTCPAddr(Type, Host + ":" + Port)
    if err != nil {
        log.Fatal("Failed to resolve tcp addr: ", err)
    }

    conn, err := net.DialTCP(Type, nil, tcpServer)
    fmt.Println("Establishing connection to:", Host)
    if err != nil {
        log.Fatal("Failed to dial server: ", err)
    }
    defer conn.Close()

    profile := &packet.Profile { Name: "Kai" }
    data, err := packet.CreatePacket(packet.Type_CSProfile, profile)
    if err != nil {
        log.Fatal("Marshal error:", err)
    }

    fmt.Println("Signing in to:", conn.RemoteAddr())
    _, err = conn.Write(data)
    if err != nil {
        log.Fatal("Failed to write data to connection: ", err)
    }
    fmt.Println("Success!")
    client := &GameClient {
        Conn: conn,
        BgColor: SkyBlue,
    }
    go client.listen()

    InitWindow(600, 400, "Game window")
    SetTargetFPS(60)
    for !WindowShouldClose() {
        BeginDrawing()
        ClearBackground(client.BgColor)

        EndDrawing()
    }

    CloseWindow()
}


type GameClient struct {
    Conn        net.Conn
    BgColor     color.RGBA
}

func (c *GameClient) listen() {
    buf := make([]byte, 2048)
    for {
        n, err := c.Conn.Read(buf)
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


func (c *GameClient) changeColor(p_color *packet.BackgroundColor) {
    rgb := p_color.Rgba
    col := NewColor(rgb[0], rgb[1], rgb[2], rgb[3])
    c.BgColor = col
}










