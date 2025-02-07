package server

import (
	//"fmt"
	"main/packet"
	"math"

	. "github.com/gen2brain/raylib-go/raylib"
)


const (
    Width   int32   = 600
    Height  int32   = 400
    PW      int32   = 10
    PH      int32   = 80
    P1X     int32   = 5
    P2X     int32   = 600 - PW - 5
    P1C     float32 = float32(P1X + PW)
    P2C     float32 = float32(P2X)
    CenterX float32 = 300
    CenterY float32 = 200
    BallS   float32 = 10
)

// pN == 0 or 1 | player 1 or 2
func NewPlayer(pN uint8) Player {
    hb := NewRectangle(0, 0, float32(PW), float32(PH))
    if pN == 0 {
        hb.X = float32(P1X)
    } else {
        hb.X = float32(P2X)
    }

    return Player { HitBox: hb }
}
type Player struct {
    Pos     int32
    Score   uint8
    HitBox  Rectangle
}
func (p *Player) Move(y int32) {
    p.Pos = y
    p.HitBox.Y = float32(y)
}

func (p *Player) CalculateHitZone(b *Ball) {
    relY := b.Pos.Y - float32(p.Pos)
    zone := relY / 10 // paddle has 6 zones, 10 pixels tall
    vel  := &b.Vel
    absX := math.Abs(float64(vel.X)) - 150

    var newY float64
    if zone >= 5 {          // southern most    +
        newY = 2 * 150 + absX
    } else if zone >= 4 {   // southern         +
        newY = 1 * 150 + absX
    } else if zone >= 3 {   // southern middle  +
        newY = 0.25 * 150 + absX
    } else if zone >= 2 {   // northern middle  -
        newY = -0.25 * 150 - absX
    } else if zone >= 1 {   // northern         -
        newY = -1 * 150 - absX
    } else {                // northern most    -
        newY = -2 * 150 - absX
    }

    // already moving in that direction
    /*
    if (newY < 0 && vel.Y < 0) || (newY > 0 && vel.Y > 0) { 
        newY += float64(vel.Y) * 0.25 // add 25% of current y velocity to new velocity
    }
    */

    vel.Y = float32(newY)
}


func NewBall() Ball {
    ball := Ball{}
    ball.Init(-1)
    return ball
}
type Ball struct {
    Pos     Vector2
    Vel     Vector2
    HitBox  Rectangle
}

func (b *Ball) Move(x, y float32) {
    b.Pos.X = x
    b.Pos.Y = y
    b.HitBox.X = x
    b.HitBox.Y = y
}
func (b *Ball) MoveY(y float32) {
    b.Pos.Y = y
    b.HitBox.Y = y
}
func (b *Ball) MoveX(x float32) {
    b.Pos.X = x
    b.HitBox.X = x
}


// dir should be -1 or 1 for left or right
func (b *Ball) Init(dir float32) {
    half := float32(math.Trunc(float64(BallS) / 2))
    b.Pos.X = CenterX - half
    b.Pos.Y = CenterY - half
    b.Vel.X = 200 * dir
    b.Vel.Y = 0
    b.HitBox.X = b.Pos.X
    b.HitBox.Y = b.Pos.Y
    b.HitBox.Width = BallS
    b.HitBox.Height = BallS
}

func (b *Ball) IncreaseVel() {
    if b.Vel.X >= 700 || b.Vel.X <= -700 { return }
    if b.Vel.X < 0 {
        b.Vel.X -= 50
        return
    }

    b.Vel.X += 50
}

func (b *Ball) CheckScore() bool {
    return CheckCollisionRecs(b.HitBox, LeftScoreBox) || CheckCollisionRecs(b.HitBox, RightScoreBox) 
}

// first bool is paddle collision
func (b *Ball) CheckPaddleCol(p *Player) bool {
    if !CheckCollisionRecs(b.HitBox, p.HitBox) { return false }

    var safeX float32
    if b.Vel.X < 0 { // left collision
        safeX = p.HitBox.Width + 6
    } else {
        safeX = p.HitBox.X - 1
    }

    b.MoveX(safeX)
    return true
}

// checks for collision, but also moves ball in case of clipping
func (b *Ball) CheckYCol() {
    if b.Pos.Y <= 0 {
        b.MoveY(1)
        b.Vel.Y *= -1
    } else if b.Pos.Y >= float32(Height) - BallS {
        b.MoveY(float32(Height) - BallS - 1)
        b.Vel.Y *= -1
    }
}

func (b *Ball) ApplyVelocity(dt float32) {
    b.Move(b.Pos.X + b.Vel.X * dt, b.Pos.Y + b.Vel.Y * dt)
}

var LeftScoreBox = NewRectangle(0, 0, 10, float32(Height))
var RightScoreBox = NewRectangle(float32(Width) - 10, 0, 10, float32(Height))
type GameState struct {
    P1      Player
    P2      Player
    Ball    Ball
    Started bool
}

func NewGame() *GameState {
    return &GameState {
        P1: NewPlayer(0),
        P2: NewPlayer(1),
        Ball: NewBall(),
        Started: false,
    }
}

func (state *GameState) ScoreAgainst(p *Player) {
    if *p == state.P1 {
        state.P2.Score += 1
        state.Ball.Init(1)
    } else {
        state.P1.Score += 1
        state.Ball.Init(-1)
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

    ball.CheckYCol()
    // paddle collision
    if ball.CheckPaddleCol(p) {
        ball.IncreaseVel()
        ball.Vel.X *= -1
        p.CalculateHitZone(ball)
    } else if ball.CheckScore() {
        s.State.ScoreAgainst(p)
        score := &packet.Score {
            P1: s.State.P1.Score,
            P2: s.State.P2.Score,
        }

        s.SendPacket(score)
    }

    ball.ApplyVelocity(s.DeltaTime)
    return Ok
}, 1)



