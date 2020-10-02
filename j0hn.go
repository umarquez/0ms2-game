package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/ungerik/go3d/float64/vec2"
	"image"
	"math"
	"path/filepath"
)

type PlayerDirection float64

const playerTick = 50
const (
	DirectionLeft PlayerDirection = iota - 1
	DirectionTop
	DirectionRight
)

var imgJ0hn *ebiten.Image
var leftOffsetRotation = vec2.T{-13, 33}
var rightOffsetRotation = vec2.T{34, -12}

func init() {
	imgJ0hn = loadSprite(filepath.Join(spritesPath, j0hnSpriteFile))
}

type J0hn struct {
	//direction        PlayerDirection
	rotation         float64
	position         *vec2.T
	acceleration     *vec2.T
	velocity         *vec2.T
	relativePosition *vec2.T
	upPosition       *vec2.T
	leftPosition     *vec2.T
	rightPosition    *vec2.T
	//friction     float64
	animationFrame int
	totalFrames    int
	isAccelerating bool
	timeAcumulator int64
	currentTileId  vec2.T

	flying bool
}

func NewJ0hn() *J0hn {
	w, _ := imgJ0hn.Size()

	return &J0hn{
		totalFrames:      w / playerSize,
		position:         new(vec2.T),
		acceleration:     new(vec2.T),
		velocity:         new(vec2.T),
		relativePosition: new(vec2.T),
	}
}

func (j0hn *J0hn) SetPosition(newPosition vec2.T) *J0hn {
	scrPosition := copyVector(newPosition)
	j0hn.upPosition = &scrPosition

	l := copyVector(scrPosition)
	l.Add(&leftOffsetRotation)
	j0hn.leftPosition = &l

	r := copyVector(scrPosition)
	r.Add(&rightOffsetRotation)
	j0hn.rightPosition = &r

	t := copyVector(newPosition)
	j0hn.position = &t
	return j0hn
}

func (j0hn *J0hn) Accelerate(amount *vec2.T) *J0hn {
	j0hn.flying = true
	j0hn.acceleration = j0hn.acceleration.Add(amount)
	j0hn.isAccelerating = true
	j0hn.animationFrame++
	if j0hn.animationFrame == 8 || j0hn.animationFrame == 14 || j0hn.animationFrame == 20 {
		j0hn.animationFrame = j0hn.animationFrame - 6
	}
	return j0hn
}

func (j0hn *J0hn) Steady() *J0hn {
	j0hn.acceleration = new(vec2.T)
	j0hn.isAccelerating = false
	j0hn.animationFrame = 1

	if j0hn.velocity.IsZero() {
		return j0hn
	}

	j0hn.velocity = j0hn.velocity.Scale(0.99)

	if j0hn.velocity.Length() < .5 {
		if j0hn.velocity[0] > j0hn.velocity[1] {

		}
	}

	//fmt.Printf("%#v\n", j0hn.velocity)

	return j0hn
}

func (j0hn *J0hn) StandUp() *J0hn {
	j0hn.velocity = new(vec2.T)
	j0hn.acceleration = new(vec2.T)
	j0hn.isAccelerating = false
	j0hn.animationFrame = 0

	return j0hn
}

func (j0hn *J0hn) Draw(screen *ebiten.Image) {
	op := ebiten.DrawImageOptions{}
	//op.GeoM.Rotate(1)
	op.GeoM.Rotate(j0hn.rotation)
	op.GeoM.Translate(float64(j0hn.position[0]), float64(j0hn.position[1]))
	op.GeoM.Scale(j0hnScale, j0hnScale)

	x1, y1 := j0hn.animationFrame*playerSize, 0
	x2, y2 := x1+playerSize, y1+playerSize

	screen.DrawImage(imgJ0hn.SubImage(image.Rect(x1, y1, x2, y2)).(*ebiten.Image), &op)
}

func (j0hn *J0hn) Update(_ *ebiten.Image, delta int64) {
	j0hn.timeAcumulator += delta

	if j0hn.timeAcumulator >= playerTick {
		j0hn.timeAcumulator = 0

		var direction = 0.0
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			direction = 1
			j0hn.rotation = -(45 * math.Pi) / 180
			j0hn.position = j0hn.leftPosition
		} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
			direction = -1
			j0hn.rotation = (45 * math.Pi) / 180
			j0hn.position = j0hn.rightPosition
		} else {
			j0hn.position = j0hn.upPosition
			j0hn.rotation = 0
		}

		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			j0hn.Accelerate(&vec2.T{direction * 2, 2})
		} else if !j0hn.flying {
			j0hn.StandUp()
		} else {
			j0hn.Steady()
		}

		if !j0hn.acceleration.IsZero() {
			j0hn.velocity = j0hn.velocity.Add(j0hn.acceleration)
		}

		var limit float64 = 20
		if j0hn.velocity[0] > limit {
			j0hn.velocity[0] = limit
		} else if j0hn.velocity[0] < -limit {
			j0hn.velocity[0] = -limit
		}

		if j0hn.velocity[1] > limit {
			j0hn.velocity[1] = limit
		} else if j0hn.velocity[1] < -limit {
			j0hn.velocity[1] = -limit
		}

		if !j0hn.velocity.IsZero() {
			j0hn.relativePosition = j0hn.relativePosition.Add(j0hn.velocity)
		}
	}
}
