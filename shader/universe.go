//go:build ignore

//kage:unit pixel

// https://www.shadertoy.com/view/lscczl

package shader

// uniform variables
var IResolution vec2
var IMouse vec2
var ITime float

func N21(p vec2) float {
	a := fract(vec3(p.xyx) * vec3(213.897, 653.453, 253.098))
	a += dot(a, a.yzx+79.76)
	return fract((a.x + a.y) * a.z)
}

func GetPos(id vec2, offs vec2, t float) vec2 {
	n := N21(id + offs)
	n1 := fract(n * 10.0)
	n2 := fract(n * 100.0)
	a := t + n
	return offs + vec2(sin(a*n1), cos(a*n2))*0.4
}

func GetT(ro vec2, rd vec2, p vec2) float {
	return dot(p-ro, rd)
}

func LineDist(a vec3, b vec3, p vec3) float {
	return length(cross(b-a, p-a)) / length(p-a)
}

func df_line(a vec2, b vec2, p vec2) float {
	pa := p - a
	ba := b - a
	h := clamp(dot(pa, ba)/dot(ba, ba), 0.0, 1.0)
	return length(pa - ba*h)
}

func line(a vec2, b vec2, uv vec2) float {
	r1 := 0.04
	r2 := 0.01

	d := df_line(a, b, uv)
	d2 := length(a - b)
	fade := smoothstep(1.5, 0.5, d2)

	fade += smoothstep(0.05, 0.02, abs(d2-0.75))
	return smoothstep(r1, r2, d) * fade
}

func NetLayer(st vec2, n float, t float) float {
	id := floor(st) + n
	st = fract(st) - 0.5

	p := [9]vec2{vec2(0, 0), vec2(0, 0), vec2(0, 0), vec2(0, 0), vec2(0, 0), vec2(0, 0), vec2(0, 0), vec2(0, 0), vec2(0, 0)}
	i := 0
	for y := -1.0; y <= 1.0; y++ {
		for x := -1.0; x <= 1.0; x++ {
			p[i] = GetPos(id, vec2(x, y), t)
			i++
		}
	}

	m := 0.0
	sparkle := 0.0
	for i = 0; i < 9; i++ {
		m += line(p[4], p[i], st)
		d := length(st - p[i])

		s := (0.005 / (d * d))
		s *= smoothstep(1.0, 0.7, d)
		pulse := sin((fract(p[i].x)+fract(p[i].y)+t)*5.0)*0.4 + 0.6
		pulse = pow(pulse, 20.0)

		s *= pulse
		sparkle += s
	}

	m += line(p[1], p[3], st)
	m += line(p[1], p[5], st)
	m += line(p[7], p[5], st)
	m += line(p[7], p[3], st)

	sPhase := (sin(t+n)+sin(t*.1))*0.25 + 0.5
	sPhase += pow(sin(t*0.1)*0.5+0.5, 50.0) * 5.0
	m += sparkle * sPhase
	return m
}

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	uv := (texCoord - IResolution.xy*0.5) / IResolution.y
	M := IMouse.xy/IResolution.xy - 0.5

	t := ITime * 0.1

	s := sin(t)
	c := cos(t)
	rot := mat2(c, -s, s, c)
	st := uv * rot
	M *= rot * 2.0

	m := 0.0
	for i := 0.0; i < 1.0; i += 1.0 / 4 {
		z := fract(t + i)
		size := mix(15.0, 1.0, z)
		fade := smoothstep(0.0, 0.6, z) * smoothstep(1.0, 0.8, z)
		m += fade * NetLayer(st*size-M*z, i, ITime)
	}

	glow := -uv.y * 2.0

	baseCol := vec3(s, cos(t*0.4), -sin(t*0.24))*0.4 + 0.6
	col := baseCol * m
	col += baseCol * glow

	col *= 1.0 - dot(uv, uv)
	t = mod(ITime, 230.0)
	col *= smoothstep(0.0, 20.0, t) * smoothstep(224.0, 200.0, t)

	return vec4(col, 1)
}
