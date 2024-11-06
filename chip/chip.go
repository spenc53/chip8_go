package chip

import (
	"image/color"
	"log"
	"math/rand"
	"time"
)

type KeysPressed func() []uint8

const (
	ZERO     = 0x0
	JUMP     = 0x1
	CALL     = 0x2
	SKP_EQN  = 0x3
	SKP_NEQN = 0x4
	SKP_EQ   = 0x5
	SET_VX   = 0x6
	ADD_VX   = 0x7
	LOG_ARTH = 0x8
	SKP_NEQ  = 0x9
	SET_I    = 0xA
	JMP_OFF  = 0xB
	RAND     = 0xC
	DRAW     = 0xD
	KEY      = 0xE
	AUX      = 0xF // controls AUX things (timers, chars)

	// LOGICAL AND ARITH OPS
	SET       = 0x0
	OR        = 0x1
	AND       = 0x2
	XOR       = 0x3
	ADD       = 0x4
	SUB_VX    = 0x5
	SUB_VY    = 0x7
	SHFT_LFT  = 0xE
	SHFT_RGHT = 0x6

	// AUX (0xF instructions)
	SAV_DEL_TMR = 0x07
	SET_DEL_TMR = 0x15
	SET_SND_TMR = 0x18
	ADD_TO_IDX  = 0x1E
	GET_KEY     = 0x0A
	FONT_CHR    = 0x29
	BIN_TO_DEC  = 0x33
	STR_MEM     = 0x55
	LOD_MEM     = 0x65
)

type Chip8 struct {
	screen      *Screen
	cpu         *CPU
	mem         *Memory
	keysPressed KeysPressed
}

type instruction struct {
	Ins uint16
	I   uint8
	X   uint8
	Y   uint8
	N   uint8
	NN  uint8
	NNN uint16
}

func InitializeChip(keysPressed KeysPressed) *Chip8 {
	screen := NewScreen(64, 32)
	chip := &Chip8{
		screen: screen,
		cpu: &CPU{
			PC: 0x200,
		},
		mem:         InitializeMemory(),
		keysPressed: keysPressed,
	}
	chip.mem.LoadRom()

	return chip
}

func (chip *Chip8) Run() {
	go func() {
		for {
			timer := time.NewTimer(time.Second / 60)
			if chip.cpu.DelayTimer > 0 {
				chip.cpu.DelayTimer--
			}
			// chip.cpu.SoundTimer--
			<-timer.C
		}
	}()

	timer := time.NewTimer(time.Second / 1000)
	for {
		ins := chip.FetchDecode()
		chip.Execute(ins)
		<-timer.C
		timer.Reset(time.Second / 1000)
	}
}

func (chip *Chip8) FetchDecode() instruction {
	ins := chip.fetch()
	return decode(ins)
}

func (chip *Chip8) Execute(ins instruction) {
	// log.Printf("RUNNING INSTRUCTION: 0x%04x", ins.Ins)
	switch ins.I {
	case ZERO:
		if ins.NN == 0xE0 { // clear
			chip.screen.clear()
		} else if ins.NN == 0xEE { // return
			chip.cpu.PC = chip.cpu.Stack[len(chip.cpu.Stack)-1]
			chip.cpu.Stack = chip.cpu.Stack[:len(chip.cpu.Stack)-1]
		}
	case CALL:
		chip.cpu.Stack = append(chip.cpu.Stack, chip.cpu.PC)
		chip.cpu.PC = ins.NNN
	case JUMP:
		chip.cpu.PC = ins.NNN
	case SET_VX:
		chip.cpu.Reg[ins.X] = ins.NN
	case ADD_VX:
		chip.cpu.Reg[ins.X] = chip.cpu.Reg[ins.X] + ins.NN
	case SET_I:
		chip.cpu.I = ins.NNN
	case SKP_EQN:
		if chip.cpu.Reg[ins.X] == ins.NN {
			chip.cpu.PC++
			chip.cpu.PC++
		}
	case SKP_NEQN:
		if chip.cpu.Reg[ins.X] != ins.NN {
			chip.cpu.PC++
			chip.cpu.PC++
		}
	case SKP_EQ:
		if chip.cpu.Reg[ins.X] == chip.cpu.Reg[ins.Y] {
			chip.cpu.PC++
			chip.cpu.PC++
		}
	case SKP_NEQ:
		if chip.cpu.Reg[ins.X] != chip.cpu.Reg[ins.Y] {
			chip.cpu.PC++
			chip.cpu.PC++
		}
	case JMP_OFF:
		chip.cpu.PC = ins.NNN + uint16(chip.cpu.Reg[0x0])
	case RAND:
		r := rand.Int()
		chip.cpu.Reg[ins.X] = ins.NN & uint8(r)
	case LOG_ARTH:
		x := chip.cpu.Reg[ins.X]
		y := chip.cpu.Reg[ins.Y]
		switch ins.N {
		case SET:
			chip.cpu.Reg[ins.X] = y
		case OR:
			chip.cpu.Reg[ins.X] = x | y
		case AND:
			chip.cpu.Reg[ins.X] = x & y
		case XOR:
			chip.cpu.Reg[ins.X] = x ^ y
		case ADD:
			chip.cpu.Reg[ins.X] = x + y
			chip.cpu.Reg[0xF] = 0
			if int(x)+int(y) >= 255 {
				chip.cpu.Reg[0xF] = 1
			}
		case SUB_VX:
			chip.cpu.Reg[ins.X] = x - y
			chip.cpu.Reg[0xF] = 1
			if y > x {
				chip.cpu.Reg[0xF] = 0
			}
		case SUB_VY:
			chip.cpu.Reg[ins.X] = y - x
			chip.cpu.Reg[0xF] = 1
			if x > y {
				chip.cpu.Reg[0xF] = 0
			}
		case SHFT_LFT:
			chip.cpu.Reg[ins.X] = x << 1
			chip.cpu.Reg[0xF] = (x >> 0x7)
		case SHFT_RGHT:
			chip.cpu.Reg[ins.X] = x >> 1
			chip.cpu.Reg[0xF] = x & 0x1
		default:
			log.Printf("UNKNOWN INSTRUCTION: 0x%04x", ins.Ins)
		}
	case DRAW:
		xPos := chip.cpu.Reg[ins.X] % 64
		yPos := chip.cpu.Reg[ins.Y] % 32
		height := ins.N
		chip.cpu.Reg[0xF] = 0

		for row := uint8(0); row < height; row++ {
			spriteByte := chip.mem.data[chip.cpu.I+uint16(row)]

			for col := uint8(0); col < 8; col++ {
				spritePixel := spriteByte & (0x80 >> col)
				pixRow := int(yPos + row)
				pixCol := int(xPos + col)
				if pixRow >= 32 || pixRow < 0 {
					continue
				}
				if pixCol >= 64 || pixCol < 0 {
					continue
				}

				if spritePixel != 0 {
					if chip.screen.IsPixelOn(pixRow, pixCol) {
						chip.cpu.Reg[0xF] = 1
						chip.screen.SetPixel(pixRow, pixCol, false)
					} else {
						chip.screen.SetPixel(pixRow, pixCol, true)
					}
				}

			}
		}
	case KEY:
		switch ins.NN { // these don't wait for key to be pressed.
		case 0x9E: // skip if key == VX
			keys := chip.keysPressed()
			for _, key := range keys {
				if key == chip.cpu.Reg[ins.X] {
					chip.cpu.PC++
					chip.cpu.PC++
					break
				}
			}
		case 0xA1: // skip if key != VX
			keys := chip.keysPressed()
			for _, key := range keys {
				if key != chip.cpu.Reg[ins.X] {
					chip.cpu.PC++
					chip.cpu.PC++
				}
			}
		default:
		}
	case AUX:
		switch ins.NN {
		case SAV_DEL_TMR:
			chip.cpu.Reg[ins.X] = chip.cpu.DelayTimer
		case SET_DEL_TMR:
			chip.cpu.DelayTimer = chip.cpu.Reg[ins.X]
		case SET_SND_TMR:
			chip.cpu.SoundTimer = chip.cpu.Reg[ins.X]
		case ADD_TO_IDX:
			chip.cpu.I = chip.cpu.I + uint16(chip.cpu.Reg[ins.X])
		case GET_KEY:
			keys := chip.keysPressed()
			if len(keys) == 0 {
				chip.cpu.PC--
				chip.cpu.PC--
				return
			}
			chip.cpu.Reg[ins.X] = keys[0]
		case FONT_CHR:
			hex := uint16(chip.cpu.Reg[ins.X] & 0xF)
			chip.cpu.I = 5*(hex) + FONT_START
		case BIN_TO_DEC:
			x := chip.cpu.Reg[ins.X]
			digs := []uint8{}

			// ex 100
			// 100 -> 10 -> 1 -> 0
			// [] -> [0] -> [0, 0] -> [0, 0, 1]
			for x != 0 {
				digs = append(digs, x%10)
				x = x / 10
			}

			for idx := range digs {
				d := digs[len(digs)-1-idx]
				chip.mem.data[int(chip.cpu.I)+idx] = d
			}

		case STR_MEM:
			i := chip.cpu.I
			for r := uint8(0); r <= ins.X; r++ {
				chip.mem.data[i+uint16(r)] = chip.cpu.Reg[r]
			}
		case LOD_MEM:
			i := chip.cpu.I
			for r := uint8(0); r <= ins.X; r++ {
				chip.cpu.Reg[r] = chip.mem.data[i+uint16(r)]
			}
		}
	default:
		log.Printf("UNKNOWN INSTRUCTION: 0x%04x", ins.Ins)
	}
}

func (chip *Chip8) fetch() uint16 {
	ins := uint16(chip.mem.data[chip.cpu.PC])<<8 | uint16(chip.mem.data[chip.cpu.PC+1])
	chip.cpu.PC++
	chip.cpu.PC++
	return ins
}

func decode(ins uint16) instruction {
	return instruction{
		Ins: ins,
		I:   uint8((ins >> 12) & 0xF),
		X:   uint8((ins >> 8) & 0xF),
		Y:   uint8((ins >> 4) & 0xF),
		N:   uint8((ins) & 0xF),
		NN:  uint8((ins) & 0xFF),
		NNN: uint16((ins) & 0xFFF),
	}
}

func (chip *Chip8) UpdateScreen() {
	chip.screen.UpdateScreen()
}

func (chip *Chip8) GetPixels() [][]color.Color {
	return chip.screen.Pixels
}
