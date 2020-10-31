package main

import rl "github.com/gen2brain/raylib-go/raylib"

// Entity allows custom objects to be grouped together through this interface
type Entity interface {
	CanReap() bool
	Draw()
	GetHitbox() rl.Rectangle
	Update()
}
