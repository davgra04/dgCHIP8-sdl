package sdlio

import (
	"fmt"
	"math/bits"
	"strconv"

	"github.com/davgra04/dgCHIP8-go/chip8"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

////////////////////////////////////////////////////////////////////////////////
// SDLAppContext
////////////////////////////////////////////////////////////////////////////////

// SDLAppContext holds references to everything needed for the SDL version of the
// emulator (configuration, window, chip8 machine, etc)
type SDLAppContext struct {
	WinCfg    *WindowConfig
	Window    *sdl.Window
	Renderer  *sdl.Renderer
	Font      *ttf.Font
	BeepChunk *mix.Chunk
	Chip8     *chip8.CHIP8
}

////////////////////////////////////////////////////////////////////////////////
// WindowConfig
////////////////////////////////////////////////////////////////////////////////

// WindowConfig represents the SDL window configuration
type WindowConfig struct {
	RefreshRate float32 // display refresh rate in Hz
	Width       int32
	Height      int32
	PixelSize   int32
	MainColor   sdl.Color
	TextColor   sdl.Color
}

// GetDefaultWindowConfig returns the default SDL window configuration
func GetDefaultWindowConfig() *WindowConfig {

	var pixelSize int32 = 10
	// var w int32 = 64*(pixelSize-1) + 63
	var w int32 = 655
	var h int32 = 715

	fmt.Printf("window size: %v x %v\n", w, h)

	return &WindowConfig{
		RefreshRate: 60,
		Width:       w,
		Height:      h,
		PixelSize:   pixelSize,
		MainColor:   sdl.Color{R: 0x00, G: 0xc0, B: 0xd3},
		TextColor:   sdl.Color{R: 0xff, G: 0xff, B: 0xff},
	}
}

////////////////////////////////////////////////////////////////////////////////
// drawing functions
////////////////////////////////////////////////////////////////////////////////

// DrawWindow draws the entire window, including chip8 display and machine state
func DrawWindow(ctx *SDLAppContext) {
	// clear screen
	ctx.Renderer.SetDrawColor(0, 0, 0, 255)
	ctx.Renderer.Clear()

	// draw things
	DrawCHIP8Display(ctx)
	DrawCHIP8MachineState(ctx)

	// update
	ctx.Window.UpdateSurface()
}

func drawKeypad(ctx *SDLAppContext, xOffset, yOffset int32, qwerty bool) {

	size := int32(30)
	gap := int32(4)

	label := "KEYPAD:"
	if qwerty {
		label = "QWERTY:"
	}
	RenderText(ctx, label, xOffset, yOffset, ctx.WinCfg.TextColor)

	for ix := int32(0); ix < 4; ix++ {
		for iy := int32(0); iy < 4; iy++ {
			x := xOffset + 10 + ix*(size+gap)
			y := yOffset + 24 + iy*(size+gap)
			key := string(chip8Keys[ix+4*iy])
			keyIdx, _ := strconv.ParseInt(key, 16, 64)
			keyLabel := key
			if qwerty {
				keyLabel = chip8KeyToQWERTY[key]
			}
			rect := sdl.Rect{
				X: x,
				Y: y,
				W: size,
				H: size,
			}

			if ctx.Chip8.Keys[keyIdx] {
				ctx.Renderer.SetDrawColor(0x00, 0xc0, 0xd3, 0xff)
				ctx.Renderer.FillRect(&rect)
			} else {
				ctx.Renderer.SetDrawColor(0xff, 0xff, 0xff, 0xff)
				ctx.Renderer.DrawRect(&rect)
			}
			RenderText(ctx, keyLabel, x+10, y+4, ctx.WinCfg.TextColor)
		}
	}
}

func drawCHIP8Keypad(ctx *SDLAppContext, xOffset, yOffset int32) {
	drawKeypad(ctx, xOffset, yOffset, false)
}

func drawQWERTYKeypad(ctx *SDLAppContext, xOffset, yOffset int32) {
	drawKeypad(ctx, xOffset, yOffset, true)
}

// DrawCHIP8MachineState draws the state of the registers/stack/memory.
func DrawCHIP8MachineState(ctx *SDLAppContext) {

	// // temp draw bounding box
	// ////////////////////////////////////////
	// boundRect := sdl.Rect{
	// 	X: 8,
	// 	Y: 335,
	// 	W: 639,
	// 	H: 372,
	// }

	// ctx.Renderer.SetDrawColor(0x00, 0xc0/2, 0xd3/2, 0xff/2)
	// ctx.Renderer.DrawRect(&boundRect)

	// ctx.Renderer.DrawLine(152, boundRect.Y, 152, boundRect.Y+boundRect.H)
	// ctx.Renderer.DrawLine(288, boundRect.Y, 288, boundRect.Y+boundRect.H)
	// ctx.Renderer.DrawLine(462, boundRect.Y, 462, boundRect.Y+boundRect.H)

	// segments := int32(4)
	// for i := int32(1); i < segments; i++ {
	// 	// vertical line
	// 	ctx.Renderer.DrawLine(int32(boundRect.X+boundRect.W*i/segments), boundRect.Y,
	// 		int32(boundRect.X+boundRect.W*i/segments), boundRect.Y+boundRect.H)
	// }

	// draw keypad
	////////////////////////////////////////
	xOffset := int32(8)
	yOffset := int32(330)
	drawCHIP8Keypad(ctx, xOffset, yOffset)
	drawQWERTYKeypad(ctx, xOffset, yOffset+170)

	// draw cycle count and frequency
	////////////////////////////////////////
	// RenderText(ctx, keyLabel, x+10, y+4, ctx.WinCfg.TextColor)
	RenderText(ctx, fmt.Sprintf("CYCLE %d", ctx.Chip8.Cycle), xOffset, yOffset+342, ctx.WinCfg.TextColor)
	RenderText(ctx, fmt.Sprintf("FREQ %.1f HZ", ctx.Chip8.Cfg.ClockFreq), xOffset, yOffset+360, ctx.WinCfg.TextColor)

	// draw registers
	////////////////////////////////////////
	xOffset = 160
	RenderText(ctx, "REGISTERS:", xOffset, yOffset, ctx.WinCfg.TextColor)
	for i := 0; i < 16; i++ {
		text := fmt.Sprintf("REG_%x 0x%02x", i, ctx.Chip8.Reg[i])
		RenderText(ctx, text, xOffset+10, int32(yOffset+18*int32(i+1)), ctx.WinCfg.TextColor)
	}
	RenderText(ctx, fmt.Sprintf("REG_I 0x%04x", ctx.Chip8.RegI), xOffset+10, yOffset+324, ctx.WinCfg.TextColor)
	RenderText(ctx, fmt.Sprintf("DELAY 0x%02x", ctx.Chip8.RegDelay), xOffset+10, yOffset+342, ctx.WinCfg.TextColor)
	RenderText(ctx, fmt.Sprintf("SOUND 0x%02x", ctx.Chip8.RegSound), xOffset+10, yOffset+360, ctx.WinCfg.TextColor)

	// draw stack
	////////////////////////////////////////
	xOffset = 294
	RenderText(ctx, "STACK:", xOffset, yOffset, ctx.WinCfg.TextColor)
	for i := 0; i < 16; i++ {
		text := fmt.Sprintf("0x%x 0x%04x", i, ctx.Chip8.Stack[i])
		if i == int(ctx.Chip8.StackPtr) {
			text += " ←HEAD"
			RenderText(ctx, text, xOffset+10, int32(yOffset+18*int32(i+1)), ctx.WinCfg.MainColor)
		} else {
			RenderText(ctx, text, xOffset+10, int32(yOffset+18*int32(i+1)), ctx.WinCfg.TextColor)
		}
	}

	// draw program
	////////////////////////////////////////
	xOffset = 468
	numWords := 20
	RenderText(ctx, "PROGRAM:", xOffset, yOffset, ctx.WinCfg.TextColor)
	for i := 0; i < numWords; i++ {
		curAddr := int(ctx.Chip8.PC) + 2*(i-numWords/2)

		if curAddr >= 0 {

			text := fmt.Sprintf("0x%04x 0x%04x", curAddr, ctx.Chip8.ReadShort(uint16(curAddr)))
			if i == numWords/2 {
				text += " ←PC"
				RenderText(ctx, text, xOffset+10, int32(yOffset+18*int32(i+1)), ctx.WinCfg.MainColor)
			} else {
				RenderText(ctx, text, xOffset+10, int32(yOffset+18*int32(i+1)), ctx.WinCfg.TextColor)
			}
		}
	}

}

// RenderText draws the given text to the window at the specified x, y coordinate
func RenderText(ctx *SDLAppContext, msg string, x, y int32, color sdl.Color) {

	// var color uint32 = 0xff00c0d3

	// colorMain := ctx.WinCfg.MainColor
	// colorWhite := sdl.Color{R: 0xff, G: 0xff, B: 0xff}

	textSurface, err := ctx.Font.RenderUTF8Blended(msg, color)
	if err != nil {
		fmt.Printf("RenderUTF8Solid Error: %v\n", err)
		return
	}
	defer textSurface.Free()

	textTexture, err := ctx.Renderer.CreateTextureFromSurface(textSurface)
	if err != nil {
		fmt.Printf("CreateTextureFromSurface Error: %v\n", err)
		return
	}
	defer textTexture.Destroy()

	textRect := sdl.Rect{X: x, Y: y, W: textSurface.W, H: textSurface.H}

	ctx.Renderer.Copy(textTexture, nil, &textRect)
}

// DrawCHIP8Display reads the CHIP8 display memory and draws into the SDL window.
func DrawCHIP8Display(ctx *SDLAppContext) {

	// var pixelSize int32 = 14
	pixelSize := ctx.WinCfg.PixelSize

	// var color uint32 = 0xff000000 + uint32(ctx.WinCfg.MainColor.R)<<16 + uint32(ctx.WinCfg.MainColor.G)<<8 + uint32(ctx.WinCfg.MainColor.B)

	mainColor := ctx.WinCfg.MainColor

	for byteIdx, byte := range ctx.Chip8.Display {
		// fmt.Printf("    byte %03d:  %v\n", byteIdx, byte)
		byte = bits.Reverse8(byte)

		for i := 0; i < 8; i++ {
			x := int32((byteIdx*8 + i) % ctx.Chip8.Cfg.ResolutionX)
			y := int32((byteIdx*8 + i) / ctx.Chip8.Cfg.ResolutionX)
			// rect := sdl.Rect{X: x * pixelSize, Y: y * pixelSize, W: pixelSize, H: pixelSize}
			rect := sdl.Rect{X: x*pixelSize + 8, Y: y*pixelSize + 8, W: pixelSize - 1, H: pixelSize - 1}
			if byte%2 == 1 {
				// surface.FillRect(&rect, color)
				ctx.Renderer.SetDrawColor(mainColor.R, mainColor.G, mainColor.B, mainColor.A)
			} else {
				// surface.FillRect(&rect, 0xff222222)
				ctx.Renderer.SetDrawColor(0x22, 0x22, 0x22, 0xff)
			}
			ctx.Renderer.FillRect(&rect)
			byte /= 2
		}

	}

}
