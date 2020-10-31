package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	GRAVITY = 9.81
)

// Coin represents a coin drop
type Coin struct {
	Active         bool
	Counter        float64
	Done           bool
	Dosh           float64
	Engine         *Engine
	LevelX, LevelY float32
	Velocity       float64
	Sound          rl.Sound
	Sprite         *Sprite
}

// NewCoin generates a new coin at the coordinates provided
func NewCoin(engine *Engine, dosh float64, levelX float32, levelY float32) *Coin {
	var sound rl.Sound
	soundPath := "assets/sounds/coin1.mp3"
	if rand.Intn(2) == 1 {
		soundPath = "assets/sounds/coin2.mp3"
	}
	sound = rl.LoadSound(soundPath)

	sprite := &Sprite{}
	sprite.Init("assets/sprites/mega.png", 96, 32, 32, 32)
	sprite.LevelX = levelX
	sprite.LevelY = levelY

	return &Coin{
		Active:   true,
		Dosh:     dosh,
		Engine:   engine,
		LevelX:   levelX,
		LevelY:   levelY,
		Sound:    sound,
		Sprite:   sprite,
		Velocity: -10,
	}
}

// CanReap returns Coin.Done
func (coin *Coin) CanReap() bool {
	return coin.Done
}

// Draw satisfies the entity interface
func (coin *Coin) Draw() {
	coin.Sprite.Draw()
}

// Update performs the coin annimation
func (coin *Coin) Update() {
	if coin.Active {
		if coin.Dosh > 0 {
			coin.Engine.Dosh += coin.Dosh
			coin.Dosh = 0
		}
		coin.Velocity++
		// If the coin has landed, play the coin sound and transact the coin Dosh to the player,
		// Then deactivate the coin
		if float64(coin.Sprite.LevelY)+coin.Velocity >= float64(GroundLevel) {
			coin.Sprite.LevelY = float32(GroundLevel)
			rl.PlaySound(coin.Sound)
			coin.Engine.Dosh += coin.Dosh
			coin.Active = false
		} else {
			coin.Sprite.LevelY += float32(coin.Velocity)
		}
	} else {
		// Once the coin is disactive, animate a fade out, then Reap
		if coin.Counter < 20 {
			coin.Sprite.Color = rl.NewColor(255, 255, 255, uint8(int(255*(1/coin.Counter))))
		} else {
			coin.Sprite.Color = rl.NewColor(255, 255, 255, 0)
			coin.Sprite.Deleted = true
			coin.Done = true
		}
		coin.Counter += 0.25
	}
}

// GetHitbox returns a rectangle to represent the entity hitbox
func (coin *Coin) GetHitbox() rl.Rectangle {
	return rl.NewRectangle(coin.Sprite.XPos, coin.Sprite.YPos, coin.Sprite.Width, coin.Sprite.Height)
}
