package sdlio

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// LoadFont returns a usable SDL Font from the bundled TTF
func LoadFont(size int) (*ttf.Font, error) {
	fontRWops, err := sdl.RWFromMem(fontData)
	if err != nil {
		return nil, err
	}

	font, err := ttf.OpenFontRW(fontRWops, 1, size)
	if err != nil {
		return nil, err
	}

	return font, nil
}
