package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Building provides an abstraction for buildings. Give it a stamp, or a collection of brushes
//	and it's coordinate pairing,
type Building struct {
	Cost        float64
	Counter     int
	Deleted     bool
	Decorations []Decoration
	Effects     []func(*Building)
	Engine      *Engine
	Filepath    string
	Palette     *Palette
	Population  int
	Stamp       *Stamp
}

// CanReap returns building.Deleted, designed to be toggled if a building is demolished
func (building *Building) CanReap() bool {
	return building.Deleted
}

// Draw renders the stamp in the X,Y coordinates given
func (building *Building) Draw() {
	building.Stamp.Draw()
	for _, d := range building.Decorations {
		d.Draw()
	}
}

// Update runs any building effects. This could be used to build levels over time
// Decorate if the building is still plain
func (building *Building) Update() {
	if len(building.Decorations) < 3 {
		Decorate(building)
	}

	// Run FX
	for _, e := range building.Effects {
		e(building)
	}
}

// GetHitbox returns a rectangle to represent the entity hitbox
func (building *Building) GetHitbox() rl.Rectangle {
	return rl.NewRectangle(building.Stamp.LevelX, building.Stamp.LevelY, building.Stamp.Width, building.Stamp.Height)
}

// GetHitboxLeft returns a rectangle to represent the entity hitbox with an offset to the left to detect anything nearby
func (building *Building) GetHitboxLeft() rl.Rectangle {
	return rl.NewRectangle(building.Stamp.LevelX-16, building.Stamp.LevelY, building.Stamp.Width+16, building.Stamp.Height)
}

// GetHitboxRight returns a rectangle to represent the entity hitbox with an offset to the right to detect anything nearby
func (building *Building) GetHitboxRight() rl.Rectangle {
	return rl.NewRectangle(building.Stamp.LevelX, building.Stamp.LevelY, building.Stamp.Width+16, building.Stamp.Height)
}

// GetHouse puts together a 2 story building with a door
func GetHouse(engine *Engine, palette *Palette) *Building {
	door := rand.Intn(3)
	window := rand.Intn(3)
	for window == door {
		window = rand.Intn(3)
	}
	buildingStamp := &Stamp{Palette: palette, Width: 48, Height: 32}
	buildingStamp.DrawCoords = append(buildingStamp.DrawCoords, DrawCoord{75, 0, 0})                    // top left
	buildingStamp.DrawCoords = append(buildingStamp.DrawCoords, DrawCoord{107, 0, 16})                  // left
	buildingStamp.DrawCoords = append(buildingStamp.DrawCoords, DrawCoord{77, 16, 0})                   // top center
	buildingStamp.DrawCoords = append(buildingStamp.DrawCoords, DrawCoord{109, 16, 16})                 // filler
	buildingStamp.DrawCoords = append(buildingStamp.DrawCoords, DrawCoord{79, 32, 0})                   // top right
	buildingStamp.DrawCoords = append(buildingStamp.DrawCoords, DrawCoord{111, 32, 16})                 // right
	buildingStamp.DrawCoords = append(buildingStamp.DrawCoords, DrawCoord{103, float32(door) * 16, 16}) // door randomized
	if rand.Intn(3) == 1 {
		buildingStamp.DrawCoords = append(buildingStamp.DrawCoords, DrawCoord{20, float32(window) * 16, 16}) // window randomized
	}

	return &Building{
		Cost:       1,
		Engine:     engine,
		Population: 1,
		Stamp:      buildingStamp,
	}
}

// GetSlum puts together a 3 story slum building with a door
func GetSlum(engine *Engine, palette *Palette) *Building {
	filepath := "assets/buildings/slum/2.json"

	buildingStamp, err := GetStampFromTiledFile(filepath)
	if err != nil {
		panic("at the disco")
	}

	return &Building{
		Cost:       10,
		Effects:    []func(*Building){Decorate},
		Engine:     engine,
		Population: 6,
		Stamp:      buildingStamp,
	}
}

// GetApartment puts together a 3 story slum building with a door
func GetApartment(engine *Engine, palette *Palette) *Building {
	filepath := "assets/buildings/slum/1.json"
	buildingStamp, _ := GetStampFromTiledFile(filepath)

	return &Building{
		Cost:       100,
		Effects:    []func(*Building){Decorate},
		Engine:     engine,
		Population: 12,
		Stamp:      buildingStamp,
	}
}

// GetChurch puts together a 3 story slum building with a door
func GetChurch(engine *Engine, palette *Palette) *Building {
	filepath := "assets/buildings/slum/church.json"
	buildingStamp, _ := GetStampFromTiledFile(filepath)

	return &Building{
		Cost:       300,
		Effects:    []func(*Building){Decorate},
		Engine:     engine,
		Population: 0,
		Stamp:      buildingStamp,
	}
}

// Decoration  represents decore on buildings
type Decoration struct {
	AttachedTo *Building
	Palette    *Palette
	Stamp      *Stamp
	Removed    bool
}

// CanReap returns true if its removed
func (d *Decoration) CanReap() bool {
	return d.Removed
}

// Draw draws
func (d *Decoration) Draw() {
	d.Stamp.Draw()
}

// Update should clean up decorations if it starts to collide
func (d *Decoration) Update() {

}

// GetHitbox returns a rectange for the building
func (d *Decoration) GetHitbox() rl.Rectangle {
	return rl.NewRectangle(d.Stamp.LevelX, d.Stamp.LevelY, d.Stamp.Width, d.Stamp.Height)
}

// Decorate adds some decore to buildings in certain parameters
func Decorate(building *Building) {
	decorBrush := 139 + rand.Intn(5)
	decorX := 16 * rand.Intn(int(building.Stamp.Width/16))
	decorY := 16 * rand.Intn(int(building.Stamp.Height/16))
	decorStamp := &Stamp{Palette: building.Stamp.Palette, LevelX: building.Stamp.LevelX, LevelY: building.Stamp.LevelY, Width: 16, Height: 16}
	decorStamp.DrawCoords = append(decorStamp.DrawCoords, DrawCoord{decorBrush, float32(decorX), float32(decorY)}) // top left

	decoration := Decoration{AttachedTo: building, Palette: building.Palette, Stamp: decorStamp}
	building.Decorations = append(building.Decorations, decoration)
}

// Layer represents a portion of the JSON file saved from Tiled
type Layer struct {
	Data    []int `json:"data"`
	Height  int   `json:"height"`
	Width   int   `json:"width"`
	Opacity int   `json:"opacity"`
	X       int   `json:"x"`
	Y       int   `json:"y"`
}

// Tiled represents the tiled file, or the JSON file exported from Tiled
type Tiled struct {
	Layers     []Layer `json:"layers"`
	TileHeight int     `json:"tileheight"`
	TileWidth  int     `json:"tilewidth"`
	Width      int     `json:"width"`
	Height     int     `json:"height"`
}

// GetStampFromTiledFile takes a filepath to a file saved from the popular tile map program, tiled:
// https://www.mapeditor.org/
func GetStampFromTiledFile(filepath string) (*Stamp, error) {
	tiled, err := loadTiled(filepath)
	if err != nil {
		return &Stamp{}, err
	}

	stamp := ParseTiled(tiled)

	return stamp, nil
}

// ParseTiled takes the data from a tiled layer, and parses it into a stamp
// We have an odd task of taking a flat 1D slice of ints, representing the tile # of the tileset
// file, and figuring out how to render that.
//
// The approach we've taken is to load in tilesets with a known height/width, and iterate through each
// possible tile, starting with 1 and incrementing until the file is read. This is how we populate our
// brush map
func ParseTiled(t *Tiled) *Stamp {
	stamp := &Stamp{}
	stamp.Height = float32(t.Height * t.TileHeight)
	stamp.Width = float32(t.Width * t.TileWidth)

	for _, layer := range t.Layers {
		y := 0
		x := 0
		for counter, tile := range layer.Data {
			// This silly looking thing is because Tiled saves the tile number with + 1.
			tileCorrection := tile - 1
			if counter == 0 {
				x = 0
			} else {
				x = counter % layer.Width
			}

			// Skip if tile is 0 (blank)
			if tile != 0 {
				stamp.DrawCoords = append(stamp.DrawCoords, DrawCoord{tileCorrection, float32(x * t.TileWidth), float32(y * t.TileHeight)})
			}
			if (counter+1)%layer.Width == 0 && counter > 0 {
				y++
			}
			counter++
		}
	}
	fmt.Printf("Generated stamp from tiled:\n%v", stamp.DrawCoords)

	return stamp
}

// loadTiled is a helper that unmarshals a file into a Tiled struct
func loadTiled(filepath string) (*Tiled, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return &Tiled{}, err
	}

	tiled := &Tiled{}
	err = json.Unmarshal(file, tiled)
	if err != nil {
		return &Tiled{}, err
	}
	return tiled, nil
}
