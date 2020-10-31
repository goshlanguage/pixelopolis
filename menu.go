package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ryanhartje/raylib-go/raygui"
)

var (
	backgroundMusic rl.Music
)

// Menu is an abstraction for menus
type Menu struct {
	Title                        string
	BackgroundImagePath          string
	Buttons                      map[int]string
	ButtonXOffset, ButtonYOffset int
	ButtonFunctions              map[int]func()
	Effects                      []func()
	ScreenX, ScreenY             int32
}

// Loop starts a loop for the main menu
func (menu *Menu) Loop() {
	MenuInit()
	if Music {
		rl.PlayMusicStream(backgroundMusic)
	}

	var background rl.Texture2D
	bgDefined := menu.BackgroundImagePath != ""
	if bgDefined {
		background = rl.LoadTexture(menu.BackgroundImagePath)
	}

	rain := &Rain{Color: rl.LightGray}
	menu.Effects = append(menu.Effects, rain.Draw, rain.Update)
	for !rl.WindowShouldClose() {
		if Music {
			rl.UpdateMusicStream(backgroundMusic)
		}
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		if bgDefined {
			rl.DrawTextureRec(background, rl.NewRectangle(500, 0, float32(menu.ScreenX), float32(menu.ScreenY)), rl.NewVector2(0, 0), rl.White)
		}

		for _, e := range menu.Effects {
			e()
		}

		for i, b := range menu.Buttons {
			button := raygui.Button(rl.NewRectangle(float32(menu.ScreenX/2)-40, float32(menu.ScreenY)-float32(menu.ScreenY/5)*float32(i+1), 100, 30), b)
			if button {
				menu.ButtonFunctions[i]()
			}
		}

		titleWidth := rl.MeasureText(menu.Title, 64) / 2
		rl.DrawText(menu.Title, (menu.ScreenX/2)-titleWidth, 100, 64, rl.White)

		rl.EndDrawing()
	}
	//MenuClose()
}

// MenuInit initializes assets created for the menu
func MenuInit() {
	rl.InitAudioDevice()
	backgroundMusic = rl.LoadMusicStream("assets/music/menudrumroll.mp3")
}

// MenuClose cleans up menu assets
func MenuClose() {
	rl.UnloadMusicStream(backgroundMusic)
	rl.CloseAudioDevice()
}
