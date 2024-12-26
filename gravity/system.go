package gravity

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hikitani/blueprint/ecs"
)

const TPS = 60

type RenderSystem struct{}

func (b *RenderSystem) Draw(ent ecs.Composer, screen *ebiten.Image) {
	var (
		render   *Render
		position *Position
	)
	for _, comp := range ent.Components() {
		if render == nil {
			render, _ = comp.(*Render)
		}

		if position == nil {
			position, _ = comp.(*Position)
		}
	}

	if render == nil || position == nil {
		return
	}

	if render.Image == nil {
		return
	}

	var draw ebiten.DrawImageOptions
	bnd := render.Image.Bounds()
	shiftX, shiftY := bnd.Dx()/2, bnd.Dy()/2
	x, y := int(position.X)-shiftX, int(position.Y)-shiftY
	draw.GeoM.Translate(float64(x), float64(y))

	screen.DrawImage(render.Image, &draw)
}

type BlockSpawnerByClick struct {
	ecs.WorldInjector
}

func (b *BlockSpawnerByClick) Handle(ent ecs.Composer) {
	var mouse *MouseEvents

	for _, comp := range ent.Components() {
		if mouse == nil {
			mouse, _ = comp.(*MouseEvents)
		}
	}

	if mouse == nil {
		return
	}

	if mouse.LeftSideState == MouseStateClicked {
		b.World().AddEntity(&Block{
			Position: Position{
				X: float64(mouse.X),
				Y: float64(mouse.Y),
			},
			Render: Render{
				Image: blockImage,
			},
		})
	}

	if mouse.RightSideState == MouseStateClicked {
		b.World().AddEntity(&GravityBlock{
			Block: Block{
				Position: Position{
					X: float64(mouse.X),
					Y: float64(mouse.Y),
				},
				Render: Render{
					Image: gravityBlockImage,
				},
			},
			GravityAttraction: GravityAttraction{
				Radious:      150,
				Acceleration: 10,
			},
		})
	}
}

type InputChecker struct{}

func (inp *InputChecker) Handle(ent ecs.Composer) {
	var mouse *MouseEvents

	for _, comp := range ent.Components() {
		if mouse == nil {
			mouse, _ = comp.(*MouseEvents)
		}
	}

	if mouse == nil {
		return
	}

	mouse.X, mouse.Y = ebiten.CursorPosition()
	mouse.LeftSideState = inp.handleState(mouse.LeftSideState, ebiten.MouseButtonLeft)
	mouse.RightSideState = inp.handleState(mouse.RightSideState, ebiten.MouseButtonRight)
}

func (inp *InputChecker) handleState(current MouseState, btn ebiten.MouseButton) MouseState {
	switch current {
	case MouseStateNone:
		if ebiten.IsMouseButtonPressed(btn) {
			return MouseStatePressed
		}
	case MouseStatePressed:
		if !ebiten.IsMouseButtonPressed(btn) {
			return MouseStateClicked
		}
	case MouseStateClicked:
		return MouseStateNone
	default:
		panic("unreachable")
	}

	return current
}

type GravitySystem struct {
	ecs.WorldInjector

	gravityBlocks []struct {
		Gravity  *GravityAttraction
		Position *Position
	}
}

func (g *GravitySystem) OnNewEntity(ent ecs.Composer) {
	var (
		gravity  *GravityAttraction
		position *Position
	)
	for _, comp := range ent.Components() {
		if gravity == nil {
			gravity, _ = comp.(*GravityAttraction)
		}

		if position == nil {
			position, _ = comp.(*Position)
		}
	}

	if gravity == nil || position == nil {
		return
	}

	g.gravityBlocks = append(g.gravityBlocks, struct {
		Gravity  *GravityAttraction
		Position *Position
	}{
		Gravity:  gravity,
		Position: position,
	})
}

func (g *GravitySystem) Handle(ent ecs.Composer) {
	var (
		position *Position
		velocity *Velocity
		static   *IsStaticMarker
	)
	for _, comp := range ent.Components() {
		if position == nil {
			position, _ = comp.(*Position)
		}

		if velocity == nil {
			velocity, _ = comp.(*Velocity)
		}

		if static == nil {
			static, _ = comp.(*IsStaticMarker)
		}
	}

	if static != nil {
		return
	}

	if position == nil || velocity == nil {
		return
	}

	for _, gravity := range g.gravityBlocks {
		dist := position.DistanceTo(gravity.Position)

		if gravity.Gravity.Radious < dist.Distance {
			continue
		}

		velocity.X += gravity.Gravity.Acceleration * (dist.Dx) / dist.Distance * 1 / TPS
		velocity.Y += gravity.Gravity.Acceleration * (dist.Dy) / dist.Distance * 1 / TPS
	}
}

type MovementSystem struct{}

func (m *MovementSystem) Handle(ent ecs.Composer) {
	var (
		position *Position
		velocity *Velocity
	)
	for _, comp := range ent.Components() {
		if position == nil {
			position, _ = comp.(*Position)
		}

		if velocity == nil {
			velocity, _ = comp.(*Velocity)
		}
	}

	if position == nil || velocity == nil {
		return
	}

	position.X += velocity.X
	position.Y += velocity.Y
}
