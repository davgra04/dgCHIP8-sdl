package sdlio

import (
	"strconv"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

// chip8KeyToQWERTY maps CHIP8 keys to QWERTY keyboard keys
var chip8KeyToQWERTY map[string]string = map[string]string{
	"1": "1",
	"2": "2",
	"3": "3",
	"C": "4",
	"4": "Q",
	"5": "W",
	"6": "E",
	"D": "R",
	"7": "A",
	"8": "S",
	"9": "D",
	"E": "F",
	"A": "Z",
	"0": "X",
	"B": "C",
	"F": "V",
}

var chip8Keys string = "123C456D789EA0BF"

// QWERTYToChip8Key maps QWERTY keyboard keys to CHIP8 keys
var QWERTYToChip8Key map[string]string = map[string]string{
	"1": "1",
	"2": "2",
	"3": "3",
	"4": "C",
	"Q": "4",
	"W": "5",
	"E": "6",
	"R": "D",
	"A": "7",
	"S": "8",
	"D": "9",
	"F": "E",
	"Z": "A",
	"X": "0",
	"C": "B",
	"V": "F",
}

var qwertyKeys string = "1234QWERASDFZXCV"

// HandleKey sets the CHIP8 input state according to the event
func HandleKey(ctx *SDLAppContext, t *sdl.KeyboardEvent) {
	// fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
	// 	t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)

	// handle CHIP8 keypresses
	if chipKey, ok := QWERTYToChip8Key[strings.ToUpper(string(t.Keysym.Sym))]; ok {
		keyIdx, _ := strconv.ParseInt(chipKey, 16, 64)
		if t.Type == sdl.KEYDOWN {
			// ctx.Chip8.Keys[keyIdx] = true
			ctx.Chip8.SetKeyState(uint8(keyIdx), true)
		} else {
			// ctx.Chip8.Keys[keyIdx] = false
			ctx.Chip8.SetKeyState(uint8(keyIdx), false)
		}
	} else { // handle pause/step emulation
		if t.Type == sdl.KEYDOWN {
			if t.Keysym.Sym == sdl.K_k && t.Repeat == 0 {
				// pause emulation
				ctx.Chip8.Paused = !ctx.Chip8.Paused
			} else if t.Keysym.Sym == sdl.K_l {
				// step emulation
				if ctx.Chip8.Paused {
					ctx.Chip8.StepEmulation()
				}
			}
		}
	}

}
