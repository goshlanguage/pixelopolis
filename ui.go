package main

import (
	"fmt"
	"reflect"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Button is meant to be embedded in a UI or Menu.
type Button struct {
	Text          string
	XPos, YPos    float32
	Width, Height float32
}

// GraphicButton represents a button represented by an image
type GraphicButton struct {
	Onclick func()
	Stamp   Stamp
}

// UI is meants to hold buttons and be overlayed on top of the screen
type UI struct {
	BuildingCache *Building // Stores building to be passed between Update and Draw
	Buttons       map[string]*Button
	ButtonValues  map[string]bool
	// Used to update if our preview cursor is collided
	CursorCollided bool
	Decorations    []Decoration
	// DrawFuncs store funcs that can be stored into the UI object to iterate through in the draw phase
	DrawFuncs []func()
	// Store a pointer to the engine to render its values as they Update
	Engine      *Engine
	Events      []*Event
	GroundLevel int32
	Halt        bool
	// Palettes give us the ability to toggle through texture maps
	// 0 - UI
	// 1 - City Tileset in Blue
	// 2 - City Tileset in Green
	// 3 - City Tileset in Yellow
	// 4 - City Tileset in Red
	Palettes         map[int]*Palette
	ScreenX, ScreenY int32
	SoundConfirm     rl.Sound
	SoundSelect      rl.Sound
	SoundCancel      rl.Sound
	Toggles          map[string]bool // Allows us to toggle things off and on
	XPos, YPos       float32
	Width, Height    float32
}

// Init initializes the UI, namely loads in assets to render the ui
func (ui *UI) Init() {
	ui.Palettes = make(map[int]*Palette)
	ui.Palettes[0] = GetUIPalette()
	ui.Palettes[1] = GetProjectMegaPalette("assets/sprites/projectmute.png")
	ui.Palettes[2] = GetProjectMegaPalette("assets/sprites/projectmuteG.png")
	ui.Palettes[3] = GetProjectMegaPalette("assets/sprites/projectmuteY.png")
	ui.Palettes[4] = GetProjectMegaPalette("assets/sprites/projectmuteR.png")

	ui.SoundConfirm = rl.LoadSound("assets/sounds/confirm.mp3")
	ui.SoundSelect = rl.LoadSound("assets/sounds/select.mp3")
	ui.SoundCancel = rl.LoadSound("assets/sounds/cancel.mp3")
	ui.Toggles = make(map[string]bool)
	ui.BuildingCache = &Building{}
}

// Draw renders the UI. Draw doesn't render the buttons, as they are used in the update process to map them to the ButtonValues map
func (ui *UI) Draw() {
	// Draw UI box at bottom of screen
	rl.DrawRectangleRec(rl.NewRectangle(ui.XPos, ui.YPos, ui.Width, ui.Height), rl.DarkGray)
	for _, f := range ui.DrawFuncs {
		f()
	}

	// if we're in building preview mode, look for collisions, then print whichever building stamp is in
	// the building cache
	if len(ui.Toggles) > 0 && ui.Toggles["drawPreview"] {
		if ui.CursorCollided {
			ui.BuildingCache.Stamp.Palette = ui.Palettes[3]
		} else if ui.Engine.Dosh < ui.BuildingCache.Cost {
			ui.BuildingCache.Stamp.Palette = ui.Palettes[4]
		} else {
			ui.BuildingCache.Stamp.Palette = ui.Palettes[2]
		}
		ui.BuildingCache.Draw()
	}

	fpsOffset := ui.ScreenX - rl.MeasureText("FPS: 000  ", 18)
	rl.DrawText(fmt.Sprintf("FPS: %v", rl.GetFPS()), fpsOffset, 20, 18, rl.Gold)
}

// Update renders the UI buttons so that it can store the values of the button bools to the ButtonValues map
func (ui *UI) Update() {
	for k, v := range ui.Buttons {
		ui.ButtonValues[k] = raygui.Button(rl.NewRectangle(v.XPos, v.YPos, v.Width, v.Height), v.Text)
	}

	// If the house button is clicked, render  the appropriate preview
	// The code that determines this is in the ui
	if ui.ButtonValues["house"] || ui.ButtonValues["slum"] || ui.ButtonValues["apartment"] || ui.ButtonValues["church"] {
		rl.PlaySound(ui.SoundSelect)
		// House button hook here
		ui.Toggles["drawPreview"] = !ui.Toggles["drawPreview"]
		if ui.ButtonValues["house"] {
			// build our bespoke building and store it in cache
			ui.BuildingCache = GetHouse(ui.Engine, ui.Palettes[2])
		}
		if ui.ButtonValues["slum"] {
			ui.BuildingCache = GetSlum(ui.Engine, ui.Palettes[2])
		}
		if ui.ButtonValues["apartment"] {
			ui.BuildingCache = GetApartment(ui.Engine, ui.Palettes[2])
		}
		if ui.ButtonValues["church"] {
			ui.BuildingCache = GetChurch(ui.Engine, ui.Palettes[2])
		}
	}

	for i, e := range ui.Events {
		if e.Trigger() && !e.Done {
			e.Execute()
			ui.Halt = true
			if !e.Triggered {
				rl.PlaySound(e.Sound)
				e.Triggered = true
			}
			if rl.IsKeyPressed(rl.KeySpace) && !e.Done {
				e.Done = true
				ui.Halt = false
			}
			if e.Done {
				ui.Events = append(ui.Events[:i], ui.Events[i+1:]...)
			}
		}
	}

	if len(ui.Toggles) > 0 && ui.Toggles["drawPreview"] {
		mouseX := rl.GetMouseX()
		ui.BuildingCache.Stamp.LevelX = float32(mouseX) - (ui.BuildingCache.Stamp.Width / 2)
		ui.BuildingCache.Stamp.LevelY = float32(ui.GroundLevel) - ui.BuildingCache.Stamp.Height + 16

		ui.CursorCollided = ui.Engine.IsCollidedWithType(ui.BuildingCache, reflect.TypeOf(&Building{}))

		// enable right click to exit
		if rl.IsMouseButtonPressed(rl.MouseRightButton) {
			rl.PlaySound(ui.SoundCancel)
			ui.Toggles["drawPreview"] = !ui.Toggles["drawPreview"]
		}

		// left click to place buildings
		if rl.IsMouseButtonDown(rl.MouseLeftButton) && !ui.CursorCollided && ui.Engine.Dosh >= ui.BuildingCache.Cost && rl.GetMouseY() <= ui.GroundLevel+100 {
			rl.PlaySound(ui.SoundConfirm)
			ui.Engine.Dosh -= ui.BuildingCache.Cost
			ui.Engine.PopulationMax += ui.BuildingCache.Population
			ui.Toggles["drawPreview"] = !ui.Toggles["drawPreview"]

			// Create new instance so pointer doesn't change all buildings all the time
			// TODO - Refactor this and other invocations to use a helper to make sure its setup
			building := &Building{}
			building.Stamp = ui.BuildingCache.Stamp
			building.Stamp.Palette = ui.Palettes[1]
			building.Stamp.LevelX = float32(mouseX) - (building.Stamp.Width / 2)
			building.Stamp.LevelY = float32(GroundLevel) - building.Stamp.Height + 16
			building.Engine = ui.Engine

			// Buildings to the front so they are rendered in the back
			ui.Engine.Entities = append([]Entity{building}, ui.Engine.Entities...)
		}
	}

}

// GetMainUI composes the UI for the main game loop
func GetMainUI(engine *Engine, groundLevel int, ScreenX, ScreenY int32) *UI {
	//height := (ScreenY / 5)
	height := 150

	ui := &UI{
		Engine:      engine,
		GroundLevel: int32(groundLevel),
		XPos:        0,
		YPos:        float32(ScreenY) - float32(height),
		Width:       float32(ScreenX),
		Height:      float32(height),
		ScreenX:     ScreenX,
		ScreenY:     ScreenY,
	}

	ui.Buttons = make(map[string]*Button)
	ui.ButtonValues = make(map[string]bool)

	ui.Buttons["house"] = &Button{"$1 - house", 10, float32(ScreenY - 130), 80, 40}
	ui.Buttons["slum"] = &Button{"$10 - slum", 10, float32(ScreenY - 80), 80, 40}
	ui.Buttons["apartment"] = &Button{"$100 - apt", 100, float32(ScreenY - 130), 80, 40}
	ui.Buttons["church"] = &Button{"$300 - church", 100, float32(ScreenY - 80), 80, 40}

	padding := rl.MeasureText("Population: 100000 / 100000", 18)
	yOffset := (ui.ScreenY / 12)

	ui.DrawFuncs = append(ui.DrawFuncs, func() {
		rl.DrawText(fmt.Sprintf("Population: %v / %v", ui.Engine.Population, ui.Engine.PopulationMax), ui.ScreenX-(padding), ui.ScreenY-yOffset, 18, rl.RayWhite)
	})
	ui.DrawFuncs = append(ui.DrawFuncs, func() {
		rl.DrawText(fmt.Sprintf("Dosh: $%.2f", ui.Engine.Dosh), ui.ScreenX-(padding), ui.ScreenY-(yOffset+18), 18, rl.RayWhite)
	})
	return ui
}

// GetUIPalette returns a pallet from a preconfigured, or swappable asset
func GetUIPalette() *Palette {
	palette := NewPalette("assets/sprites/ui.png", 16, 16)
	palette.Update()
	return palette
}
