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
	"log"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Point struct {
	x, y float64
}

const (
	screenWidth  = 1024
	screenHeight = 768
)

var (
	world   = &Point{.93, .75}
	θ       = -math.Pi / 2.0
	fovHalf = math.Pi / 4.0
	near    = .005
	far     = .03

	pixels  *image.RGBA
	texture *image.RGBA
)

func init() {
	pixels = image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))

	// Convert the texture to RGBA
	_, tmp, _ := ebitenutil.NewImageFromFile("mk1.png", ebiten.FilterDefault)
	texture = image.NewRGBA(image.Rect(0, 0, tmp.Bounds().Size().X, tmp.Bounds().Size().Y))
	for y := 0; y < tmp.Bounds().Size().Y; y++ {
		for x := 0; x < tmp.Bounds().Size().X; x++ {
			texture.Set(x, y, tmp.At(x, y))
		}
	}
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
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		near += 0.002
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		near -= 0.002
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		far += 0.002
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		far -= 0.002
	}
	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		fovHalf += 0.002
	}
	if ebiten.IsKeyPressed(ebiten.KeyX) {
		fovHalf -= 0.002
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		fovHalf = math.Pi / 4.0
		near = .005
		far = .03
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

			// Sample color from the texture
			sx := int(sample.x * float64(texture.Bounds().Size().X))
			sy := int(sample.y * float64(texture.Bounds().Size().Y-1.0))
			var r, g, b, a uint8
			if sx < 0 || sx >= texture.Bounds().Size().X || sy < 0 || sy >= texture.Bounds().Size().Y {
				r, g, b, a = 0, 0, 0, 0
			} else {
				t := texture.PixOffset(sx, sy)
				r, g, b, a = texture.Pix[t], texture.Pix[t+1], texture.Pix[t+2], texture.Pix[t+3]
			}

			// Draw the pixel
			i := pixels.PixOffset(x, y+screenHeight/2)
			pixels.Pix[i] = r
			pixels.Pix[i+1] = g
			pixels.Pix[i+2] = b
			pixels.Pix[i+3] = a
		}
	}

	// Draw the pixels
	_ = screen.ReplacePixels(pixels.Pix)

	// Draw the message
	msg := fmt.Sprintf("X: %f\nY: %f\nT: %f\n", world.x, world.y, θ)
	msg += fmt.Sprintf("near: %f\nfar: %f\nfov: %f\n", near, far, fovHalf)
	msg += fmt.Sprintf("FPS: %f, TPS: %f\n", ebiten.CurrentFPS(), ebiten.CurrentTPS())
	msg += fmt.Sprintf("Use arrows to move around. q/w for near, a/s for far, z/x for fov.")
	_ = ebitenutil.DebugPrint(screen, msg)
	return nil
}

func main() {
	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "Mode7"); err != nil {
		log.Fatal(err)
	}
}
