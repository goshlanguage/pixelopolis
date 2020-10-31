package main

import (
	"math"
	"reflect"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Engine holds the game state
type Engine struct {
	BuildingBoxes []rl.Rectangle
	Counter       int
	Dosh          float64
	Effects       []func(*Engine)
	Entities      []Entity
	Lightcycle    rl.Color
	Pi            float64
	Population    int
	PopulationMax int
	Tax           float64
	UI            *UI
}

// Draw renders any all entities stored in the engine
func (e *Engine) Draw() {
	for _, e := range e.Entities {
		e.Draw()
	}
}

// Update updates all entities stored in the engine
func (e *Engine) Update() {
	houses := 0.0
	population := 0
	e.BuildingBoxes = []rl.Rectangle{}
	for i, entity := range e.Entities {
		if e.UI.Halt == false {
			entity.Update()
		} else {
			if reflect.TypeOf(entity) != reflect.TypeOf(&Person{}) && reflect.TypeOf(entity) != reflect.TypeOf(&Taxi{}) {
				entity.Update()
			}
		}
		// TODO - Refactor this. We don't need to run it every update most likely
		typeOfEntity := reflect.TypeOf(entity)
		if typeOfEntity == reflect.TypeOf(&Building{}) {
			e.BuildingBoxes = append(e.BuildingBoxes, entity.GetHitbox())
			houses++
		}
		if typeOfEntity == reflect.TypeOf(&Person{}) {
			population++
		}
		if entity.CanReap() {
			if len(e.Entities) > i {
				e.Entities = append(e.Entities[0:i], e.Entities[i+1:]...)
			} else {
				e.Entities = e.Entities[0:i]
			}
		}
	}
	e.Population = population

	// Here's the economy part
	// TODO - break this out into a package and design some tests to make the economy more iterable, and long term fun
	if e.Dosh == 0 {
		e.Dosh += 0.01
	}
	e.Dosh += float64(e.Population)*e.Tax*(0.0001*houses) + 0.0001

	// LightCycle Effects
	// A good timespan is around 2000 cycles. Cycles 0-200 Should be sun up - 800-1000 sun down - and times between at the peaks of the Pi
	// TODO - These constraints are poorly designed and results in an improper lightcycle. Resolve buggy lightcycle effects
	e.Counter++
	if e.Counter > 2000 {
		e.Counter = 0
	}
	if e.Counter <= 200 || (e.Counter >= 1000 && e.Counter <= 1200) {
		e.Pi += math.Pi / 400
	}
	if e.Pi > math.Pi {
		e.Pi = 0
	}

	R := uint8(49 + (191 * math.Sin(e.Pi)))
	G := uint8(52 + (187 * math.Sin(e.Pi)))
	B := uint8(49 + (191 * math.Sin(e.Pi)))
	A := uint8(200 + (55 * math.Sin(e.Pi)))

	e.Lightcycle = rl.Color{R: R, G: G, B: B, A: A}
	// fmt.Printf("Counter: %v\t Alpha: %v\tSinPi: %v\tPi: %v\n", e.Counter, (255 * math.Sin(e.Pi)), math.Sin(e.Pi), e.Pi)
}

// IsCollidedWithType takes a target entity and tells you if it's collided with another entity of the given type
func (e *Engine) IsCollidedWithType(target Entity, targetType reflect.Type) bool {
	for _, e := range e.Entities {
		if reflect.TypeOf(e) == targetType {
			if rl.CheckCollisionRecs(target.GetHitbox(), e.GetHitbox()) {
				return true
			}
		}
	}
	return false
}
