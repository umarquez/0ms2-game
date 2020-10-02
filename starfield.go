package main

import (
	"errors"
	"github.com/hajimehoshi/ebiten"
	log "github.com/sirupsen/logrus"
	"github.com/ungerik/go3d/float64/vec2"
	"image/color"
	"math/rand"
)

type Tile struct {
	startingPoint *vec2.T
	img           *ebiten.Image
	position      *vec2.T
	op            *ebiten.DrawImageOptions
}

func NewTile(origin, position vec2.T) *Tile {
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(position[0]), float64(position[1]))

	max := copyVector(origin)

	max.Add(&vec2.T{
		windowWidth,
		windowHeight,
	})

	tile := Tile{
		startingPoint: &origin,
		position:      &position,
		op:            &op,
	}

	log.WithFields(map[string]interface{}{
		"starting_point": tile.startingPoint,
		"position":       tile.position,
	}).Debugf("new tile")

	tile.img, _ = ebiten.NewImage(windowWidth, windowHeight, ebiten.FilterNearest)
	for v := float64(windowWidth*windowHeight) * starsProportion; v > 0; v-- {
		tile.img.Set(rand.Intn(windowWidth), rand.Intn(windowHeight), color.White)
	}

	return &tile
}

type Starfield struct {
	player         *J0hn
	tiles          map[string]*Tile
	timeAcumulator int64
}

func NewStarfield(player *J0hn) *Starfield {
	p := &Starfield{
		player: player,
		tiles:  make(map[string]*Tile),
	}
	tile := NewTile(vec2.T{}, vec2.T{})
	p.tiles[vec2.Zero.String()] = tile
	return p
}

func (stars *Starfield) Update(screen *ebiten.Image, delta int64) {
	stars.timeAcumulator += delta
	//fmt.Println(stars.timeAcumulator)

	if stars.timeAcumulator >= playerTick {
		stars.timeAcumulator = 0

		var currentTile *vec2.T

		for id, tile := range stars.tiles {
			max := copyVector(*tile.position)
			min := copyVector(*tile.position)
			max.Add(&vec2.T{
				windowWidth,
				windowHeight,
			})
			absArea := vec2.Rect{
				Min: min,
				Max: max,
			}

			if absArea.ContainsPoint(stars.player.position) {
				log.WithFields(map[string]interface{}{
					"player position":       stars.player.position,
					"current tile position": tile.position,
				}).Traceln("player position")
				currentTile = tile.startingPoint
			}

			vel := copyVector(*stars.player.velocity)
			vel.Scale(.1)
			tile.position.Add(&vel)
			//fmt.Printf("%#v\n", tile.position)

			tile.op.GeoM.Reset()
			tile.op.GeoM.Translate(float64(tile.position[0]), float64(tile.position[1]))
			tile.op.GeoM.Scale(1, 1)

			// offscreen? then delete instance
			if tile.position[1] > windowHeight {
				log.WithField("id", id).Debugln("Killing Tile")
				delete(stars.tiles, id)
			}
		}

		if currentTile == nil {
			Panic("Starfield.Update", map[string]interface{}{
				"position": stars.player.position,
			}, errors.New("current tile does not exists, there must be a very good reason for this... revert whatever you touch"))
		}

		if stars.player.currentTileId != *currentTile {
			/*log.WithFields(map[string]interface{}{
				"tilesTotal": len(stars.tiles),
				"tile":     currentTile,
				"position": stars.player.position,
			}).Debugln(currentTile.ContainsPoint(stars.player.position))*/

			stars.player.currentTileId = *currentTile
		}

		for v := -1.0; v <= 1; v += 1 {
			stepTop := &vec2.T{
				windowWidth * v,
				windowHeight,
			}
			//log.Debugf("v: %#v", stepTop)
			tileId := copyVector(*currentTile)
			tileId.Add(stepTop)

			if _, ok := stars.tiles[tileId.String()]; !ok {
				pos := copyVector(*stars.tiles[currentTile.String()].position)
				pos.Sub(stepTop)
				stars.tiles[tileId.String()] = NewTile(copyVector(tileId), pos)
			}
		}

		switch {
		case stars.player.velocity[0] < 0:
			stepSide := &vec2.T{-windowWidth, 0}
			tileId := copyVector(*currentTile)
			tileId.Add(stepSide)

			if _, ok := stars.tiles[tileId.String()]; !ok {
				pos := copyVector(*stars.tiles[currentTile.String()].position)
				pos.Sub(stepSide)
				stars.tiles[tileId.String()] = NewTile(copyVector(tileId), pos)
			}
		case stars.player.velocity[0] > 0:
			stepSide := &vec2.T{windowWidth, 0}
			tileId := copyVector(*currentTile)
			tileId.Add(stepSide)

			if _, ok := stars.tiles[tileId.String()]; !ok {
				pos := copyVector(*stars.tiles[currentTile.String()].position)
				pos.Sub(stepSide)
				stars.tiles[tileId.String()] = NewTile(copyVector(tileId), pos)
			}
		}

		log.WithFields(map[string]interface{}{
			"tiles_counter": len(stars.tiles),
		}).Traceln("Active tiles")
	}
}

func (stars *Starfield) Draw(screen *ebiten.Image) {
	for k, t := range stars.tiles {
		//if k.Min.Eq(image.Point{
		//	X: 0,
		//	Y: 0,
		//}) {
		//	continue
		//}
		_ = k
		screen.DrawImage(t.img, t.op)
		//fmt.Printf("%#v: %#v\n", k, t.op.GeoM)
	}
}
