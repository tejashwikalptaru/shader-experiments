package asset

import (
	_ "embed"
	"image"
	"image/png"
)

var (
	//go:embed rain.png
	Rain []byte

	//go:embed pirates-channel0.png
	Pirates0 []byte

	//go:embed pirates-channel1.png
	Pirates1 []byte
)

func init() {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
}
