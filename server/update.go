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
    if !s.State.Running { return Ok }
    if s.Players[0].Started && s.Players[1].Started {
        return Ok
    }
    var resend []*Profile = nil
    if !s.Players[0].Started && !s.Players[1].Started {
        resend = s.Players[:]
    } else if !s.Players[0].Started {
        resend = s.Players[0:1]
    } else {
        resend = s.Players[1:2]
    }

    for _, p := range resend {
        fmt.Println("Resending start packet to:", p.Conn.RemoteAddr())
        s.SendPacketTo(&packet.GameStart{}, p) 
    }
    return Ok
}, 15)

var ConfirmStop = NewUpdate(func(s *GameServer) UpStatus {
    if s.State.Running { return Ok }
    if s.Players[0] == nil && s.Players[1] == nil { return Ok }
    if s.Players[0] != nil && s.Players[1] != nil { // potentially future issue, shouldnt happen now
        fmt.Println("Both players exist but game isnt running???")
        return Ok
    }
    var p *Profile
    if s.Players[0] != nil {
        p = s.Players[0]
    } else {
        p = s.Players[1]
    }

    if p.Started {
        fmt.Println("Resending stop packet to:", p.Conn.RemoteAddr())
        s.SendPacketTo(&packet.GameStop{}, p) 
    }
    return Ok
}, 15)

var SendBallMove = NewUpdate(func(s *GameServer) UpStatus {
    if !s.State.Running { return Ok }
    ball := s.State.Ball
    s.SendPacket(&packet.BallMove {
        X: float32(ball.Pos.X),
        Y: float32(ball.Pos.Y),
    })

    return Ok
}, 1)

var SendScoreUpdate = NewUpdate(func(s *GameServer) UpStatus {
    if !s.State.Running { return Ok }
    score := &packet.Score {
        P1: s.State.P1.Score,
        P2: s.State.P2.Score,
    }

    s.SendPacket(score)
    return Ok
}, 60)




const (
    TargetTPS = 60
    TargetTickTime = 1.0 / float64(TargetTPS)
    Cooldown time.Duration = time.Millisecond * 5
)

func (s *GameServer) UpdateClients() {
    var elapsed float64 = 0.0
    last := time.Now().UnixNano() / int64(time.Millisecond)

    for {
        current := time.Now().UnixNano() / int64(time.Millisecond)
        s.DeltaTime = float64(current - last) / 1000 // Delta time seconds
        last = current
        elapsed += s.DeltaTime

        s.updateTicks(&elapsed)
        time.Sleep(Cooldown)
    }
}

func (s *GameServer) updateTicks(elapsed *float64) {
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
