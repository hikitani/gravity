package ecs

import (
	"errors"
	"iter"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	ErrEntityAlreadyInitialized = errors.New("entity already initialized")
	ErrEntityNotInitialized     = errors.New("entity not initialized")
)

type OnNewEntityObserver interface {
	OnNewEntity(ent Composer)
}

type Identifier interface {
	ID() (int, error)
	assignID(id int) error
}

type Composer interface {
	Identifier
	Components() []any
}

type Drawer interface {
	Draw(ent Composer, screen *ebiten.Image)
}

type Logic interface {
	Handle(ent Composer)
}

type WorldInjector struct {
	world *World
}

func (i *WorldInjector) inject(w *World) {
	i.world = w
}

func (i *WorldInjector) World() *World {
	return i.world
}

func tryInjectWorld(w *World, v any) {
	if i, ok := v.(interface {
		inject(w *World)
	}); ok {
		i.inject(w)
	}
}

type Entity struct {
	id   int
	init bool
}

func (ent *Entity) assignID(id int) error {
	if ent.init {
		return ErrEntityAlreadyInitialized
	}

	ent.id = id
	ent.init = true
	return nil
}

func (ent *Entity) ID() (int, error) {
	if !ent.init {
		return 0, ErrEntityNotInitialized
	}
	return ent.id, nil
}

var _ ebiten.Game = &World{}

type World struct {
	width, height int

	idNum int
	ents  []Composer

	draws []Drawer
	logic []Logic

	newEntityObservers []func(ent Composer)
}

func New(width, height int) *World {
	return &World{
		width:  width,
		height: height,
	}
}

func (w *World) AddObserverOnNewEntity(observer func(ent Composer)) {
	w.newEntityObservers = append(w.newEntityObservers, observer)
}

func (w *World) tryAddObserverOnNewEntity(v any) {
	if v, ok := v.(OnNewEntityObserver); ok {
		w.AddObserverOnNewEntity(v.OnNewEntity)
	}
}

func (w *World) AddDrawer(drawer Drawer) *World {
	tryInjectWorld(w, drawer)
	w.tryAddObserverOnNewEntity(drawer)
	w.draws = append(w.draws, drawer)
	return w
}

func (w *World) AddLogic(logic Logic) *World {
	tryInjectWorld(w, logic)
	w.tryAddObserverOnNewEntity(logic)
	w.logic = append(w.logic, logic)
	return w
}

// Draw implements ebiten.Game.
func (w *World) Draw(screen *ebiten.Image) {
	for _, drawer := range w.draws {
		for _, ent := range w.ents {
			drawer.Draw(ent, screen)
		}
	}
}

// Layout implements ebiten.Game.
func (w *World) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return w.width, w.height
}

// Update implements ebiten.Game.
func (w *World) Update() error {
	for _, logic := range w.logic {
		for _, ent := range w.ents {
			logic.Handle(ent)
		}
	}

	return nil
}

func (w *World) Get(id int) (Composer, bool) {
	if id < 0 || id >= w.idNum {
		return nil, false
	}

	return w.ents[id], true
}

func (w *World) AddEntity(ent Composer) {
	err := ent.assignID(w.idNum)
	if errors.Is(err, ErrEntityAlreadyInitialized) {
		return
	}
	if err != nil {
		panic("Unknown error when assigning ID to entity: " + err.Error())
	}
	w.idNum++
	w.ents = append(w.ents, ent)
	for _, observer := range w.newEntityObservers {
		observer(ent)
	}
}

func ComponentsFromEntities[Component any](w *World) iter.Seq[Component] {
	return func(yield func(Component) bool) {
		for _, ent := range w.ents {
			for _, comp := range ent.Components() {
				comp, ok := comp.(Component)
				if !ok {
					continue
				}
				if !yield(comp) {
					return
				}
				break
			}
		}
	}
}

// ValidateEntity checks if the entity has components and if they are pointers.
func ValidateEntity[E Composer]() error {
	var e E
	comps := e.Components()
	if len(comps) == 0 {
		return errors.New("entity has no components")
	}

	for _, comp := range comps {
		if !checkIsPointer(comp) {
			return errors.New("component is not a pointer, it will be copy")
		}
	}

	return nil
}

func checkIsPointer(v any) bool {
	return reflect.TypeOf(v).Kind() == reflect.Pointer
}
