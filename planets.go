package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	log "github.com/sirupsen/logrus"
	"github.com/ungerik/go3d/float64/vec2"
	"math/rand"
	"path/filepath"
	"sort"
)

const planetScale = 3

type Planet struct {
	id              uint
	position        vec2.T
	velocity        vec2.T
	sprite          *ebiten.Image
	op              *ebiten.DrawImageOptions
	playerInfluence float64
}

func (planet *Planet) Update(image *ebiten.Image, delta int64) {
	planet.position.Add(&planet.velocity)
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

const planetsUpdateInterval = 20
const velocityScale = 2
const newPlanetProbability = .2
const maxPlayerInfluence = .1

var planetsSprites []*ebiten.Image

func init() {
	for i := 1; i < 5; i++ {
		planetFile := fmt.Sprintf("planet-%v.png", i)
		img, _, err := ebitenutil.NewImageFromFile(filepath.Join(spritesPath, planetFile), ebiten.FilterNearest)
		if err != nil {
			log.WithField("sprite", planetFile).Error(err)
			continue
		}

		planetsSprites = append(planetsSprites, img)
	}
}

type PlanetsSpawner struct {
	activePlanets    map[uint]*Planet
	drawablePlanets  []*Planet
	lastId           uint
	timerAccumulator int64
	player           *J0hn
}

func NewPlanetSpawner(player *J0hn) *PlanetsSpawner {
	planets := new(PlanetsSpawner)
	planets.player = player
	planets.activePlanets = make(map[uint]*Planet)

	return planets
}

func (spawner *PlanetsSpawner) Update(image *ebiten.Image, delta int64) {
	spawner.timerAccumulator += delta

	if spawner.timerAccumulator >= planetsUpdateInterval {
		newDrawables := []*Planet{}
		for _, item := range spawner.activePlanets {
			v := copyVector(*spawner.player.velocity)
			item.position.Add(v.Scale(item.playerInfluence))
			item.Update(image, spawner.timerAccumulator)
			if item.position[1] > windowHeight/planetScale {
				log.WithField("planetId", item.id).Debug("killing planet")
				delete(spawner.activePlanets, item.id)
			}

			if item.position[0] > -32*planetScale &&
				item.position[1] > -32*planetScale &&
				item.position[0] < windowWidth &&
				item.position[1] < windowHeight {
				newDrawables = append(newDrawables, item)
			}
		}

		sort.Slice(newDrawables, func(i, j int) bool {
			return newDrawables[i].id < newDrawables[j].id
		})
		spawner.drawablePlanets = newDrawables

		if rand.Float64() < (float64(spawner.timerAccumulator)/1000)*newPlanetProbability {
			fx := (rand.Float64() * 2) - .5
			px := fx * ((windowWidth - 32) / planetScale)
			fy := rand.Float64()
			if fx > 0 && fx < 1 {
				fy *= .5
			}
			fy -= .5

			py := fy * ((windowHeight - 32) / planetScale)

			initPos := vec2.T{px, py}
			initVel := copyVector(*spawner.player.position)
			initVel.Sub(&initPos)
			initVel.Normalize()
			initVel.Scale((rand.Float64() * (velocityScale - 0.5)) + 0.5)

			p := Planet{
				id:              spawner.lastId + 1,
				sprite:          planetsSprites[rand.Intn(len(planetsSprites))],
				op:              &ebiten.DrawImageOptions{},
				position:        initPos,
				velocity:        initVel,
				playerInfluence: rand.Float64() * maxPlayerInfluence,
			}

			log.WithFields(map[string]interface{}{
				"position": initPos,
				"velocity": initVel,
			}).Debug("spawning new planet.")

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
