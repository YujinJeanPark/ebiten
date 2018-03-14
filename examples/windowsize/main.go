// Copyright 2015 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build example jsgo

package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/examples/resources/images"
	"github.com/hajimehoshi/ebiten/inpututil"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	initScreenWidth  = 320
	initScreenHeight = 240
	initScreenScale  = 2
)

var (
	gophersImage *ebiten.Image
	count        = 0
)

func createRandomIconImage() image.Image {
	const size = 32

	r := uint8(rand.Intn(0x100))
	g := uint8(rand.Intn(0x100))
	b := uint8(rand.Intn(0x100))
	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	for j := 0; j < size; j++ {
		for i := 0; i < size; i++ {
			img.Pix[j*img.Stride+4*i] = r
			img.Pix[j*img.Stride+4*i+1] = g
			img.Pix[j*img.Stride+4*i+2] = b
			img.Pix[j*img.Stride+4*i+3] = uint8(float64(i+j) / float64(2*size) * 0xff)
		}
	}

	return img
}

func update(screen *ebiten.Image) error {
	screenScale := ebiten.ScreenScale()
	d := int(32 / screenScale)
	screenWidth, screenHeight := screen.Size()
	fullscreen := ebiten.IsFullscreen()
	runnableInBackground := ebiten.IsRunnableInBackground()
	cursorVisible := ebiten.IsCursorVisible()

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		screenHeight += d
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		if 16 < screenHeight && d < screenHeight {
			screenHeight -= d
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		if 16 < screenWidth && d < screenWidth {
			screenWidth -= d
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		screenWidth += d
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		switch screenScale {
		case 1:
			screenScale = 1.5
		case 1.5:
			screenScale = 2
		case 2:
			screenScale = 1
		default:
			panic("not reached")
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		fullscreen = !fullscreen
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		runnableInBackground = !runnableInBackground
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		cursorVisible = !cursorVisible
	}
	ebiten.SetScreenSize(screenWidth, screenHeight)
	ebiten.SetScreenScale(screenScale)
	ebiten.SetFullscreen(fullscreen)
	ebiten.SetRunnableInBackground(runnableInBackground)
	ebiten.SetCursorVisible(cursorVisible)

	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		ebiten.SetWindowIcon([]image.Image{createRandomIconImage()})
	}

	count++

	if ebiten.IsRunningSlowly() {
		return nil
	}

	screen.Fill(color.RGBA{0x80, 0x80, 0xc0, 0xff})
	w, h := gophersImage.Size()
	w2, h2 := screen.Size()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(-w+w2)/2, float64(-h+h2)/2)
	dx := math.Cos(2*math.Pi*float64(count)/360) * 10
	dy := math.Sin(2*math.Pi*float64(count)/360) * 10
	op.GeoM.Translate(dx, dy)
	screen.DrawImage(gophersImage, op)

	x, y := ebiten.CursorPosition()
	msg := fmt.Sprintf(`Press arrow keys to change the window size
Press S key to change the window scale
Press F key to switch the fullscreen state
Press B key to switch the run-in-background state
Press C key to switch the cursor visibility
Press I key to change the window icon
Cursor: (%d, %d)
FPS: %0.2f`, x, y, ebiten.CurrentFPS())
	ebitenutil.DebugPrint(screen, msg)
	return nil
}

func main() {
	fmt.Printf("Device scale factor: %0.2f\n", ebiten.DeviceScaleFactor())

	img, _, err := image.Decode(bytes.NewReader(images.Gophers_jpg))
	if err != nil {
		log.Fatal(err)
	}
	gophersImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	ebiten.SetWindowIcon([]image.Image{createRandomIconImage()})

	if err := ebiten.Run(update, initScreenWidth, initScreenHeight, initScreenScale, "Window Size (Ebiten Demo)"); err != nil {
		log.Fatal(err)
	}
}
