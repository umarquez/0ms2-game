package main

import (
	"github.com/hajimehoshi/ebiten"
	log "github.com/sirupsen/logrus"
	"github.com/ungerik/go3d/float64/vec2"
)

var (
	positionCurrent vec2.T
	positionN       = vec2.T{0, -windowHeight}
	positionNE      = vec2.T{windowWidth, -windowHeight}
	//positionE       = vec2.T{windowWidth, 0}
	//positionW       = vec2.T{-windowWidth, 0}
	positionNW = vec2.T{-windowWidth, -windowHeight}
)

type Background struct {
	player          *J0hn
	tiles           map[Tile]bool
	timeAccumulator int64
	FirstTile       Tile
}

func NewBackgroundSystem(player *J0hn) *Background {
	t := NewInitTile(windowWidth, windowHeight*3, 1)

	p := &Background{
		player: player,
		tiles:  map[Tile]bool{t: true},
	}

	p.FirstTile = t
	p.player.currentTile = t
	return p
}

func (bg *Background) Update(_ *ebiten.Image, delta int64) {
	bg.timeAccumulator += delta

	if float64(bg.timeAccumulator) >= playerTick {
		bg.timeAccumulator = 0
		tilesCollection := map[vec2.T]Tile{
			positionCurrent: nil,
			positionN:       nil,
			positionNE:      nil,
			positionNW:      nil,
		}
		for tile := range bg.tiles {
			vel := copyVector(*bg.player.velocity)
			vel.Scale(float64(delta) / 300)

			tile.Update(vel)

			if !bg.player.isLifting {
				for k := range tilesCollection {
					pos := copyVector(bg.player.currentTile.GetPosition().Min)
					pos.Add(&k).Add(&vec2.T{windowWidth / 2, windowHeight / 2})

					if k == positionCurrent {
						pos = copyVector(*bg.player.position)
						pos.Scale(j0hnScale)
					}

					if tile.GetPosition().ContainsPoint(&pos) {
						tilesCollection[k] = tile
						break
					}
				}
			}

			if tile == bg.FirstTile {
				pos := *tile.GetPosition()
				if pos.Min[1] > 0 {
					bg.player.isLifting = false
				}
			}

			// offscreen? then delete instance
			if tile.IsOffscreen() {
				log.WithField("id", tile.GetId()).Debugln("Killing Tile")
				delete(bg.tiles, tile)
			}
		}

		if !bg.player.isLifting && tilesCollection[positionCurrent] == nil {
			log.WithFields(map[string]interface{}{
				"player_position": bg.player.collitionBox,
				"last_tile":       bg.player.currentTile.GetPosition(),
			}).Error("impossible to locate the current tile")
		} else if bg.player.isLifting {
		} else {
			if bg.player.currentTile.GetId() != tilesCollection[positionCurrent].GetId() {
				log.WithFields(map[string]interface{}{
					"player_position": bg.player.collitionBox,
					"last_tile":       tilesCollection[positionCurrent].GetPosition(),
				}).Tracef("%v contains player", tilesCollection[positionCurrent].GetId())
				bg.player.currentTile = tilesCollection[positionCurrent]
			}

			for k, v := range tilesCollection {
				if k == positionCurrent || v != nil {
					continue
				}

				pos := copyVector(bg.player.currentTile.GetPosition().Min)
				pos.Add(&k)
				t := NewStarsTile(pos)
				bg.tiles[t] = true
			}
		}

		log.WithFields(map[string]interface{}{
			"tiles_counter": len(bg.tiles),
		}).Traceln("Active tiles")
	}
}

func (bg *Background) Draw(screen *ebiten.Image) {
	for t := range bg.tiles {
		t.Draw(screen)
		//fmt.Printf("%#v: %#v\n", k, t.op.GeoM)
	}
}
