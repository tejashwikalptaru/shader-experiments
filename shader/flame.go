//go:build ignore

//kage:unit pixel

// https://www.shadertoy.com/view/DssyDj

package shader

// uniform variables
var IResolution vec2
var ITime float

func noise(p vec3) float {
	i := floor(p)
	a := dot(i, vec3(1.0, 57.0, 21.0)) + vec4(0.0, 57.0, 21.0, 78.0)
	f := cos((p-i)*acos(-1.0))*(-0.5) + 0.5
	a = mix(sin(cos(a)*a), sin(cos(1.0+a)*(1.0+a)), f.x)
	a.xy = mix(a.xz, a.yw, f.y)
	return mix(a.x, a.y, f.z)
}

func sphere(p vec3, spr vec4) float {
	return length(spr.xyz-p) - spr.w
}

func flame(p vec3) float {
	d := sphere(p*vec3(0.0, 0.5, 1.0), vec4(0.0, -1.0, 0.0, 1.0))
	return d + (noise(p+vec3(0.0, ITime*2.0, 0.0))+noise(p*3.0)*0.5)*0.25*(p.y)
}

func scene(p vec3) float {
	return min(100.0-length(p), abs(flame(p)))
}

func raymarch(org vec3, dir vec3) vec4 {
	d := 0.0
	glow := 0.0
	eps := 0.02
	p := org

	for i := 0; i < 256; i++ {
		d = scene(p) + eps
		p += d * dir
		if d > eps && flame(p) < 0.0 {
			glow = float(i) / 64.0
		}
	}
	return vec4(p, glow)
}

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	v := -1.0 + 2.0*texCoord.xy/IResolution.xy
	v.y *= -1
	v.x *= IResolution.x / IResolution.y

	org := vec3(0.0, -2.0, 4.0)
	dir := normalize(vec3(v.x*1.6, -v.y, -1.5))

	p := raymarch(org, dir)
	glow := p.w

	col := mix(vec4(1.0, 0.5, 0.1, 1.0), vec4(0.1, 0.5, 1.0, 1.0), p.y*0.02+0.4)

	return mix(vec4(0.0), col, pow(glow*2.0, 4.0))
}
