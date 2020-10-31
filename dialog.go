package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// NewDialog is a helper to spawn a new bit of game dialog
func NewDialog(text string, fontSize float32) {
	rl.DrawRectangleRec(
		rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight()/5)),
		rl.Black,
	)
	rl.DrawRectangleLinesEx(
		rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight()/5)),
		4,
		rl.White,
	)
	rl.DrawTextRecEx(
		rl.GetFontDefault(),
		text,
		rl.NewRectangle(20, 20, float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight()/5)),
		fontSize,
		1,
		true,
		rl.RayWhite,
		0,
		int32(rl.GetScreenWidth()),
		rl.White,
		rl.Black,
	)
}
