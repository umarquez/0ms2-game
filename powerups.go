package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	log "github.com/sirupsen/logrus"
	"github.com/ungerik/go3d/float64/vec2"
	"image/color"
	"math/rand"
	"path/filepath"
	"sort"
)

const powerupScale = 2

type PowerupType string

const FuelType PowerupType = "fuel"
const O2Type PowerupType = "o2"

type Powerup struct {
	id              uint
	position        vec2.T
	velocity        vec2.T
	sprite          *ebiten.Image
	op              *ebiten.DrawImageOptions
	playerInfluence float64
	puType          PowerupType
	collitionBox    vec2.Rect
}

func (powerup *Powerup) UpdatePosition(playerVelocity vec2.T) {
	v := copyVector(powerup.velocity)
	//v.Scale(2)
	v.Add(&playerVelocity)
	powerup.position.Add(&v)
}

func (powerup *Powerup) Draw(screen *ebiten.Image) {
	powerup.op.GeoM.Reset()
	powerup.op.GeoM.Translate(powerup.position[0], powerup.position[1])
	powerup.op.GeoM.Scale(powerupScale, powerupScale)

	// [Drawing collition box]
	if DrawCollitionBoxes {
		max := copyVector(powerup.collitionBox.Max)
		min := copyVector(powerup.collitionBox.Min)

		max.Sub(&min)
		size := []float64{max[0], max[1]}

		ebitenutil.DrawRect(screen, powerup.collitionBox.Min[0], powerup.collitionBox.Min[1], size[0], size[1], color.RGBA{
			R: 0x70,
			G: 0x70,
			B: 0xFF,
			A: 0x90,
		})
	}

	err := screen.DrawImage(powerup.sprite, powerup.op)
	if err != nil {
		log.Error(err)
	}
}

const PowerupsUpdateInterval = (1 / 60) * 1000
const puVelocityScale = .5
const newPowerupProbability = .5
const puMaxPlayerInfluence = .05
const powerupSize = 32

var PowerupsSprites map[PowerupType]*ebiten.Image

func init() {
	PowerupsSprites = make(map[PowerupType]*ebiten.Image)
	var err error
	PowerupsSprites[O2Type], _, err = ebitenutil.NewImageFromFile(filepath.Join(spritesPath, "o2.png"), ebiten.FilterNearest)
	if err != nil {
		log.Error(err)
	}

	PowerupsSprites[FuelType], _, err = ebitenutil.NewImageFromFile(filepath.Join(spritesPath, "gas.png"), ebiten.FilterNearest)
	if err != nil {
		log.Error(err)
	}
}

type PowerupsSpawner struct {
	activePowerups     map[uint]*Powerup
	drawablePowerups   []*Powerup
	lastId             uint
	timerAccumulator   int64
	player             *J0hn
	lastPlayerPosition vec2.T
}

func NewPowerupSpawner(player *J0hn) *PowerupsSpawner {
	Powerups := new(PowerupsSpawner)
	Powerups.player = player
	Powerups.activePowerups = make(map[uint]*Powerup)

	return Powerups
}

func (spawner *PowerupsSpawner) Update(_ *ebiten.Image, delta int64) {
	spawner.timerAccumulator += delta

	if spawner.timerAccumulator >= PowerupsUpdateInterval {
		newDrawables := []*Powerup{}
		for _, item := range spawner.activePowerups {
			v := copyVector(*spawner.player.velocity)
			v.Scale(item.playerInfluence)
			item.UpdatePosition(v)

			if item.position[1] > windowHeight/powerupScale {
				log.WithField("PowerupId", item.id).Debug("killing Powerup")
				delete(spawner.activePowerups, item.id)
			}

			if item.position[0] > -powerupSize*powerupScale &&
				item.position[1] > -powerupSize*powerupScale &&
				item.position[0] < windowWidth &&
				item.position[1] < windowHeight {
				newDrawables = append(newDrawables, item)
			}

			min := copyVector(item.position)
			min.Mul(&vec2.T{powerupScale, powerupScale})
			max := copyVector(min)
			max.Add(&vec2.T{planetSize * powerupScale, planetSize * powerupScale})

			vPos := &vec2.Rect{min, max}
			item.collitionBox = *vPos
			if spawner.player.Collition(vPos) {
				switch item.puType {
				case FuelType:
					spawner.player.AddFuel(100)
				case O2Type:
					spawner.player.AddO2(100)
				}
				delete(spawner.activePowerups, item.id)
			}
		}

		sort.Slice(newDrawables, func(i, j int) bool {
			return newDrawables[i].id < newDrawables[j].id
		})
		spawner.drawablePowerups = newDrawables

		if spawner.lastPlayerPosition != *spawner.player.position && spawner.player.flying && !spawner.player.isLifting && rand.Float64() < (float64(spawner.timerAccumulator)/500)*newPowerupProbability {
			spawner.lastPlayerPosition = *spawner.player.position
			fx := (rand.Float64() * 2) - .5
			px := fx * ((windowWidth - powerupSize) / powerupScale)
			fy := rand.Float64()
			if fx > 0 && fx < 1 {
				fy *= .5
			}
			fy -= .5

			py := fy * ((windowHeight - powerupSize) / powerupScale)

			initPos := vec2.T{px, py}
			initVel := copyVector(*spawner.player.position)
			initVel.Sub(&initPos)
			initVel.Normalize()
			initVel.Scale(rand.Float64() * puVelocityScale)

			puType := FuelType

			if rand.Float64() < .5 {
				puType = O2Type
			}

			p := Powerup{
				id:              spawner.lastId + 1,
				sprite:          PowerupsSprites[puType],
				op:              &ebiten.DrawImageOptions{},
				position:        initPos,
				velocity:        initVel,
				playerInfluence: (rand.Float64() * (puMaxPlayerInfluence / 2)) + (puMaxPlayerInfluence / 2),
				puType:          puType,
			}

			log.WithFields(map[string]interface{}{
				"position": initPos,
				"velocity": initVel,
			}).Debug("spawning new Powerup.")

			spawner.activePowerups[spawner.lastId+1] = &p
			spawner.lastId++
		}

		spawner.timerAccumulator = 0
	}
}

func (spawner *PowerupsSpawner) Draw(screen *ebiten.Image) {
	for _, item := range spawner.drawablePowerups {
		item.Draw(screen)
	}
}
