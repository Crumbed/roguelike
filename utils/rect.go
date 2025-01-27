package utils

import . "github.com/gen2brain/raylib-go/raylib"

func Rect(x float32, y float32, w float32, h float32) Rectangle {
    return Rectangle {
        X: x,
        Y: y,
        Width: w,
        Height: h,
    }
}

func Vec2(x float32, y float32) Vector2 {
    return Vector2 { X: x, Y: y }
}
