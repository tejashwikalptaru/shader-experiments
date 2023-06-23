package shader

import _ "embed"

var (
	//go:embed stars.go
	Stars []byte

	//go:embed color.blindness.go
	ColorBlindness []byte

	//go:embed rain.go
	Rain []byte

	//go:embed flame.go
	Flame []byte

	//go:embed universe.go
	Universe []byte
)
