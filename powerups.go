package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	log "github.com/sirupsen/logrus"
	"github.com/ungerik/go3d/float64/vec2"
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
}

func (Powerup *Powerup) Update(image *ebiten.Image, delta int64) {
	Powerup.position.Add(&Powerup.velocity)
}

func (Powerup *Powerup) Draw(screen *ebiten.Image) {
	Powerup.op.GeoM.Reset()
	Powerup.op.GeoM.Translate(Powerup.position[0], Powerup.position[1])
	Powerup.op.GeoM.Scale(powerupScale, powerupScale)

	err := screen.DrawImage(Powerup.sprite, Powerup.op)
	if err != nil {
		log.Error(err)
	}
}

const PowerupsUpdateInterval = 20
const puVelocityScale = 2
const newPowerupProbability = .1
const puMaxPlayerInfluence = .5

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
	activePowerups   map[uint]*Powerup
	drawablePowerups []*Powerup
	lastId           uint
	timerAccumulator int64
	player           *J0hn
}

func NewPowerupSpawner(player *J0hn) *PowerupsSpawner {
	Powerups := new(PowerupsSpawner)
	Powerups.player = player
	Powerups.activePowerups = make(map[uint]*Powerup)

	return Powerups
}

func (spawner *PowerupsSpawner) Update(image *ebiten.Image, delta int64) {
	spawner.timerAccumulator += delta

	if spawner.timerAccumulator >= PowerupsUpdateInterval {
		newDrawables := []*Powerup{}
		for _, item := range spawner.activePowerups {
			v := copyVector(*spawner.player.velocity)
			item.position.Add(v.Scale(item.playerInfluence))
			item.Update(image, spawner.timerAccumulator)
			if item.position[1] > windowHeight/powerupScale {
				log.WithField("PowerupId", item.id).Debug("killing Powerup")
				delete(spawner.activePowerups, item.id)
			}

			if item.position[0] > -32*powerupScale &&
				item.position[1] > -32*powerupScale &&
				item.position[0] < windowWidth &&
				item.position[1] < windowHeight {
				newDrawables = append(newDrawables, item)
			}

			min := copyVector(item.position)
			min.Mul(&vec2.T{powerupScale, powerupScale})
			max := copyVector(min)
			max.Add(&vec2.T{planetSize, planetSize})

			vPos := &vec2.Rect{min, max}
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

		if rand.Float64() < (float64(spawner.timerAccumulator)/1000)*newPowerupProbability {
			fx := (rand.Float64() * 2) - .5
			px := fx * ((windowWidth - 32) / powerupScale)
			fy := rand.Float64()
			if fx > 0 && fx < 1 {
				fy *= .5
			}
			fy -= .5

			py := fy * ((windowHeight - 32) / powerupScale)

			initPos := vec2.T{px, py}
			initVel := copyVector(*spawner.player.position)
			initVel.Sub(&initPos)
			initVel.Normalize()
			initVel.Scale((rand.Float64() * (puVelocityScale - 0.5)) + 0.5)

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
				playerInfluence: rand.Float64() * puMaxPlayerInfluence,
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
