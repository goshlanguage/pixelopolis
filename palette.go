package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Brush represents a paintbrush's coordinates and w,h in an image full of textures
type Brush struct {
	XPos, YPos, Width, Height float32
}

// Palette represents a bunch of brushes in a texture
type Palette struct {
	Brushes map[int]Brush
	// Filepath to sprite sheet
	Spritesheet string
	Texture     rl.Texture2D
	// Set width and Height so we can automatically parse our Texture into a tilesheet
	Width, Height         int
	TileWidth, TileHeight int
}

// NewPalette is a factory that takes a filepath to a tilesheet, and the tilesheet's tile width and height
func NewPalette(filepath string, tileHeight, tileWidth int) *Palette {
	img := rl.LoadImage(filepath)

	// Declare brushMap, cover any sort of possible edge case here. We start at 1 instead of 0 due to the
	// same behavior being embedded in Tiled in its tileset data handling
	brushMap := make(map[int]Brush)
	brushMap[0] = Brush{0, 0, 0, 0}

	xCycles := img.Width / int32(tileWidth)
	yCycles := img.Height / int32(tileHeight)
	length := xCycles * yCycles
	x := 0
	y := 0
	for i := 0; i < int(length); i++ {
		if i > 0 {
			x = i % tileWidth
		}
		brushMap[i] = Brush{float32(x * tileWidth), float32(y * tileHeight), float32(tileWidth), float32(tileHeight)}

		if x > 0 && (x+1)%tileWidth == 0 {
			y++
		}
		x++
	}

	return &Palette{
		Brushes:     brushMap,
		Spritesheet: filepath,
		Texture:     rl.LoadTextureFromImage(img),
		Width:       int(img.Width),
		Height:      int(img.Height),
		TileHeight:  tileHeight,
		TileWidth:   tileWidth,
	}
}

// Draw uses the given brush at an X,Y point
func (p *Palette) Draw(brush, x, y int) {
	rectangle := rl.NewRectangle(p.Brushes[brush].XPos, p.Brushes[brush].YPos, p.Brushes[brush].Width, p.Brushes[brush].Height)
	position := rl.NewVector2(float32(x), float32(y))
	rl.DrawTextureRec(p.Texture, rectangle, position, rl.White)
}

// Update loads the brushes into textures so they can be drawn to the sceen
func (p *Palette) Update() {}

// DrawCoord represents a coordinate pair
type DrawCoord struct {
	Brush            int
	XOffset, YOffset float32
}

// Stamp is a grouping of brushes preassembled to represent objects
type Stamp struct {
	DrawCoords                    []DrawCoord
	Palette                       *Palette
	LevelX, LevelY, Width, Height float32
}

// Draw renders the stamp from the given x, y coordinates
func (s *Stamp) Draw() {
	for _, i := range s.DrawCoords {
		s.Palette.Draw(i.Brush, int(s.LevelX+i.XOffset), int(s.LevelY+i.YOffset))
	}
}

// GetProjectMegaPalette takes a filepath, so we can sellect which file we want to create the palette from
func GetProjectMegaPalette(filepath string) *Palette {
	/** Refacor this... EEK
	brushMap[0] = Brush{176, 160, 16, 16} // ground
	brushMap[1] = Brush{192, 160, 16, 16} // grass
	brushMap[2] = Brush{208, 160, 16, 16} // grass
	brushMap[3] = Brush{224, 160, 16, 16} // grass
	brushMap[4] = Brush{240, 160, 16, 16} // grass
	brushMap[5] = Brush{176, 64, 16, 16}  // building top left
	brushMap[6] = Brush{208, 64, 16, 16}  // building top center
	brushMap[7] = Brush{240, 64, 16, 16}  // building top right
	brushMap[8] = Brush{176, 96, 16, 16}  // building left
	brushMap[9] = Brush{208, 96, 16, 16}  // building center (fill)
	brushMap[10] = Brush{240, 96, 16, 16} // building right
	brushMap[11] = Brush{128, 96, 16, 16} // door
	brushMap[12] = Brush{0, 96, 16, 16}   // small buildings 1
	brushMap[13] = Brush{16, 96, 16, 16}  // small buildings 2
	brushMap[14] = Brush{24, 96, 16, 16}  // small buildings 3
	brushMap[15] = Brush{64, 16, 16, 16}  // window
	brushMap[16] = Brush{32, 16, 16, 16}  // big window 1 left
	brushMap[17] = Brush{48, 16, 16, 16}  // big window 1 right

	// Building Decorations
	brushMap[18] = Brush{0, 80, 16, 16}    // banner right
	brushMap[19] = Brush{16, 80, 16, 16}   // banner left
	brushMap[20] = Brush{144, 112, 16, 16} // Building Texture 1
	brushMap[21] = Brush{160, 112, 16, 16} // Building Texture 1
	brushMap[22] = Brush{176, 112, 16, 16} // Building Texture 1
	brushMap[23] = Brush{192, 112, 16, 16} // Building Texture 1
	brushMap[24] = Brush{208, 112, 16, 16} // Building Texture 1
	brushMap[25] = Brush{224, 112, 16, 16} // Building Texture 1
	brushMap[26] = Brush{240, 112, 16, 16} // Building Texture 1
	brushMap[27] = Brush{144, 128, 16, 16} // Building Texture 1
	brushMap[28] = Brush{160, 128, 16, 16} // Building Texture 1
	brushMap[29] = Brush{176, 128, 16, 16} // Building Texture 1
	brushMap[30] = Brush{192, 128, 16, 16} // Building Texture 1
	brushMap[31] = Brush{208, 128, 16, 16} // Building Texture 1
	brushMap[32] = Brush{224, 128, 16, 16} // Building Texture 1
	brushMap[33] = Brush{240, 128, 16, 16} // Building Texture 1
	brushMap[34] = Brush{144, 144, 16, 16} // Building Texture 1
	brushMap[35] = Brush{160, 144, 16, 16} // Building Texture 1
	brushMap[36] = Brush{176, 144, 16, 16} // Building Texture 1
	brushMap[37] = Brush{192, 144, 16, 16} // Building Texture 1
	brushMap[38] = Brush{208, 144, 16, 16} // Building Texture 1
	brushMap[39] = Brush{224, 144, 16, 16} // Building Texture 1
	brushMap[40] = Brush{240, 144, 16, 16} // Building Texture 1
	**/

	palette := NewPalette(filepath, 16, 16)
	palette.Update()
	return palette
}
