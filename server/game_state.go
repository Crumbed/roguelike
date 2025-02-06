package server


import (
	. "github.com/gen2brain/raylib-go/raylib"
)


const (
    Width   int32   = 600
    Height  int32   = 400
    PW      int32   = 20
    PH      int32   = 80
    P1X     int32   = 5
    P2X     int32   = 600 - PW - 5
    CenterX float32 = 300
    CenterY float32 = 200
)

type Player struct {
    Pos     int32
    Score   uint8
}

type Ball struct {
    Pos Vector2
    Vel Vector2    
}

type GameState struct {
    P1      Player
    P2      Player
    Ball    Ball
    Started bool
}

func NewGame() *GameState {
    return &GameState {
        Ball: Ball { 
            Pos: NewVector2(CenterX, CenterY),
            Vel: NewVector2(-10, -10),
        },
        Started: false,
    }
}


