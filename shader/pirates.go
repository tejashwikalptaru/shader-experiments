//go:build ignore

//kage:unit pixel

// https://www.shadertoy.com/view/ldXXDj

package shader

// uniform variables
var IResolution vec2
var ITime float

func texture(channel int, uv vec2) vec4 {
	origin, size := imageSrcRegionOnTexture()
	coords := uv*size + origin
	if channel == 0 {
		return imageSrc0At(coords)
	} else if channel == 1 {
		return imageSrc1At(coords)
	} else if channel == 2 {
		return imageSrc2At(coords)
	} else if channel == 3 {
		return imageSrc3At(coords)
	}
	return vec4(0, 0, 0, 0)
}

func fbm(p vec2) float {
	return 0.5000*texture(1, p*1.00).x + 0.2500*texture(1, p*2.02).x + 0.1250*texture(1, p*4.03).x + 0.0625*texture(1, p*8.04).x
}

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	time := mod(ITime, 60.0)
	p := (2.0*texCoord - IResolution.xy) / IResolution.y
	p.y *= -1
	i := p

	// camera
	p += vec2(1.0, 3.0) * 0.001 * 2.0 * cos(ITime*5.0+vec2(0.0, 1.5))
	p += vec2(1.0, 3.0) * 0.001 * 1.0 * cos(ITime*9.0+vec2(1.0, 4.5))
	an := 0.3 * sin(0.1*time)
	co := cos(an)
	si := sin(an)
	p = mat2(co, -si, si, co) * p * 0.85

	// water
	q := vec2(p.x, 1.0) / p.y
	q.y -= 0.9 * time
	off := texture(0, 0.1*q*vec2(1.0, 2.0)-vec2(0.0, 0.007*ITime)).xy
	q += 0.4 * (-1.0 + 2.0*off)
	col := 0.2 * sqrt(texture(0, 0.05*q*vec2(1.0, 4.0)+vec2(0.0, 0.01*ITime)).zyx)
	re := 1.0 - smoothstep(0.0, 0.7, abs(p.x-0.6)-abs(p.y)*0.5+0.2)
	col += 1.0 * vec3(1.0, 0.9, 0.73) * re * 0.2 * (0.1 + 0.9*off.y) * 5.0 * (1.0 - col.x)
	re2 := 1.0 - smoothstep(0.0, 2.0, abs(p.x-0.6)-abs(p.y)*0.85)
	col += 0.7 * re2 * smoothstep(0.35, 1.0, texture(1, 0.075*q*vec2(1.0, 4.0)).x)

	// sky
	sky := vec3(0.0, 0.05, 0.1) * 1.4

	// stars
	sky += 0.5 * smoothstep(0.95, 1.00, texture(1, 0.25*p).x)
	sky += 0.5 * smoothstep(0.85, 1.0, texture(1, 0.25*p).x)
	sky += 0.2 * pow(1.0-max(0.0, p.y), 2.0)

	// clouds
	f := fbm(0.002 * vec2(p.x, 1.0) / p.y)
	cloud := vec3(0.3, 0.4, 0.5) * 0.7 * (1.0 - 0.85*smoothstep(0.4, 1.0, f))
	sky = mix(sky, cloud, 0.95*smoothstep(0.4, 0.6, f))
	sky = mix(sky, vec3(0.33, 0.34, 0.35), pow(1.0-max(0.0, p.y), 2.0))
	col = mix(col, sky, smoothstep(0.0, 0.1, p.y))

	// horizon
	col += 0.1 * pow(clamp(1.0-abs(p.y), 0.0, 1.0), 9.0)

	// moon
	d := length(p - vec2(0.6, 0.5))
	moon := vec3(0.98, 0.97, 0.95) * (1.0 - 0.1*smoothstep(0.2, 0.5, f))
	col += 0.8 * moon * exp(-4.0*d) * vec3(1.1, 1.0, 0.8)
	col += 0.2 * moon * exp(-2.0*d)
	moon *= 0.85 + 0.15*smoothstep(0.25, 0.7, fbm(0.05*p+0.3))
	col = mix(col, moon, 1.0-smoothstep(0.2, 0.22, d))

	// postprocess
	col = pow(1.4*col, vec3(1.5, 1.2, 1.0))
	col *= clamp(1.0-0.3*length(i), 0.0, 1.0)

	// fade
	col *= smoothstep(3.0, 6.0, time)
	col *= 1.0 - smoothstep(44.0, 50.0, time)

	return vec4(col, 1.0)
}
