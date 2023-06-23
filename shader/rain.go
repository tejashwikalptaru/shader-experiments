//go:build ignore

// credit to https://www.shadertoy.com/view/DdSczW

package shader

// uniform variables
var IResolution vec2
var ITime float
var ICursor vec2

func N13(p float) vec3 {
	p3 := fract(vec3(p) * vec3(0.1031, 0.11369, 0.13787))
	p3 += dot(p3, p3.yzx+19.19)
	return fract(vec3((p3.x+p3.y)*p3.z, (p3.x+p3.z)*p3.y, (p3.y+p3.z)*p3.x))
}

func N14(t float) vec4 {
	return fract(sin(t*vec4(123.0, 1024.0, 1456.0, 264.0)) * vec4(6547.0, 345.0, 8799.0, 1564.0))
}

func N(t float) float {
	return fract(sin(t*12345.564) * 7658.76)
}

func Saw(b float, t float) float {
	return smoothstep(0.0, b, t) * smoothstep(1.0, b, t)
}

func DropLayer2(uv vec2, t float) vec2 {
	UV := uv

	uv.y += t * 0.75
	a := vec2(6.0, 1.0)
	grid := a * 2.0
	id := floor(uv * grid)

	colShift := N(id.x)
	uv.y += colShift

	id = floor(uv * grid)
	n := N13(id.x*35.2 + id.y*2376.1)
	st := fract(uv*grid) - vec2(0.5, 0)

	x := n.x - 0.5

	y := UV.y * 20.0
	wiggle := sin(y + sin(y))
	x += wiggle * (0.5 - abs(x)) * (n.z - 0.5)
	x *= 0.7
	ti := fract(t + n.z)
	y = (Saw(0.85, ti)-0.5)*0.9 + 0.5
	p := vec2(x, y)

	d := length((st - p) * a.yx)

	mainDrop := smoothstep(0.4, 0.0, d)

	r := sqrt(smoothstep(1.0, y, st.y))
	cd := abs(st.x - x)
	trail := smoothstep(0.23*r, 0.15*r*r, cd)
	trailFront := smoothstep(-0.02, 0.02, st.y-y)
	trail *= trailFront * r * r

	y = UV.y
	trail2 := smoothstep(0.2*r, 0.0, cd)
	droplets := max(0.0, (sin(y*(1.0-y)*120.0)-st.y)) * trail2 * trailFront * n.z
	y = fract(y*10.0) + (st.y - 0.5)
	dd := length(st - vec2(x, y))
	droplets = smoothstep(0.3, 0.0, dd)
	m := mainDrop + droplets*r*trailFront

	return vec2(m, trail)
}

func StaticDrops(uv vec2, t float) float {
	uv *= 40.0

	id := floor(uv)
	uv = fract(uv) - 0.5
	n := N13(id.x*107.45 + id.y*3543.654)
	p := (n.xy - 0.5) * 0.7
	d := length(uv - p)

	fade := Saw(0.025, fract(t+n.z))
	return smoothstep(0.3, 0.0, d) * fract(n.z*10.0) * fade
}

func Drops(uv vec2, t float, l0 float, l1 float, l2 float) vec2 {
	s := StaticDrops(uv, t) * l0
	m1 := DropLayer2(uv, t) * l1
	m2 := DropLayer2(uv*1.85, t) * l2

	c := s + m1.x + m2.x
	c = smoothstep(0.3, 1.0, c)

	return vec2(c, max(m1.y*l0, m2.y*l1))
}

func imageColorAtPixel(pixelCoords vec2) vec4 {
	sizeInPixels := imageSrcTextureSize()
	offsetInTexels, _ := imageSrcRegionOnTexture()
	adjustedTexelCoords := pixelCoords/sizeInPixels + offsetInTexels
	return imageSrc0At(adjustedTexelCoords)
}

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	uv := (texCoord.xy - .5*IResolution.xy) / IResolution.y
	UV := texCoord.xy / IResolution.xy
	M := ICursor.xy / IResolution.xy
	T := ITime + M.x*2.0

	T = mod(ITime, 102.0)
	T = mix(T, M.x*102.0, 0.0)

	t := T * 0.2
	rainAmount := sin(T*0.05)*0.3 + 0.7

	//maxBlur := mix(3.0, 6.0, rainAmount)
	//minBlur := 2.0

	story := 0.0
	heart := 0.0

	staticDrops := smoothstep(-0.5, 1.0, rainAmount) * 2.0
	layer1 := smoothstep(0.25, 0.75, rainAmount)
	layer2 := smoothstep(0.0, 0.5, rainAmount)

	c := Drops(uv, t, staticDrops, layer1, layer2)

	e := vec2(0.001, 0.0)
	cx := Drops(uv+e, t, staticDrops, layer1, layer2).x
	cy := Drops(uv+e.yx, t, staticDrops, layer1, layer2).x
	n := vec2(cx-c.x, cy-c.x) // expensive normals

	n *= 1.0 - smoothstep(60.0, 85.0, T)
	c.y *= 1.0 - smoothstep(80.0, 100.0, T)*0.8

	//focus := mix(maxBlur-c.y, minBlur, smoothstep(0.1, 0.2, c.x))
	//iCol := imageSrc0At(vec2(UV+n, focus))
	iCol := imageColorAtPixel(position.xy)
	col := vec3(iCol.r, iCol.g, iCol.b)

	t = (T + 3.0) * 0.5 // make time sync with first lightnoing
	colFade := sin(t*0.2)*0.5 + 0.5 + story
	col *= mix(vec3(1.0), vec3(0.8, 0.9, 1.3), colFade)    // subtle color shift
	fade := smoothstep(0.0, 10.0, T)                       // fade in at the start
	lightning := sin(t * sin(t*10.0))                      // lighting flicker
	lightning *= pow(max(0.0, sin(t+sin(t))), 10.0)        // lightning flash
	col *= 1.0 + lightning*fade*mix(1.0, 0.1, story*story) // composite lightning
	UVm := UV - 0.5
	col *= 1.0 - dot(UVm, UV)

	col = mix(pow(col, vec3(1.2)), col, heart)
	fade *= smoothstep(102.0, 97.0, T)
	col *= fade

	return vec4(col, 1.0)
}
