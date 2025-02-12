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

const (
    IpSize  float32 = 8
    IpLen   float32 = IpSize * 21
    EntryW  float32 = IpLen + 10 
    EntryH  float32 = IpSize + 20
    EntryX  float32 = float32(CenterX) - EntryW / 2
    EntryY  float32 = float32(CenterY) - EntryH

    PlayFS  float32 = 16
    PlayW   float32 = EntryW
    PlayH   float32 = PlayFS + 5
    PlayX   float32 = float32(CenterX) - PlayW / 2
    PlayY   float32 = EntryY + EntryH + 5
)

var TextBox = rl.NewRectangle(EntryX, EntryY, EntryW, EntryH)
var PlayBtn = rl.NewRectangle(PlayX, PlayY, PlayW, PlayH)
var ConnectError = ""

func InitScreen() Screen {
    rl.SetConfigFlags(rl.FlagWindowResizable)
    rl.InitWindow(Width, Height, "Pong Online")
    rl.SetWindowMinSize(int(Width), int(Height))
    rl.SetTargetFPS(60)

    font := rl.LoadFont("assets/font.fnt")
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
    rl.DrawTextEx(c.screen.font, "pong online", rl.NewVector2(float32(CenterX) - 264, 50), 48, 0, rl.White)
    rl.DrawRectangleRec(TextBox, rl.DarkGray)
    if len(c.serverIp) == 0 {
        rl.DrawTextEx(c.screen.font, "ip : port", rl.NewVector2(EntryX + 5, EntryY + 10), IpSize, 0, rl.Gray)
    } else {
        rl.DrawTextEx(c.screen.font, string(c.serverIp), rl.NewVector2(EntryX + 5, EntryY + 10), IpSize, 0, rl.White)
    }

    if ConnectError != "" {
        rl.DrawTextEx(c.screen.font, ConnectError, rl.NewVector2(EntryX, EntryY - 10), 8, 0, rl.Red)
    }

    rl.DrawRectangleRec(PlayBtn, rl.DarkGreen)
    rl.DrawTextEx(c.screen.font, "play", rl.NewVector2(float32(CenterX) - 32, PlayY + 2), PlayFS, 0, rl.White)

    /*
    rl.DrawTextEx(c.screen.font, "111.111.111.111:65535", rl.NewVector2(0, 0), 8, 0, rl.White)
    le := rl.MeasureTextEx(c.screen.font, "111.111.111.111:65535", 8, 0)
    fmt.Println(le)
    */
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

func ipInput(c *Client) {
    key := rl.GetCharPressed()

    for key > 0 {
        if len(c.serverIp) >= 21 { break }
        if key != 46 && (key < 48 || key > 58) { break }
        c.serverIp = append(c.serverIp, byte(key))
        key = rl.GetCharPressed()
    }

    if rl.IsKeyPressed(rl.KeyBackspace) && len(c.serverIp) > 0 {
        c.serverIp = c.serverIp[:len(c.serverIp)-1]
    }

    clickedPlay := rl.CheckCollisionPointRec(rl.GetMousePosition(), PlayBtn) && rl.IsMouseButtonPressed(rl.MouseButtonLeft)
    if rl.IsKeyPressed(rl.KeyEnter) || clickedPlay {
        err := c.Connect()
        if err != nil {
            ConnectError = "connection failed."
            return
        }
    }
}
