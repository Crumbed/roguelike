package main

import (
    "github.com/gen2brain/raylib-go/raylib"
    "fmt"
)

func render() {
    rl.BeginDrawing()
    rl.ClearBackground(rl.RayWhite)
    rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)
    rl.EndDrawing()
}

func main() {
    fmt.Println("Hello, world")

    rl.InitWindow(800, 450, "raylib [core] example - basic window")
    rl.SetTargetFPS(60)

    for !rl.WindowShouldClose() {
        render()
    }

    rl.CloseWindow()
}
