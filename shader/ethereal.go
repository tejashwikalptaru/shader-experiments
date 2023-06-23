//go:build ignore

//kage:unit pixel

//https://www.shadertoy.com/view/DsjcRR

package shader

// uniform variables
var IResolution vec2
var ITime float

func easeOutSine(x float) float {
	PI := acos(-1.0)
	return sin((x * PI) / 2.0)
}

func hsv2rgb(c vec3) vec3 {
	rgb := clamp(abs(mod(c.x*6.0+vec3(0.0, 4.0, 2.0), 6.0)-3.0)-1.0, 0.0, 1.0)
	return c.z * mix(vec3(1.0), rgb, c.y)
}

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	PI := acos(-1.0)
	TWO_PI := PI * 2.0
	const STEPS = 8.0
	ANG_DIFF := TWO_PI / STEPS
	CIR_RAD := 0.1
	TIME_SCALER := 4.0
	LINE_W := 1.0 / max(IResolution.x, IResolution.y)

	uv := (2.0*texCoord.xy - IResolution.xy) / IResolution.y * 1.1

	interval := TWO_PI * (STEPS)
	scaledTime := ITime * TIME_SCALER
	intervalN := floor(scaledTime/interval) + 2.0
	time := mod(scaledTime, interval)
	totalColor := vec3(0.0)

	lastStart := (STEPS-1.0)*ANG_DIFF*0.5 + (STEPS-1.)*PI

	j := STEPS - 1.0
	for i := 0.0; i < STEPS; i++ {
		startTime := i * ANG_DIFF * 0.5
		endTime := lastStart + TWO_PI + j*PI - j*ANG_DIFF*0.5
		s := sin(time - startTime)
		cntr := vec2(0.0, (abs(s))*sign(s))
		a := startTime
		m := mat2(cos(a), -sin(a), sin(a), cos(a))
		cntr *= m

		start := startTime + i*PI
		rad := CIR_RAD * pow(clamp(time-start, 0.0, 1.0), 2.0) * pow(1.0-clamp(time-endTime, 0.0, 1.0), 4.0)

		circClr := hsv2rgb(vec3(fract((startTime+intervalN*37.17)*0.1), 1.0, 1.0)) * step(startTime*1.5, time)
		circleColor := LINE_W * 25.0 / abs(distance(uv, cntr)-rad)

		lineColor := 0.0

		even := sign(mod(i, 2.0) - 0.5)
		rotUV := vec2(uv.x*even, uv.y*-even) * m
		end := pow(1.0-clamp(time-endTime, 0.0, 1.0), 8.0)
		head := 0.9 * clamp(time-start, 0.0, 1.0)
		tail := 0.9 * clamp(time-start-PI, 0.0, 1.0)
		lineColor = LINE_W * end * (atan((rotUV.y+tail)/rotUV.x) - atan((rotUV.y-head)/rotUV.x)) / rotUV.x

		totalColor += circClr * max(circleColor, lineColor)

		j--
	}
	return vec4(totalColor, 1.0)
}
