package server

import (
	"math"
)


type Direction uint8 
const (
    Left Direction = iota
    Right
    None
)

func Radians(angle float64) float64 { return angle * (math.Pi / 180) }

type Position struct { X, Y float64 }
func (p *Position) ApplyVelocity(vel *Velocity, deltaTime float64) {
    p.X += vel.x * deltaTime
    p.Y += vel.y * deltaTime
}

type HitBox struct {
    pos     *Position   // Top left of hitbox
    width   int32
    height  int32
}

func (hb *HitBox) CollidesWith(other *HitBox) bool {
    cx := hb.pos.X <= other.pos.X + float64(other.width) && other.pos.X <= hb.pos.X + float64(hb.width) 
    cy := hb.pos.Y <= other.pos.Y + float64(other.height) && other.pos.Y <= hb.pos.Y + float64(hb.height)

    return cx && cy
}



// Just a fancy Vector3 where the z is the length of the vector
type Velocity struct { 
    x, y    float64 //  x & y of vector
    units   float64 //  length of vector
}
func NewVelocity(x, y float64) Velocity {
    v := Velocity { x: x, y: y }
    v.CalculateUnits()
    return v
}

func (v *Velocity) Set(x, y float64) {
    v.x = x
    v.y = y
    v.CalculateUnits()
}

// Rotate the vector direction relative to the current direction
func (v *Velocity) Rotate(angle float64) {
    cosres := math.Cos(angle)
    sinres := math.Sin(angle)

    x, y := v.x, v.y
    v.x = x*cosres - y*sinres
    v.y = x*sinres + y*cosres
}

// Sets the angle of rotation
func (v *Velocity) SetRotation(angle float64) {
	x := math.Cos(angle)
	y := math.Sin(angle)

	v.x = x * v.units
	v.y = y * v.units
}

// Inverts the X direction
func (v *Velocity) InvertX() { v.x *= -1 }
// Inverts the Y direction
func (v *Velocity) InvertY() { v.y *= -1 }

// Scales the vector to a given value. 
// !!!DOES NOT RECALCULATE UNIT LENGTH!!!
func (v *Velocity) Scale(scale float64) {
    v.x *= scale
    v.y *= scale
}

// Normalizes the vector
func (v *Velocity) Normalize() { v.Scale(1/v.units) }

// Add X & Y
func (v *Velocity) Add(x, y float64) {
    v.x += x
    v.y += y
    v.CalculateUnits()
}
// Add another vector
func (v *Velocity) Addv(other *Velocity) {
    v.x += other.x
    v.y += other.y
    v.CalculateUnits()
}
// Add a flat value to X & Y
func (v *Velocity) Addf(flat float64) {
    v.x += flat
    v.y += flat
    v.CalculateUnits()
}
func (v *Velocity) Len() float64 { return v.units }
func (v *Velocity) CalculateUnits() {
    v.units = math.Sqrt((v.x * v.x) + (v.y * v.y))
}

// Change the unit length of the vector
func (v *Velocity) SetUnitLength(units float64) {
    v.Scale(1/v.units)
    v.Scale(units)
    v.units = units
}

func (v *Velocity) GetDir() Direction {
    if v.x < 0 {
        return Left 
    } else if v.x > 0 {
        return Right
    }

    return None
}
