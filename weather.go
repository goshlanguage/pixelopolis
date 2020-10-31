package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Droplet stores the coordinates and how to draw it
type Droplet struct {
	Color rl.Color
	// Allows for droplets to fall at varying rates
	Speed      int
	XPos, YPos int
}

// CanReap returns true when the droplet is off screen
func (d *Droplet) CanReap() bool {
	return d.YPos > int(rl.GetScreenHeight())
}

// Draw renders the droplet to the screen
func (d *Droplet) Draw() {
	rl.DrawRectangleRec(rl.NewRectangle(float32(d.XPos), float32(d.YPos), 4, 4), d.Color)
}

// Update imposes gravity on the droplet
func (d *Droplet) Update() {
	d.YPos += d.Speed
}

// Rain satisfies entity so it can be used on menus and in game
type Rain struct {
	Color    rl.Color
	Done     bool
	Droplets []*Droplet
	Music    rl.Music
}

// CanReap returns true when the droplet is off screen ( > rl.GetScreenHeight )
func (r *Rain) CanReap() bool {
	return r.Done
}

// Draw renders each droplet to the screen
func (r *Rain) Draw() {
	for _, droplet := range r.Droplets {
		droplet.Draw()
	}
}

// Init loads our SFX in
func (r *Rain) Init() {
	r.Music = rl.LoadMusicStream("assets/music/rain.mp3")
	rl.PlayMusicStream(r.Music)
	rl.SetMusicVolume(r.Music, .08)
}

// Update manages the rain inventory for droplets
// TODO - Rain is just a faster snow. Refactor to take a base speed for reuse between rain/snow
func (r *Rain) Update() {
	rl.UpdateMusicStream(r.Music)

	for i, droplet := range r.Droplets {
		droplet.Update()
		// Once the droplet is offscreen, clean it up
		if droplet.YPos > rl.GetScreenHeight() {
			if i == 0 {
				r.Droplets = r.Droplets[1:]
				continue
			}
			r.Droplets = append(r.Droplets[0:i], r.Droplets[i+1:]...)
		}
	}

	if rand.Intn(2) == 1 {
		droplet := &Droplet{Speed: 6 + rand.Intn(10), XPos: rand.Intn(rl.GetScreenWidth()), YPos: -4, Color: r.Color}
		r.Droplets = append(r.Droplets, droplet)
	}
}

// GoAway stops our music
func (r *Rain) GoAway() {
	rl.StopMusicStream(r.Music)
}

// GetHitbox returns a rectangle to represent the entity hitbox
// Kind of a hack to satisfy the entity interface
// TODO - Refactor into an Effect
func (r *Rain) GetHitbox() rl.Rectangle {
	return rl.NewRectangle(0, 0, 0, 0)
}
