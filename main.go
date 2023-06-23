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
	CustomShader *ebiten.Shader
	Time         int
	RainImage    *ebiten.Image
	Pirates0     *ebiten.Image
	Pirates1     *ebiten.Image
	CursorX      int
	CursorY      int
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
	op.Images[0] = g.Pirates0
	op.Images[1] = g.Pirates1
	screen.DrawRectShader(w, h, g.CustomShader, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Shader Test")

	shaderToRender, shaderErr := ebiten.NewShader(shader.Pirates)
	if shaderErr != nil {
		fmt.Println("Shader Err: \n" + shaderErr.Error())
		return
	}
	rainImage, _, imageDecodeErr := image.Decode(bytes.NewReader(asset.Rain))
	if imageDecodeErr != nil {
		fmt.Println("Image load Err: \n" + imageDecodeErr.Error())
		return
	}
	pirates0, _, imageDecodeErr := image.Decode(bytes.NewReader(asset.Pirates0))
	if imageDecodeErr != nil {
		fmt.Println("Image load Err: \n" + imageDecodeErr.Error())
		return
	}
	pirates1, _, imageDecodeErr := image.Decode(bytes.NewReader(asset.Pirates1))
	if imageDecodeErr != nil {
		fmt.Println("Image load Err: \n" + imageDecodeErr.Error())
		return
	}
	game := &Game{
		CustomShader: shaderToRender,
		RainImage:    ebiten.NewImageFromImage(rainImage),
		Pirates0:     ebiten.NewImageFromImage(pirates0),
		Pirates1:     ebiten.NewImageFromImage(pirates1),
	}

	if err := ebiten.RunGame(game); err != nil {
		fmt.Println(err)
	}
}
