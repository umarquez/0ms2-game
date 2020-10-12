package main

import (
	"github.com/ungerik/go3d/float64/vec3"
	"image/color"
	"sort"
)

type GradientColor struct {
	percentPosition float64
	color           *vec3.T
}

type Gradient []GradientColor

func (grad *Gradient) AddColor(position float64, c *vec3.T) {
	*grad = append(*grad, GradientColor{
		percentPosition: position,
		color:           c,
	})
}

func (grad *Gradient) GetColor(percent float64) color.Color {
	g := *grad
	sort.Slice(g, func(i, j int) bool {
		a := g[i]
		b := g[j]
		return a.percentPosition < b.percentPosition
	})

	grad = &g

	lastColor := GradientColor{}
	for _, c := range *grad {
		if c.percentPosition < percent {
			lastColor = c
			continue
		}

		o := vec3.T{c.color[0], c.color[1], c.color[2]}
		o.Sub(lastColor.color)
		p := percent - lastColor.percentPosition
		d := c.percentPosition - lastColor.percentPosition
		o.Scale(p / d)
		o.Add(lastColor.color)

		return color.RGBA{
			uint8(o[0]),
			uint8(o[1]),
			uint8(o[2]),
			255,
		}
	}
	g = *grad

	return color.RGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}
}
