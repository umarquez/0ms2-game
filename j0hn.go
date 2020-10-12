package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	log "github.com/sirupsen/logrus"
	"github.com/ungerik/go3d/float64/vec2"
	"image"
	"image/color"
	"math"
	"path/filepath"
)

//type PlayerDirection float64

const playerTick = (1 / 60.0) * 1000
const frictionFactor = 0.99
const frameTime = 100

/*const (
	DirectionLeft PlayerDirection = iota - 1
	DirectionTop
	DirectionRight
)*/

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
	currentTile    Tile
	collitionBox   vec2.Rect
	frameStep      int64
	isLifting      bool

	o2, fuel float64

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
		fuel:             100,
		o2:               100,
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
	j0hn.acceleration.Add(amount)
	j0hn.isAccelerating = true

	if j0hn.velocity[0] > 50 {
		j0hn.velocity[0] = 50
	}

	if j0hn.velocity[1] > 50 {
		j0hn.velocity[1] = 50
	} else if j0hn.velocity[1] < -50 {
		j0hn.velocity[1] = -50
	}

	j0hn.velocity.Add(j0hn.acceleration)

	return j0hn
}

func (j0hn *J0hn) Steady() *J0hn {
	j0hn.acceleration = new(vec2.T)
	j0hn.isAccelerating = false
	j0hn.animationFrame = 1

	/*
		if j0hn.velocity.Length() < .5 {
			if j0hn.velocity[0] > j0hn.velocity[1] {

			}
		}*/

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

	// [ Drawing collition box behind J0hn, if enabled
	if DrawCollitionBoxes {
		max := copyVector(j0hn.collitionBox.Max)
		min := copyVector(j0hn.collitionBox.Min)

		max.Sub(&min)
		size := []float64{max[0], max[1]}

		ebitenutil.DrawRect(screen, j0hn.collitionBox.Min[0], j0hn.collitionBox.Min[1], size[0], size[1], color.RGBA{0xFF, 0x70, 0x70, 0x90})
	}

	_ = screen.DrawImage(imgJ0hn.SubImage(image.Rect(x1, y1, x2, y2)).(*ebiten.Image), &op)
	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf(
			"Fuel: %0.2f\nO2: %0.2f\nVelocity:\n  x: %0.2f\n  y: %0.2f",
			j0hn.fuel,
			j0hn.o2,
			j0hn.velocity[0],
			j0hn.velocity[1],
		),
		0,
		windowHeight/2,
	)
}

func (j0hn *J0hn) Update(_ *ebiten.Image, delta int64) {
	j0hn.timeAcumulator += delta
	j0hn.frameStep += delta

	if j0hn.isAccelerating && j0hn.fuel > 0 && j0hn.frameStep >= frameTime {
		j0hn.frameStep = 0
		j0hn.animationFrame++
		if j0hn.animationFrame == 8 {
			j0hn.animationFrame -= 6
		}
	}

	if float64(j0hn.timeAcumulator) >= playerTick {
		j0hn.timeAcumulator = 0

		if j0hn.isLifting {
			j0hn.velocity = &vec2.T{0, 100}
		} else {
			if j0hn.flying {
				j0hn.o2 -= float64(delta) / 1000
			}

			if j0hn.o2 < 0 {
				j0hn.o2 = 0
				j0hn.velocity = &vec2.Zero
			}

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

			if ebiten.IsKeyPressed(ebiten.KeySpace) && j0hn.fuel > 0 {
				if !j0hn.flying {
					j0hn.isLifting = true
				}
				amount := vec2.T{direction, 1}
				amount.Scale(1 / float64(delta))
				j0hn.Accelerate(&amount)

				if j0hn.fuel < 0 {
					j0hn.fuel = 0
				} else if j0hn.fuel > 0 {
					j0hn.fuel -= float64(delta) / 100
				}
			} else if !j0hn.flying {
				j0hn.StandUp()
			} else {
				j0hn.Steady()
				j0hn.velocity.Scale(frictionFactor)
			}
		}
		/*const limit float64 = 50
		if j0hn.velocity[0] > limit {
			j0hn.velocity[0] = limit
		} else if j0hn.velocity[0] < -limit {
			j0hn.velocity[0] = -limit
		}

		if j0hn.velocity[1] > limit {
			j0hn.velocity[1] = limit
		} else if j0hn.velocity[1] < -limit {
			j0hn.velocity[1] = -limit
		}*/

		if math.IsNaN(j0hn.velocity[0]) {
			j0hn.velocity[0] = 0
		}

		if math.IsNaN(j0hn.velocity[1]) {
			j0hn.velocity[1] = 0
		}

		v := copyVector(*j0hn.velocity)
		v.Scale(float64(delta) / 1000)
		j0hn.relativePosition = j0hn.relativePosition.Add(&v)
	}
}

func (j0hn *J0hn) AddO2(amount float64) {
	total := j0hn.o2 + amount
	if total > 100 {
		total = 100
	}

	j0hn.o2 = total
}

func (j0hn *J0hn) AddFuel(amount float64) {
	total := j0hn.fuel + amount
	if total > 100 {
		total = 100
	}

	j0hn.fuel = total
}

func (j0hn *J0hn) Collition(obj *vec2.Rect) bool {
	position := copyVector(*j0hn.upPosition)
	position.Scale(j0hnScale)
	max := copyVector(position)
	max.Add(&vec2.T{playerSize * j0hnScale, playerSize * j0hnScale})

	position.Add(&vec2.T{(playerSize * j0hnScale) / 4, 0})
	max.Sub(&vec2.T{(playerSize * j0hnScale) / 4, 0})

	playerArea := vec2.Rect{
		Min: position,
		Max: max,
	}

	j0hn.collitionBox = playerArea

	log.WithFields(map[string]interface{}{
		"player": playerArea,
		"obj":    obj,
	}).Trace("")

	return playerArea.ContainsPoint(&obj.Min) ||
		playerArea.ContainsPoint(&obj.Max) ||
		playerArea.ContainsPoint(&vec2.T{obj.Min[0], obj.Max[1]}) ||
		playerArea.ContainsPoint(&vec2.T{obj.Max[0], obj.Min[1]})
}
