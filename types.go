package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Level should be an exportable/importable data model for the game
type Level struct {
	Texture rl.Texture2D
}

// Init takes a background file, converts it to texture for rendering the level background
func (level *Level) Init(filepath string, width, height int32) {
	image := rl.LoadImage(filepath)
	rl.ImageResize(image, width, height)
	level.Texture = rl.LoadTextureFromImage(image)
}

// Draw renders the level texture to screen
func (level *Level) Draw() {
	bgRect := rl.NewRectangle(0, 0, 800, 500)
	bgPos := rl.NewVector2(0, 0)
	rl.DrawTextureRec(level.Texture, bgRect, bgPos, rl.White)
}

// Taxi represents a taxi that spawn people, bringing them into the city
type Taxi struct {
	Dosh       int
	Effects    []func(*Taxi)
	Engine     *Engine
	Entities   *[]Entity
	Passengers int
	Sound      rl.Sound
	SpawnRate  int
	Sprite     Sprite
}

// CanReap can't touch my taxi.
func (taxi *Taxi) CanReap() bool {
	return false
}

// Draw renders the taxi sprite to the screen
func (taxi *Taxi) Draw() {
	taxi.Sprite.Draw()
}

// Update drives the taxi along the X axis
func (taxi *Taxi) Update() {
	if taxi.Sprite.LevelX >= -taxi.Sprite.Width && taxi.Sprite.LevelX <= float32(rl.GetScreenWidth())+taxi.Sprite.Width {
		taxi.Sprite.LevelX += 4
	} else {
		// If we're not in motion, respawn if we get a random 1
		// rate is the rate of respawn. If we make the odds 1 in 60, we should expect to trigger this
		// once a second. We instead want to trigger it once a minute
		if taxi.SpawnRate == 0 {
			taxi.SpawnRate = 60 * 10
		}
		// Randomly spawn a taxi to drop off a person, assuming we have the population allowance
		if rand.Intn(taxi.SpawnRate) == 1 && taxi.Engine.PopulationMax > taxi.Engine.Population {
			taxi.Sprite.LevelX = -taxi.Sprite.Width
			rl.PlaySound(taxi.Sound)
		}

		if taxi.Engine.Population > 1 && float64(taxi.Engine.PopulationMax) >= (float64(taxi.Engine.Population)*2) {
			taxi.Passengers = 2
		} else {
			taxi.Passengers = 1
		}
	}

	if taxi.Sprite.LevelX == float32(rl.GetScreenWidth()/2) {
		randomizer := rand.Intn(4)

		// Randomly pick between the available choices of characters on the sprite sheet
		s := &Sprite{}
		s.Init("assets/sprites/mega.png", 0, 864+float32(32*randomizer), 32, 32)
		s.FrameCount = 4

		for i := 0; i < taxi.Passengers; i++ {
			p := &Person{Dosh: rand.Intn(100)}
			p.Init(taxi.Engine)
			p.Sprite = *s
			p.Sprite.LevelX = taxi.Sprite.LevelX
			p.Sprite.LevelY = taxi.Sprite.LevelY
			p.Effects = append(p.Effects, Wander)
			// Spawn a money bags passenger about 1 in 10 times
			if rand.Intn(10) == 1 {
				p.Effects = append(p.Effects, MoneyBags)
			}
			taxi.Engine.Entities = append(taxi.Engine.Entities, p)
		}
		// fmt.Println("SPAWN at: %d, %d", p.Sprite.LevelX, p.Sprite.LevelY)

		// Taxis give a flat fare per delivery, great for early game
		coin := NewCoin(taxi.Engine, 2.5, taxi.Sprite.LevelX, taxi.Sprite.LevelY)

		// Add the taxi and the fare tax to the entity bag
		taxi.Engine.Entities = append(taxi.Engine.Entities, coin)
	}
}

// GetHitbox returns a rectangle to represent the entity hitbox
func (taxi *Taxi) GetHitbox() rl.Rectangle {
	return rl.NewRectangle(taxi.Sprite.XPos, taxi.Sprite.YPos, taxi.Sprite.Width, taxi.Sprite.Height)
}
