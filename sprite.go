package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Sprite represents a sprite in a spritesheet
type Sprite struct {
	Animated bool

	Color      rl.Color
	Counter    int
	Deleted    bool
	Effects    map[string]func()
	Frame      int // Tracks which cycle of animation the object should be in
	FrameCount int
	// LevelX, LevelY represent the X,Y screen coords
	LevelX, LevelY float32
	// Toggle true to render the sprite in reverse
	Reversed bool
	// Scale is for rendering at a different scale
	Scale float32
	// Speed represents a Sprite's step distance along the X axis
	Speed                     int
	Texture                   rl.Texture2D
	XPos, YPos, Width, Height float32 // Represents the sprite's parameters on a spritesheet
}

// Init sets the sprite's initial position on the provided spritesheet
func (s *Sprite) Init(filepath string, x, y, w, h float32) {
	s.Color = rl.White
	s.Texture = rl.LoadTexture(filepath)
	s.XPos = x
	s.YPos = y
	s.Width = w
	s.Height = h
	s.FrameCount = 1
	s.Scale = 1
	s.Speed = 1
}

// CanReap returns Sprite's Deleted bool
func (s *Sprite) CanReap() bool {
	return s.Deleted
}

// Draw renders the sprite to the screen in its frame of animation
func (s *Sprite) Draw() {
	rectangle := rl.NewRectangle(s.XPos+(float32(s.Frame)*s.Width), s.YPos, s.Width*s.Scale, s.Height)
	if s.Reversed {
		rectangle = rl.NewRectangle(s.XPos+(float32(s.Frame)*s.Width), s.YPos, -s.Width*s.Scale, s.Height)
	}
	position := rl.NewVector2(s.LevelX, s.LevelY)
	rl.DrawTextureRec(s.Texture, rectangle, position, s.Color)
}

// Update cycles the animation
// Counter should increment roughly 60 times a second at 60 FPS, so 60 % 10 should occur about 6 times a second
func (s *Sprite) Update() {
	if s.Counter%10 == 0 && s.Animated {
		s.Frame++
		if s.Frame > s.FrameCount {
			s.Frame = 1
		}
		if s.Counter == 60 {
			s.Counter = 0
		}
	}

	if !s.Animated {
		s.Frame = 0
	}
	s.Counter++
}

// GetHitbox returns a rectangle to represent the entity hitbox
func (s *Sprite) GetHitbox() rl.Rectangle {
	return rl.NewRectangle(s.XPos, s.YPos, s.Width, s.Height)
}
