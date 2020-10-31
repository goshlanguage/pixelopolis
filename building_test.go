package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	filePath := "assets/buildings/slum/1.json"
	tiled, err := loadTiled(filePath)
	assert.NoError(t, err)

	assert.NotEqual(t, tiled, &Tiled{}, "Got an empty struct")
	assert.Equal(t, tiled.TileHeight, 16)
	assert.Equal(t, tiled.TileWidth, 16)
}

func TestPaseTiled(t *testing.T) {
	filePath := "assets/buildings/slum/1.json"
	tiled, err := loadTiled(filePath)
	assert.NoError(t, err)

	assert.NotEqual(t, tiled, &Tiled{}, "Got an empty struct")
	assert.Equal(t, tiled.TileHeight, 16)
	assert.Equal(t, tiled.TileWidth, 16)

	stamp := ParseTiled(tiled)
	assert.Equal(t, DrawCoord{75, 0, 0}, stamp.DrawCoords[0])
	assert.Equal(t, DrawCoord{77, 16, 0}, stamp.DrawCoords[1])
	assert.Equal(t, DrawCoord{77, 32, 0}, stamp.DrawCoords[2])

	assert.Equal(t, stamp.DrawCoords[8], DrawCoord{107, 0, 16})
	t.Errorf("%v", stamp.DrawCoords)
}
