package main

import (
	. "main/utils"
	. "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
    Texture     Texture2D
    DestRect    Rectangle
    Vel         Vector2
}


func (p *Player) move() {
    p.Vel.X = 0
    p.Vel.Y = 0

    if IsKeyDown(KeyD) {
        p.Vel.X = 100
    }
    if IsKeyDown(KeyA) {
        p.Vel.X = -100
    }

    if IsKeyDown(KeyW) {
        p.Vel.Y = -100
    }
    if IsKeyDown(KeyS) {
        p.Vel.Y = 100
    }
}

func (p *Player) applyVelocity() {
    p.DestRect.X += p.Vel.X * GetFrameTime()
    p.DestRect.Y += p.Vel.Y * GetFrameTime()
}

func (self *Player) draw() {
    DrawTexturePro(
        self.Texture,
        Rect(0, 0, 16, 16),
        self.DestRect,
        Vec2(0, 0),
        0,
        RayWhite,
    )
}

func main() {
    InitWindow(600, 400, "Game window")
    SetTargetFPS(60)

    player := Player {
        Texture: LoadTexture("assets/player.png"),
        DestRect: Rect(10, 100, 100, 100),
    }

    for !WindowShouldClose() {
        player.move()
        player.applyVelocity()

        BeginDrawing()
        ClearBackground(SkyBlue)

        //rl.DrawTexture(playerTexture, 10, 100, rl.RayWhite)
        player.draw()

        EndDrawing()
    }

    UnloadTexture(player.Texture)
    CloseWindow()
}
