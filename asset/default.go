package asset

import (
	_ "embed"
	"image"
	"image/png"
)

var (
	//go:embed rain.png
	Rain []byte
)

func init() {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
}
