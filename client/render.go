package client

import (
	"fmt"
	"math"

	"github.com/gen2brain/raylib-go/raylib"
)


type Display uint8
const (
    StartMenu Display = iota
    Game 
)


func InitScreen() Screen {
    rl.SetConfigFlags(rl.FlagWindowResizable)
    rl.InitWindow(Width, Height, "Pong Online")
    rl.SetWindowMinSize(int(Width), int(Height))
    rl.SetTargetFPS(60)

    font := rl.LoadFont("assets/font.fnt")
    //v := rl.MeasureTextEx(font, "Pong Online", 50, 0)
    //fmt.Println(v)

    return Screen { 
        w: Width,
        h: Height,
        scale: 1,
        disp: StartMenu,
        buffer: rl.LoadRenderTexture(Width, Height),
        cPos: rl.GetMousePosition(),
        font: font,
    }
}
type Screen struct {
    w, h    int32
    scale   float32
    disp    Display
    buffer  rl.RenderTexture2D
    cPos    rl.Vector2
    font    rl.Font
}

func (s *Screen) UpdateSize() {
    w, h := int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight())
    if s.w == w && s.h == h { return }
    s.w = w
    s.h = h
    s.scale = float32(math.Min(float64(w) / float64(Width), float64(h) / float64(Height)))
    /*
    rl.UnloadFont(s.font)
    s.font = rl.LoadFontEx("assets/joystix_mono.otf", int32(float32(FSize) * s.scale), nil)
    rl.SetTextureFilter(s.font.Texture, rl.FilterBilinear)
    */
}

func (s *Screen) StartRendering() {
    rl.BeginTextureMode(s.buffer)
}

func (s *Screen) FinalizeRender() {
    rl.EndTextureMode()
    rl.BeginDrawing()
    rl.ClearBackground(rl.Black)

    sw, sh := float32(s.w), float32(s.h)
    w, h := float32(Width), float32(Height)
    scaleRect := rl.NewRectangle(0, 0, w * s.scale, h * s.scale)
    scaledW := w * s.scale
    scaledH := h * s.scale
    if scaledW == sw { // center y
        scaleRect.Y = (sh - scaledH) * 0.5
    } else if scaledH == sh { // center x
        scaleRect.X = (sw - scaledW) * 0.5
    } else { // center both
        scaleRect.X = (sw - scaledW) * 0.5
        scaleRect.Y = (sh - scaledH) * 0.5
    }

    //rl.BeginBlendMode(rl.BlendAddColors)
    rl.DrawTexturePro(s.buffer.Texture, 
        rl.NewRectangle(0, 0, float32(s.buffer.Texture.Width), float32(-s.buffer.Texture.Height)),
        scaleRect,
        rl.NewVector2(0, 0), 0, rl.White)
    //rl.EndBlendMode()

    rl.EndDrawing()
}



func (c *Client) drawMenu() {
    //rl.DrawTextEx(c.screen.font, "Pong Online", rl.NewVector2(0, 5), 50, 0, rl.White)
    //rl.DrawTexture(texture, 0, 0, rl.White)
    rl.DrawTextEx(c.screen.font, "pong online", rl.NewVector2(5, 5), 32, 0, rl.White)
}


func (c *Client) drawGame() {
    p1, p2 := &c.Players[0], &c.Players[1]
    rl.ClearBackground(rl.Black)
    rl.DrawRectangle(
        CenterX - 1, 0,
        2, Height,
        rl.Gray)

    if !c.Started { return }
    c.drawScore(p1, p2)

    // player 1 paddle
    p1.render(0)
    // player 2 paddle
    p2.render(1)

    // ball
    c.Ball.render()
}

func (c *Client) render() {
    m := rl.GetMousePosition()
    if !rl.Vector2Equals(c.screen.cPos, m) && rl.IsCursorHidden() {
        rl.ShowCursor()
    }
    c.screen.cPos = m

    c.screen.UpdateSize()
    c.screen.StartRendering()
    defer c.screen.FinalizeRender()
    //rl.BeginDrawing()
    //defer rl.EndDrawing()

    switch c.screen.disp {
    case StartMenu: c.drawMenu()
    case Game: c.drawGame()
    }
}


func (c *Client) drawScore(p1, p2 *Player) {
    p1str := fmt.Sprintf("%d", p1.Score)
    p2str := fmt.Sprintf("%d", p2.Score)
    p1pos := rl.NewVector2(float32(CenterX) - 20 - float32(len(p1str)) * 69, 0)
    p2pos := rl.NewVector2(float32(CenterX) + 20, 0)
    
    rl.DrawTextEx(c.screen.font, p1str, p1pos, 100, 0, rl.Gray)
    rl.DrawTextEx(c.screen.font, p2str, p2pos, 100, 0, rl.Gray)
}
