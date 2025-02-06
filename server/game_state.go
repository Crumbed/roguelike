package server

import (
	"fmt"
	"math"

	. "github.com/gen2brain/raylib-go/raylib"
)


const (
    Width   int32   = 600
    Height  int32   = 400
    PW      int32   = 10
    PH      int32   = 60
    P1X     int32   = 5
    P2X     int32   = 600 - PW - 5
    P1C     float32 = float32(P1X + PW)
    P2C     float32 = float32(P2X)
    CenterX float32 = 300
    CenterY float32 = 200
    BallS   float32 = 10
)

type Player struct {
    Pos     int32
    Score   uint8
    ColX    float32
}



type Ball struct {
    Pos     Vector2
    Vel     Vector2
}

// first bool is paddle collision, second is x collision (meaning a point was scored)
func (b *Ball) CheckPaddleCol(p *Player) (bool, bool) {
    xCol := false
    var safeX float32
    if p.ColX == P1C && b.Pos.X <= P1C {        // player 1
        xCol = true
        safeX = P1C + 1
    } else if p.ColX == P2C && b.Pos.X >= P2C - BallS { // player 2
        xCol = true
        safeX = P2C - BallS - 1
    }

    pY := float32(p.Pos)
    if xCol && (b.Pos.Y <= pY + float32(PH) && b.Pos.Y >= pY - BallS) {
        b.Pos.X = safeX
        return true, true
    }

    return false, xCol
}

// checks for collision, but also moves ball in case of clipping
func (b *Ball) CheckYCol() bool {
    col := false
    if b.Pos.Y <= 0 {
        col = true 
        b.Pos.Y = 1
    } else if b.Pos.Y >= float32(Height) - BallS {
        col = true
        b.Pos.Y = float32(Height) - BallS - 1
    }

    return col
}

func (b *Ball) CheckXCol() bool { 
    return b.Pos.X <= 0 || b.Pos.X >= float32(Width) - BallS 
}

func (b *Ball) ApplyVelocity(dt float32) {
    b.Pos.X += b.Vel.X * dt
    b.Pos.Y += b.Vel.Y * dt
}

func (b *Ball) CheckCollision(p *Player) {
    if b.CheckYCol() { // top or bottom border
        b.Vel.Y *= -1 // invert y velocity
        //fmt.Println("Top or bottom collision:", b.Vel, b.Pos)
    }

    if b.CheckXCol() { // side border
        b.Vel.X *= -1 // invert x velocity
        //fmt.Println("Side collision:", b.Vel, b.Pos)
    }

    // paddle collision
    pc, xc := b.CheckPaddleCol(p)
    if pc {
        b.Vel.X *= -1
        //fmt.Println("Paddle collision:", b.Vel, b.Pos)
    } else if xc {
        fmt.Println("Point scored")
    }
}

type GameState struct {
    P1      Player
    P2      Player
    Ball    Ball
    Started bool
}

func NewGame() *GameState {
    half := float32(math.Trunc(float64(BallS) / 2))
    return &GameState {
        P1: Player { ColX: P1C },
        P2: Player { ColX: P2C },
        Ball: Ball { 
            Pos: NewVector2(CenterX - half, CenterY - half),
            Vel: NewVector2(-100, -100),
        },
        Started: false,
    }
}

var UpdateBall = NewUpdate(func(s *GameServer) UpStatus {
    ball := &s.State.Ball
    //lastPos := ball.Pos
    var p *Player
    if ball.Vel.X < 0 {         // ball moving left >> p1
        p = &s.State.P1
    } else if ball.Vel.X > 0 {  // ball moving right >> p2
        p = &s.State.P2
    } else { return Ok }        // ball isnt moving left or right (shouldnt be moving at all)

    ball.CheckCollision(p)
    ball.ApplyVelocity(s.DeltaTime)

    return Ok
}, 1)



