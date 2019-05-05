package main

// Super Nintendo Mode7 demo
//
// Ported from Javidx9
//  - Source code: https://github.com/OneLoneCoder/videos/blob/master/OneLoneCoder_Pseudo3DPlanesMode7.cpp
//  - Video tutorial: https://www.youtube.com/watch?v=ybLZyY655iY
//
// Using Ebiten for graphics: https://ebiten.org/

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"golang.org/x/image/colornames"
)

type Point struct {
	x, y float64
}

const (
	screenWidth  = 1024
	screenHeight = 768
	mapSize      = 1024
	fovHalf      = math.Pi / 4.0
)

var (
	world = &Point{.5, .5}
	θ     = math.Pi / 4.0

	near = .005
	far  = .03

	pixels  *image.RGBA
	texture *image.RGBA
)

func init() {
	pixels = image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))
	texture = image.NewRGBA(image.Rect(0, 0, mapSize, mapSize))
	for x := 0; x <= mapSize; x += 32 {
		for y := 0; y < mapSize; y++ {
			texture.Set(x, y, colornames.Magenta)
			texture.Set(x-1, y, colornames.Magenta)
			texture.Set(x+1, y, colornames.Magenta)
			texture.Set(y, x, colornames.Blue)
			texture.Set(y, x-1, colornames.Blue)
			texture.Set(y, x+1, colornames.Blue)
		}
	}
}

func SampleColor(p *Point) color.Color {
	sx := int(p.x * float64(mapSize))
	sy := int(p.y * float64(mapSize-1.0))

	if sx < 0 || sx >= mapSize || sy < 0 || sy >= mapSize {
		return colornames.Black
	}
	return texture.At(sx%mapSize, sy%mapSize)
}

func update(screen *ebiten.Image) error {
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		θ += 0.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		θ -= 0.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		world.x += math.Cos(θ) * 0.002
		world.y += math.Sin(θ) * 0.002
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		world.x -= math.Cos(θ) * 0.002
		world.y -= math.Sin(θ) * 0.002
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// Create Frustum corner points
	far1 := &Point{
		world.x + math.Cos(θ-fovHalf)*far,
		world.y + math.Sin(θ-fovHalf)*far,
	}
	near1 := &Point{
		world.x + math.Cos(θ-fovHalf)*near,
		world.y + math.Sin(θ-fovHalf)*near,
	}
	far2 := &Point{
		world.x + math.Cos(θ+fovHalf)*far,
		world.y + math.Sin(θ+fovHalf)*far,
	}
	near2 := &Point{
		world.x + math.Cos(θ+fovHalf)*near,
		world.y + math.Sin(θ+fovHalf)*near,
	}

	// Starting with furthest away line and work towards the camera Point
	for y := 0; y < screenHeight/2; y++ {
		// Take a sample Point for depth linearly related to rows down screen
		sampleDepth := float64(y) / float64(screenHeight/2.0)

		// Use sample Point in non-linear (1/x) way to enable perspective
		// and grab start and end points for lines across the screen
		start := &Point{
			(far1.x-near1.x)/sampleDepth + near1.x,
			(far1.y-near1.y)/sampleDepth + near1.y,
		}
		end := &Point{
			(far2.x-near2.x)/sampleDepth + near2.x,
			(far2.y-near2.y)/sampleDepth + near2.y,
		}

		// Linearly interpolate lines across the screen
		for x := 0; x < screenWidth; x++ {
			sampleWidth := float64(x) / float64(screenWidth)
			sample := &Point{
				(end.x-start.x)*sampleWidth + start.x,
				(end.y-start.y)*sampleWidth + start.y,
			}

			// Sample pixel from the texture and draw the pixel
			pixels.Set(x, y+screenHeight/2, SampleColor(sample))
		}
	}

	// Draw the pixels
	_ = screen.ReplacePixels(pixels.Pix)

	// Draw the message
	msg := fmt.Sprintf("X: %f\nY: %f\nT: %f\n", world.x, world.y, θ)
	msg += fmt.Sprintf("FPS: %f\n", ebiten.CurrentFPS())
	msg += fmt.Sprintf("Use arrows to move around")
	_ = ebitenutil.DebugPrint(screen, msg)
	return nil
}

func main() {
	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "Mode7"); err != nil {
		log.Fatal(err)
	}
}
