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
    //client.Reset()
    ConnectError = "game is full."
}

func SCGameStartListener(context *packet.PacketContext, p packet.Packet) {
    client := context.Handler.(*Client)
    client.Started = true
    client.SendPacket(p)
    fmt.Println("Starting game...")
}

func CCPaddleMoveListener(context *packet.PacketContext, data packet.Packet) {
    client := context.Handler.(*Client)
    update := data.(*packet.PaddleMove)

    client.Players[update.PlayerN].Target = float32(update.Pos)
    //client.Players[update.PlayerN].Pos = update.Pos
}

func SCBallMoveListener(context *packet.PacketContext, data packet.Packet) {
    client := context.Handler.(*Client)
    ballPos := data.(*packet.BallMove)

    client.Ball.NewPos.X = ballPos.X
    client.Ball.NewPos.Y = ballPos.Y
}

func SCScoreListener(context *packet.PacketContext, data packet.Packet) {
    client := context.Handler.(*Client)
    score := data.(*packet.Score)

    client.Players[0].Score = score.P1
    client.Players[1].Score = score.P2
}

