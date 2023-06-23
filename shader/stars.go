//go:build ignore

//kage:unit pixel

// credit to https://www.shadertoy.com/view/ttKGDt

package main

// uniform variables
var IResolution vec2
var ITime float

func rot(a float) mat2 {
	c := cos(a)
	s := sin(a)
	return mat2(c, s, -s, c)
}

func pmod(p vec2, r float) vec2 {
	pi := acos(-1.0)
	pi2 := pi * 2.0

	a := atan2(p.x, p.y) + pi/r
	n := pi2 / r
	a = floor(a/n) * n
	return p * rot(-a)
}

func box(p vec3, b vec3) float {
	d := abs(p) - b
	return min(max(d.x, max(d.y, d.z)), 0.0) + length(max(d, 0.0))
}

func ifsBox(p vec3) float {
	for i := 0; i < 5; i++ {
		p = abs(p) - 1.0
		p.xy *= rot(ITime * 0.3)
		p.xz *= rot(ITime * 0.1)
	}
	p.xz *= rot(ITime)
	return box(p, vec3(0.4, 0.8, 0.3))
}

func makeMap(p vec3, cPos vec3) float {
	p1 := p
	p1.x = mod(p1.x-5, 10) - 5
	p1.y = mod(p1.y-5, 10) - 5
	p1.z = mod(p1.z, 16) - 8
	p1.xy = pmod(p1.xy, 5.0)
	return ifsBox(p1)
}

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	p := (texCoord.xy*2.0 - IResolution.xy) / min(IResolution.x, IResolution.y)
	cPos := vec3(0.0, 0.0, -3.0*ITime)
	cDir := normalize(vec3(0.0, 0.0, -1.0))
	cUp := vec3(sin(ITime), 1.0, 0.0)
	cSide := cross(cDir, cUp)

	ray := normalize(cSide*p.x + cUp*p.y + cDir)

	acc := 0.0
	acc2 := 0.0
	t := 0.0

	for i := 0; i < 99; i++ {
		pos := cPos + ray*t
		dist := makeMap(pos, cPos)
		dist = max(abs(dist), 0.02)
		a := exp(-dist * 3.0)
		if mod(length(pos)+24.0*ITime, 30.0) < 3.0 {
			a *= 2.0
			acc2 += a
		}
		acc += a
		t += dist * 0.5
	}
	col := vec3(acc*0.01, acc*0.011+acc2*0.002, acc*0.012+acc2*0.005)
	return vec4(col, 1.0-t*0.03)
}
