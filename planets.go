package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	log "github.com/sirupsen/logrus"
	"github.com/ungerik/go3d/float64/vec2"
	"image"
	"math/rand"
	"path/filepath"
	"sort"
)

const planetScale = 4

type Planet struct {
	id              uint
	position        vec2.T
	velocity        vec2.T
	sprite          *ebiten.Image
	op              *ebiten.DrawImageOptions
	playerInfluence float64
}

func (planet *Planet) UpdatePosition(playerVelocity vec2.T) {
	v := copyVector(planet.velocity)
	//v.Scale(2)
	v.Add(&playerVelocity)
	planet.position.Add(&v)
}

func (planet *Planet) Draw(screen *ebiten.Image) {
	planet.op.GeoM.Reset()
	planet.op.GeoM.Translate(planet.position[0], planet.position[1])
	planet.op.GeoM.Scale(planetScale, planetScale)

	err := screen.DrawImage(planet.sprite, planet.op)
	if err != nil {
		log.Error(err)
	}
}

const planetSize = 32
const planetsUpdateInterval = (1 / 60) * 1000
const newPlanetProbability = .5
const planetVelocityScale = .05
const maxPlayerInfluence = .01

var planetsSprites []*ebiten.Image

func init() {
	planetFile := "planets.png"
	img, _, err := ebitenutil.NewImageFromFile(filepath.Join(spritesPath, planetFile), ebiten.FilterNearest)
	if err != nil {
		log.WithField("sprite", planetFile).Error(err)
	}

	for i := 0; i < img.Bounds().Max.X/planetSize; i++ {
		sprite := img.SubImage(image.Rect(planetSize*i, 0, planetSize*(i+1), planetSize)).(*ebiten.Image)
		planetsSprites = append(planetsSprites, sprite)
	}
}

type PlanetsSpawner struct {
	activePlanets      map[uint]*Planet
	drawablePlanets    []*Planet
	lastId             uint
	timerAccumulator   int64
	player             *J0hn
	lastPlayerPosition vec2.T
}

func NewPlanetSpawner(player *J0hn) *PlanetsSpawner {
	planets := new(PlanetsSpawner)
	planets.player = player
	planets.activePlanets = make(map[uint]*Planet)

	return planets
}

func (spawner *PlanetsSpawner) Update(_ *ebiten.Image, delta int64) {
	spawner.timerAccumulator += delta

	if spawner.timerAccumulator >= planetsUpdateInterval {
		newDrawables := []*Planet{}
		for _, item := range spawner.activePlanets {
			v := copyVector(*spawner.player.velocity)
			v.Scale(item.playerInfluence)
			item.UpdatePosition(v)

			if item.position[1] > windowHeight {
				log.WithField("planetId", item.id).Trace("killing planet")
				delete(spawner.activePlanets, item.id)
			}

			if item.position[0] > -planetSize*planetScale &&
				item.position[1] > -planetSize*planetScale &&
				item.position[0] < windowWidth &&
				item.position[1] < windowHeight {
				newDrawables = append(newDrawables, item)
			}
		}

		sort.Slice(newDrawables, func(i, j int) bool {
			return newDrawables[i].id < newDrawables[j].id
		})
		spawner.drawablePlanets = newDrawables

		if spawner.lastPlayerPosition != *spawner.player.position && spawner.player.flying && !spawner.player.isLifting && rand.Float64() < (float64(spawner.timerAccumulator)/100)*newPlanetProbability {
			spawner.lastPlayerPosition = *spawner.player.position
			fx := (rand.Float64() * 2) - .5
			px := fx * ((windowWidth - planetSize) / planetScale)
			fy := rand.Float64()
			if fx > 0 && fx < 1 {
				fy *= .5
			}
			fy -= .5

			py := fy * ((windowHeight - planetSize) / planetScale)

			initPos := vec2.T{px, py}
			initVel := copyVector(*spawner.player.position)
			initVel.Sub(&initPos)
			initVel.Normalize()
			initVel.Scale(rand.Float64() * planetVelocityScale)

			p := Planet{
				id:              spawner.lastId + 1,
				sprite:          planetsSprites[rand.Intn(len(planetsSprites))],
				op:              &ebiten.DrawImageOptions{},
				position:        initPos,
				velocity:        initVel,
				playerInfluence: (rand.Float64() * (maxPlayerInfluence / 2)) + (maxPlayerInfluence / 2),
			}

			log.WithFields(map[string]interface{}{
				"position": initPos,
				"velocity": initVel,
			}).Trace("spawning new planet.")

			spawner.activePlanets[spawner.lastId+1] = &p
			spawner.lastId++
		}

		spawner.timerAccumulator = 0
	}
}

func (spawner *PlanetsSpawner) Draw(screen *ebiten.Image) {
	for _, item := range spawner.drawablePlanets {
		item.Draw(screen)
	}
}
