package sdlio

import (
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

// LoadBeep returns a usable SDL_mix chunk from the bundled WAV
func LoadBeep() (*mix.Chunk, error) {
	beepRWops, err := sdl.RWFromMem(beepWav)
	if err != nil {
		return nil, err
	}

	// chunk, err := mix.QuickLoadWAV(beepWav)
	chunk, err := mix.LoadWAVRW(beepRWops, true)
	if err != nil {
		return nil, err
	}

	return chunk, nil
}

// HandleBeepEvent starts or stops playing the CHIP8's beep
func HandleBeepEvent(ctx *SDLAppContext, beep bool) {
	if beep {
		if mix.Playing(-1) == 0 {
			ctx.BeepChunk.Play(0, -1)
		}
	} else {
		if mix.Playing(-1) != 0 {
			// mix.Pause(-1)
			mix.HaltChannel(0)
		}
	}
}
