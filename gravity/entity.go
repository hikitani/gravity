package gravity

import "github.com/hikitani/blueprint/ecs"

type Block struct {
	ecs.Entity
	Position
	Velocity
	Render
}

func (b *Block) Components() []any {
	return []any{&b.Position, &b.Render, &b.Velocity}
}

type Input struct {
	ecs.Entity
	MouseEvents
}

func (i *Input) Components() []any {
	return []any{&i.MouseEvents}
}

type GravityBlock struct {
	Block
	GravityAttraction
	IsStaticMarker
}

func (g *GravityBlock) Components() []any {
	return append(g.Block.Components(), &g.GravityAttraction, &g.IsStaticMarker)
}
