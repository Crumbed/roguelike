package client

import (
	"fmt"
	"main/packet"
)



func SCJoinResponseListener(context *packet.PacketContext, p packet.Packet) {
    client := context.Handler.(*Client)
    response := p.(*packet.JoinResponse)
    if response.IsOk() {
        fmt.Println("Successfully connected to game")
        client.Iam = PlayerN(response.PlayerN)
        return 
    }

    fmt.Println("Game is full, connection refused...")
    client.Conn.Close()
}

func SCGameStartListener(context *packet.PacketContext, _ packet.Packet) {
    client := context.Handler.(*Client)
    client.Started = true
    fmt.Println("Starting game...")
}


