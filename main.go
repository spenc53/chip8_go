package main

import (
	"chip_8/chip"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var KeyMap = map[ebiten.Key]uint8{
	ebiten.Key1: 0x1,
	ebiten.Key2: 0x2,
	ebiten.Key3: 0x3,
	ebiten.Key4: 0xC,

	ebiten.KeyQ: 0x4,
	ebiten.KeyW: 0x5,
	ebiten.KeyE: 0x6,
	ebiten.KeyR: 0xD,

	ebiten.KeyA: 0x7,
	ebiten.KeyS: 0x8,
	ebiten.KeyD: 0x9,
	ebiten.KeyF: 0xE,

	ebiten.KeyZ: 0xA,
	ebiten.KeyX: 0x0,
	ebiten.KeyC: 0xB,
	ebiten.KeyV: 0xF,

	ebiten.KeyArrowUp:   0xE,
	ebiten.KeyArrowDown: 0xF,
	ebiten.KeyEnter:     0xA,
}

type Game struct {
	chip *chip.Chip8
}

func (g *Game) Update() error {
	// keys := []ebiten.Key{}
	// keys = inpututil.AppendPressedKeys(keys)
	// if len(keys) != 0 {
	// 	key := keys[0]
	// 	if val, ok := KeyMap[key]; ok {
	// 		g.chip.SetKeyPressed(val)
	// 	} else {
	// 		g.chip.ClearKeyPress()
	// 	}
	// } else {
	// 	g.chip.ClearKeyPress()
	// }
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// ebitenutil.DebugPrint(screen, "Hello, World!")
	pixels := g.chip.GetPixels()
	for row, pixelRow := range pixels {
		for col, pix := range pixelRow {
			vector.DrawFilledRect(screen, float32(col)*10, float32(row)*10, 10, 10, pix, false)
			g.chip.UpdateScreen()
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 320
}

func main() {

	chip := chip.InitializeChip(func() []uint8 {
		ks := []ebiten.Key{}
		ks = inpututil.AppendPressedKeys(ks)
		keys := []uint8{}
		for _, k := range ks {
			if val, ok := KeyMap[k]; ok {
				keys = append(keys, val)
			}
		}
		return keys
	})

	go func() {
		chip.Run()
	}()

	ebiten.SetWindowSize(640, 320)
	ebiten.SetWindowTitle("CHIP-8")
	if err := ebiten.RunGame(&Game{
		chip: chip,
	}); err != nil {
		log.Fatal(err)
	}

}
