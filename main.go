package main

import (
	"math/rand"
	"os"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Keybindings tracks the user's key configuration
var (
	Music       = false
	Keybindings map[string]int32
	GroundLevel int
	ScreenX     = int32(rl.GetMonitorWidth(0))
	ScreenY     = int32(rl.GetMonitorHeight(0))
	// calculate the rate of increase. If we want it to happen once a minute, you have to consider the rate
	// that we're moving at is 60 FPS, so 60 * 60 = the rate of 1 minute.
	rate = 60.0 * 60.0
)

// main initializes raylib, and drops into the Main Menu
func main() {
	Init()
	ScreenX = int32(rl.GetScreenWidth())
	ScreenY = int32(rl.GetScreenHeight())
	GroundLevel = int(ScreenY - (ScreenY / 4))

	// Run the menu loop for the user
	mainMenu := &Menu{Title: "_ Pixelopolis _", Buttons: make(map[int]string), ScreenX: ScreenX, ScreenY: ScreenY}
	mainMenu.ButtonFunctions = make(map[int]func())
	mainMenu.ButtonFunctions[0] = func() {
		Run()
		os.Exit(0)
	}
	mainMenu.Buttons[0] = "Start"
	mainMenu.Loop()
}

// Run runs our game loop
func Run() {
	engine := &Engine{Dosh: 1, Tax: 1.05, Lightcycle: rl.RayWhite}

	// Group together some ground, grass, and skyline brushes to draw onto the screen for our background
	type tile struct {
		brush, x, y int
	}
	bgTiles := []tile{}

	for x := 0; x < rl.GetScreenWidth(); x += 16 {
		bgTiles = append(bgTiles, tile{171 + rand.Intn(4), x, GroundLevel})
	}

	// Setup our Taxi
	taxi := &Taxi{Passengers: 1, Engine: engine}
	taxi.Sound = rl.LoadSound("assets/sounds/taxi.mp3")
	taxi.Sprite.Init("assets/sprites/mega.png", 0, 1072, 96, 32)
	taxi.Sprite.Speed = 4
	// Spawn this off screen
	taxi.Sprite.LevelX = float32(ScreenX + 96)
	taxi.Sprite.LevelY = float32(GroundLevel)
	engine.Entities = append(engine.Entities, taxi)

	// rl.InitAudioDevice()
	backgroundMusic := rl.LoadMusicStream("assets/music/gameloop.mp3")
	if Music {
		rl.PlayMusicStream(backgroundMusic)
	}

	greet := "Welcome to Pixelopolis. Construct a slum or hovel to attract new citizens.\nPress space to continue..."

	ui := GetMainUI(engine, GroundLevel, ScreenX, ScreenY)
	ui.Init()
	ui.Events = append(ui.Events, NewEventDuration(7*time.Second, func() {
		NewDialog(greet, 48)
	}))
	ui.Toggles["drawPreview"] = false
	engine.UI = ui

	rain := &Rain{Color: rl.NewColor(57, 16, 90, 200)}
	rain.Init()
	engine.Entities = append(engine.Entities, rain)

	engine.Dosh = 300

	for !rl.WindowShouldClose() {
		rl.UpdateMusicStream(backgroundMusic)

		rl.BeginDrawing()
		rl.ClearBackground(engine.Lightcycle)

		for _, t := range bgTiles {
			if engine.Counter > 1000 {
				ui.Palettes[3].Draw(t.brush, t.x, t.y)
				rain.Color = rl.RayWhite
			} else {
				ui.Palettes[1].Draw(t.brush, t.x, t.y)
				rain.Color = rl.NewColor(57, 16, 90, 200)
			}
		}
		// Engine entities are triggered through this call
		engine.Draw()
		engine.Update()
		ui.Draw()
		ui.Update()

		rl.EndDrawing()
	}

	rl.CloseWindow()
}

// Init sets up our keybindings
func Init() {
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(ScreenX, ScreenY, "Pixelopolis")
	rl.SetTargetFPS(60)

	Keybindings = make(map[string]int32)
	Keybindings["forward"] = rl.KeyW
	Keybindings["backward"] = rl.KeyS
	Keybindings["left"] = rl.KeyA
	Keybindings["right"] = rl.KeyD
	Keybindings["space"] = rl.KeySpace
	Keybindings["exit"] = rl.KeyEscape
}
