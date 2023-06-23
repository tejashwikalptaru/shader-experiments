package main

import (
	"bytes"
	"fmt"
	"game/asset"
	"game/shader"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type Game struct {
	StarShader *ebiten.Shader
	Time       int
	RainImage  *ebiten.Image
	CursorX    int
	CursorY    int
}

func (g *Game) Update() error {
	g.Time++
	g.CursorX, g.CursorY = ebiten.CursorPosition()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	op := &ebiten.DrawRectShaderOptions{
		Uniforms: map[string]any{
			"ITime": float32(g.Time) / 60,
			"IResolution": []float32{
				float32(w),
				float32(h),
			},
			"ICursor": []float32{
				float32(g.CursorX),
				float32(g.CursorY),
			},
		},
	}
	//op.Images[0] = g.RainImage
	screen.DrawRectShader(w, h, g.StarShader, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Shader Test")

	starShader, shaderErr := ebiten.NewShader(shader.Universe)
	if shaderErr != nil {
		fmt.Println("Shader Err: \n" + shaderErr.Error())
		return
	}
	rainImage, _, imageDecodeErr := image.Decode(bytes.NewReader(asset.Rain))
	if imageDecodeErr != nil {
		fmt.Println("Image load Err: \n" + imageDecodeErr.Error())
		return
	}
	game := &Game{
		StarShader: starShader,
		RainImage:  ebiten.NewImageFromImage(rainImage),
	}

	if err := ebiten.RunGame(game); err != nil {
		fmt.Println(err)
	}
}
