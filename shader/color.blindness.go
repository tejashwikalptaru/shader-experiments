//go:build ignore

//kage:unit pixel

// credits https://www.shadertoy.com/view/ddBcRW

package shader

// uniform variables
var IResolution vec2
var ITime float

func palette(t float) vec3 {
	a := vec3(0.5, 0.5, 0.5)
	b := vec3(0.5, 0.5, 0.5)
	c := vec3(1.0, 1.0, 1.0)
	d := vec3(0.263, 0.416, 0.557)

	return a + b*cos(6.28318*(c*t+d))
}

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	uv := (texCoord*2.0 - IResolution.xy) / IResolution.y
	uv0 := uv
	finalColor := vec3(0.0)

	for i := 0.0; i < 4.0; i++ {
		uv = fract(uv*1.5) - 0.5
		d := length(uv) * exp(-length(uv0))
		col := palette(length(uv0) + i*0.4 + ITime*0.4)
		d = sin(d*8.0+ITime) / 8.0
		d = abs(d)
		d = pow(0.01/d, 2.0)
		finalColor += col * d
	}
	return vec4(finalColor, 1.0)
}
