package server

import (
	"fmt"
	"main/packet"
	"time"
)

type UpStatus uint8
const (
    Ok UpStatus = iota
    Remove
)

func NewUpdate(fn func(*GameServer) UpStatus, rate uint8) UpdateFn {
    return UpdateFn { Rate: rate, Update: fn }
}
type UpdateFn struct {
    Rate    uint8                       // How many ticks between calls
    Passed  uint8                       // How many ticks have passed since last call
    Update  func(*GameServer) UpStatus  // Update function -> status of UpdateFn
}
func (up *UpdateFn) Check(server *GameServer) UpStatus {
    up.Passed += 1
    if up.Passed != up.Rate { return Ok }
    up.Passed = 0
    return up.Update(server)
}

var ConfirmReady = NewUpdate(func(s *GameServer) UpStatus {
    var p *Profile = nil
    if s.Players[0].Started && s.Players[1].Started {
        return Remove
    } else if s.Players[0].Started {
        p = s.Players[1]
    } else {
        p = s.Players[0]
    }

    fmt.Println("Resending start packet to:", p.Conn.RemoteAddr())
    s.SendPacketTo(&packet.GameStart{}, p) 
    return Ok
}, 15)

var SendBallMove = NewUpdate(func(s *GameServer) UpStatus {
    ball := &s.State.Ball
    s.SendPacket(&packet.BallMove {
        X: ball.Pos.X,
        Y: ball.Pos.Y,
    })

    return Ok
}, 2)





const (
    TargetTPS = 60
    TargetTickTime = 1.0 / float32(TargetTPS)
    Cooldown time.Duration = time.Millisecond * 10
)

func (s *GameServer) UpdateClients() {
    var elapsed float32 = 0.0
    last := time.Now().UnixNano() / int64(time.Millisecond)

    for {
        current := time.Now().UnixNano() / int64(time.Millisecond)
        s.DeltaTime = float32(current - last) / 1000 // Delta time seconds
        last = current
        elapsed += s.DeltaTime

        s.updateTicks(&elapsed)
        time.Sleep(Cooldown)
    }
}


func (s *GameServer) updateTicks(elapsed *float32) {
    for *elapsed >= TargetTickTime {
        at := s.updateFns.Front()
        for at != nil {
            status := at.Value.(*UpdateFn).Check(s)
            if status == Remove {
                s.updateFns.Remove(at)
            }

            at = at.Next()
        }

        *elapsed -= TargetTickTime
    }
}
