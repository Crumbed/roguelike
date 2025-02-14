package server

import (
	//"fmt"
	"main/packet"
	"math"
	//. "github.com/gen2brain/raylib-go/raylib"
)


const (
    Width   int32   = 600
    Height  int32   = 400
    PW      int32   = 10
    PH      int32   = 80
    P1X     int32   = 5
    P2X     int32   = 600 - PW - 5
    P1C     float64 = float64(P1X + PW)
    P2C     float64 = float64(P2X)
    CenterX float64 = 300
    CenterY float64 = 200
    BallS   float64 = 10
    InitVel float64 = 300
)



// pN == 0 or 1 | player 1 or 2
func NewPlayer(pN uint8) *Player {
    p := &Player {
        hb: HitBox { width: PW, height: PH },
    }
    switch pN {
    case 0: p.pos.X = float64(P1X)
    case 1: p.pos.X = float64(P2X)
    }
    p.hb.pos = &p.pos

    return p
}
type Player struct {
    pos     Position
    Score   uint8
    hb      HitBox
}

func (p *Player) CalculateHitZone(b *Ball) {
    relY := float64(b.Pos.Y) - float64(p.pos.Y)
    zone := relY / 10 // paddle has 8 zones, 10 pixels tall
    vel  := &b.Vel

    var angle float64
    switch {
    case zone >= 7: angle = 45      // southern most    +
    case zone >= 6: angle = 35      // south            +
    case zone >= 5: angle = 25      // southern mid     +
    case zone >= 4: angle = 15      // center           +
    case zone >= 3: angle = 345     // center           -
    case zone >= 2: angle = 335     // northern mid     -
    case zone >= 1: angle = 325     // north            -
    default:        angle = 315     // northern most    -
    }

    // if ball is going right, we need to invert this angle
    if vel.GetDir() == Right { angle = 180 - angle }
    //fmt.Println("hit zone:", zone)
    //fmt.Println("angle:", angle)
    vel.SetRotation(Radians(angle))
}


func NewBall() *Ball {
    ball := &Ball{ Vel: Velocity{InitVel, 0, InitVel} }
    ball.Box.width = int32(BallS)
    ball.Box.height = int32(BallS)
    ball.Box.pos = &ball.Pos
    ball.Init(-1)
    return ball
}
type Ball struct {
    Pos     Position
    Vel     Velocity
    Box     HitBox
}

// dir should be -1 or 1 for left or right
func (b *Ball) Init(dir float64) {
    half := math.Trunc(BallS / 2)
    b.Pos.X = CenterX - half
    b.Pos.Y = CenterY - half
    b.Vel.Set(InitVel * dir, 0)
    //b.Vel.Set(InitVel, 0)
    //b.Vel.SetRotation(Radians(angle))
}

func (b *Ball) IncreaseVel() {
    vlen := b.Vel.Len()
    //if vlen >= 700 { return }
    b.Vel.SetUnitLength(vlen + 100)
}

// first bool is paddle collision
func (b *Ball) CheckPaddleCol(p *Player) bool {
    if !p.hb.CollidesWith(&b.Box) { return false }

    var safeX int32
    if b.Vel.GetDir() == Left {
        safeX = P1X + p.hb.width + 1
    } else {
        safeX = int32(p.hb.pos.X - BallS) - 1
    }

    b.Pos.X = float64(safeX)
    return true
}

// checks for collision, but also moves ball in case of clipping
func (b *Ball) CheckYCol() {
    if b.Box.CollidesWith(TopBorder) {
        b.Pos.Y = 1
        b.Vel.InvertY()
    } else if b.Box.CollidesWith(BottomBorder) {
        b.Pos.Y = float64(Height) - BallS - 1
        b.Vel.InvertY()
    }
}

func (b *Ball) CheckScore() bool {
    return b.Box.CollidesWith(LeftScoreBox) || b.Box.CollidesWith(RightScoreBox)
}

var TopBorder       = &HitBox{ &Position{0,0}, Width, 0 }
var BottomBorder    = &HitBox{ &Position{0,float64(Height)}, Width, 0 }
var LeftScoreBox    = &HitBox{ &Position{0,0}, 10, Height }
var RightScoreBox   = &HitBox{ &Position{float64(Width) - 10,0}, 10, Height }
type GameState struct {
    P1      *Player
    P2      *Player
    Ball    *Ball
    Running bool
}

func NewGame() *GameState {
    return &GameState {
        P1: NewPlayer(0),
        P2: NewPlayer(1),
        Ball: NewBall(),
        Running: false,
    }
}

func (state *GameState) ScoreAgainst(p *Player) {
    if p == state.P1 {
        state.P2.Score += 1
        state.Ball.Init(1)
    } else {
        state.P1.Score += 1
        state.Ball.Init(-1)
    }
}

var UpdateBall = NewUpdate(func(s *GameServer) UpStatus {
    if !s.State.Running { return Ok }
    ball := s.State.Ball
    ball.Pos.ApplyVelocity(&ball.Vel, s.DeltaTime)

    var p *Player
    switch ball.Vel.GetDir() {
    case Left: p = s.State.P1
    case Right: p = s.State.P2
    default: return Ok
    }

    ball.CheckYCol()
    // paddle collision
    if ball.CheckPaddleCol(p) {
        p.CalculateHitZone(ball)
        ball.IncreaseVel()
        //ball.Vel.X *= -1
        //fmt.Println("Paddle collision:", ball.Vel)
    } else if ball.CheckScore() {
        s.State.ScoreAgainst(p)
        score := &packet.Score {
            P1: s.State.P1.Score,
            P2: s.State.P2.Score,
        }

        s.SendPacket(score)
    }

    return Ok
}, 1)



