package client

import (
	"fmt"
	"main/packet"
	"net"
)



func SCJoinResponseListener(context *packet.PacketContext, p packet.Packet) {
    response := p.(*packet.JoinResponse)
    if response.Is(true) {
        fmt.Println("Successfully connected to game")
        return 
    }

    fmt.Println("Game is full, connection refused...")
    context.Sender.(net.Conn).Close()
}



