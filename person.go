package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Person is an abstration for a person in the city
type Person struct {
	Counter  float32
	Deceased bool
	Dragged  bool
	Dosh     int
	Effects  []func(*Person)
	Engine   *Engine
	// Moving flags to animate the sprite
	OnTask bool
	Sprite Sprite
	// Step tracks the animation
	Step int
	// 0 - clicked
	Sounds    map[int]rl.Sound
	WaypointX float32
}

// CanReap returns when the Person is deceased.
func (person *Person) CanReap() bool {
	return person.Deceased
}

// IsFalling is a simple helper to stop other animations when falling
func (person *Person) IsFalling() bool {
	return int(person.Sprite.LevelY) < GroundLevel
}

// Init innit?
// TODO
// - Refactor sounds to be loaded in and played by detections in the engine. The engine should handle what assets to hold in memory and
//   when to queue theme
func (person *Person) Init(engine *Engine) {
	person.Engine = engine
	person.Sounds = make(map[int]rl.Sound)
	person.Sounds[0] = rl.LoadSound("assets/sounds/jump.mp3")
	person.Sounds[1] = rl.LoadSound("assets/sounds/arrived.mp3")
	rl.PlaySound(person.Sounds[1])
}

// Draw renders a person's sprite to the screen
func (person *Person) Draw() {
	person.Sprite.Draw()
}

// Update updates the sprites and runs any effects (like Person wandering etc)
func (person *Person) Update() {
	for _, e := range person.Effects {
		e(person)
	}
	if person.IsClicked() {
		person.Dragged = true
		// TODO - This worked ok except if the cursor exited the screen. Perhaps reimplement later
		/*
			if person.Dragged {
				rl.DisableCursor()
			} else {
				rl.EnableCursor()
			}
		*/
	}

	if person.Dragged {
		person.OnTask = false
		person.Counter = 0

		person.Sprite.LevelX = float32(rl.GetMouseX())
		if int(rl.GetMouseY()) <= GroundLevel {
			person.Sprite.LevelY = float32(rl.GetMouseY())
		} else {
			person.Sprite.LevelY = float32(GroundLevel)
		}
		rl.PlaySound(person.Sounds[0])

		if rl.IsMouseButtonDown(rl.MouseRightButton) || rl.IsMouseButtonUp(rl.MouseLeftButton) {
			person.Dragged = !person.Dragged
		}
	} else {
		if person.Sprite.LevelY < float32(GroundLevel) {
			if int(person.Sprite.LevelY+2.5+person.Counter) > GroundLevel {
				person.Sprite.LevelY = float32(GroundLevel)
			} else {
				person.Sprite.LevelY += 2.5 + person.Counter
			}
			person.Counter++
		}
		if !person.IsFalling() {
			person.Sprite.Update()
		}
	}
}

// GetHitbox returns a rectangle to represent the entity hitbox
func (person *Person) GetHitbox() rl.Rectangle {
	return rl.NewRectangle(person.Sprite.LevelX, person.Sprite.LevelY, person.Sprite.Width, person.Sprite.Height)
}

// IsClicked returns true when a person is clicked on
func (person *Person) IsClicked() bool {
	return rl.IsMouseButtonDown(rl.MouseLeftButton) && rl.CheckCollisionPointRec(rl.Vector2{X: float32(rl.GetMouseX()), Y: float32(rl.GetMouseY())}, person.GetHitbox())
}

// Wander is an effect intended to set a waypoint for a Person, then walk them to it.
func Wander(person *Person) {
	// rate 60 is roughly once a second
	rate := 60 * 5
	if !person.OnTask && !person.IsFalling() {
		person.Sprite.Animated = false
		if rand.Intn(rate) == 1 {
			waypoint := rand.Intn(rl.GetScreenWidth())
			remainder := waypoint % 4
			person.OnTask = true
			person.WaypointX = float32(waypoint - remainder)
			if person.WaypointX < 0 {
				person.WaypointX += 4
			}
		}
	} else {
		if person.Sprite.LevelX > person.WaypointX {
			person.Sprite.Animated = true
			person.Sprite.Reversed = true
			person.Sprite.LevelX -= float32(person.Sprite.Speed)
		}
		if person.Sprite.LevelX < person.WaypointX {
			person.Sprite.Animated = true
			person.Sprite.Reversed = false
			person.Sprite.LevelX += float32(person.Sprite.Speed)
		}
	}

	if person.Sprite.LevelX == person.WaypointX {
		person.OnTask = false
	}
}

// MoneyBags is a modifier that makes the person its effecting drop their dosh
func MoneyBags(person *Person) {
	// rate 60 * 60 is roughly once a minute
	rate := 60 * 60
	if person.Dosh > 0 && rand.Intn(rate) == 1 {
		// DropRate means we will drop X times where X=dropRate
		dropRate := 10
		coin := NewCoin(person.Engine, float64(person.Dosh/dropRate), person.Sprite.LevelX, person.Sprite.LevelY)
		person.Engine.Entities = append(person.Engine.Entities, coin)
	}
}
